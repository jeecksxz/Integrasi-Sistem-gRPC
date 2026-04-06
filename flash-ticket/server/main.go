package main

import (
	"context"
	"fmt"
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
	mu         sync.Mutex
	stock      map[string]int32
	queueUsers []string // Daftar ID user yang sedang antre secara global
}

type catalogServer struct {
	pb.UnimplementedCatalogServiceServer
	state *SharedState
}
type bookingServer struct {
	pb.UnimplementedBookingServiceServer
	state *SharedState
}
type queueServer struct {
	pb.UnimplementedQueueServiceServer
	state *SharedState
}

// 1. Catalog Service - Server Streaming
func (s *catalogServer) WatchLiveStock(req *pb.EventRequest, stream pb.CatalogService_WatchLiveStockServer) error {
	for {
		s.state.mu.Lock()
		val := s.state.stock[req.EventId]
		s.state.mu.Unlock()

		err := stream.Send(&pb.StockUpdate{RemainingStock: val, Status: "LIVE"})
		if err != nil {
			return err
		}
		if val <= 0 {
			break
		}
		time.Sleep(2 * time.Second)
	}
	return nil
}

// 2. Booking Service - Unary
func (s *bookingServer) BookTicket(ctx context.Context, req *pb.BookingRequest) (*pb.BookingResponse, error) {
	s.state.mu.Lock()
	defer s.state.mu.Unlock()

	if req.Quantity <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Jumlah harus > 0")
	}
	if s.state.stock[req.EventId] < req.Quantity {
		return nil, status.Error(codes.ResourceExhausted, "Tiket habis!")
	}

	s.state.stock[req.EventId] -= req.Quantity
	log.Printf("[BOOKING] %s memesan %d tiket. Sisa: %d", req.UserId, req.Quantity, s.state.stock[req.EventId])

	return &pb.BookingResponse{
		Success:   true,
		Message:   "Berhasil mengamankan tiket!",
		BookingId: fmt.Sprintf("BK-%d", time.Now().Unix()),
	}, nil
}

// 3. Queue Service - Bi-directional Streaming (Real Global Queue)
func (s *queueServer) JoinPaymentQueue(stream pb.QueueService_JoinPaymentQueueServer) error {
	req, err := stream.Recv()
	if err != nil {
		return err
	}
	currentID := req.BookingId

	// Tambahkan ke antrean global
	s.state.mu.Lock()
	s.state.queueUsers = append(s.state.queueUsers, currentID)
	log.Printf("[QUEUE] %s masuk antrean. Total: %d", currentID, len(s.state.queueUsers))
	s.state.mu.Unlock()

	// Loop untuk update posisi
	for {
		select {
		case <-stream.Context().Done(): // Jika client exit/cancel context
			goto cleanup
		default:
			s.state.mu.Lock()
			posisi := -1
			for i, id := range s.state.queueUsers {
				if id == currentID {
					posisi = i + 1
					break
				}
			}
			total := len(s.state.queueUsers)
			s.state.mu.Unlock()

			if posisi == -1 {
				goto cleanup
			}

			msg := fmt.Sprintf("Menunggu... (Antrean Aktif: %d)", total)
			if posisi == 1 {
				msg = "SILAKAN BAYAR SEKARANG!"
			}

			err := stream.Send(&pb.QueueUpdate{Position: int32(posisi), Message: msg})
			if err != nil {
				goto cleanup
			}
			time.Sleep(2 * time.Second)
		}
	}

cleanup:
	s.state.mu.Lock()
	for i, id := range s.state.queueUsers {
		if id == currentID {
			s.state.queueUsers = append(s.state.queueUsers[:i], s.state.queueUsers[i+1:]...)
			break
		}
	}
	log.Printf("[QUEUE] %s keluar antrean. Sisa: %d", currentID, len(s.state.queueUsers))
	s.state.mu.Unlock()
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	state := &SharedState{
		stock:      map[string]int32{"KONSER_ROCK": 100},
		queueUsers: []string{},
	}

	s := grpc.NewServer()
	pb.RegisterCatalogServiceServer(s, &catalogServer{state: state})
	pb.RegisterBookingServiceServer(s, &bookingServer{state: state})
	pb.RegisterQueueServiceServer(s, &queueServer{state: state})

	fmt.Println("=== FLASH-TICKET SERVER STARTED (PORT 50051) ===")
	s.Serve(lis)
}