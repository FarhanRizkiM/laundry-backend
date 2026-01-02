# API Contract - Laundry Management System

## General Rules

- Base URL: /api/v1
- All requests and responses use JSON
- Authentication via Authorization header
- Timestamps use ISO 8601 format

## Authentication

Authenticated endpoints require:
Authorization: Bearer <token>

## Roles:

- owner
- kasir
- staff
- courier

## Roles & Access Rules

- Owner: read-only access to summaries
- Kasir: create orders, manage payments
- Staff: update production status
- Courier: handle delivery status
- Customer: public access (no auth) for order status lookup

## Endpoints

### Auths

- POST /auth/login
<!-- Fungsi: Menukarkan username & password dengan Token Akses (JWT) | Kenapa: Pintu masuk utama. Backend harus memvalidasi apakah user ada di tabel users dan statusnya is_active = 1. -->

- POST /auth/logout
<!-- Fungsi: Menghanguskan token atau menghapus sesi aktif. | Kenapa: Keamanan standar. Memastikan akun tidak bisa digunakan lagi setelah user selesai bekerja/pulang. -->

- POST /auth/refresh-token
<!-- Fungsi: Meminta token akses baru tanpa harus input password lagi. | Kenapa: User Experience (UX). Agar Kurir yang sedang di jalan tidak tiba-tiba "logout paksa" saat token utamanya kadaluarsa (expired). -->

- GET /auth/me
<!-- Fungsi: Mengambil data profil (ID, Nama, Role) milik user yang sedang login (berdasarkan Token). | Kenapa: SANGAT PENTING. Frontend butuh data ini untuk membedakan tampilan. Jika response-nya role: courier, maka menu "Laporan Keuangan" harus disembunyikan. -->

### Users

- POST /users
<!-- Fungsi: Menambahkan karyawan baru (Create). | Kenapa: Karena sistem ini Private/Internal, Kurir, Kasir dan staff tidak bisa Register sendiri. Owner harus input manual di sini. | Logic DB: Wajib mengisi role ('cashier'/'courier'/'staff') dan password awal. -->

- GET /users
<!-- Fungsi: Mengambil daftar semua karyawan. | Kenapa: Digunakan di halaman "Kelola Staff". | Catatan Penting: Endpoint ini WAJIB mendukung Query Param ?role=courier. | Skenario: Saat Admin membuat Order Delivery, dia butuh dropdown yang isinya hanya nama-nama Kurir (bukan Kasir). Filter ini solusinya. -->

- GET /users/{id}
<!-- Fungsi: Melihat detail profil satu karyawan. | Kenapa: Untuk admin mengecek detail spesifik, misal no HP kurir saat darurat, atau cek sejak kapan bergabung. -->

- PUT /users/{id}
<!-- Fungsi: Mengupdate data karyawan. | Kenapa: Digunakan jika karyawan ganti nomor HP, ganti nama (typo), atau Admin ingin mereset password karyawan yang lupa sandi. -->

- DELETE /users/{id}
<!-- Fungsi: Menonaktifkan karyawan (Soft Delete). | Kenapa: SANGAT KRUSIAL. Kita tidak boleh melakukan DELETE FROM users (Hard Delete). | Alasan: Jika Jono (ID 4) dihapus permanen, maka ribuan data transaksi sejarah yang collected_by = 4 akan error/hilang relasinya. Solusinya: Update kolom is_active menjadi false atau isi kolom deleted_at. -->

### Service Categories

- POST /categories
<!-- Fungsi: Membuat kategori layanan baru. | Contoh: "Kiloan", "Satuan", "Karpet/Boneka", "Express". | Kenapa: Agar di menu Kasir, layanan tidak berantakan. Kasir bisa klik tab "Kiloan" untuk melihat harga cuci kiloan saja. -->

- GET /categories
<!-- Fungsi: Mengambil semua daftar kategori. | Kenapa: Digunakan oleh Frontend untuk membuat Tab atau Dropdown filter di halaman Kasir/List Harga. -->

- GET /categories/{id}
<!-- Fungsi: Melihat detail satu kategori. | Kenapa: Standar CRUD, untuk memastikan data sebelum diedit. -->

- PUT /categories/{id}
<!-- Fungsi: Mengganti nama kategori. | Kenapa: Misal owner ingin mengubah nama kategori "Satuan" menjadi "Dry Clean/Satuan". -->

- DELETE /categories/{id}
<!-- Fungsi: Menghapus kategori. | ⚠️ PERINGATAN LOGIC: Backend harus punya validasi ketat. | Aturan: Kategori TIDAK BOLEH dihapus jika masih ada Layanan (Services) di dalamnya. | Contoh: Jika kategori "Kiloan" dihapus, nasib layanan "Cuci Komplit" bagaimana? Error. Jadi harus kosongkan dulu isinya, baru boleh hapus kategorinya. -->

### Services

- POST /services
<!-- Fungsi: Menambahkan layanan baru (Menu Baru). | Contoh: Input "Cuci Komplit Kilat", Harga 10.000, Satuan 'kg', Kategori 'Express'. | Kenapa: Owner butuh input layanan baru saat ekspansi bisnis. | Logic DB: Wajib menyertakan category_id agar terkelompok dengan rapi. -->

- GET /services
<!-- Fungsi: Mengambil daftar semua layanan (Menu). | Kenapa: Endpoint ini yang paling sering dipanggil. | Skenario Kasir: Saat Kasir klik Tab "Kiloan", Frontend akan request GET /services?category_id=1. Jadi endpoint ini WAJIB support filter category_id. | Skenario Owner: Menampilkan tabel master data harga untuk dikelola. -->

- GET /services/{id}
<!-- Fungsi: Melihat detail satu layanan. | Kenapa: Standard CRUD. -->

- PUT /services/{id}
<!-- Fungsi: Update layanan (Terutama ganti harga/nama). | Logic Bisnis: | Pertanyaan: "Kalau harga diubah, order lama ikut berubah tidak?" | Jawabannya: TIDAK. Karena di tabel order_items kita menyimpan price saat transaksi terjadi (snapshot). Jadi Owner aman mengubah harga di sini kapan saja tanpa merusak laporan keuangan masa lalu. -->

- DELETE /services/{id}
<!-- Fungsi: Menghapus layanan (Menu yang sudah tidak dijual). | ⚠️ PERINGATAN LOGIC (Soft Delete): | Jangan gunakan DELETE FROM services (Hard Delete) jika layanan ini pernah dipesan orang. | Gunakan Soft Delete (is_active = 0 atau deleted_at) agar riwayat order tahun lalu yang mencuci "Paket Promo Lebaran" (yang sekarang sudah dihapus) datanya tetap bisa dibuka/tidak error. -->

### Orders

- POST /orders
<!-- Fungsi: Membuat transaksi order baru. | Input: customer_id (jika member), service_id (jenis cuci), weight/qty (berat/jumlah). | Logic Coding: | 1. Backend otomatis set status_internal = 'pending'. | 2. Backend otomatis set payment_status = 'unpaid'. | 3. Backend menghitung total_price berdasarkan harga service saat ini (Snapshot Price). | 4. Generate invoice_number unik (misal: INV/2025/01/001). -->

- GET /orders
<!-- Fungsi: Melihat daftar semua order (Dashboard). | Logic Filter: Endpoint ini harus canggih. | - ?status=pending (Admin cari order baru). | - ?payment_status=unpaid (Admin cari yang belum bayar). | - ?date=today (Owner cek orderan hari ini). | Kenapa: Tanpa filter ini, Admin akan pusing melihat ribuan data campur aduk. -->

- GET /orders/{id}
<!-- Fungsi: Melihat detail lengkap satu order (Faktur). | Output: Data Customer + Data Cucian (Item) + Total Harga + Status Terakhir. | Kenapa: Saat customer datang bawa struk, Kasir scan/input ID untuk melihat detail ini sebelum menyerahkan baju. -->

- PUT /orders/{id}
<!-- Fungsi: Mengedit Data Order (Revisi). | Bukan Ganti Status: Hanya untuk edit data yang salah input. | Contoh: Kasir salah input berat (Harusnya 5kg tertulis 4kg) atau salah pilih parfum. | Restriction: Sebaiknya dibatasi. Jika status sudah 'completed', data tidak boleh diedit lagi demi keamanan audit. -->

- PATCH /orders/{id}
<!-- Fungsi: Khusus Update STATUS Order (Workflow). | Input: status (misal: 'washing', 'ready-delivery', 'cancelled'). | Logic "Side Effect" (Sangat Penting): | Saat endpoint ini dipanggil, Backend melakukan 2 hal sekaligus (Database Transaction): | 1. Update kolom status_internal di tabel orders. | 2. Auto-Insert baris baru ke tabel status_history (mencatat siapa yang ubah & jam berapa). -->

### Payments

- POST /payments
<!-- Fungsi: Mencatat pembayaran baru yang diterima. | User: Kasir (saat customer bayar di toko). | Input: order_id, amount (nominal), payment_method (cash/qris/transfer), bank_name (opsional). | Logic Coding (Automation): | Backend WAJIB melakukan update otomatis ke tabel orders: | 1. Tambahkan nominal ini ke kolom total_paid di tabel orders. | 2. Cek matematika: Jika total_paid >= total_price, maka update payment_status order menjadi 'paid'. | 3. Jika masih kurang, biarkan 'unpaid' (atau 'partial' jika ada). -->

- GET /payments
<!-- Fungsi: Melihat riwayat daftar pembayaran yang masuk. | User: Owner (untuk Audit) & Kasir (untuk Re-check pembukuan hari ini). | Filter: Wajib ada filter tanggal (?date=today) dan metode bayar (?method=cash). | Kenapa: Owner butuh ini untuk mencocokkan fisik uang di laci kasir dengan data di sistem setiap malam (Closingan). Jika fisik uang beda dengan data di sini, berarti ada selisih. -->

- GET /payments/{id}
<!-- Fungsi: Menampilkan detail lengkap satu transaksi pembayaran spesifik. | User: Owner (Utamanya saat hendak melakukan Audit detail atau Edit). | Skenario Penggunaan: Saat Owner melihat ada angka aneh di laporan, lalu menekan tombol "Edit". Frontend WAJIB memanggil endpoint ini dulu untuk mengambil data lama (pre-fill form) agar muncul di layar sebelum diubah. | Informasi Penting: Endpoint ini harus me-return data collected_by (nama kasir) dan created_at (jam transaksi) dengan jelas untuk bukti audit. | Kenapa: Tanpa endpoint ini, fitur Edit (PUT) akan sulit berjalan karena form editnya kosong (blind edit), dan Owner tidak bisa memverifikasi detail siapa kasir yang menginput transaksi tersebut. -->

- PUT /payments/{id}
<!-- Fungsi: Mengoreksi data pembayaran (Edit). | User: KHUSUS OWNER (Kasir dilarang akses demi keamanan keuangan). | Skenario: Human Error. Misal Kasir salah input, customer bayar Rp 50.000 tapi tertulis Rp 500.000, atau salah pilih metode bayar (Tunai jadi Transfer). | Logic Side Effect (Automation): | Saat nominal diedit, Backend harus hitung ulang total_paid di tabel orders secara real-time. | 1. Kurangi total dengan nominal lama. | 2. Tambah total dengan nominal baru. | 3. Update status lunas/belum lunasnya lagi. -->

- DELETE /payments/{id}
<!-- Fungsi: Membatalkan/Menghapus pembayaran. | User: KHUSUS OWNER. | Skenario: Transaksi double input (kasir input 2x) atau customer minta refund uang (batal total). | Logic Penting (Soft Delete): | 1. Jangan hapus permanen dari database, gunakan kolom deleted_at. | 2. PENTING: Kurangi nominal yang dihapus tadi dari total_paid di tabel orders. | 3. Jika total paid jadi 0 atau kurang dari harga, kembalikan status order ke 'unpaid'. -->

### Delivery

- POST /deliveries
<!-- Fungsi: Kurir Mengambil/Klaim Order (Pick Up). | User: KHUSUS KURIR. | Logic: 1. Validasi status order harus ready-delivery. 2. Insert ke tabel deliveries. 3. Update status order jadi being-delivered. 4. Catat history: "Kurir Udin mulai pengantaran". -->

- GET /deliveries
<!-- Fungsi: Melihat daftar tugas kurir. | User: KHUSUS KURIR. | Filter: ?status=process (Yang sedang dibawa) & ?status=delivered (Riwayat kerja). -->

- GET /deliveries/{id}
<!-- Fungsi: Halaman Detail Pengiriman. | User: KHUSUS KURIR. | Kenapa: Untuk menampilkan Link Google Maps, Tombol WA Customer, dan Info Tagihan (COD). -->

- PATCH /deliveries/{id}
<!-- Fungsi: Menyelesaikan pengiriman (Finish). | User: KHUSUS KURIR (Di lokasi customer). | Logic: Upload bukti foto & Update status order jadi finished-delivery. -->

- DELETE /deliveries/{id}
<!-- Fungsi: Membatalkan tugas (Cancel Pickup). | User: KHUSUS KURIR (Hanya bisa dilakukan jika status masih 'process' / belum selesai). | Skenario: Kurir salah pencet tombol "Ambil" atau motor mogok tidak jadi antar. | Logic Coding (Revert Status): 1. Soft delete data di tabel deliveries. 2. PENTING: Kembalikan status order menjadi ready-delivery 3. Agar order tersebut muncul lagi di dashboard kurir lain dan bisa diambil orang lain. -->

### Customer (Endpoint Public Tracking)

- GET /orders/track/{invoice_number}
<!-- Fungsi: Melacak status laundry dan melihat detail tagihan secara publik. | Akses: PUBLIC (Tanpa Login / Tanpa Token JWT). | User: Customer (Pak Doni) yang memegang struk/nota. | Input: Nomor Invoice unik (Contoh: INV/2025/01/999). | Logic Backend (Data Fetching): | Backend harus mengambil paket data lengkap (All-in-One): 1. Header Order: Nama Customer, Total Harga, Status Pembayaran (Lunas/Belum). 2. Items: Baju apa saja yang dicuci (agar customer bisa cek kelengkapan). 3. Timeline: Ambil data dari tabel status_history (Dibuat -> Dicuci -> Diantar). | Logic Frontend (Mapping): Seperti diskusi kita sebelumnya, Backend mengirim status mentah (in-progress), nanti Frontend yang mengubahnya jadi kalimat "Sedang Dicuci" + Icon Mesin Cuci. -->

### Reports

- GET /reports/dashboard
- GET /reports/revenue
- GET /reports/employees
