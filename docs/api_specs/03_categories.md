# LAUNDRY MANAGEMENT SYSTEM — API SPECIFICATION

## CATEGORIES MODULE SPECIFICATION

---

## Endpoint : `POST /categories`

### Headers :

- `Authorization: Bearer <access_token>` (Required)

#### Description :

Endpoint ini digunakan untuk membuat Kategori Layanan baru (misal: "Satuan", "Kiloan", "Dry Clean"). Kategori ini nanti akan menjadi pengelompokan untuk layanan-layanan laundry.

### Request Body :

Mengirimkan object data kategori baru. Field category_name wajib diisi (Max 150 karakter) dan harus unik. Field description bersifat opsional (Max 255 karakter).

```json
{
  "category_name": "Laundry Kiloan",
  "description": "Kategori untuk layanan cuci kiloan"
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
    "description": "Kategori untuk layanan cuci kiloan",
    "is_active": true,
    "created_at": "2024-02-20T10:00:00Z",
    "updated_at": "2024-03-15T14:30:00Z"
  }
}
```

#### ⚠️ 400 Bad Request (Validation Error)

Input tidak valid (misal: nama kosong atau melebihi batas karakter).

```json
{
  "success": false,
  "message": "Input Validation Failed",
  "data": {
    "category_name": "Category name is required"
  }
}
```

#### 🚫 401 Unauthorized (Invalid Token)

Token tidak valid atau expired.

```json
{
  "success": false,
  "message": "Invalid or missing access token",
  "data": null
}
```

#### ⛔ 403 Forbidden (Insufficient Permissions)

User bukan Owner mencoba membuat kategori.

```json
{
  "success": false,
  "message": "You do not have permission to perform this action",
  "data": null
}
```

#### 💥 409 Conflict (Duplicate Data)

Nama kategori sudah terdaftar sebelumnya.

```json
{
  "success": false,
  "message": "Category name already exists",
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

## Endpoint : `GET /categories`

### Headers :

- `Authorization: Bearer <access_token>` (Required)

#### Description :

Endpoint ini digunakan untuk mengambil daftar semua kategori layanan. Dapat diakses oleh Semua User yang Login (Owner, Kasir, Staff) karena dibutuhkan untuk operasional transaksi. Mendukung pagination dan filter pencarian.

### Parameters :

Parameter ini ditempel di URL (Query String).

| Key       | Tipe   | Default     | Deskripsi                                                      |
| --------- | ------ | ----------- | -------------------------------------------------------------- |
| page      | Int    | 1           | Halaman ke berapa yang mau diambil.                            |
| limit     | Int    | 10, Max 100 | Jumlah data per halaman.                                       |
| name      | String | -           | (Opsional) Cari kategori berdasarkan nama.                     |
| is_active | String | Show all    | (Opsional) Filter status: "true" (Aktif) atau "false" (Arsip). |

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
      "description": "Kategori untuk layanan cuci kiloan",
      "is_active": true,
      "created_at": "2024-02-20T10:00:00Z",
      "updated_at": "2024-03-15T14:30:00Z"
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
    "total_pages": 5,
    "total_items": 12,
    "per_page": 10
  }
}
```

#### ⚠️ 400 Bad Request (Invalid Parameters)

Terjadi jika parameter pagination salah format.

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

## Endpoint : `GET /categories/{id}`

### Headers :

- `Authorization: Bearer <access_token>` (Required)

#### Description :

Endpoint ini digunakan untuk melihat detail data satu kategori spesifik berdasarkan ID-nya. Dapat diakses oleh semua user yang login (Owner, Kasir, Staff).

### Parameters :

Parameter ini adalah Path Variable.

| Key | Tipe | Default | Deskripsi                                      |
| --- | ---- | ------- | ---------------------------------------------- |
| id  | Int  | -       | ID Unik kategori yang ingin dilihat detailnya. |

### Request Body :

```json
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK (Success)

Data kategori ditemukan. Mengembalikan object detail kategori (Tanpa pagination).

```json
{
  "success": true,
  "message": "Category detail retrieved successfully",
  "data": {
    "id": 1,
    "category_name": "Laundry Kiloan",
    "description": "Kategori untuk layanan cuci kiloan",
    "is_active": true,
    "created_at": "2024-02-20T10:00:00Z",
    "updated_at": "2024-03-15T14:30:00Z"
  }
}
```

#### ⚠️ 400 Bad Request (Invalid ID Format)

Terjadi jika ID di URL bukan angka (misal: /categories/abc).

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

Token tidak valid, expired, atau user belum login.

```json
{
  "success": false,
  "message": "Invalid or missing access token",
  "data": null
}
```

#### 🔍 404 Not Found (Data Missing)

ID kategori valid (angka), tapi datanya tidak ditemukan di database.

```json
{
  "success": false,
  "message": "Category not found",
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

## Endpoint : `PUT /categories/{id}`

### Headers :

- `Authorization: Bearer <access_token>` (Required)

#### Description :

Endpoint ini digunakan untuk memperbarui data kategori. Mendukung Partial Update. Hanya Owner yang memiliki izin untuk mengakses endpoint ini.

### Parameters :

Parameter ini adalah Path Variable.

| Key | Tipe | Default | Deskripsi                           |
| --- | ---- | ------- | ----------------------------------- |
| id  | Int  | -       | ID Unik kategori yang ingin diedit. |

### Request Body :

Semua field bersifat Optional (Boleh dikirim salah satu saja).

```json
{
  "category_name": "Laundry Satuan",
  "description": "Kategori untuk layanan cuci satuan",
  "is_active": true
}
```

### Responses Body :

#### ✅ 200 OK (Success)

Data berhasil diperbarui.

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

Format input salah atau ID bukan angka.

```json
{
  "success": false,
  "message": "Input Validation Failed",
  "data": {
    "id": "ID must be a number",
    "category_name": "Category name is required"
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

User bukan Owner mencoba mengedit kategori.

```json
{
  "success": false,
  "message": "You do not have permission to perform this action",
  "data": null
}
```

#### 404 Not Found (Category Not Found)

ID kategori tidak ditemukan di database.

```json
{
  "success": false,
  "message": "Category not found",
  "data": null
}
```

#### 💥 409 Conflict (Duplicate Data)

User mencoba mengganti nama kategori menjadi nama yang sudah dipakai kategori lain.

```json
{
  "success": false,
  "message": "Category name already exists",
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

## Endpoint : `DELETE /categories/{id}`

### Headers :

- `Authorization: Bearer <access_token>` (Required)

#### Description :

Endpoint ini digunakan untuk Menonaktifkan (Soft Delete) kategori layanan. Catatan Penting: Kategori tidak dihapus permanen agar riwayat transaksi masa lalu tetap aman. Namun, kategori ini akan hilang dari menu pemilihan saat Kasir membuat transaksi baru. Hanya Owner yang boleh melakukan ini.

### Parameters :

Parameter ini adalah Path Variable.

| Key | Tipe | Default | Deskripsi                                  |
| --- | ---- | ------- | ------------------------------------------ |
| id  | Int  | -       | ID Unik kategori yang ingin dinonaktifkan. |

### Request Body :

```json
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK (Success)

Kategori berhasil dinonaktifkan (Status is_active menjadi false).

```json
{
  "success": true,
  "message": "Category deactivated successfully",
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

User bukan Owner mencoba menghapus kategori.

```json
{
  "success": false,
  "message": "You do not have permission to perform this action",
  "data": null
}
```

#### 🔍 404 Not Found (Category Missing)

ID kategori tidak ditemukan di database.

```json
{
  "success": false,
  "message": "Category not found",
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
