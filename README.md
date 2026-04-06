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

---

# Penjelasan Fitur
### 1. Fitur Utama 

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

### 2. Teknis gRPC

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
  
Dalam sistem **Flash-Ticket** yang kita bangun, fitur yang menggunakan **Unary gRPC** (Request-Response) adalah fitur **Booking Tiket (Pemesanan Tiket)**.

Berikut adalah penjelasan detailnya untuk bahan presentasi atau videomu:

---

### 1. Nama Fitur: **Instant Booking (Reservasi Tiket)**
Ini adalah fitur utama di mana pengguna melakukan aksi "Klik Beli" untuk mengamankan tiket mereka.

*   **Cara Kerja Unary:** 
    Sama seperti protokol HTTP biasa (Request-Response), Client mengirimkan **satu** permintaan, dan Server memberikan **satu** jawaban.
    *   **Client Mengirim:** Data berupa `UserId`, `EventId` (Nama Konser), dan `Quantity` (Jumlah Tiket).
    *   **Server Membalas:** Pesan sukses/gagal, serta `BookingID` unik jika berhasil.

*   **Kenapa Harus Unary?**
    Karena transaksi pembelian bersifat **final dan atomik**. Kamu tidak butuh aliran data terus-menerus untuk satu kali aksi "Beli". User hanya butuh kepastian: *"Apakah saya dapat tiketnya atau tidak?"* di detik itu juga.

---

### 2. Fitur Pendukung: **Validasi & Error Handling**
Masih di dalam fungsi `BookTicket`, Unary gRPC digunakan untuk memberikan respon error yang instan jika terjadi kesalahan.

*   **Skenario:**
    *   Jika User memasukkan jumlah tiket **negatif (-1)** atau **nol (0)**.
    *   Jika User meminta tiket lebih banyak dari sisa stok (misal: minta 500, stok sisa 90).
*   **Respon gRPC Unary:**
    Server akan langsung memutus permintaan dan mengirimkan status error (seperti `InvalidArgument` atau `ResourceExhausted`).
*   **Kenapa Penting?**
    Agar sistem tidak membuang-buang sumber daya (bandwidth) untuk memproses permintaan yang sudah pasti salah sejak awal.

---

### 3. Peran Unary dalam "State Management"
Unary gRPC di fitur ini adalah **pemicu (trigger)** perubahan data di memori server.

*   **Prosesnya:**
    1. Client melakukan panggilan Unary `BookTicket`.
    2. Server menerima, lalu melakukan **Locking (Mutex)** pada In-memory State.
    3. Stok dikurangi.
    4. Server membalas (Response) ke Client bahwa stok berhasil dikurangi.
    5. Setelah Unary ini sukses, barulah fitur **Server-streaming** (Menu 1) mendeteksi perubahan dan mengirimkan angka stok baru ke semua user lain.

---

Agar dosenmu mudah paham, gunakan analogi ini:

| Tipe gRPC | Fitur di Projek | Analogi Dunia Nyata |
| :--- | :--- | :--- |
| **Unary** | **Beli Tiket (BookTicket)** | Kamu memesan 1 burger di kasir dan langsung mendapat struk bukti bayar. |
| **Server-Streaming** | **Pantau Stok (WatchStock)** | Kamu melihat papan menu digital yang harganya terus berubah otomatis. |
| **Bi-directional** | **Antrean (JoinQueue)** | Kamu ngobrol bolak-balik dengan satpam di pintu masuk tentang posisi antreanmu. |

---

### Ringkasan Fitur Unary untuk di Video:
> *"Fitur **Unary gRPC** kami implementasikan pada **Booking Service**. Fungsinya adalah untuk menangani transaksi pemesanan tiket secara cepat dan aman. Ketika user menekan tombol beli, client mengirimkan satu request ke server, dan server akan melakukan validasi stok di memori secara atomik, lalu mengirimkan satu response balik berupa status keberhasilan dan ID booking. Di sini kami juga menerapkan **Error Handling** untuk mencegah input yang tidak valid."*

---

### 1. Server-side Streaming (Pesan dari Server ke Client)
**Implementasi:** Fitur **Live Stock Monitoring** (`WatchLiveStock`).

*   **Apa itu?** 
    Model di mana Client mengirimkan **satu** permintaan (*Request*), dan Server membalas dengan **aliran data terus-menerus** (*Stream of Responses*) selama koneksi terbuka.
*   **Fungsi dalam Fitur:**
    Dalam fitur Pantau Stok, Client (User) hanya perlu mengeklik "Pantau Stok" sekali. Server kemudian akan mengirimkan update sisa tiket setiap kali ada perubahan atau setiap 2 detik. 
*   **Mengapa Menggunakan Ini?**
    Daripada Client harus melakukan *Refresh* atau *Polling* (bertanya berulang-ulang: "Masih ada tiket? Masih ada tiket?"), Server secara aktif "membisiki" Client data terbaru. Ini jauh lebih hemat kuota data dan beban kerja server jadi lebih ringan.

---

### 2. Bi-directional Streaming (Komunikasi Dua Arah)
**Implementasi:** Fitur **Real-time Global Waiting Room** (`JoinPaymentQueue`).

*   **Apa itu?** 
    Jenis komunikasi yang paling kompleks, di mana Client dan Server bisa saling mengirimkan aliran data secara **bersamaan** dan **mandiri** dalam satu koneksi yang sama (*Full Duplex*).
*   **Fungsi dalam Fitur:**
    Dalam fitur Antrean Pembayaran, terjadi percakapan dua arah:
    *   **Dari Client ke Server:** Client mengirimkan sinyal "Heartbeat" (Status: WAITING) untuk membuktikan bahwa user tersebut masih aktif berada di depan layar dan tidak *afk* atau *close* aplikasi.
    *   **Dari Server ke Client:** Server membalas dengan posisi antrean terbaru secara dinamis (misal: "Posisi Anda: 3... sekarang 2... sekarang 1").
*   **Mengapa Menggunakan Ini?**
    Ini sangat krusial untuk fitur antrean. Jika Andi (Posisi 1) menutup aplikasinya, Server langsung mendeteksi koneksi Andi terputus, lalu Server secara otomatis mengirimkan update posisi terbaru ke Budi (Posisi 2) agar naik menjadi Posisi 1. Komunikasi dua arah ini membuat antrean jadi sangat responsif dan adil.

---

### Perbandingan dengan Unary (Sebagai Kontras)
Sebagai perbandingan, fitur **Beli Tiket** (`BookTicket`) menggunakan **Unary gRPC**:
*   **Fungsi:** Transaksi cepat sekali jalan.
*   **Cara kerja:** Client kirim pesanan -> Server proses -> Selesai.
*   **Alasan:** Transaksi bersifat final dan tidak butuh *update* terus-menerus, jadi cukup pakai satu kali *Request-Response*.

---

| Jenis Streaming | Fitur | Fungsi Utama |
| :--- | :--- | :--- |
| **Server-side Streaming** | Live Stock Update | Server mendorong (*Push*) data stok terbaru ke semua User secara otomatis tanpa perlu *Refresh*. |
| **Bi-directional Streaming** | Global Waiting Room | Komunikasi dua arah untuk memantau kehadiran User dan meng-update posisi antrean secara adil dan *live*. |

