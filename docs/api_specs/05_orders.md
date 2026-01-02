# LAUNDRY MANAGEMENT SYSTEM — API SPECIFICATION

## ORDERS MODULE SPECIFICATION

---

## Endpoint : `POST /orders`

### Headers :

- `Authorization: Bearer <access_token>` (Required)
- `Content-Type: application/json` (Required)

#### Description :

Endpoint ini digunakan oleh Kasir (Staff) atau Owner untuk mencatat transaksi baru (Membuat Tagihan). Saat endpoint ini dipanggil, Backend akan:

1. Menghitung Total Harga (Berdasarkan harga layanan di database x berat/jumlah).
2. Menghitung Estimasi Selesai otomatis (Berdasarkan durasi layanan terlama).
3. Generate Invoice Number unik (Contoh: INV/2025/01/1001).
4. Set status pembayaran default menjadi unpaid.

### Request Body :

Data input yang dikirimkan oleh Kasir. Catatan: Field estimated_ready_at dihapus karena sudah otomatis.

```json
{
  "customer_name": "Bu Rina",
  "customer_phone": "081223334444",
  "customer_address": "Jl. Mawar No 10",
  "is_delivery": 0,
  "notes": "jangan dicampur baju putih",
  "items": [
    {
      "service_id": 1,
      "item_notes": "pisahkan warna putih",
      "quantity": 1,
      "weight_kg": 2.5
    },
    {
      "service_id": 3,
      "item_notes": null,
      "quantity": 1,
      "weight_kg": 0
    }
  ]
}
```

### Responses Body :

#### ✅ 201 Created (Success)

Berhasil membuat order. Backend mengembalikan ID Order (penting untuk pembayaran) dan Total Harga.

```json
{
  "success": true,
  "message": "Order created successfully",
  "data": {
    "id": 101,
    "invoice_number": "INV/2025/01/101",
    "customer_name": "Bu Rina",
    "total_price": 17500.0,
    "payment_status": "unpaid",
    "status_internal": "pending",
    "estimated_ready_at": "2025-01-06 14:00:00",
    "created_at": "2025-01-03 14:00:00"
  }
}
```

#### ⚠️ 400 Bad Request (Validation Error)

Terjadi jika input tidak lengkap atau salah logika (Misal: Item kosong, Berat negatif).

```json
{
  "success": false,
  "message": "Items list cannot be empty",
  "data": null
}
```

#### ⚠️ 401 Unauthorized (Invalid Token)

User tidak mengirim token atau token sudah kadaluwarsa (Expired).

```json
{
  "success": false,
  "message": "Invalid or missing access token",
  "data": null
}
```

#### 🚫 403 Forbidden (Insufficient Permissions)

User login tapi tidak punya hak akses (Misal: Akun Kurir mencoba input order kasir).

```json
{
  "success": false,
  "message": "You do not have permission to perform this action",
  "data": null
}
```

#### 🔍 404 Not Found (Foreign Key Error)

Terjadi jika Kasir mengirim service_id yang tidak ada di database (Misal: ID Service 99 padahal cuma ada 1-5).

```json
{
  "success": false,
  "message": "Service with ID 99 not found",
  "data": null
}
```

#### 💥 409 Conflict (Duplicate Data)

Sangat jarang terjadi di POST, tapi bisa terjadi jika sistem gagal generate Invoice Number yang unik (Duplikat Invoice).

```json
{
  "success": false,
  "message": "Failed to generate unique invoice number, please try again",
  "data": null
}
```

#### 🔥 500 Internal Server Error

Terjadi kesalahan di sisi server (Database down).

```json
{
  "success": false,
  "message": "An unexpected error occurred",
  "data": null
}
```

---

## Endpoint : `GET /orders`

### Headers :

- `Authorization: Bearer <access_token>` (Required)

#### Description :

Endpoint ini digunakan untuk menampilkan Daftar Transaksi (Dashboard). Mendukung fitur Filter (berdasarkan status/tanggal), Pencarian (Invoice/Nama), dan Pagination (Halaman). Catatan: Jika data kosong, backend tetap mengembalikan 200 OK dengan data array kosong [], bukan 404.

### Parameters :

Parameter ini ditempel di URL (Contoh: /orders?page=1&status=pending).

| Key            | Tipe    | Default | Deskripsi                                                     |
| -------------- | ------- | ------- | ------------------------------------------------------------- |
| page           | Integer | 1       | Nomor halaman (Pagination).                                   |
| limit          | Integer | 10      | Jumlah data per halaman.                                      |
| search         | String  | null    | Cari berdasarkan Invoice Number atau Nama Customer.           |
| status         | String  | null    | Filter status proses (pending, washing, ready-delivery, dll). |
| payment_status | String  | null    | Filter pembayaran (paid, unpaid).                             |
| date           | String  | null    | Filter tanggal transaksi (Format: YYYY-MM-DD).                |
| sort_by        | String  | desc    | Urutan data (desc = Terbaru, asc = Terlama).                  |

### Request Body :

```json
None (Kosong, karena method GET tidak boleh punya body).
```

### Responses Body :

#### ✅ 200 OK (Success)

Berhasil mengambil data. Perhatikan ada field meta untuk keperluan pagination Frontend.

```json
{
  "success": true,
  "message": "Orders retrieved successfully",
  "data": [
    {
      "id": 101,
      "invoice_number": "INV/2025/01/101",
      "customer_name": "Bu Rina",
      "total_price": 17500.0,
      "payment_status": "unpaid",
      "status_internal": "pending",
      "created_at": "2025-01-03 14:00:00"
    },
    {
      "id": 100,
      "invoice_number": "INV/2025/01/100",
      "customer_name": "Pak Budi",
      "total_price": 50000.0,
      "payment_status": "paid",
      "status_internal": "washing",
      "created_at": "2025-01-03 12:00:00"
    }
  ],
  "meta": {
    "current_page": 1,
    "per_page": 10,
    "total_data": 50,
    "total_pages": 5
  }
}
```

#### ⚠️ 400 Bad Request (Validation Error)

Terjadi jika parameter filter salah format.

```json
{
  "success": false,
  "message": "Invalid date format, expected YYYY-MM-DD",
  "data": null
}
```

#### ⚠️ 401 Unauthorized (Invalid Token)

Deskripsi bagian ini

```json
{
  "success": false,
  "message": "Invalid or missing access token",
  "data": null
}
```

#### 🔥 500 Internal Server Error

Kesalahan koneksi database.

```json
{
  "success": false,
  "message": "Database connection failed",
  "data": null
}
```

---

## Endpoint : `GET /orders/{id}`

### Headers :

- `Authorization: Bearer <access_token>` (Required)

#### Description :

Endpoint ini digunakan untuk melihat Detail Lengkap dari satu transaksi spesifik berdasarkan ID. Backend akan melakukan Join Table untuk mengambil data dari tabel orders, order_items, dan payments.

### Parameters :

Parameter ini adalah Path Variable.

| Key | Tipe | Default | Deskripsi                                       |
| --- | ---- | ------- | ----------------------------------------------- |
| id  | Int  | -       | ID dari Order yang ingin dilihat (Contoh: 101). |

### Request Body :

```json
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK (Success)

Perhatikan struktur data. Di sini kita menampilkan object order, yang di dalamnya ada Array items dan Array payments.

```json
{
  "success": true,
  "message": "Order details retrieved successfully",
  "data": {
    "id": 101,
    "invoice_number": "INV/2025/01/101",
    "customer_name": "Bu Rina",
    "customer_phone": "081223334444",
    "customer_address": "Jl. Mawar No 10",
    "is_delivery": 0,
    "status_internal": "pending",
    "payment_status": "unpaid",
    "total_price": 17500.0,
    "estimated_ready_at": "2025-01-06 14:00:00",
    "notes": "jangan dicampur baju putih",
    "created_at": "2025-01-03 14:00:00",
    "items": [
      {
        "id": 1,
        "service_name": "Cuci Reguler", // Backend ambil nama dari table services
        "quantity": 1,
        "weight_kg": 2.5,
        "item_notes": "pisahkan warna putih",
        "subtotal": 17500.0
      }
    ],
    "payment_history": [] // Kosong jika belum ada pembayaran (Unpaid)
  }
}
```

Contoh jika status sudah PAID (Ada data payment):

```json
"payment_history": [
      {
        "id": 50,
        "date": "2025-01-03 14:05:00",
        "amount": 17500.00,
        "method": "cash",
        "collected_by": "Kasir Siti"
      }
    ]
```

#### ⚠️ 400 Bad Request (Invalid ID Format)

Terjadi jika ID di URL bukan angka.

```json
{
  "success": false,
  "message": "Invalid ID format",
  "data": {
    "id": "ID must be a number"
  }
}
```

#### ⚠️ 401 Unauthorized (Invalid Token)

Token tidak valid atau expired.

```json
{
  "success": false,
  "message": "Invalid or missing access token",
  "data": null
}
```

#### 🔍 404 Not Found (Data Missing)

ID valid (angka), tapi data layanan tidak ada di database.

```json
{
  "success": false,
  "message": "Order with ID 9999 not found",
  "data": null
}
```

#### 🔥 500 Internal Server Error

Terjadi kesalahan di sisi server (Database down).

```json
{
  "success": false,
  "message": "An unexpected error occurred",
  "data": null
}
```

---

## Endpoint : `PUT /orders/{id}`

### Headers :

- `Authorization: Bearer <access_token>` (Required)

#### Description :

Endpoint ini digunakan untuk memperbarui data layanan/produk. Mendukung Partial Update (boleh kirim field yang mau diubah saja). PENTING: Hanya Owner yang boleh mengakses endpoint ini.

### Parameters :

Parameter ini adalah Path Variable.

| Key | Tipe | Default | Deskripsi                          |
| --- | ---- | ------- | ---------------------------------- |
| id  | Int  | -       | ID Unik layanan yang ingin diedit. |

### Request Body :

Semua field bersifat Optional. Field category_id (jika diubah) harus valid. Field code (jika diubah) harus unik.

```json
{
  "category_id": 1,
  "code": "SVC-LKR-1",
  "service_name": "Kiloan Regular",
  "unit": "kg",
  "price": 8000, // update harga
  "is_active": true
}
```

### Responses Body :

#### ✅ 200 OK (Success)

Data layanan berhasil diperbarui.

```json
{
  "success": true,
  "message": "Service updated successfully",
  "data": {
    "id": 1,
    "category_id": 1,
    "code": "SVC-LKR-1",
    "service_name": "Kiloan Regular",
    "unit": "kg",
    "price": 8000, // hasil update harga
    "is_active": true,
    "created_at": "2025-12-26T18:49:21.410Z",
    "updated_at": "2025-12-26T18:49:21.410Z"
  }
}
```

#### ⚠️ 400 Bad Request (Validation Error)

Terjadi jika format input salah, harga negatif, atau Foreign Key (Category ID) tidak ditemukan.

```json
{
  "success": false,
  "message": "Input Validation Failed",
  "data": {
    "category_id": "Category ID does not exist", // Validasi FK
    "price": "Price cannot be negative"
  }
}
```

#### ⚠️ 401 Unauthorized (Invalid Token)

Token tidak valid atau expired.

```json
{
  "success": false,
  "message": "Invalid or missing access token",
  "data": null
}
```

#### 🚫 403 Forbidden (Insufficient Permissions)

User bukan Owner mencoba mengedit harga/layanan.

```json
{
  "success": false,
  "message": "You do not have permission to perform this action",
  "data": null
}
```

#### 404 Not Found (Service Missing)

ID Layanan yang mau diedit tidak ditemukan.

```json
{
  "success": false,
  "message": "Services not found",
  "data": null
}
```

#### 💥 409 Conflict (Duplicate Code)

Kode layanan baru bentrok dengan layanan lain.

```json
{
  "success": false,
  "message": "Service code 'SVC-LKR-1' already exists",
  "data": null
}
```

#### 🔥 500 Internal Server Error

Terjadi kesalahan di sisi server (Database down).

```json
{
  "success": false,
  "message": "An unexpected error occurred",
  "data": null
}
```

## Endpoint : `PATCH /orders/{id}`

### Headers :

- `Authorization: Bearer <access_token>` (Required)

#### Description :

Endpoint ini digunakan untuk Menonaktifkan (Soft Delete) layanan laundry. Catatan: Data tidak dihapus permanen dari database agar riwayat transaksi lama tidak error. Layanan yang dinonaktifkan tidak akan muncul lagi di menu Kasir (GET /services). Hanya Owner yang boleh melakukan ini.

### Parameters :

Parameter ini adalah Path Variable.

| Key | Tipe | Default | Deskripsi                                 |
| --- | ---- | ------- | ----------------------------------------- |
| id  | Int  | -       | ID Unik layanan yang ingin dinonaktifkan. |

### Request Body :

```json
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK (Success)

Layanan berhasil dinonaktifkan (Status is_active berubah jadi false).

```json
{
  "success": true,
  "message": "Service deactivated successfully",
  "data": null
}
```

#### ⚠️ 400 Bad Request (Invalid ID Format)

Terjadi jika ID di URL bukan angka.

```json
{
  "success": false,
  "message": "Invalid ID format",
  "data": {
    "id": "ID must be a number"
  }
}
```

#### ⚠️ 401 Unauthorized (Invalid Token)

Token tidak valid, expired, atau tidak dikirim.

```json
{
  "success": false,
  "message": "Invalid or missing access token",
  "data": null
}
```

#### 🚫 403 Forbidden (Insufficient Permissions)

User bukan Owner mencoba menghapus layanan.

```json
{
  "success": false,
  "message": "You do not have permission to perform this action",
  "data": null
}
```

#### 🔍 404 Not Found (Service Missing)

ID layanan tidak ditemukan di database.

```json
{
  "success": false,
  "message": "Service not found",
  "data": null
}
```

#### 🔥 500 Internal Server Error

Terjadi kesalahan di sisi server (Database down).

```json
{
  "success": false,
  "message": "An unexpected error occurred",
  "data": null
}
```
