# LAUNDRY MANAGEMENT SYSTEM — API SPECIFICATION

## ORDERS MODULE SPECIFICATION

---

## Endpoint : `POST /orders`

### Headers :

- `Authorization`: `Bearer <access_token>`
- `Accept`: `application/json`
- `Content-Type`: `application/json`
- `Permissions`: `cashier, owner`

#### Description :

Endpoint ini digunakan untuk mencatat pesanan laundry baru ke dalam sistem. Akses dibatasi hanya untuk Cashier dan Owner. Backend akan menjalankan transaksi atomik (semua sukses atau semua gagal) untuk mengisi 5 tabel utama:

1. customers: Mengambil data pelanggan berdasarkan customer_id atau membuat data baru jika belum terdaftar.
2. orders: Menyimpan data induk pesanan (Status awal: pending, Payment: unpaid).
3. order_items: Menyimpan detail setiap layanan yang dipilih (berat, harga per unit, subtotal).
4. deliveries: Menyimpan data logistik/pengiriman jika is_delivery bernilai 1.
5. status_history: Mencatat log pertama kali pesanan dibuat sebagai riwayat pelacakan.

### Request Body :

Berisi detail pesanan, pelanggan, dan item cucian. Pada tahap ini, fokus utama adalah validasi fisik cucian (berat/jumlah) dan data pengiriman.

```json
{
  "customer_id": null,
  "customer_name": "Mpok Romlah",
  "customer_phone": "081234567890",
  "customer_address": "Jl. Merpati No. 12",
  "is_delivery": 1,
  "notes": "Jangan dicampur dengan baju luntur",
  "shipping_cost": 10000.0,
  "order_items": [
    {
      "service_id": 1, // Misal: Cuci Kiloan Reguler
      "quantity": null,
      "weight_kg": 5.0,
      "item_notes": "Pisahkan warna putih"
    }
  ]
}
```

### Responses Body :

#### ✅ 201 Created

Pesanan berhasil dibuat. Response ini mengembalikan invoice_number yang digunakan pelanggan untuk pelacakan (tracking) dan proses pembayaran selanjutnya di endpoint POST /payments.

```json
{
  "success": true,
  "message": "Order created successfully",
  "data": {
    "id": 45,
    "invoice_number": "INV-260105-001",
    "is_delivery": 1,
    "total_price": 60000.0,
    "payment_status": "unpaid",
    "status_internal": "pending",
    "estimated_ready_at": "2026-01-08 13:00:00",
    "notes": "Jangan dicampur dengan baju luntur",
    "created_by": 2,
    "created_by_name": "Siti Aminah",
    "created_at": "2026-01-05 13:00:00",
    "updated_at": null,
    "customer": {
      "id": 101,
      "name": "Mpok Romlah",
      "phone": "081234567890",
      "address": "Jl. Merpati No. 12"
    },
    "order_items": [
      {
        "id": 12,
        "service_id": 1,
        "service_name": "Cuci Kiloan Reguler",
        "item_notes": "Pisahkan warna putih",
        "quantity": null,
        "weight_kg": 5.0,
        "unit": "Kg",
        "unit_price": 10000.0,
        "subtotal": 50000.0
      }
    ],
    "delivery": {
      "id": 12,
      "shipping_cost": 10000.0,
      "courier_id": null,
      "courier_name": null,
      "courier_phone": null,
      "courier_departed_at": null,
      "courier_arrived_at": null,
      "cod_collected_amount": 0.0
    },
    "status_history": [
      {
        "id": 1,
        "previous_status": null,
        "new_status": "pending",
        "actor_name": "Siti Aminah",
        "actor_role": "cashier",
        "notes": "Initial order creation",
        "created_at": "2026-01-05 13:00:00"
      }
    ]
  }
}
```

#### ⚠️ 400 Bad Request

Terjadi ketika validasi input gagal memenuhi kriteria bisnis.

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "items.0.weight_kg": "Weight must be greater than 0 for kiloan service"
  }
}
```

#### ⚠️ 401 Unauthorized

Terjadi ketika token akses tidak valid, kedaluwarsa, atau tidak disertakan dalam header.

```json
{
  "success": false,
  "message": "Invalid or missing access token",
  "data": null
}
```

#### 🚫 403 Forbidden

Terjadi ketika user berhasil terautentikasi (token valid), namun tidak memiliki izin akses (Role) yang sesuai untuk endpoint ini.

```json
{
  "success": false,
  "message": "Your role does not have permission",
  "data": null
}
```

#### ⚔️ 409 Conflict

Terjadi jika terdapat konflik pada server, seperti duplikasi data unik atau pelanggaran aturan bisnis tertentu.

```json
{
  "success": false,
  "message": "Data already exists",
  "data": {
    "invoice_number": "Invoice number INV-260105-001 is already in use"
  }
}
```

#### 🚫 429 Too Many Requests

Audit: Karena proses pembuatan order melibatkan penulisan ke banyak tabel, sistem membatasi frekuensi request untuk mencegah spam pesanan palsu atau serangan DDoS.

```json
{
  "success": false,
  "message": "Too many requests, please try again later",
  "data": {
    "retry_after_seconds": 60
  }
}
```

#### 🔥 500 Internal Server Error

Terjadi kegagalan pada sistem internal atau kegagalan transaksi pada database.

```json
{
  "success": false,
  "message": "An unexpected server error occurred",
  "data": {
    "error_code": "DB_TRANSACTION_FAILED"
  }
}
```

---

## Endpoint : `GET /orders`

### Headers :

- `Authorization`: `Bearer <access_token>`
- `Accept`: `application/json`
- `Permissions`: `owner, cashier, staff, courier` (Internal Only)

#### Description :

Endpoint ini digunakan untuk mengambil daftar seluruh pesanan laundry. Karena Anda sedang membangun sistem yang skalabel, endpoint ini menggunakan fitur Pagination agar beban server Golang Anda tetap ringan dan aplikasi tetap responsif. Data yang ditampilkan adalah ringkasan informasi terbaru dari tabel orders tanpa menyertakan riwayat lengkap (history) untuk menghindari over-fetching.

### Parameters :

Bagian ini digunakan untuk mengontrol aliran data yang keluar dari database. Dengan menggunakan parameter ini, Anda memastikan server Golang tidak menarik data yang tidak diperlukan (Query Optimization).

| Key            | Tipe    | Default | Deskripsi                                                                  |
| -------------- | ------- | ------- | -------------------------------------------------------------------------- |
| page           | Integer | 1       | Nomor halaman (Pagination).                                                |
| per_page       | Integer | 10      | Jumlah data per halaman. Selaras dengan meta.per_page pada response.       |
| search         | String  | null    | Cari berdasarkan Invoice Number atau Nama Customer.                        |
| status         | String  | null    | Filter status internal proses (pending, in-progress, ready-delivery, dll). |
| payment_status | String  | null    | Filter pembayaran (paid, unpaid).                                          |

```
GET /api/orders?page=1&per_page=10&search=INV-260105&status=pending&payment_status=unpaid
```

### Request Body :

```json
None (Kosong, karena method GET tidak boleh punya body).
```

### Responses Body :

#### ✅ 200 OK

Data berhasil diambil. Struktur ini memisahkan antara data (daftar pesanan) dan meta (informasi halaman) agar Frontend dapat membuat komponen pagination dengan mudah.

```json
{
  "success": true,
  "message": "Orders retrieved successfully",
  "data": [
    {
      "id": 45,
      "invoice_number": "INV-260105-001",
      "is_delivery": 1,
      "total_price": 60000.0,
      "payment_status": "unpaid",
      "status_internal": "pending",
      "estimated_ready_at": "2026-01-08 13:00:00",
      "created_by": 2,
      "created_by_name": "Siti Aminah",
      "created_at": "2026-01-05 13:00:00",
      "updated_at": null,
      "customer": {
        "id": 101,
        "name": "Mpok Romlah",
        "phone": "081234567890"
      },
      "delivery": {
        "id": 12,
        "shipping_cost": 10000.0
      }
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

#### ⚠️ 400 Bad Request

Terjadi jika parameter filter salah format.

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "page": "Page must be a positive integer"
  }
}
```

#### ⚠️ 401 Unauthorized

Terjadi ketika token akses tidak valid, kedaluwarsa, atau tidak disertakan dalam header.

```json
{
  "success": false,
  "message": "Invalid or missing access token",
  "data": null
}
```

#### 🚫 403 Forbidden

Audit: Menjaga insting keamanan Anda agar Customer tidak bisa bypass data.

```json
{
  "success": false,
  "message": "Your role does not have permission",
  "data": null
}
```

#### 🚫 429 Too Many Requests

Dalam sistem profesional, Anda harus melindungi server Golang Anda dari serangan Brute Force atau Spamming.

- Skenario: Jika ada user (atau bot) yang memanggil GET /orders sebanyak 100 kali dalam 1 detik.

```json
{
  "success": false,
  "message": "Too many requests, please try again later",
  "data": {
    "retry_after_seconds": 30
  }
}
```

#### 🔥 500 Internal Server Error

Terjadi kegagalan pada sistem internal atau kesalahan saat melakukan query data ke database.

```json
{
  "success": false,
  "message": "An unexpected server error occurred",
  "data": {
    "error_code": "DB_QUERY_FAILED"
  }
}
```

---

## Endpoint : `GET /orders/{id}`

### Headers :

- `Authorization`: `Bearer <access_token>`
- `Accept`: `application/json`
- `Permissions`: `owner, cashier, staff, courier` (Internal Only)

#### Description :

Endpoint ini digunakan untuk mengambil data detail lengkap dari satu pesanan tertentu. Backend akan melakukan operasi JOIN atau pemanggilan data dari 5 tabel (orders, customers, order_items, deliveries, dan status_history) untuk memberikan gambaran utuh mengenai satu transaksi kepada user.

### Parameters :

Bagian ini menggunakan Path Parameter untuk mengidentifikasi sumber data secara spesifik di database.

| Key | Tipe    | In   | Deskripsi                                   |
| --- | ------- | ---- | ------------------------------------------- |
| id  | Integer | Path | ID unik pesanan yang ingin diambil datanya. |

```
GET /api/orders/45
```

### Request Body :

```json
None (Kosong, karena method GET tidak boleh punya body).
```

### Responses Body :

#### ✅ 200 OK

Data ditemukan dan dikembalikan secara lengkap. Struktur data tetap konsisten menggunakan objek singular untuk hubungan 1-ke-1 dan array untuk riwayat/item.

```json
{
  "success": true,
  "message": "Order retrieved successfully",
  "data": {
    "id": 45,
    "invoice_number": "INV-260105-001",
    "is_delivery": 1,
    "total_price": 60000.0,
    "payment_status": "unpaid",
    "status_internal": "pending",
    "estimated_ready_at": "2026-01-08 13:00:00",
    "notes": "Jangan dicampur dengan baju luntur",
    "created_by": 2,
    "created_by_name": "Siti Aminah",
    "created_at": "2026-01-05 13:00:00",
    "updated_at": null,
    "customer": {
      "id": 101,
      "name": "Mpok Romlah",
      "phone": "081234567890",
      "address": "Jl. Merpati No. 12"
    },
    "order_items": [
      {
        "id": 12,
        "service_id": 1,
        "service_name": "Cuci Kiloan Reguler",
        "item_notes": "Pisahkan warna putih",
        "quantity": null,
        "weight_kg": 5.0,
        "unit": "Kg",
        "unit_price": 10000.0,
        "subtotal": 50000.0
      }
    ],
    "delivery": {
      "id": 12,
      "shipping_cost": 10000.0,
      "courier_id": null,
      "courier_name": null,
      "courier_phone": null,
      "courier_departed_at": null,
      "courier_arrived_at": null,
      "cod_collected_amount": 0.0
    },
    "status_history": [
      {
        "id": 1,
        "previous_status": null,
        "new_status": "pending",
        "actor_name": "Siti Aminah",
        "actor_role": "cashier",
        "notes": "Initial order creation",
        "created_at": "2026-01-05 13:00:00"
      }
    ]
  }
}
```

#### ⚠️ 400 Bad Request

Terjadi jika format ID yang dikirimkan pada URL tidak valid (bukan angka).

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "id": "Order ID must be a valid integer"
  }
}
```

#### ⚠️ 401 Unauthorized

Terjadi ketika token akses tidak valid, kedaluwarsa, atau tidak disertakan dalam header.

```json
{
  "success": false,
  "message": "Invalid or missing access token",
  "data": null
}
```

#### 🚫 403 Forbidden

Terjadi jika role user (misal: Customer) mencoba mengakses detail pesanan milik orang lain (IDOR Protection).

```json
{
  "success": false,
  "message": "Your role does not have permission",
  "data": null
}
```

#### 🔍 404 Not Found

Terjadi jika ID pesanan yang dicari tidak ada di database MySQL.

```json
{
  "success": false,
  "message": "Order not found",
  "data": null
}
```

#### 🚫 429 Too Many Requests

Mencegah upaya pemindaian ID secara otomatis (ID Scanning/Brute Force).

```json
{
  "success": false,
  "message": "Too many requests, please try again later",
  "data": {
    "retry_after_seconds": 30
  }
}
```

#### 🔥 500 Internal Server Error

Terjadi kegagalan koneksi database atau kesalahan logika JOIN pada query Golang.

```json
{
  "success": false,
  "message": "An unexpected server error occurred",
  "data": {
    "error_code": "DB_QUERY_FAILED"
  }
}
```

---

## Endpoint : `PUT /orders/{id}`

### Headers :

- `Authorization`: `Bearer <access_token>`
- `Accept`: `application/json`
- `Content-Type`: `application/json`
- `Permissions`: `cashier, owner`

#### Description :

Endpoint ini digunakan untuk melakukan pembaruan total (Full Update) pada data pesanan. Sesuai aturan bisnis, perubahan hanya diizinkan jika status_internal masih bernilai pending. Jika pesanan sudah masuk tahap in-progress atau lebih lanjut, permintaan akan ditolak. Backend akan menghitung ulang total_price secara otomatis jika terdapat perubahan pada order_items atau shipping_cost.

### Parameters :

Bagian ini menggunakan Path Parameter untuk mengunci ID pesanan yang akan direvisi di database.

| Key | Tipe    | In   | Deskripsi                              |
| --- | ------- | ---- | -------------------------------------- |
| id  | Integer | Path | ID unik pesanan yang ingin diperbarui. |

### Request Body :

Mengirimkan struktur data lengkap untuk menggantikan data lama. Gunakan customer_id jika pelanggan sudah terdaftar.

```json
{
  "customer_id": 101,
  "customer_name": "Mpok Romlah",
  "customer_phone": "081234567890",
  "customer_address": "Jl. Merpati No. 12",
  "is_delivery": 1,
  "notes": "Jangan dicampur dengan baju luntur",
  "shipping_cost": 10000.0,
  "order_items": [
    {
      "service_id": 1,
      "quantity": null,
      "weight_kg": 10.0, // update berat laundry kiloan
      "item_notes": "Pisahkan warna putih"
    }
  ]
}
```

### Responses Body :

#### ✅ 200 OK

Pesanan berhasil diperbarui. Response mengembalikan objek data terbaru yang sudah dikalkulasi ulang oleh server.

```json
{
  "success": true,
  "message": "Order updated successfully",
  "data": {
    "id": 45,
    "invoice_number": "INV-260105-001",
    "is_delivery": 1,
    "total_price": 110000.0,
    "payment_status": "unpaid",
    "status_internal": "pending",
    "estimated_ready_at": "2026-01-08 13:00:00",
    "notes": "Jangan dicampur dengan baju luntur",
    "created_by": 2,
    "created_by_name": "Siti Aminah",
    "created_at": "2026-01-05 13:00:00",
    "updated_at": "2026-01-05 13:10:00",
    "customer": {
      "id": 101,
      "name": "Mpok Romlah",
      "phone": "081234567890",
      "address": "Jl. Merpati No. 12"
    },
    "order_items": [
      {
        "id": 12,
        "service_id": 1,
        "service_name": "Cuci Kiloan Reguler",
        "item_notes": "Pisahkan warna putih",
        "quantity": null,
        "weight_kg": 10.0, // update tambah berat layanan kiloan
        "unit": "Kg",
        "unit_price": 10000.0,
        "subtotal": 100000.0
      }
    ],
    "delivery": {
      "id": 12,
      "shipping_cost": 10000.0,
      "courier_id": null,
      "courier_name": null,
      "courier_phone": null,
      "courier_departed_at": null,
      "courier_arrived_at": null,
      "cod_collected_amount": 0.0
    },
    "status_history": [
      {
        "id": 1,
        "previous_status": null,
        "new_status": "pending",
        "actor_name": "Siti Aminah",
        "actor_role": "cashier",
        "notes": "Initial order creation",
        "created_at": "2026-01-05 13:00:00"
      }
    ]
  }
}
```

#### ⚠️ 400 Bad Request

Terjadi jika format input salah atau melanggar aturan bisnis (mencoba edit pesanan yang sudah diproses).

```json
{
  "success": false,
  "message": "Update failed",
  "data": {
    "status": "Order can only be edited when status is pending"
  }
}
```

#### ⚠️ 401 Unauthorized

Terjadi jika token akses tidak valid atau sudah kadaluwarsa.

```json
{
  "success": false,
  "message": "Invalid or missing access token",
  "data": null
}
```

#### 🚫 403 Forbidden

Terjadi jika role user tidak memiliki izin (misal: Staff atau Courier mencoba melakukan edit total).

```json
{
  "success": false,
  "message": "Your role does not have permission",
  "data": null
}
```

#### 404 Not Found

Terjadi jika ID pesanan yang ingin di-update tidak ditemukan di database.

```json
{
  "success": false,
  "message": "Order not found",
  "data": null
}
```

#### 💥 409 Conflict

Terjadi jika terdapat duplikasi data unik yang tidak sengaja tercipta saat proses update.

```json
{
  "success": false,
  "message": "Data already exists",
  "data": {
    "invoice_number": "Invoice number already in use"
  }
}
```

#### 🚫 429 Too Many Requests

Perlindungan server dari upaya perubahan data yang terlalu masif dalam waktu singkat.

```json
{
  "success": false,
  "message": "Too many requests, please try again later",
  "data": {
    "retry_after_seconds": 30
  }
}
```

#### 🔥 500 Internal Server Error

Terjadi kesalahan pada transaksi database atau kegagalan sistem internal saat pemrosesan.

```json
{
  "success": false,
  "message": "An unexpected server error occurred",
  "data": {
    "error_code": "DB_TRANSACTION_FAILED"
  }
}
```

---

## Endpoint : `PATCH /orders/{id}`

### Headers :

- `Authorization`: `Bearer <access_token>`
- `Accept`: `application/json`
- `Content-Type`: `application/json`
- `Permissions`: `owner, cashier, staff, courier` (Internal Only)

#### Description :

Endpoint ini digunakan untuk memperbarui status operasional pesanan secara parsial (State Transition). Berbeda dengan PUT yang melakukan update data fisik, PATCH berfokus pada pergerakan workflow. Backend akan menjalankan transaksi atomik untuk memperbarui dua tabel utama:

1. orders: Memperbarui kolom status_internal dan updated_at.
2. status_history: Menambahkan baris riwayat baru untuk melacak siapa yang mengubah status, kapan, dan alasan perubahannya (audit trail).

### Parameters :

Bagian ini menggunakan Path Parameter untuk menentukan pesanan mana yang akan diproses transisi statusnya.

| Key | Tipe | In   | Deskripsi                                    |
| --- | ---- | ---- | -------------------------------------------- |
| id  | Int  | Path | ID unik pesanan yang ingin diubah statusnya. |

### Request Body :

Hanya mengirimkan informasi transisi status. Field notes bersifat opsional namun sangat disarankan untuk audit internal.

```json
{
  "new_status": "in-progress",
  "notes": "Pakaian mulai dimasukkan ke mesin cuci nomor 03"
}
```

### Responses Body :

#### ✅ 200 OK

Status berhasil diperbarui. Mengikuti prinsip "The Finished Plate", response mengembalikan data lengkap agar UI dapat langsung memindahkan card pesanan ke kolom yang sesuai tanpa re-fetch.

```json
{
  "success": true,
  "message": "Order updated successfully",
  "data": {
    "id": 45,
    "invoice_number": "INV-260105-001",
    "is_delivery": 1,
    "total_price": 110000.0,
    "payment_status": "unpaid",
    "status_internal": "in-progress",
    "estimated_ready_at": "2026-01-08 13:00:00",
    "notes": "Jangan dicampur dengan baju luntur",
    "created_by": 2,
    "created_by_name": "Siti Aminah",
    "created_at": "2026-01-05 13:00:00",
    "updated_at": "2026-01-06 10:00:00",
    "customer": {
      "id": 101,
      "name": "Mpok Romlah",
      "phone": "081234567890",
      "address": "Jl. Merpati No. 12"
    },
    "order_items": [
      {
        "id": 12,
        "service_id": 1,
        "service_name": "Cuci Kiloan Reguler",
        "item_notes": "Pisahkan warna putih",
        "quantity": null,
        "weight_kg": 10.0,
        "unit": "Kg",
        "unit_price": 10000.0,
        "subtotal": 100000.0
      }
    ],
    "delivery": {
      "id": 12,
      "shipping_cost": 10000.0,
      "courier_name": null,
      "courier_phone": null,
      "courier_departed_at": null,
      "courier_arrived_at": null,
      "cod_collected_amount": 0.0
    },
    "status_history": [
      {
        "id": 1,
        "previous_status": null,
        "new_status": "pending",
        "actor_name": "Siti Aminah",
        "actor_role": "cashier",
        "notes": "Initial order creation",
        "created_at": "2026-01-05 13:00:00"
      },
      {
        "id": 2,
        "previous_status": "pending",
        "new_status": "in-progress",
        "actor_name": "Fadhillah Kurnia",
        "actor_role": "staff",
        "notes": "Pakaian mulai dimasukkan ke mesin cuci nomor 03",
        "created_at": "2026-01-06 10:00:00"
      }
    ]
  }
}
```

#### ⚠️ 400 Bad Request

Terjadi jika transisi status melanggar aturan bisnis (misal: pesanan cancelled tidak bisa diubah ke in-progress).

```json
{
  "success": false,
  "message": "Invalid status transition",
  "data": {
    "status": "Cannot change status from cancelled to in-progress"
  }
}
```

#### ⚠️ 401 Unauthorized

Terjadi ketika token akses tidak valid, kedaluwarsa, atau tidak disertakan dalam header.

```json
{
  "success": false,
  "message": "Invalid or missing access token",
  "data": null
}
```

#### 🚫 403 Forbidden

Audit: Menjamin bahwa hanya user dengan role tertentu yang bisa mengubah status spesifik (misal: hanya Courier yang bisa mengubah ke being-delivered).

```json
{
  "success": false,
  "message": "Your role does not have permission",
  "data": null
}
```

#### 🔍 404 Not Found

Terjadi jika ID pesanan yang ingin di-patch tidak ditemukan di database.

```json
{
  "success": false,
  "message": "Order not found",
  "data": null
}
```

#### 🚫 429 Too Many Requests

Mencegah spamming aksi pada tombol operasional yang berakibat pada penulisan log berlebihan di database.

```json
{
  "success": false,
  "message": "Too many requests, please try again later",
  "data": {
    "retry_after_seconds": 10
  }
}
```

#### 🔥 500 Internal Server Error

Terjadi kegagalan transaksi pada database saat mencoba menulis riwayat status baru.

```json
{
  "success": false,
  "message": "An unexpected server error occurred",
  "data": {
    "error_code": "DB_HISTORY_INSERT_FAILED"
  }
}
```
