package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"flash-ticket/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Gagal connect: %v", err)
	}
	defer conn.Close()

	catalogCli := pb.NewCatalogServiceClient(conn)
	bookingCli := pb.NewBookingServiceClient(conn)
	queueCli := pb.NewQueueServiceClient(conn)

	fmt.Println("=== FLASH-TICKET CLIENT ===")
	fmt.Print("Masukkan Nama: ")
	var name string
	fmt.Scanln(&name)

	for {
		fmt.Println("\nMENU:")
		fmt.Println("1. Pantau Stok Live")
		fmt.Println("2. Beli Tiket")
		fmt.Println("3. Masuk Antrean")
		fmt.Println("4. Keluar")
		fmt.Print("Pilih [1-4]: ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			fmt.Println("[CATALOG] Memantau stok (5 update)...")
			stream, _ := catalogCli.WatchLiveStock(context.Background(), &pb.EventRequest{EventId: "KONSER_ROCK"})
			for i := 0; i < 5; i++ {
				msg, err := stream.Recv()
				if err == io.EOF { break }
				fmt.Printf(">>> Stok Saat Ini: %d [%s]\n", msg.RemainingStock, msg.Status)
			}
		case 2:
			var qty int32
			fmt.Print("Jumlah tiket: ")
			fmt.Scanln(&qty)
			resp, err := bookingCli.BookTicket(context.Background(), &pb.BookingRequest{UserId: name, EventId: "KONSER_ROCK", Quantity: qty})
			if err != nil {
				fmt.Printf(">>> ERROR: %v\n", err)
			} else {
				fmt.Printf(">>> SUKSES: %s (ID: %s)\n", resp.Message, resp.BookingId)
			}
		case 3:
			fmt.Println("[QUEUE] Bergabung ke antrean...")
			stream, _ := queueCli.JoinPaymentQueue(context.Background())
			for i := 0; i < 3; i++ {
				stream.Send(&pb.QueueStatus{BookingId: name, Status: "WAITING"})
				msg, _ := stream.Recv()
				fmt.Printf(">>> Posisi Antrean: %d (%s)\n", msg.Position, msg.Message)
				time.Sleep(1 * time.Second)
			}
		case 4:
			fmt.Println("Bye!")
			os.Exit(0)
		}
	}
}