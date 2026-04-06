package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"flash-ticket/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SharedState struct {
	mu    sync.Mutex
	stock map[string]int32
}

// Implementasi 3 Services
type catalogServer struct { pb.UnimplementedCatalogServiceServer; state *SharedState }
type bookingServer struct { pb.UnimplementedBookingServiceServer; state *SharedState }
type queueServer   struct { pb.UnimplementedQueueServiceServer; state *SharedState }

// SERVICE 1: Catalog - WatchLiveStock (Server-side Streaming)
func (s *catalogServer) WatchLiveStock(req *pb.EventRequest, stream pb.CatalogService_WatchLiveStockServer) error {
	log.Printf("[CATALOG] User memantau stok event: %s", req.EventId)
	for {
		s.state.mu.Lock()
		currentStock := s.state.stock[req.EventId]
		s.state.mu.Unlock()

		err := stream.Send(&pb.StockUpdate{RemainingStock: currentStock, Status: "LIVE"})
		if err != nil { return err }

		if currentStock <= 0 { break }
		time.Sleep(2 * time.Second)
	}
	return nil
}

// SERVICE 2: Booking - BookTicket (Unary)
func (s *bookingServer) BookTicket(ctx context.Context, req *pb.BookingRequest) (*pb.BookingResponse, error) {
	s.state.mu.Lock()
	defer s.state.mu.Unlock()

	// Error Handling
	if req.Quantity <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Jumlah harus > 0")
	}
	if s.state.stock[req.EventId] < req.Quantity {
		return nil, status.Error(codes.ResourceExhausted, "Tiket habis!")
	}

	// Update State
	s.state.stock[req.EventId] -= req.Quantity
	log.Printf("[BOOKING] User %s beli %d tiket. Sisa: %d", req.UserId, req.Quantity, s.state.stock[req.EventId])

	return &pb.BookingResponse{
		Success: true,
		Message: "Booking Berhasil!",
		BookingId: fmt.Sprintf("BK-%d", time.Now().Unix()),
	}, nil
}

// SERVICE 3: Queue - JoinPaymentQueue (Bi-directional Streaming)
func (s *queueServer) JoinPaymentQueue(stream pb.QueueService_JoinPaymentQueueServer) error {
	posisi := 5
	for {
		_, err := stream.Recv()
		if err == io.EOF { return nil }
		if err != nil { return err }

		err = stream.Send(&pb.QueueUpdate{
			Position: int32(posisi), 
			Message: "Menunggu pembayaran...",
		})
		if err != nil { return err }

		if posisi > 1 { posisi-- }
		time.Sleep(2 * time.Second)
	}
}

func main() {
	fmt.Println("=== FLASH-TICKET SERVER (3 SERVICES) ===")
	lis, err := net.Listen("tcp", ":50051")
	if err != nil { log.Fatal(err) }

	state := &SharedState{stock: map[string]int32{"KONSER_ROCK": 100}}
	grpcServer := grpc.NewServer()

	// Registrasi 3 Services sekaligus
	pb.RegisterCatalogServiceServer(grpcServer, &catalogServer{state: state})
	pb.RegisterBookingServiceServer(grpcServer, &bookingServer{state: state})
	pb.RegisterQueueServiceServer(grpcServer, &queueServer{state: state})

	log.Println("Server running on port :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}