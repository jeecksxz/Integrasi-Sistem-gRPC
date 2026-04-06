# Project gRPC Integrasi Sistem


| Nama          | NRP          |
| ------------- | ------------ |
| Syela Zeruya Tandi Lalong | 5027231076 |
| Nayyara Ashila | 5027231083 |

## Cara Menjalankan Kode
Ini adalah **Tutorial Lengkap Panduan Menjalankan Projek Flash-Ticket** dari awal sampai demo. Ikuti langkah-langkah ini secara urut agar tidak ada error "package not found" atau "undefined".

### **Langkah 1: Generate Kontrak gRPC**
Buka Terminal/CMD (disarankan di VS Code), pastikan kamu berada di root folder `D:\Integrasi Sistem\flash-ticket`. Jalankan perintah ini:

```powershell
# Jalankan Compiler Protobuf
protoc --go_out=. --go-grpc_out=. pb/ticketing.proto
```
**Cek Folder PB:** Pastikan muncul file `ticketing.pb.go` dan `ticketing_grpc.pb.go`. Jika tidak muncul, periksa apakah instalasi `protoc` kamu sudah benar.

---

### **Langkah 2: Sinkronisasi Library (Wajib)**
Jalankan perintah ini untuk mendownload library gRPC dan Protobuf ke dalam projekmu:

```powershell
go mod tidy
```
*Tunggu sampai selesai. Jika berhasil, file `go.sum` akan muncul secara otomatis.*

---

### **Langkah 3: Menjalankan Server (Terminal 1)**
Tetap di root folder `flash-ticket`, jalankan server dengan perintah ini:

```powershell
go run server/main.go
```
*Tunggu sampai muncul tulisan: **"=== FLASH-TICKET SERVER STARTED (PORT 50051) ==="***.
**Jangan tutup terminal ini!** Server harus tetap menyala.

---

### **Langkah 4: Menjalankan Client & Simulasi (Terminal 2 & 3)**

Untuk melihat fitur **Multi-client** dan **Antrean Global**, kamu harus membuka dua terminal client.

#### **Terminal 2 (Client Andi):**
1. Buka terminal baru (Split terminal di VS Code).
2. Jalankan: `go run client/main.go`
3. Masukkan Nama: **Andi**
4. Pilih **Menu 1 (Pantau Stok)**. Andi akan melihat stok 100.

#### **Terminal 3 (Client Budi):**
1. Buka terminal baru lagi.
2. Jalankan: `go run client/main.go`
3. Masukkan Nama: **Budi**
4. Pilih **Menu 2 (Beli Tiket)**, masukkan angka **10**.
5. *Hasil:* Budi akan melihat "Sukses!". 
6. **Lihat Terminal Andi:** Di layar Andi, stok akan otomatis berubah dari **100 menjadi 90** tanpa Andi melakukan apapun.

---

### **Langkah 5: Simulasi Antrean Global (Fitur Paling Canggih)**

Ini adalah skenario untuk menunjukkan **Bi-directional Streaming** dan **Global State**:

1. **Andi (Terminal 2):** Pilih **Menu 3 (Masuk Antrean)**. 
   * Andi akan melihat: `POSISI: 1 | GILIRAN ANDA!`
2. **Budi (Terminal 3):** Pilih **Menu 3 (Masuk Antrean)**. 
   * Budi akan melihat: `POSISI: 2 | Menunggu...`
3. **Pamer Fitur:** Di terminal **Andi**, tekan tombol **ENTER**. 
   * Andi akan keluar dari antrean dan kembali ke Menu Utama.
4. **Lihat Terminal Budi:** Secara ajaib (real-time), status Budi akan berubah otomatis menjadi:
   `POSISI: 1 | GILIRAN ANDA!`
5. **Kesimpulan:** Dosen akan melihat bahwa server kamu secara cerdas mengelola urutan orang yang sedang aktif antre.

---

### **Langkah 6: Simulasi Error Handling**
1. Di salah satu terminal client, pilih **Menu 2 (Beli Tiket)**.
2. Masukkan jumlah tiket **500** (sedangkan stok cuma 90).
3. *Hasil:* Muncul pesan error **ResourceExhausted: Tiket habis!**.

---

### **Troubleshooting (Jika Error):**
*   **Error "package not in std":** Pastikan kamu menjalankan `go run` dari folder root (`flash-ticket`), bukan masuk ke folder `server` dulu.
*   **Error "undefined":** Pastikan kamu sudah menjalankan perintah `protoc` di Langkah 2 dengan benar.
*   **Port 50051 Busy:** Itu artinya ada server lama yang masih nyala. Tekan `CTRL+C` di terminal server lama atau restart VS Code.

**Tutorial ini sudah sangat lengkap untuk kamu ikuti saat rekaman video nanti. Selamat mencoba!**

# Penjelasan Fitur
### 1. Fitur Utama (Sisi Fungsional)

*   **Live Stock Monitoring (Catalog Service):**
    *   **Deskripsi:** User bisa melihat sisa kuota tiket yang terus berkurang secara otomatis di layar tanpa perlu melakukan *refresh*.
    *   **Kegunaan:** Memberikan informasi stok paling akurat saat kondisi "War Tiket" yang sangat sibuk.
*   **Atomic Ticket Reservation (Booking Service):**
    *   **Deskripsi:** Fitur untuk mengamankan (booking) tiket secara instan.
    *   **Kegunaan:** Menjamin user mendapatkan kuota tiket sebelum lanjut ke tahap pembayaran, mencegah *overselling* (tiket terjual melebihi stok).
*   **Real Global Waiting Room (Queue Service):**
    *   **Deskripsi:** Fitur antrean pembayaran yang bersifat global. Jika Andi masuk duluan, dia nomor 1. Jika Budi masuk kemudian, dia nomor 2. Jika Andi keluar/selesai, Budi otomatis naik jadi nomor 1.
    *   **Kegunaan:** Mengatur trafik transaksi agar server tidak *crash* dan adil bagi user yang datang lebih awal.

---

### 2. Pemenuhan Syarat Tugas (Sisi Teknis gRPC)

Ini adalah bagian yang paling dicari dosen. Kamu harus tunjukkan bahwa kodenya sudah memenuhi **6 Fitur Wajib**:

1.  **Minimal 3 Services (TERPENUHI):**
    *   Sistem dibagi menjadi `CatalogService`, `BookingService`, dan `QueueService`. Masing-masing punya tanggung jawab sendiri.
2.  **Request-Response / Unary (TERPENUHI):**
    *   Implementasi pada fungsi `BookTicket`. Client kirim pesanan, Server balas Sukses/Gagal.
3.  **Streaming gRPC (TERPENUHI - Kamu pakai 2 jenis):**
    *   **Server-side Streaming:** Pada fungsi `WatchLiveStock`. Server terus mengirimkan update stok ke banyak client sekaligus.
    *   **Bi-directional Streaming:** Pada fungsi `JoinPaymentQueue`. Client dan Server saling bertukar pesan (Heartbeat & Update Posisi) secara bersamaan dalam satu koneksi terbuka.
4.  **Error Handling (TERPENUHI):**
    *   Menggunakan gRPC Status Codes. Contoh: `codes.ResourceExhausted` jika tiket habis dan `codes.InvalidArgument` jika jumlah tiket salah.
5.  **State Management In-Memory (TERPENUHI):**
    *   Data stok dan daftar antrean disimpan dalam variabel global di RAM server menggunakan **`SharedState`**.
    *   Dilengkapi dengan **`sync.Mutex`** (Locking) agar data tidak rusak saat diakses oleh ribuan orang secara bersamaan (*Thread-safe*).
6.  **Multi-Client (TERPENUHI):**
    *   Sistem mendukung banyak terminal client sekaligus. Perubahan stok yang dilakukan Client A akan langsung terlihat di layar Client B secara real-time.

---

### 3. Keunggulan Proyek Kamu (Nilai Plus Kreativitas)

Sampaikan poin ini agar nilaimu semakin tinggi:
*   **Graceful Exit:** Client bisa keluar dari antrean dengan menekan tombol **ENTER** (menggunakan *Context Cancellation*), sehingga posisi antrean user lain di bawahnya langsung naik secara otomatis. Ini menunjukkan pemahaman mendalam tentang *Goroutine* di Golang.
*   **Reactive Update:** Tidak menggunakan *polling* (request berulang-ulang), tapi menggunakan *streaming* murni sehingga sangat hemat kuota data dan latensi sangat rendah.

---

### Saran untuk Presentasi Video (5 Menit):
1.  **Menit 1:** Judul (Flash-Ticket) & Kenalkan 3 Service-nya.
2.  **Menit 2:** Tunjukkan File `.proto` (Sebutkan mana yang Unary, mana yang Streaming).
3.  **Menit 3:** Demo Multi-Client (Jalankan 2 Client bersamaan, tunjukkan stok berkurang di kedua layar).
4.  **Menit 4:** Demo Antrean Global (Tunjukkan Budi naik posisi saat Andi keluar antrean).
5.  **Menit 5:** Tunjukkan Error Handling (Beli tiket 500 buah saat stok cuma 100).

**Semua sudah lengkap! Sekarang kamu tinggal rekam videonya. Semangat, proyek kamu sudah kategori "High-End" untuk tugas ini!** Ada lagi yang mau ditanyakan?
