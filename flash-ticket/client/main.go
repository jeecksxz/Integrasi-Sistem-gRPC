package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"flash-ticket/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
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

	fmt.Print("Masukkan Nama Anda: ")
	var name string
	fmt.Scanln(&name)

	for {
		fmt.Printf("\n--- DASHBOARD FLASH-TICKET (%s) ---\n", name)
		fmt.Println("1. Pantau Stok Tiket (Streaming)")
		fmt.Println("2. Beli Tiket (Unary)")
		fmt.Println("3. MASUK ANTREAN PEMBAYARAN (Bi-directional)")
		fmt.Println("4. Keluar")
		fmt.Print("Pilih [1-4]: ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			stream, _ := catalogCli.WatchLiveStock(context.Background(), &pb.EventRequest{EventId: "KONSER_ROCK"})
			fmt.Println("[INFO] Memantau 5 update stok...")
			for i := 0; i < 5; i++ {
				m, err := stream.Recv()
				if err == io.EOF { break }
				fmt.Printf(">>> STOK SAAT INI: %d [%s]\n", m.RemainingStock, m.Status)
			}
		case 2:
			fmt.Print("Jumlah tiket: ")
			var q int32
			fmt.Scanln(&q)
			r, err := bookingCli.BookTicket(context.Background(), &pb.BookingRequest{UserId: name, EventId: "KONSER_ROCK", Quantity: q})
			if err != nil {
				fmt.Printf(">>> GAGAL: %v\n", err)
			} else {
				fmt.Printf(">>> SUKSES: %s | ID: %s\n", r.Message, r.BookingId)
			}
		case 3:
			ctx, cancel := context.WithCancel(context.Background())
			stream, err := queueCli.JoinPaymentQueue(ctx)
			if err != nil {
				fmt.Println("Gagal join antrean:", err)
				cancel()
				continue
			}

			stream.Send(&pb.QueueStatus{BookingId: name, Status: "WAITING"})
			fmt.Println("\n[!] ANDA DALAM ANTREAN. Tekan ENTER untuk keluar ke menu utama.")

			go func() {
				for {
					m, err := stream.Recv()
					if err == io.EOF || status.Code(err) == codes.Canceled { return }
					if err != nil { return }
					fmt.Printf("\r>>> POSISI: %d | %s              ", m.Position, m.Message)
				}
			}()

			var exit string
			fmt.Scanln(&exit) // Menunggu tombol ENTER
			cancel()
			fmt.Println("\n[INFO] Kembali ke menu...")
			time.Sleep(1 * time.Second)
		case 4:
			return
		}
	}
}