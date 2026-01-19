# LAUNDRY MANAGEMENT SYSTEM — API SPECIFICATION

## SERVICE CATEGORIES MODULE SPECIFICATION

---

## Endpoint : `POST /categories`

### Description :

Endpoint ini digunakan oleh **Owner** untuk membuat kategori layanan baru (Contoh: "Kiloan", "Satuan", "Dry Clean"). Backend akan memastikan nama kategori belum pernah ada sebelumnya (Unique) sebelum menyimpannya ke database.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`
- `Content-Type`: `application/json`

### Parameters :

Bagian ini mendefinisikan aturan main data yang harus dikirimkan oleh klien melalui Request Body.

| Key           | Tipe   | In   | Default | Deskripsi                                              |
| ------------- | ------ | ---- | ------- | ------------------------------------------------------ |
| category_name | String | Body | -       | Nama unik kategori layanan (Contoh: "Layanan Kiloan"). |
| description   | String | Body | -       | Deskripsi singkat tentang kategori layanan (Opsional). |

```
{
  "category_name": "Layanan Kiloan",
  "description": "Cuci pakaian sehari-hari dihitung per kilogram"
}
```

### Request Body :

Objek JSON berisi nama kategori dan deskripsinya.

```json
{
  "category_name": "Layanan Kiloan",
  "description": "Cuci pakaian sehari-hari dihitung per kilogram"
}
```

### Responses Body :

#### ✅ 201 Created

Kategori layanan berhasil didaftarkan secara sukses.

```json
{
  "success": true,
  "message": "Category created successfully",
  "data": {
    "id": 1,
    "category_name": "Laundry Kiloan",
    "description": "Cuci pakaian sehari-hari dihitung per kilogram",
    "is_active": 1,
    "created_at": "2026-01-13 16:00:00",
    "updated_at": null
  }
}
```

#### ⚠️ 400 Bad Request

Input tidak lolos validasi (misal: category_name kosong).

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "category_name": "Category name is required"
    }
  }
}
```

#### ⚠️ 401 Unauthorized

Token tidak valid, expired, atau tidak disertakan.

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

Akses ditolak karena pengguna bukan Role **Owner**.

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

#### 🚫 409 Conflict

Nama kategori sudah digunakan sebelumnya di database.

```json
{
  "success": false,
  "message": "Data already exists",
  "data": {
    "error_code": "DUPLICATE_DATA",
    "errors": {
      "category_name": "Category 'Layanan Kiloan' already exists"
    }
  }
}
```

#### 🚫 429 Too Many Requests

Terjadi jika terlalu banyak permintaan yang dikirim dalam waktu singkat, memicu mekanisme rate limiting.

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

Kesalahan teknis saat penyimpanan ke database.

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

## Endpoint : `GET /categories`

### Description :

Endpoint ini digunakan untuk mengambil daftar kategori layanan secara terorganisir. Akses dibatasi hanya untuk staf administratif (Owner & Cashier) guna keperluan manajemen layanan dan pembuatan transaksi baru. Dilengkapi dengan **Pagination** dan **Filtering**.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner`, `cashier`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`

### Parameters :

Bagian ini mendefinisikan filter pencarian melalui Query String.

| Key           | Tipe   | In    | Default  | Deskripsi                                           |
| ------------- | ------ | ----- | -------- | --------------------------------------------------- |
| page          | Int    | Query | 1        | Menentukan halaman data.                            |
| limit         | Int    | Query | 10       | Jumlah data per halaman.                            |
| category_name | String | Query | -        | (Optional) Filter berdasarkan nama kategori.        |
| is_active     | Int    | Query | Show all | (Optional) Filter status: 1 (Aktif) atau 0 (Arsip). |

```
GET /api/categories?page=1&limit=10&is_active=1&category_name=Kiloan
```

### Request Body :

```json
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK

Data kategori berhasil diambil secara sukses beserta metadata.

```json
{
  "success": true,
  "message": "Data retrieved successfully",
  "data": [
    {
      "id": 1,
      "category_name": "Layanan Kiloan",
      "description": "Cuci pakaian sehari-hari dihitung per kilogram",
      "is_active": 1,
      "created_at": "2025-12-28 07:24:03",
      "updated_at": null
    },
    {
      "id": 2,
      "category_name": "Layanan Satuan",
      "description": "Kategori untuk layanan per item",
      "is_active": 1,
      "created_at": "2025-12-28 07:25:10",
      "updated_at": null
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

Parameter query tidak valid (misal: format angka salah).

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "page": "Must be a number"
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

Permintaan terlalu sering dalam waktu singkat (Rate Limiting).

```json
{
  "success": false,
  "message": "Too many requests, please try again later",
  "data": {
    "error_code": "TOO_MANY_REQUESTS",
    "errors": null
  }
}
```

#### 🔥 500 Internal Server Error

Kesalahan teknis pada internal server.

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

## Endpoint : `GET /categories/{id}`

#### Description :

Endpoint ini digunakan untuk mengambil detail data dari satu kategori spesifik berdasarkan ID uniknya. Biasanya digunakan oleh Frontend untuk mengisi data pada formulir edit atau melihat detail kategori sebelum melakukan perubahan.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner`, `cashier`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`

### Parameters :

Menargetkan ID unik melalui jalur URL (Path Parameter).

| Key | Tipe | In   | Default | Deskripsi                                                         |
| --- | ---- | ---- | ------- | ----------------------------------------------------------------- |
| id  | Int  | Path | -       | ID Unik (Primary Key) dari kategori yang ingin dilihat detailnya. |

```
GET /api/categories/1
```

### Request Body :

```json
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK (Success)

Detail kategori ditemukan dan dikembalikan dalam bentuk Object `{}`.

```json
{
  "success": true,
  "message": "Data retrieved successfully",
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

#### ⚠️ 400 Bad Request

Format ID yang dikirimkan pada URL bukan merupakan angka.

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

Token tidak valid, expired, atau header Authorization tidak disertakan.

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

Akses ditolak karena Role pengguna (Staff/Courier) tidak memiliki izin untuk melihat master data kategori.

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

Format ID benar (angka), namun data dengan ID tersebut tidak ditemukan di database.

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

#### 🚫 429 Too Many Requests

Terjadi jika terlalu banyak permintaan dalam waktu singkat, memicu mekanisme rate limiting.

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
  "message": "An unexpected error occurred",
  "data": {
    "error_code": "INTERNAL_SERVER_ERROR",
    "errors": null
  }
}
```

---

## Endpoint : `PUT /categories/{id}`

### Description :

Endpoint ini digunakan oleh **Owner** untuk memperbarui data kategori. Jika nama kategori diubah, sistem akan memvalidasi keunikannya terhadap data lain di database agar tidak terjadi duplikasi. Field yang tidak dikirim di Request Body akan tetap menggunakan nilai lama.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`
- `Content-Type`: `application/json`

### Parameters :

Bagian ini mendefinisikan ID kategori yang ingin diubah, disisipkan langsung di URL.

| Key | Tipe | In   | Default | Deskripsi                                             |
| --- | ---- | ---- | ------- | ----------------------------------------------------- |
| id  | Int  | Path | -       | ID unik (Primary Key) kategori yang ingin diperbarui. |

```
PUT /api/categories/1
```

### Request Body :

Objek JSON berisi data yang ingin diubah.

```json
{
  "category_name": "Layanan Satuan",
  "description": "Kategori untuk layanan cuci per item",
  "is_active": 1
}
```

### Responses Body :

#### ✅ 200 OK

Data kategori berhasil diperbarui secara sukses.

```json
{
  "success": true,
  "message": "Category updated successfully",
  "data": {
    "id": 1,
    "category_name": "Layanan Satuan",
    "description": "Kategori untuk layanan cuci per item",
    "is_active": 1,
    "created_at": "2025-12-28 07:24:03",
    "updated_at": "2026-01-13 17:15:00"
  }
}
```

#### ⚠️ 400 Bad Request

Input tidak lolos validasi atau format ID pada URL salah.

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "category_name": "Category name is required"
    }
  }
}
```

#### ⚠️ 401 Unauthorized

Token tidak valid, expired, atau tidak disertakan.

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

Akses ditolak karena pengguna bukan Role **Owner**.

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

ID kategori tidak ditemukan di database.

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

Gagal update karena nama kategori baru sudah digunakan oleh kategori lain.

```json
{
  "success": false,
  "message": "Data already exists",
  "data": {
    "error_code": "DUPLICATE_DATA",
    "errors": {
      "category_name": "Category name 'Layanan Satuan' already taken"
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

Kegagalan sistem saat proses pembaruan data di database.

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

## Endpoint : `DELETE /categories/{id}`

### Description :

Endpoint ini digunakan oleh **Owner** untuk menghapus kategori secara logika (**Soft Delete**). Sistem akan mengubah nilai `is_active` menjadi `0`. Hal ini memastikan data transaksi lama yang menggunakan kategori ini tetap memiliki referensi yang valid, namun kategori tersebut tidak akan muncul lagi sebagai pilihan di menu operasional.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`

### Parameters :

Identitas kategori yang akan dinonaktifkan ditentukan melalui Path Parameter.

| Key | Tipe | In   | Default | Deskripsi                                  |
| --- | ---- | ---- | ------- | ------------------------------------------ |
| id  | Int  | Path | -       | ID unik kategori yang ingin dinonaktifkan. |

```
DELETE /api/categories/1
```

### Request Body :

```
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK

Kategori berhasil dinonaktifkan secara sukses.

```json
{
  "success": true,
  "message": "Category deleted successfully",
  "data": {
    "id": 1
  }
}
```

#### ⚠️ 400 Bad Request

Format ID pada URL tidak valid (bukan angka).

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

Akses ditolak karena Anda bukan **Owner**.

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

Data kategori tidak ditemukan di database.

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

Kegagalan teknis saat proses pembaruan status di database.

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
