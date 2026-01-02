# LAUNDRY MANAGEMENT SYSTEM — API SPECIFICATION

## SERVICE CATEGORIES MODULE SPECIFICATION

---

## Endpoint : `POST /categories`

### Headers :

- `Authorization: Bearer <access_token>` (Required) - Hanya Owner yang diizinkan.
- `Content-Type`: `application/json` (Required)

#### Description :

Endpoint ini digunakan untuk membuat Kategori Layanan baru. Kategori digunakan untuk mengelompokkan jenis layanan agar lebih rapi di laporan dan aplikasi kasir (Contoh: "Kiloan", "Satuan", "Dry Clean", "Karpet"). Hanya Owner yang memiliki hak akses untuk menambah kategori.

### Request Body :

Objek JSON berisi nama kategori dan deskripsinya.

```json
{
  "category_name": "Laundry Kiloan",
  "description": "Cuci pakaian sehari-hari dihitung per kilogram"
}
```

### Responses Body :

#### ✅ 201 Created (Success)

Kategori berhasil dibuat.

```json
{
  "success": true,
  "message": "Category created successfully",
  "data": {
    "id": 1,
    "category_name": "Laundry Kiloan",
    "description": "Cuci pakaian sehari-hari dihitung per kilogram",
    "is_active": true,
    "created_at": "2024-02-20T10:00:00Z",
    "updated_at": null
  }
}
```

#### ⚠️ 400 Bad Request (Validation Error)

Validasi input gagal (nama kategori kosong).

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "category_name": "Category name is required"
  }
}
```

#### 🚫 401 Unauthorized (Invalid Token)

Token tidak valid atau tidak dikirim.

```json
{
  "success": false,
  "message": "Invalid or missing access token",
  "data": null
}
```

#### ⛔ 403 Forbidden (Insufficient Permissions)

Token valid, tetapi role pengguna bukan Owner.

```json
{
  "success": false,
  "message": "You do not have permission to perform this action",
  "data": null
}
```

#### 💥 409 Conflict (Duplicate Data)

Nama kategori sudah ada di database.

```json
{
  "success": false,
  "message": "Data already exists",
  "data": {
    "category_name": "Category 'Laundry Kiloan' already exists"
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

## Endpoint : `GET /categories`

### Headers :

- `Authorization: Bearer <access_token>` (Required) - Semua Role (Owner, Staff, Cashier) boleh mengakses ini (penting untuk dropdown menu di Frontend).

#### Description :

Endpoint ini digunakan untuk mengambil daftar kategori layanan. Karena ini adalah endpoint "List", respon dilengkapi dengan Pagination dan Filtering (berdasarkan nama atau status aktif). Endpoint ini biasanya dipanggil saat memuat halaman daftar layanan atau saat Kasir memilih kategori di menu transaksi.

### Parameters :

Parameter ini dikirim melalui URL (contoh: /categories?page=1&limit=10&is_active=true).

| Key       | Tipe   | Default  | Deskripsi                                                      |
| --------- | ------ | -------- | -------------------------------------------------------------- |
| page      | Int    | 1        | Halaman ke berapa yang mau diambil.                            |
| limit     | Int    | 10       | Jumlah data per halaman.                                       |
| name      | String | -        | (Opsional) Cari kategori berdasarkan nama.                     |
| is_active | String | Show all | (Opsional) Filter status: "true" (Aktif) atau "false" (Arsip). |

### Request Body :

```json
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK (Success)

Data kategori berhasil diambil.

```json
{
  "success": true,
  "message": "Categories retrieved successfully",
  "data": [
    {
      "id": 1,
      "category_name": "Laundry Kiloan",
      "description": "Cuci pakaian sehari-hari dihitung per kilogram",
      "is_active": true,
      "created_at": "2024-02-20T10:00:00Z",
      "updated_at": null
    },
    {
      "id": 2,
      "category_name": "Laundry Satuan",
      "description": "Kategori untuk layanan per item",
      "is_active": true,
      "created_at": "2024-02-20T10:00:00Z",
      "updated_at": "2024-03-15T14:30:00Z"
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

Parameter query tidak valid (misal: page berupa huruf).

```json
{
  "success": false,
  "message": "Invalid query parameters",
  "data": {
    "page": "Must be a number"
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

## Endpoint : `GET /categories/{id}`

### Headers :

- `Authorization: Bearer <access_token>` (Required) - Semua Role diperbolehkan mengakses (biasanya untuk pengecekan sebelum edit).

#### Description :

Endpoint ini digunakan untuk mengambil detail data dari satu kategori spesifik berdasarkan ID. Biasanya digunakan oleh Frontend saat menampilkan halaman "Edit Kategori" untuk mengisi formulir dengan data yang sedang aktif di database.

### Parameters :

Parameter ini disisipkan langsung di URL (contoh: /categories/1).

| Key | Tipe | Required | Deskripsi                                      |
| --- | ---- | -------- | ---------------------------------------------- |
| id  | Int  | Yes      | ID Unik kategori yang ingin dilihat detailnya. |

### Request Body :

```json
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK (Success)

Detail kategori ditemukan.

```json
{
  "success": true,
  "message": "Category detail retrieved successfully",
  "data": {
    "id": 1,
    "category_name": "Laundry Kiloan",
    "description": "Cuci pakaian sehari-hari dihitung per kilogram",
    "is_active": true,
    "created_at": "2024-02-20T10:00:00Z",
    "updated_at": null
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

ID valid secara format, tetapi data kategori tidak ditemukan di database.

```json
{
  "success": false,
  "message": "Category not found",
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

## Endpoint : `PUT /categories/{id}`

### Headers :

- `Authorization: Bearer <access_token>` (Required) - Hanya Owner yang diizinkan.
- `Content-Type`: `application/json` (Required)

#### Description :

Endpoint ini digunakan untuk memperbarui data kategori layanan. Owner dapat mengubah nama, deskripsi, atau menonaktifkan kategori (is_active: false) jika kategori tersebut sudah tidak digunakan lagi. Sistem akan memvalidasi ulang keunikan nama kategori jika ada perubahan pada nama.

### Parameters :

Parameter ini disisipkan langsung di URL (contoh: /categories/1).

| Key | Tipe | Required | Deskripsi                           |
| --- | ---- | -------- | ----------------------------------- |
| id  | Int  | Yes      | ID Unik kategori yang ingin diedit. |

### Request Body :

Objek JSON berisi data yang ingin diubah.

```json
{
  "category_name": "Laundry Satuan", // Required, Unique
  "description": "Kategori untuk layanan cuci satuan", // Optional
  "is_active": true // Required
}
```

### Responses Body :

#### ✅ 200 OK (Success)

Data kategori berhasil diperbarui.

```json
{
  "success": true,
  "message": "Category updated successfully",
  "data": {
    "id": 1,
    "category_name": "Laundry Satuan",
    "description": "Kategori untuk layanan cuci satuan",
    "is_active": true,
    "created_at": "2024-02-20T10:00:00Z",
    "updated_at": "2025-03-15T14:30:00Z"
  }
}
```

#### ⚠️ 400 Bad Request (Validation Error)

Format ID salah atau input tidak valid.

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "id": "ID must be a number",
    "category_name": "Category name is required"
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

#### ❓ 404 Not Found (Category Not Found)

ID valid, tetapi data kategori tidak ditemukan di database.

```json
{
  "success": false,
  "message": "Category not found",
  "data": null
}
```

#### 💥 409 Conflict (Duplicate Data)

Gagal update karena nama kategori baru sudah digunakan oleh kategori lain.

```json
{
  "success": false,
  "message": "Data already exists",
  "data": {
    "category_name": "Category name 'Laundry Satuan' already exists"
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

## Endpoint : `DELETE /categories/{id}`

### Headers :

- `Authorization: Bearer <access_token>` (Required) - Hanya Owner yang diizinkan.

#### Description :

Endpoint ini digunakan untuk menonaktifkan Kategori Layanan (Soft Delete). Data kategori TIDAK DIHAPUS permanen, melainkan statusnya diubah menjadi tidak aktif (is_active = false). Catatan: Ketika kategori dinonaktifkan, biasanya layanan-layanan di dalamnya akan ikut disembunyikan dari menu kasir (tergantung logika frontend), namun riwayat transaksi lama yang menggunakan kategori ini tetap aman.

### Parameters :

Parameter ini disisipkan langsung di URL (contoh: /categories/2).

| Key | Tipe | Required | Deskripsi                                  |
| --- | ---- | -------- | ------------------------------------------ |
| id  | Int  | Yes      | ID Unik kategori yang ingin dinonaktifkan. |

### Request Body :

```json
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK (Success)

Kategori berhasil dinonaktifkan.

```json
{
  "success": true,
  "message": "Category deactivated successfully",
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

#### 🔍 404 Not Found (Category Missing)

Kategori dengan ID tersebut tidak ditemukan di database.

```json
{
  "success": false,
  "message": "Category not found",
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
