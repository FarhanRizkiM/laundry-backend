# LAUNDRY MANAGEMENT SYSTEM — API SPECIFICATION

## USERS MODULE SPECIFICATION

---

## Endpoint : `POST /users`

### Headers :

- `Authorization: Bearer <access_token>` (Required)

#### Description :

Endpoint ini digunakan untuk pendaftaran pengguna baru (Sign Up). Data yang dikirim akan divalidasi, dan jika sukses, user baru akan disimpan ke database dengan status aktif.

### Request Body :

Mengirimkan data diri lengkap untuk pembuatan akun.

```json
{
  "full_name": "Farhan Rizki Maulana", // Required, Max 100 chars
  "username": "farhanrizkimln", // Required, Unique, No spaces
  "email": "farhanrizki@gmail.com", // Required, Valid Email format
  "password": "rahasia123", // Required, Min 8 chars
  "phone_number": "081234567890", // Required, Numeric only
  "role": "owner" // Optional (Default: "staff" if empty)
}
```

### Responses Body :

#### ✅ 201 Created (Success)

User berhasil dibuat. Mengembalikan data user tanpa password.

```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": 3,
    "full_name": "Farhan Rizki Maulana",
    "username": "farhanrizkimln",
    "email": "farhanrizki@gmail.com",
    "role": "owner",
    "phone_number": "081234567890",
    "is_active": true,
    "created_at": "2024-01-30T10:00:00Z"
  }
}
```

#### ⚠️ 400 Bad Request (Validation Error)

Format input salah (Email tidak valid, Password kependekan, atau field wajib kosong).

```json
{
  "success": false,
  "message": "Input Validation Failed",
  "data": {
    "email": "Invalid email format",
    "password": "Password must be at least 8 characters"
  }
}
```

#### 🚫 401 Unauthorized (Invalid Token)

Terjadi jika Token tidak valid, Expired, atau tidak dikirim di Header.

```json
{
  "success": false,
  "message": "Invalid or missing access token",
  "data": null
}
```

#### ⛔ 403 Forbidden (Insufficient Permissions)

User token valid, tapi Role-nya bukan Owner. (Misal: Staff iseng coba buat akun baru).

```json
{
  "success": false,
  "message": "You do not have permission to perform this action",
  "data": null
}
```

#### 🚫 409 Conflict (Duplicate Data)

Username atau Email sudah terdaftar di database.

```json
{
  "success": false,
  "message": "Username or Email already exists",
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

## Endpoint : `GET /users`

### Headers :

- `Authorization: Bearer <access_token>` (Required)

#### Description :

Endpoint ini digunakan untuk mengambil daftar semua user (karyawan) yang terdaftar. Hanya bisa diakses oleh Owner. Mendukung fitur Pagination agar data tidak dimuat sekaligus (berat).

### Parameters :

Parameter ini ditempel di URL (contoh: /users?page=1&role=cashier&is_active=true).

| Key       | Tipe   | Default     | Deskripsi                                                                          |
| --------- | ------ | ----------- | ---------------------------------------------------------------------------------- |
| page      | Int    | 1           | Halaman ke berapa yang mau diambil                                                 |
| limit     | Int    | 10, Max 100 | Jumlah data per halaman                                                            |
| name      | String | -           | (Opsional) Filter cari nama users                                                  |
| role      | String | -           | (Opsional) Filter jabatan (e.g. "cashier", "courier").                             |
| is_active | String | Show all    | (Opsional) Isi "true" (Aktif), "false" (Non-aktif). Jika kosong = Tampilkan Semua. |

### Request Body :

```json
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK (Success)

Berhasil mengambil daftar user. Disertai info pagination di field meta.

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
      "full_name": "Fakhri Aslam Muzhaffar",
      "username": "fakhriaslam123",
      "email": "fakhriaslam@gmail.com",
      "role": "staff",
      "phone_number": "082345678901",
      "is_active": true,
      "last_login_at": null, // <--- INFO: Belum pernah login
      "created_at": "2024-01-31T10:00:00Z",
      "updated_at": "2024-01-31T10:00:00Z"
    }
  ],
  "meta": {
    "current_page": 1, // "Kamu sekarang ada di Halaman 1"
    "total_pages": 5, // "Totalnya ada 5 halaman lho"
    "total_items": 50, // "Total karyawanmu ada 50 orang"
    "per_page": 10 // "Di halaman ini, aku nampilin 10 orang"
  }
}
```

#### ⚠️ 400 Bad Request (Invalid Parameters)

Terjadi jika parameter pagination/filter tidak sesuai format (misal: page berupa huruf, atau limit terlalu besar).

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

Token tidak valid atau expired.

```json
{
  "success": false,
  "message": "Invalid or missing access token",
  "data": null
}
```

#### 🚫 403 Forbidden (Insufficient Permissions)

User sudah login (punya token), tapi role-nya bukan Owner (misal: Kasir coba intip data user lain).

```json
{
  "success": false,
  "message": "You do not have permission to access this resource",
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

## Endpoint : `GET /users/{id}`

### Headers :

- `Authorization: Bearer <access_token>` (Required)

#### Description :

Endpoint ini digunakan untuk melihat detail data satu user spesifik berdasarkan ID-nya. Hanya Owner yang boleh melihat detail user lain.

### Parameters :

Parameter ini adalah Path Variable, artinya ditempel langsung sebagai bagian dari URL (contoh: /users/1).

| Key | Tipe | Default | Deskripsi                                  |
| --- | ---- | ------- | ------------------------------------------ |
| id  | Int  | -       | ID Unik user yang ingin dilihat detailnya. |

### Request Body :

```json
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK (Success)

Data user ditemukan. Mengembalikan object detail user.

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
    "last_login_at": "2025-02-01T08:00:00Z", // <--- PENTING: Owner sedang online
    "created_at": "2024-01-30T10:00:00Z",
    "updated_at": "2025-02-01T08:00:00Z" // <--- PENTING: Data terakhir berubah
  }
}
```

#### ⚠️ 400 Bad Request (Invalid ID Format)

Terjadi jika ID yang dimasukkan di URL bukan angka (misal: /users/abc).

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

#### 🚫 403 Forbidden (Insufficient Permissions)

User login (bukan Owner) mencoba melihat detail orang lain.

```json
{
  "success": false,
  "message": "You do not have permission to access this resource",
  "data": null
}
```

#### 404 Not Found (Data Missing)

[PENTING] ID valid secara format (angka), tapi datanya tidak ditemukan di database.

```json
{
  "success": false,
  "message": "User not found",
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

## Endpoint : `PUT /users/{id}`

### Headers :

- `Authorization: Bearer <access_token>` (Required)

#### Description :

Endpoint ini digunakan untuk memperbarui data user yang sudah ada. Mendukung Partial Update (Artinya: Client boleh hanya mengirim field yang ingin diubah saja, tidak perlu mengirim semua data).

### Parameters :

Parameter ini adalah Path Variable.

| Key | Tipe | Default | Deskripsi                       |
| --- | ---- | ------- | ------------------------------- |
| id  | Int  | -       | ID Unik user yang ingin diedit. |

### Request Body :

Semua field di bawah ini bersifat Optional (Boleh dikirim, boleh tidak).

```json
{
  "full_name": "Farhan Rizki Maulana",
  "username": "farhanrizkimln", // Harus unik (jika diubah)
  "email": "farhanrizki@gmail.com", // Harus unik (jika diubah)
  "password": "rahasia123", // Min 8 chars
  "phone_number": "081234567890",
  "role": "owner", // hanya owners yang bisa edit
  "is_active": true // hanya owners yang bisa edit
}
```

### Responses Body :

#### ✅ 200 OK (Success)

Data berhasil diperbarui. Mengembalikan data terbaru.

```json
{
  "success": true,
  "message": "User updated successfully",
  "data": {
    "id": 1,
    "full_name": "Farhan Rizki Maulana",
    "username": "farhanrizki123",
    "email": "farhanrizki@gmail.com",
    "role": "owner",
    "phone_number": "083456789012",
    "is_active": true,
    "last_login_at": "2025-02-01T08:00:00Z",
    "created_at": "2024-01-30T10:00:00Z",
    "updated_at": "2025-02-01T09:00:00Z" // Jam berubah
  }
}
```

#### ⚠️ 400 Bad Request (Validation Error)

Format input salah, ID bukan angka, atau password kependekan.

```json
{
  "success": false,
  "message": "Input Validation Failed",
  "data": {
    "password": "Password must be at least 8 characters",
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

#### 🚫 403 Forbidden (Insufficient Permissions)

User bukan Owner mencoba mengedit data.

```json
{
  "success": false,
  "message": "You do not have permission to perform this action",
  "data": null
}
```

#### 404 Not Found (User not found)

ID user yang mau diedit tidak ditemukan di database.

```json
{
  "success": false,
  "message": "User not found",
  "data": null
}
```

#### 💥 409 Conflict (Duplicate Data)

User mencoba mengganti Email/Username menjadi nama yang sudah dipakai orang lain.

```json
{
  "success": false,
  "message": "Username or Email already exists",
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

## Endpoint : `DELETE /users/{id}`

### Headers :

- `Authorization: Bearer <access_token>` (Required)

#### Description :

Endpoint ini digunakan untuk Menonaktifkan (Soft Delete) user. Data tidak benar-benar dihapus dari database demi menjaga riwayat transaksi, namun user tersebut tidak akan bisa login lagi. Hanya Owner yang boleh melakukan ini.

### Parameters :

Parameter ini adalah Path Variable.

| Key | Tipe | Default | Deskripsi                                      |
| --- | ---- | ------- | ---------------------------------------------- |
| id  | Int  | -       | ID Unik user yang ingin dihapus/dinonaktifkan. |

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

[WAJIB ADA] Token tidak valid, expired, atau tidak dikirim.

```json
{
  "success": false,
  "message": "Invalid or missing access token",
  "data": null
}
```

#### 🚫 403 Forbidden (Insufficient Permissions)

[WAJIB ADA] User bukan Owner mencoba menghapus data.

```json
{
  "success": false,
  "message": "You do not have permission to perform this action",
  "data": null
}
```

#### 404 Not Found (User Missing)

ID user tidak ditemukan di database.

```json
{
  "success": false,
  "message": "User not found",
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
