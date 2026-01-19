# LAUNDRY MANAGEMENT SYSTEM — API SPECIFICATION

## SERVICES MODULE SPECIFICATION

---

## Endpoint : `POST /services`

### Description :

Endpoint ini digunakan oleh **Owner** untuk menambahkan layanan laundry baru ke dalam sistem. Setiap layanan wajib dikaitkan dengan `category_id` yang valid dan memiliki kode unik (`code`) sebagai SKU layanan.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`
- `Content-Type`: `application/json`

### Parameters :

Bagian ini mendefinisikan aturan data yang dikirim melalui Request Body.

| Key            | Tipe   | In   | Deskripsi                                           |
| -------------- | ------ | ---- | --------------------------------------------------- |
| category_id    | Int    | Body | ID Kategori yang sudah terdaftar di database.       |
| code           | String | Body | Kode unik layanan (Contoh: `SVC-LKR-1`).            |
| service_name   | String | Body | Nama layanan (Contoh: `Kiloan Regular`).            |
| unit           | Enum   | Body | Satuan layanan (`kg` atau `pcs`).                   |
| price          | Int    | Body | Harga layanan dalam nominal penuh (Contoh: `7000`). |
| duration_hours | Int    | Body | Durasi pengerjaan dalam satuan jam (Contoh: `72`).  |

```
{
  "category_id": 1,
  "code": "SVC-LKR-1",
  "service_name": "Layanan Cuci Kiloan Regular",
  "unit": "kg",
  "price": 7000,
  "duration_hours": 72
}
```

### Request Body :

Objek JSON berisi detail layanan.

```json
{
  "category_id": 1,
  "code": "SVC-LKR-1",
  "service_name": "Kiloan Regular",
  "unit": "kg",
  "price": 7000,
  "duration_hours": 72
}
```

### Responses Body :

#### ✅ 201 Created

Layanan baru berhasil didaftarkan secara sukses.

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
    "duration_hours": 72,
    "is_active": 1,
    "created_at": "2026-01-13 18:00:00",
    "updated_at": null
  }
}
```

#### ⚠️ 400 Bad Request

Input tidak valid (misal: harga negatif atau field wajib kosong).

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "price": "Price cannot be negative",
      "unit": "Unit must be one of: kg, pcs"
    }
  }
}
```

#### ⚠️ 401 Unauthorized

Token tidak valid, kadaluarsa, atau tidak disertakan.

```json
{
  "success": false,
  "message": "Invalid or missing access token",
  "data": {
    "error_code": "UNAUTHORIZED_ACCESS",
    "errors": null
  }
}
```

#### 🚫 403 Forbidden

Akses ditolak karena role Anda bukan **Owner**.

```json
{
  "success": false,
  "message": "Your role does not have permission",
  "data": {
    "error_code": "FORBIDDEN_ACCESS",
    "errors": null
  }
}
```

#### 🚫 404 Not Found

404 terjadi ketika `category_id` yang dikirimkan tidak ditemukan di database.

```json
{
  "success": false,
  "message": "Category not found",
  "data": {
    "error_code": "RESOURCE_NOT_FOUND",
    "errors": null
  }
}
```

#### 🚫 409 Conflict

Kode layanan atau nama layanan sudah digunakan sebelumnya.

```json
{
  "success": false,
  "message": "Data already exists",
  "data": {
    "error_code": "DUPLICATE_DATA",
    "errors": {
      "code": "Service code 'SVC-LKR-1' already exists"
    }
  }
}
```

#### 🚫 429 Too Many Requests

Terlalu banyak permintaan dalam waktu singkat.

```json
{
  "success": false,
  "message": "Too many requests, please try again later",
  "data": {
    "error_code": "RATE_LIMIT_EXCEEDED",
    "errors": null
  }
}
```

#### 🔥 500 Internal Server Error

Kegagalan sistem saat proses penyimpanan data layanan ke database.

```json
{
  "success": false,
  "message": "An unexpected server error occurred",
  "data": {
    "error_code": "INTERNAL_SERVER_ERROR",
    "errors": null
  }
}
```

---

## Endpoint : `GET /services`

#### Description :

Endpoint ini digunakan untuk mengambil daftar seluruh layanan laundry. Data dikirimkan secara Nested (Objek Kategori ada di dalam Objek Layanan) untuk memudahkan Kasir melihat jenis layanan tanpa perlu melakukan request tambahan ke endpoint kategori. Dilengkapi dengan **Pagination** dan **Filtering** yang komprehensif.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner, cashier`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`

### Parameters :

bagan ini mendefinisikan parameter query yang dapat digunakan untuk pagination dan filtering.

| Key          | Tipe   | In    | Default     | Deskripsi                                           |
| ------------ | ------ | ----- | ----------- | --------------------------------------------------- |
| page         | Int    | Query | 1           | Halaman data yang ingin ditampilkan.                |
| limit        | Int    | Query | 10, Max 100 | Jumlah data per halaman.                            |
| category_id  | Int    | Query | -           | (Optional) Filter berdasarkan ID Kategori tertentu. |
| code         | String | Query | -           | (Optional) Cari berdasarkan Kode SVC Layanan.       |
| service_name | String | Query | -           | (Optional) Cari berdasarkan nama (Partial search).  |
| is_active    | Int    | Query | 1           | (Optional) Filter status: 1 (Aktif) atau 0 (Arsip). |

```
GET /api/services?page=1&limit=10&is_active=1&category_id=1
```

### Request Body :

```json
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK

Data layanan berhasil diambil secara sukses beserta metadata dan info kategori.

```json
{
  "success": true,
  "message": "Data retrieved successfully",
  "data": [
    {
      "id": 1,
      "code": "SVC-LKR-1",
      "service_name": "Layanan Cuci Kiloan Regular",
      "unit": "kg",
      "price": 7000,
      "duration_hours": 72,
      "is_active": 1,
      "created_at": "2025-12-28 07:24:03",
      "updated_at": null,
      "category": {
        "id": 1,
        "category_name": "Layanan Kiloan",
        "description": "Cuci pakaian sehari-hari dengan sistem kiloan"
      }
    },
    {
      "id": 2,
      "code": "SVC-LJSR-1",
      "service_name": "Layanan Cuci Jas Satuan Regular",
      "unit": "pcs",
      "price": 10000,
      "duration_hours": 72,
      "is_active": 1,
      "created_at": "2025-12-28 07:25:10",
      "updated_at": null,
      "category": {
        "id": 2,
        "category_name": "Layanan Satuan",
        "description": "Cuci pakaian sehari-hari dengan sistem satuan"
      }
    }
  ],
  "meta": {
    "current_page": 1,
    "per_page": 10,
    "total_items": 12,
    "total_pages": 2
  }
}
```

#### ⚠️ 400 Bad Request

Parameter query tidak valid.

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "category_id": "Must be a number"
    }
  }
}
```

#### ⚠️ 401 Unauthorized

Token tidak valid, kadaluarsa, atau tidak disertakan.

```json
{
  "success": false,
  "message": "Invalid or missing access token",
  "data": {
    "error_code": "UNAUTHORIZED_ACCESS",
    "errors": null
  }
}
```

#### 🚫 403 Forbidden

Jika staff atau courier mencoba memanggil endpoint ini, sistem akan menolak.

```json
{
  "success": false,
  "message": "Your role does not have permission",
  "data": {
    "error_code": "FORBIDDEN_ACCESS",
    "errors": null
  }
}
```

#### 🚫 429 Too Many Requests

Terlalu banyak permintaan dalam waktu singkat (Rate Limit).

```json
{
  "success": false,
  "message": "Too many requests, please try again later",
  "data": {
    "error_code": "RATE_LIMIT_EXCEEDED",
    "errors": null
  }
}
```

#### 🔥 500 Internal Server Error

Kegagalan sistem saat pengambilan data dari database.

```json
{
  "success": false,
  "message": "An unexpected server error occurred",
  "data": {
    "error_code": "INTERNAL_SERVER_ERROR",
    "errors": null
  }
}
```

---

## Endpoint : `GET /services/{id}`

#### Description :

Endpoint ini digunakan untuk mengambil detail data dari satu layanan spesifik secara mendalam. Data dikirimkan secara nested dengan objek kategori di dalamnya. Digunakan oleh Admin/Owner untuk proses validasi data sebelum melakukan pembaruan (Edit).

### Role Based Access Control (RBAC) :

- `Permissions`: `owner, cashier`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`

### Parameters :

Parameter ini disisipkan langsung di URL (contoh: /services/1).

| Key | Tipe | In   | Default | Deskripsi                                     |
| --- | ---- | ---- | ------- | --------------------------------------------- |
| id  | Int  | Path | -       | ID unik layanan yang ingin dilihat detailnya. |

```
GET /api/services/1
```

### Request Body :

```json
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK

Detail layanan berhasil diambil secara sukses.

```json
{
  "success": true,
  "message": "Data retrieved successfully",
  "data": {
    "id": 1,
    "code": "SVC-LKR-1",
    "service_name": "Layanan Cuci Kiloan Regular",
    "unit": "kg",
    "price": 7000,
    "duration_hours": 72,
    "is_active": 1,
    "created_at": "2025-12-28 07:24:03",
    "updated_at": null,
    "category": {
      "id": 1,
      "category_name": "Layanan Kiloan",
      "description": "Cuci pakaian sehari-hari dengan sistem kiloan"
    }
  }
}
```

#### ⚠️ 400 Bad Request

Format ID pada URL tidak valid (misal: /services/abc).

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "id": "ID must be a valid number"
    }
  }
}
```

#### ⚠️ 401 Unauthorized

Token tidak valid, expired, atau tidak dikirim.

```json
{
  "success": false,
  "message": "Invalid or missing access token",
  "data": {
    "error_code": "UNAUTHORIZED_ACCESS",
    "errors": null
  }
}
```

#### 🚫 403 Forbidden

Akses ditolak karena Role pengguna (Staff/Courier) tidak memiliki izin untuk melihat master data layanan.

```json
{
  "success": false,
  "message": "Your role does not have permission",
  "data": {
    "error_code": "FORBIDDEN_ACCESS",
    "errors": null
  }
}
```

#### 🚫 404 Not Found

ID layanan tidak ditemukan di database.

```json
{
  "success": false,
  "message": "Service not found",
  "data": {
    "error_code": "RESOURCE_NOT_FOUND",
    "errors": null
  }
}
```

#### 🚫 429 Too Many Requests

Terlalu banyak permintaan dalam waktu singkat (Rate Limit).

```json
{
  "success": false,
  "message": "Too many requests, please try again later",
  "data": {
    "error_code": "RATE_LIMIT_EXCEEDED",
    "errors": null
  }
}
```

#### 🔥 500 Internal Server Error

Terjadi kesalahan di sisi server.

```json
{
  "success": false,
  "message": "An unexpected server error occurred",
  "data": {
    "error_code": "INTERNAL_SERVER_ERROR",
    "errors": null
  }
}
```

---

## Endpoint : `PUT /services/{id}`

#### Description :

Endpoint ini digunakan oleh Owner untuk memperbarui data layanan laundry. Sistem akan memvalidasi keunikan `code` (SKU) jika terjadi perubahan kode, serta memastikan `category_id` yang baru (jika dipindah kategori) benar-benar tersedia di database. Field yang tidak dikirim di request body akan tetap menggunakan nilai lama.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`
- `Content-Type`: `application/json`

### Parameters :

bagian ini mendefinisikan parameter path yang diperlukan untuk mengidentifikasi layanan yang akan diperbarui.

| Key | Tipe | In   | Default | Deskripsi                                                 |
| --- | ---- | ---- | ------- | --------------------------------------------------------- |
| id  | Int  | Path | -       | ID unik (Primary Key) dari layanan yang ingin diperbarui. |

```
PUT /api/services/1
```

### Request Body :

Objek JSON berisi data yang ingin diubah

```json
{
  "category_id": 1,
  "code": "SVC-LKR-1",
  "service_name": "Kiloan Regular",
  "unit": "kg",
  "price": 7500,
  "duration_hours": 72,
  "is_active": 1
}
```

### Responses Body :

#### ✅ 200 OK

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
    "price": 7500,
    "duration_hours": 72,
    "is_active": 1,
    "created_at": "2025-12-28 07:24:03",
    "updated_at": "2026-01-15 17:15:00"
  }
}
```

#### ⚠️ 400 Bad Request

Input tidak valid atau format salah.

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "price": "Price cannot be negative",
      "duration_hours": "Duration must be greater than 0"
    }
  }
}
```

#### ⚠️ 401 Unauthorized

Token tidak valid, expired, atau tidak dikirim.

```json
{
  "success": false,
  "message": "Invalid or missing access token",
  "data": {
    "error_code": "UNAUTHORIZED_ACCESS",
    "errors": null
  }
}
```

#### 🚫 403 Forbidden

Token valid, tetapi role pengguna bukan Owner.

```json
{
  "success": false,
  "message": "Your role does not have permission",
  "data": {
    "error_code": "FORBIDDEN_ACCESS",
    "errors": null
  }
}
```

#### 🚫 404 Not Found

Terjadi jika ID Service tidak ditemukan, atau `category_id` baru tidak terdaftar.

```json
{
  "success": false,
  "message": "Service not found",
  "data": {
    "error_code": "RESOURCE_NOT_FOUND",
    "errors": null
  }
}
```

#### 🚫 409 Conflict

Gagal update karena kode layanan baru sudah dipakai layanan lain.

```json
{
  "success": false,
  "message": "Data already exists",
  "data": {
    "error_code": "DUPLICATE_DATA",
    "errors": {
      "code": "Service code 'SVC-LKR-1' already exists"
    }
  }
}
```

#### 🚫 429 Too Many Requests

Terlalu banyak permintaan dalam waktu singkat.

```json
{
  "success": false,
  "message": "Too many requests, please try again later.",
  "data": {
    "error_code": "RATE_LIMIT_EXCEEDED",
    "errors": null
  }
}
```

#### 🔥 500 Internal Server Error

Terjadi kesalahan di sisi server.

```json
{
  "success": false,
  "message": "An unexpected server error occurred",
  "data": {
    "error_code": "INTERNAL_SERVER_ERROR",
    "errors": null
  }
}
```

## Endpoint : `DELETE /services/{id}`

### Description :

Endpoint ini digunakan oleh Owner untuk menonaktifkan layanan laundry secara logika (**Soft Delete**). Status layanan diubah menjadi `is_active = 0`. Hal ini memastikan integritas data pada transaksi lama tetap terjaga, namun layanan tersebut tidak akan muncul lagi sebagai pilihan pada pembuatan transaksi baru.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`

### Parameters :

Identifikasi layanan yang akan dinonaktifkan melalui Path Parameter.

| Key | Tipe | In   | Default | Deskripsi                                               |
| --- | ---- | ---- | ------- | ------------------------------------------------------- |
| id  | Int  | Path | -       | ID unik (Primary Key) layanan yang ingin dinonaktifkan. |

```
DELETE /api/services/1
```

### Request Body :

```json
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK

Layanan berhasil dinonaktifkan.

```json
{
  "success": true,
  "message": "Service deleted successfully",
  "data": {
    "id": 1
  }
}
```

#### ⚠️ 400 Bad Request

Format ID pada URL bukan angka.

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "id": "ID must be a valid number"
    }
  }
}
```

#### ⚠️ 401 Unauthorized

Token tidak valid, expired, atau tidak dikirim.

```json
{
  "success": false,
  "message": "Invalid or missing access token",
  "data": {
    "error_code": "UNAUTHORIZED_ACCESS",
    "errors": null
  }
}
```

#### 🚫 403 Forbidden

Token valid, tapi yang melakukan request bukan Owner.

```json
{
  "success": false,
  "message": "Your role does not have permission",
  "data": {
    "error_code": "FORBIDDEN_ACCESS",
    "errors": null
  }
}
```

#### 🚫 404 Not Found

Layanan dengan ID tersebut tidak ditemukan di database.

```json
{
  "success": false,
  "message": "Service not found",
  "data": {
    "error_code": "RESOURCE_NOT_FOUND",
    "errors": null
  }
}
```

#### 🚫 429 Too Many Requests

Terlalu banyak permintaan dalam waktu singkat.

```json
{
  "success": false,
  "message": "Too many requests, please try again later",
  "data": {
    "error_code": "RATE_LIMIT_EXCEEDED",
    "errors": null
  }
}
```

#### 🔥 500 Internal Server Error

Terjadi kesalahan tak terduga di server.

```json
{
  "success": false,
  "message": "An unexpected server error occurred",
  "data": {
    "error_code": "INTERNAL_SERVER_ERROR",
    "errors": null
  }
}
```
