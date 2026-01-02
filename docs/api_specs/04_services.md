# LAUNDRY MANAGEMENT SYSTEM — API SPECIFICATION

## SERVICES MODULE SPECIFICATION

---

## Endpoint : `POST /services`

### Headers :

- `Authorization: Bearer <access_token>` (Required) - Hanya Owner yang diizinkan.
- `Content-Type`: `application/json` (Required)

#### Description :

Endpoint ini digunakan untuk menambahkan layanan laundry baru ke dalam database. Penting: Layanan harus dikaitkan dengan Kategori yang valid (melalui category_id). Selain itu, code (SKU) harus unik untuk memudahkan pencarian cepat saat transaksi.

### Request Body :

Objek JSON berisi detail layanan.

```json
{
  "category_id": 1, // Foreign Key (Wajib Valid)
  "code": "SVC-LKR-1", // Unique, SKU/Kode Barang
  "service_name": "Kiloan Regular", // Singular (bukan services_name)
  "unit": "kg", // Enum: "kg", "pcs"
  "price": 7000, // Decimal/Int. Hindari float 7.0 jika bisa, gunakan nominal penuh.
  "duration_hours": 72 // Integer. (72 = 3 hari, 24 = 1 hari, 12 = setengah hari, 6 = 6 jam saja).
}
```

### Responses Body :

#### ✅ 201 Created (Success)

Layanan berhasil dibuat.

```json
{
  "success": true,
  "message": "Service created successfully",
  "data": {
    "id": 1,
    "category_id": 1,
    "code": "SVC-LKR-1",
    "service_name": "Kiloan Regular",
    "unit": "kg",
    "price": 7000,
    "duration_hours": 72, // <--- Data durasi tersimpan
    "is_active": true,
    "created_at": "2025-12-26T18:49:21.410Z",
    "updated_at": null
  }
}
```

#### ⚠️ 400 Bad Request (Validation Error)

Input tidak valid (misal: harga negatif, satuan salah, atau field wajib kosong).

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "duration_hours": "Duration must be a positive integer (greater than 0)",
    "price": "Price cannot be negative",
    "unit": "Unit must be one of: kg, pcs"
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

Token valid, tetapi role pengguna bukan Owner.

```json
{
  "success": false,
  "message": "You do not have permission to perform this action",
  "data": null
}
```

#### 🔍 404 Not Found (Foreign Key Error)

ID Kategori yang dikirim tidak ditemukan di database.

```json
{
  "success": false,
  "message": "Category with ID 1 not found",
  "data": null
}
```

#### 💥 409 Conflict (Duplicate Data)

Gagal simpan karena Kode Layanan atau Nama Layanan sudah ada sebelumnya.

```json
{
  "success": false,
  "message": "Data already exists",
  "data": {
    "code": "Service code 'SVC-LKR-1' already exists"
  }
}
```

#### 🔥 500 Internal Server Error

Terjadi kesalahan di sisi server.

```json
{
  "success": false,
  "message": "An unexpected error occurred",
  "data": null
}
```

---

## Endpoint : `GET /services`

### Headers :

- `Authorization: Bearer <access_token>` (Required) - Semua Role (Owner, Staff, Kasir) boleh akses.

#### Description :

Endpoint ini digunakan untuk mengambil daftar layanan laundry. Response dilengkapi dengan Pagination dan Filtering yang lengkap. Endpoint ini sangat krusial karena akan digunakan di halaman "Kasir" (POS) untuk memilih layanan. Data response menyertakan detail Kategori (Nested JSON) dan Durasi pengerjaan. Endpoint ini digunakan untuk mengambil daftar layanan laundry. Poin Penting:

1. Respon menggunakan format Nested JSON (Objek Kategori ada di dalam Objek Layanan) agar Frontend bisa langsung menampilkan nama kategori tanpa request tambahan.
2. Menyertakan duration_hours agar sistem Kasir bisa menghitung estimasi waktu selesai secara otomatis.
3. Mendukung Pagination dan Filtering.

### Parameters :

Parameter ini dikirim melalui URL (contoh: /services?page=1&category_id=1).

| Key          | Tipe    | Default     | Deskripsi                                                              |
| ------------ | ------- | ----------- | ---------------------------------------------------------------------- |
| page         | Int     | 1           | Halaman ke berapa yang mau diambil.                                    |
| limit        | Int     | 10, Max 100 | Jumlah data per halaman.                                               |
| category_id  | Int     | -           | Filter layanan berdasarkan ID Kategori (misal: cuma tampilkan Kiloan). |
| code         | String  | -           | Filter berdasarkan Kode Layanan (misal: hasil scan barcode).           |
| service_name | String  |             | Cari layanan berdasarkan nama (Partial Search).                        |
| is_active    | Boolean | true        | Filter status. Kasir biasanya default true.                            |

### Request Body :

```json
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK (Success)

Daftar layanan berhasil diambil dengan struktur data bersarang.

```json
{
  "success": true,
  "message": "Services retrieved successfully",
  "data": [
    {
      "id": 1,
      "code": "SVC-LKR-1",
      "service_name": "Kiloan Regular",
      "unit": "kg",
      "price": 7000,
      "duration_hours": 72, // 3 Hari
      "is_active": true,
      "created_at": "2025-12-26T18:49:21.410Z",
      "updated_at": null,
      "category": {
        // <--- STRUKTUR BERSARANG (NESTED)
        "id": 1,
        "category_name": "Laundry Kiloan",
        "description": "Cuci pakaian sehari-hari"
      }
    },
    {
      "id": 2,
      "code": "SVC-LKK-1",
      "service_name": "Kiloan Kilat",
      "unit": "kg",
      "price": 10000,
      "duration_hours": 6, // 6 Jam
      "is_active": true,
      "created_at": "2025-12-26T18:49:21.410Z",
      "updated_at": null,
      "category": {
        // <--- Kategori ikut terbawa
        "id": 1,
        "category_name": "Laundry Kiloan",
        "description": "Cuci pakaian sehari-hari"
      }
    }
  ],
  "meta": {
    "current_page": 1,
    "total_pages": 2,
    "total_items": 12,
    "per_page": 10
  }
}
```

#### ⚠️ 400 Bad Request (Invalid Parameters)

Parameter query tidak valid.

```json
{
  "success": false,
  "message": "Invalid query parameters",
  "data": {
    "page": "Must be a number",
    "category_id": "Must be a number"
  }
}
```

#### ⚠️ 401 Unauthorized (Invalid Token)

Token tidak valid atau tidak dikirim.

```json
{
  "success": false,
  "message": "Invalid or missing access token",
  "data": null
}
```

#### 🔥 500 Internal Server Error

Terjadi kesalahan di sisi server.

```json
{
  "success": false,
  "message": "An unexpected error occurred",
  "data": null
}
```

---

## Endpoint : `GET /services/{id}`

### Headers :

- `Authorization: Bearer <access_token>` (Required) - Semua Role diperbolehkan mengakses (biasanya untuk pengecekan sebelum edit atau validasi kasir).

#### Description :

Endpoint ini digunakan untuk mengambil detail data dari satu layanan spesifik berdasarkan ID. Biasanya digunakan oleh Frontend saat:

1. Menampilkan halaman "Edit Layanan" (Pre-fill form).
2. Kasir melakukan scan barcode (mencari detail harga & durasi).

### Parameters :

Parameter ini disisipkan langsung di URL (contoh: /services/1).

| Key | Tipe | Default | Deskripsi                                     |
| --- | ---- | ------- | --------------------------------------------- |
| id  | Int  | -       | ID Unik layanan yang ingin dilihat detailnya. |

### Request Body :

```json
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK (Success)

Detail layanan ditemukan (Lengkap dengan Durasi & Info Kategori).

```json
{
  "success": true,
  "message": "Service detail retrieved successfully",
  "data": {
    "id": 1,
    "code": "SVC-LKR-1",
    "service_name": "Kiloan Regular",
    "unit": "kg",
    "price": 7000,
    "duration_hours": 72,
    "is_active": true,
    "created_at": "2025-12-26T18:49:21.410Z",
    "updated_at": null,
    "category": {
      "id": 1,
      "category_name": "Laundry Kiloan",
      "description": "Cuci pakaian sehari-hari"
    }
  }
}
```

#### ⚠️ 400 Bad Request (Invalid ID Format)

Format ID pada URL bukan angka.

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

#### 🔍 404 Not Found (Data Missing)

ID valid secara format, tetapi data layanan tidak ditemukan di database.

```json
{
  "success": false,
  "message": "Service not found",
  "data": null
}
```

#### 🔥 500 Internal Server Error

Terjadi kesalahan di sisi server.

```json
{
  "success": false,
  "message": "An unexpected error occurred",
  "data": null
}
```

---

## Endpoint : `PUT /services/{id}`

### Headers :

- `Authorization: Bearer <access_token>` (Required) - Hanya Owner yang diizinkan.
- `Content-Type`: `application/json` (Required)

#### Description :

Endpoint ini digunakan untuk memperbarui data layanan laundry. Owner dapat mengubah harga, nama, kategori, atau durasi pengerjaan (duration_hours). Validasi: Sistem akan mengecek kembali keunikan code (SKU) jika ada perubahan pada kode tersebut.

### Parameters :

Parameter ini disisipkan langsung di URL (contoh: /services/1).

| Key | Tipe | Required | Deskripsi                          |
| --- | ---- | -------- | ---------------------------------- |
| id  | Int  | Yes      | ID Unik layanan yang ingin diedit. |

### Request Body :

Objek JSON berisi data yang ingin diubah

```json
{
  "category_id": 1, // Foreign Key (Bisa dipindah kategori)
  "code": "SVC-LKR-1", // Unique (Cek duplikasi jika berubah)
  "service_name": "Kiloan Regular",
  "unit": "kg",
  "price": 7000, // Update harga
  "duration_hours": 72, // Update durasi (misal ubah estimasi hari)
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
    "price": 7000,
    "duration_hours": 72, // Data durasi terbaru
    "is_active": true,
    "created_at": "2025-12-26T18:49:21.410Z",
    "updated_at": "2025-12-27T10:00:00.000Z" // Terisi tanggal baru
  }
}
```

#### ⚠️ 400 Bad Request (Validation Error)

Input tidak valid atau format salah.

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "category_id": "Category ID does not exist",
    "price": "Price cannot be negative",
    "duration_hours": "Duration must be greater than 0"
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

Token valid, tetapi role pengguna bukan Owner.

```json
{
  "success": false,
  "message": "You do not have permission to perform this action",
  "data": null
}
```

#### 404 Not Found (Service Missing)

Terjadi jika:

1. ID Service yang diedit tidak ditemukan.
2. ID Category baru (category_id) tidak valid.

```json
{
  "success": false,
  "message": "Service not found", // atau "Category not found"
  "data": null
}
```

#### 💥 409 Conflict (Duplicate Code)

Gagal update karena kode layanan baru sudah dipakai layanan lain.

```json
{
  "success": false,
  "message": "Data already exists",
  "data": {
    "code": "Service code 'SVC-LKR-1' already exists"
  }
}
```

#### 🔥 500 Internal Server Error

Terjadi kesalahan di sisi server.

```json
{
  "success": false,
  "message": "An unexpected error occurred",
  "data": null
}
```

## Endpoint : `DELETE /services/{id}`

### Headers :

- `Authorization: Bearer <access_token>` (Required) - Hanya Owner yang diizinkan.

#### Description :

Endpoint ini digunakan untuk menonaktifkan layanan laundry (Soft Delete). Penting: Data layanan TIDAK DIHAPUS permanen dari database. Statusnya hanya diubah menjadi tidak aktif (is_active = false). Hal ini dilakukan untuk menjaga integritas data laporan. Transaksi lama yang menggunakan layanan ini tidak akan error, tetapi layanan ini tidak akan muncul lagi di menu Kasir untuk transaksi baru.

### Parameters :

Parameter ini disisipkan langsung di URL (contoh: /services/5).

| Key | Tipe | Required | Deskripsi                                 |
| --- | ---- | -------- | ----------------------------------------- |
| id  | Int  | Yes      | ID Unik layanan yang ingin dinonaktifkan. |

### Request Body :

```json
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK (Success)

Layanan berhasil dinonaktifkan.

```json
{
  "success": true,
  "message": "Service deactivated successfully",
  "data": null
}
```

#### ⚠️ 400 Bad Request (Invalid ID Format)

Format ID pada URL bukan angka.

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

Token valid, tapi yang melakukan request bukan Owner.

```json
{
  "success": false,
  "message": "You do not have permission to perform this action",
  "data": null
}
```

#### 🔍 404 Not Found (Service Missing)

Layanan dengan ID tersebut tidak ditemukan di database.

```json
{
  "success": false,
  "message": "Service not found",
  "data": null
}
```

#### 🔥 500 Internal Server Error

Terjadi kesalahan tak terduga di server.

```json
{
  "success": false,
  "message": "An unexpected error occurred",
  "data": null
}
```
