# LAUNDRY MANAGEMENT SYSTEM — API SPECIFICATION

## USERS MODULE SPECIFICATION

---

## Endpoint : `POST /users`

### Headers :

- `Authorization: Bearer <access_token>` (Required) - Hanya Role Owner yang diizinkan.
- `Content-Type`: `application/json` (Required)

#### Description :

Endpoint ini digunakan untuk menambahkan pengguna baru (Staff, Cashier, Courier). Endpoint ini dilindungi dan hanya bisa diakses oleh Owner. Data yang dikirim akan divalidasi keunikannya (Email & Username), di-hash password-nya, lalu disimpan ke database.

### Request Body :

Data lengkap pengguna baru.

```json
{
  "full_name": "Siti Aminah", // Required, Max 100 chars
  "username": "sitiaminah", // Required, Unique, No spaces
  "email": "sitiaminah@gmail.com", // Required, Valid Email format
  "password": "rahasia123", // Required, Min 8 chars
  "phone_number": "082345678901", // Required, Numeric only
  "role": "cashier" // Optional (Default: "staff" if empty)
}
```

### Responses Body :

#### ✅ 201 Created (Success)

User berhasil dibuat. Mengembalikan data user yang baru dibuat (tanpa password).

```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": 2,
    "full_name": "Siti Aminah",
    "username": "sitiaminah",
    "email": "sitiaminah@gmail.com",
    "role": "cashier",
    "phone_number": "082345678901",
    "is_active": true,
    "created_at": "2024-02-01T10:00:00Z"
  }
}
```

#### ⚠️ 400 Bad Request (Validation Error)

Format input salah (Email tidak valid, Password kependekan, field wajib kosong, atau format Enum role salah).

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "email": "Invalid email format",
    "password": "Password must be at least 8 characters",
    "role": "Invalid role selected"
  }
}
```

#### 🚫 401 Unauthorized (Invalid Token)

Token tidak valid, expired, atau tidak dikirim di Header.

```json
{
  "success": false,
  "message": "Invalid or missing access token",
  "data": null
}
```

#### ⛔ 403 Forbidden (Insufficient Permissions)

Token valid, tetapi Role pengguna bukan Owner. (Contoh: Staff mencoba mendaftarkan user baru).

```json
{
  "success": false,
  "message": "You do not have permission to perform this action",
  "data": null
}
```

#### 🚫 409 Conflict (Duplicate Data)

Data gagal disimpan karena Username atau Email sudah terdaftar sebelumnya.

```json
{
  "success": false,
  "message": "Data already exists",
  "data": {
    "username": "Username 'sitiaminah' is already taken",
    "email": "Email already registered"
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

## Endpoint : `GET /users`

### Headers :

- `Authorization: Bearer <access_token>` (Required) - Hanya Role Owner yang diizinkan.

#### Description :

Endpoint ini digunakan oleh Owner untuk melihat daftar seluruh akun karyawan (Staff, Cashier, Courier). Karena data pengguna bisa berjumlah banyak, endpoint ini menerapkan sistem Pagination dan mendukung Filtering (pencarian berdasarkan nama, role, atau status aktif).

### Parameters :

Parameter ini dikirim melalui URL (contoh: /users?page=2&limit=10&role=cashier).

| Key       | Tipe   | Default  | Deskripsi                                                                        |
| --------- | ------ | -------- | -------------------------------------------------------------------------------- |
| page      | Int    | 1        | Menentukan halaman data yang ingin diambil.                                      |
| limit     | Int    | 10       | Jumlah data yang ditampilkan per halaman.                                        |
| name      | String | -        | (Optional) Filter pencarian berdasarkan nama lengkap atau username.              |
| role      | String | -        | (Optional) Filter berdasarkan jabatan (contoh: cashier, staff, courier).         |
| is_active | String | Show all | (Optional) Filter status. Isi true untuk yang aktif, false untuk yang non-aktif. |

### Request Body :

```json
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK (Success)

Daftar pengguna berhasil diambil. Response mencakup array data pengguna dan objek meta untuk informasi halaman.

```json
{
  "success": true,
  "message": "User retrieved successfully",
  "data": [
    {
      "id": 1,
      "full_name": "Farhan Rizki Maulana",
      "username": "farhanrizkimln",
      "email": "farhanrizki@gmail.com",
      "role": "owner",
      "phone_number": "081234567890",
      "is_active": true,
      "last_login_at": "2025-02-01T08:00:00Z", // <--- PENTING: Owner sedang online
      "created_at": "2024-01-30T10:00:00Z",
      "updated_at": "2025-02-01T08:00:00Z" // <--- PENTING: Data terakhir berubah
    },
    {
      "id": 2,
      "full_name": "Siti Aminah",
      "username": "sitiaminah",
      "email": "sitiaminah@gmail.com",
      "role": "cashier",
      "phone_number": "082345678901",
      "is_active": true,
      "last_login_at": "2025-02-01T08:00:00Z",
      "created_at": "2024-01-31T10:00:00Z",
      "updated_at": null // <--- INFO: Belum pernah Update Data
    }
  ],
  "meta": {
    "current_page": 1, // "Kamu sekarang ada di Halaman 1"
    "total_pages": 1, // "Totalnya ada 5 halaman lho"
    "total_items": 4, // "Total karyawanmu ada 50 orang"
    "per_page": 10 // "Di halaman ini, aku nampilin 10 orang"
  }
}
```

#### ⚠️ 400 Bad Request (Invalid Parameters)

Parameter query yang dikirim tidak valid (misalnya page berupa huruf, atau limit melebihi batas).

```json
{
  "success": false,
  "message": "Invalid query parameters",
  "data": {
    "page": "Must be a number",
    "limit": "Max limit is 100"
  }
}
```

#### ⚠️ 401 Unauthorized (Invalid Token)

Token tidak valid, expired, atau header Authorization tidak dikirim.

```json
{
  "success": false,
  "message": "Invalid or missing access token",
  "data": null
}
```

#### 🚫 403 Forbidden (Insufficient Permissions)

Token valid, tetapi role pengguna bukan `Owner`.

```json
{
  "success": false,
  "message": "You do not have permission to access this resource",
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

## Endpoint : `GET /users/{id}`

### Headers :

- `Authorization: Bearer <access_token>` (Required) - Hanya Role Owner yang diizinkan.

#### Description :

Endpoint ini digunakan untuk mengambil data lengkap satu pengguna secara spesifik berdasarkan id. Endpoint ini biasanya digunakan saat Owner mengklik tombol "Detail" atau "Edit" pada tabel daftar pengguna.

### Parameters :

Parameter ini disisipkan langsung di URL (contoh: /users/5).

| Key | Tipe | Default | Deskripsi                                                         |
| --- | ---- | ------- | ----------------------------------------------------------------- |
| id  | Int  | -       | ID Unik (Primary Key) dari pengguna yang ingin dilihat detailnya. |

### Request Body :

```json
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK (Success)

Data pengguna ditemukan dan dikembalikan.

```json
{
  "success": true,
  "message": "User retrieved successfully",
  "data": {
    "id": 1,
    "full_name": "Farhan Rizki Maulana",
    "username": "farhanrizkimln",
    "email": "farhanrizki@gmail.com",
    "role": "owner",
    "phone_number": "081234567890",
    "is_active": true,
    "last_login_at": "2025-02-01T08:00:00Z",
    "created_at": "2024-01-30T10:00:00Z",
    "updated_at": "2025-02-01T08:00:00Z"
  }
}
```

#### ⚠️ 400 Bad Request (Invalid ID Format)

Format ID yang dikirim bukan angka (misal: /users/abc).

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

Token tidak valid, expired, atau header Authorization tidak dikirim.

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
  "message": "You do not have permission to access this resource",
  "data": null
}
```

#### 404 Not Found (Data Missing)

ID memiliki format yang benar (angka), tetapi data pengguna dengan ID tersebut tidak ditemukan di database.

```json
{
  "success": false,
  "message": "User not found",
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

## Endpoint : `PUT /users/{id}`

### Headers :

- `Authorization: Bearer <access_token>` (Required)

  - Owner: Bisa mengakses ID siapapun.
  - Staff/Courier/Cashier: Hanya bisa mengakses ID miliknya sendiri (sesuai token).

- `Content-Type`: `application/json` (Required)

#### Description :

Endpoint ini digunakan untuk memperbarui data pengguna. Logika akses terbagi dua:

1. Owner: Bisa mengedit data seluruh karyawan, termasuk mengubah jabatan (role) dan menonaktifkan akun (is_active).
2. User Biasa: Hanya bisa mengedit data dirinya sendiri. Jika User Biasa mencoba mengubah role atau is_active, sistem akan mengabaikan perubahan tersebut (tetap menggunakan nilai lama) demi keamanan.

### Parameters :

Parameter ini disisipkan langsung di URL (contoh: /users/5).

| Key | Tipe | Default | Deskripsi                       |
| --- | ---- | ------- | ------------------------------- |
| id  | Int  | -       | ID Unik user yang ingin diedit. |

### Request Body :

Objek JSON berisi data yang ingin diubah.

```json
{
  "full_name": "Farhan Rizki Maulana", // Required
  "username": "farhanrizkimln", // Required, Unique
  "email": "farhanrizki@gmail.com", // Required, Valid Email
  "password": "barurahasia", // Optional (Isi jika ganti pass, kosong jika tetap)
  "phone_number": "081234567890", // Optional
  "role": "owner", // Optional (Hanya diproses jika requester adalah OWNER)
  "is_active": true // Optional (Hanya diproses jika requester adalah OWNER)
}
```

### Responses Body :

#### ✅ 200 OK (Success)

Data pengguna berhasil diperbarui.

```json
{
  "success": true,
  "message": "User updated successfully",
  "data": {
    "id": 1,
    "full_name": "Farhan Rizki Maulana",
    "username": "farhanrizkimln",
    "email": "farhanrizki@gmail.com",
    "role": "owner",
    "phone_number": "081234567890",
    "is_active": true,
    "last_login_at": "2025-02-01T08:00:00Z",
    "created_at": "2024-01-30T10:00:00Z",
    "updated_at": "2025-02-01T09:00:00Z" // Waktu berubah sesuai saat update
  }
}
```

#### ⚠️ 400 Bad Request (Validation Error)

Format input salah, password baru terlalu pendek, atau ID di URL bukan angka.

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "password": "Password must be at least 8 characters",
    "email": "Invalid email format"
  }
}
```

#### ⚠️ 401 Unauthorized (Invalid Token)

Token tidak valid, expired, atau header Authorization tidak dikirim.

```json
{
  "success": false,
  "message": "Invalid or missing access token",
  "data": null
}
```

#### 🚫 403 Forbidden (Access Denied)

Terjadi jika:

1. User biasa mencoba mengedit User ID orang lain.
2. User mencoba mengubah role/is_active (opsional, bisa juga di-ignore diam-diam).

```json
{
  "success": false,
  "message": "You are not allowed to edit this user's profile",
  "data": null
}
```

#### ❓ 404 Not Found (User not found)

ID valid secara format, tetapi data pengguna tidak ditemukan di database.

```json
{
  "success": false,
  "message": "User not found",
  "data": null
}
```

#### 💥 409 Conflict (Duplicate Data)

Gagal update karena Username atau Email baru yang dimasukkan sudah digunakan oleh pengguna lain.

```json
{
  "success": false,
  "message": "Data already exists",
  "data": {
    "email": "Email 'farhanrizki@gmail.com' is already taken by another user"
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

## Endpoint : `DELETE /users/{id}`

### Headers :

- `Authorization: Bearer <access_token>` (Required) - Hanya Owner yang diizinkan.

#### Description :

Endpoint ini digunakan untuk menonaktifkan pengguna (Soft Delete). Data pengguna TIDAK DIHAPUS dari database, melainkan statusnya diubah menjadi tidak aktif (is_active = 0). Ini memastikan riwayat transaksi yang pernah dilakukan oleh pengguna tersebut tetap tersimpan aman untuk laporan.

### Parameters :

Parameter ini disisipkan langsung di URL (contoh: /users/5).

| Key | Tipe    | Required | Deskripsi                              |
| --- | ------- | -------- | -------------------------------------- |
| id  | Integer | Yes      | ID Unik user yang ingin dinonaktifkan. |

### Request Body :

```json
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK (Success)

User berhasil dinonaktifkan.

```json
{
  "success": true,
  "message": "User deactivated successfully",
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

#### ❓ 404 Not Found (User Missing)

User dengan ID tersebut tidak ditemukan di database.

```json
{
  "success": false,
  "message": "User not found",
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
