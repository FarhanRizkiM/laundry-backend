# LAUNDRY MANAGEMENT SYSTEM — API SPECIFICATION

## SERVICES MODULE SPECIFICATION

---

## Endpoint : `POST /services`

### Headers :

- `Authorization: Bearer <access_token>` (Required)

#### Description :

Endpoint ini digunakan untuk membuat Layanan/Produk baru. PENTING: Wajib menyertakan category_id yang valid. Hanya Owner yang boleh mengakses endpoint ini.

### Request Body :

Field category_id harus berupa ID Kategori yang sudah ada di database. unit dibatasi pada pilihan tertentu (enum).

```json
{
  "category_id": 1, // Foreign Key (Wajib Valid)
  "code": "SVC-LKR-1", // Unique, SKU/Kode Barang
  "service_name": "Kiloan Regular", // Singular (bukan services_name)
  "unit": "kg", // Enum: "kg", "pcs"
  "price": 7000 // Decimal/Int. Hindari float 7.0 jika bisa, gunakan nominal penuh.
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
    "is_active": true,
    "created_at": "2025-12-26T18:49:21.410Z",
    "updated_at": "2025-12-26T18:49:21.410Z"
  }
}
```

#### ⚠️ 400 Bad Request (Validation Error)

Validasi input gagal (Harga negatif, unit tidak dikenal, atau format salah).

```json
{
  "success": false,
  "message": "Input Validation Failed",
  "data": {
    "price": "Price cannot be negative",
    "unit": "Unit must be one of: kg, pcs"
  }
}
```

#### ⚠️ 401 Unauthorized (Invalid Token)

Token tidak valid.

```json
{
  "success": false,
  "message": "Invalid or missing access token",
  "data": null
}
```

#### 🚫 403 Forbidden (Insufficient Permissions)

User bukan Owner mencoba membuat harga/layanan.

```json
{
  "success": false,
  "message": "You do not have permission to perform this action",
  "data": null
}
```

#### 🔍 404 Not Found (Foreign Key Error)

Terjadi jika category_id yang dikirim tidak ditemukan di database.

```json
{
  "success": false,
  "message": "Category with ID 1 not found",
  "data": null
}
```

#### 💥 409 Conflict (Duplicate Data)

Kode layanan (service_code) sudah terdaftar.

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

---

## Endpoint : `GET /services`

### Headers :

- `Authorization: Bearer <access_token>` (Required)

#### Description :

Endpoint ini digunakan untuk menampilkan daftar layanan/menu laundry. Sangat penting untuk operasional Kasir (memilih menu) dan Owner (melihat stok layanan). Mendukung filter kategori, pencarian kode, dan nama.

### Parameters :

Parameter ini ditempel di URL (Query String).

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

Data layanan berhasil diambil.

```json
{
  "success": true,
  "message": "Services retrieved successfully",
  "data": [
    {
      "id": 1,
      "category_id": 1,
      "code": "SVC-LKR-1", // Sesuai DB
      "service_name": "Kiloan Regular", // Sesuai DB
      "unit": "kg",
      "price": 7000,
      "is_active": true,
      "created_at": "2025-12-26T18:49:21.410Z",
      "updated_at": "2025-12-26T18:49:21.410Z"
    },
    {
      "id": 2,
      "category_id": 1,
      "code": "SVC-LKK-1",
      "service_name": "Kiloan Kilat",
      "unit": "kg",
      "price": 10000,
      "is_active": true,
      "created_at": "2025-12-26T18:49:21.410Z",
      "updated_at": "2025-12-26T18:49:21.410Z"
    }
  ],
  "meta": {
    "current_page": 1,
    "total_pages": 5,
    "total_items": 12,
    "per_page": 10
  }
}
```

#### ⚠️ 400 Bad Request (Invalid Parameters)

Terjadi jika parameter salah format (misal page=abc).

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

Token tidak valid, expired, atau user belum login.

```json
{
  "success": false,
  "message": "Invalid or missing access token",
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

## Endpoint : `GET /services/{id}`

### Headers :

- `Authorization: Bearer <access_token>` (Required)

#### Description :

Endpoint ini digunakan untuk melihat detail lengkap satu layanan laundry berdasarkan ID-nya. Dapat diakses oleh Owner, Kasir, dan Staff (untuk kebutuhan pengecekan harga atau detail item).

### Parameters :

Parameter ini adalah Path Variable.

| Key | Tipe | Default | Deskripsi                                     |
| --- | ---- | ------- | --------------------------------------------- |
| id  | Int  | -       | ID Unik layanan yang ingin dilihat detailnya. |

### Request Body :

```json
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK (Success)

Data layanan ditemukan.

```json
{
  "success": true,
  "message": "Services detail retrieved successfully",
  "data": {
    "id": 1,
    "category_id": 1,
    "code": "SVC-LKR-1", // Sesuai DB
    "service_name": "Kiloan Regular", // Sesuai DB
    "unit": "kg",
    "price": 7000,
    "is_active": true,
    "created_at": "2025-12-26T18:49:21.410Z",
    "updated_at": "2025-12-26T18:49:21.410Z"
  }
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

---

## Endpoint : `PUT /services/{id}`

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

## Endpoint : `DELETE /services/{id}`

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
