# LAUNDRY MANAGEMENT SYSTEM — API SPECIFICATION

## USERS MODULE SPECIFICATION

---

## Endpoint : `POST /users`

### Description :

Endpoint ini digunakan oleh **Owner** untuk mendaftarkan karyawan baru ke dalam sistem. Backend akan melakukan validasi ganda (keaslian email & keunikan username), melakukan hashing password (Bcrypt), dan menyimpannya ke tabel `users`.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`
- `Content-Type`: `application/json`

### Parameters :

Bagian ini mendefinisikan aturan main data yang harus dikirimkan oleh klien.

| Key          | Tipe   | In   | Deskripsi                                        |
| ------------ | ------ | ---- | ------------------------------------------------ |
| full_name    | String | Body | Nama lengkap user (Max. 150 karakter).           |
| username     | String | Body | Username unik, tanpa spasi (Max. 100 karakter).  |
| email        | String | Body | Email unik dan format valid (Max. 150 karakter). |
| password     | String | Body | Kata sandi minimal 8 karakter.                   |
| phone_number | String | Body | Nomor telepon aktif (Max. 30 karakter).          |
| role         | Enum   | Body | Pilihan: `owner`, `cashier`, `staff`, `courier`. |

```
{
  "full_name": "Siti Aminah",
  "username": "sitiaminah",
  "email": "sitiaminah@gmail.com",
  "password": "rahasia123",
  "phone_number": "082345678901",
  "role": "cashier"
}
```

### Request Body :

Objek JSON yang dikirimkan untuk proses registrasi karyawan baru.

```json
{
  "full_name": "Siti Aminah",
  "username": "sitiaminah",
  "email": "sitiaminah@gmail.com",
  "password": "rahasia123",
  "phone_number": "082345678901",
  "role": "cashier"
}
```

### Responses Body :

#### ✅ 201 Created

Karyawan baru berhasil didaftarkan secara sukses.

```json
{
  "success": true,
  "message": "User created successfully",
  "data": {
    "id": 2,
    "full_name": "Siti Aminah",
    "username": "sitiaminah",
    "email": "sitiaminah@gmail.com",
    "role": "cashier",
    "phone_number": "082345678901",
    "is_active": 1,
    "created_at": "2025-12-28 07:24:03",
    "updated_at": null
  }
}
```

#### ⚠️ 400 Bad Request

Input tidak lolos validasi (format salah atau field kosong).

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "email": "Invalid email format",
      "password": "Password must be at least 8 characters"
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

User memiliki token, tapi rolenya bukan `owner`.

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

Username atau Email sudah terdaftar sebelumnya di database.

```json
{
  "success": false,
  "message": "Data already exists",
  "data": {
    "error_code": "DUPLICATE_DATA",
    "errors": {
      "username": "Username 'sitiaminah' is already taken"
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

Kegagalan teknis pada server atau database.

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

## Endpoint : `GET /users`

### Description :

Endpoint ini digunakan oleh **Owner** untuk mendapatkan daftar seluruh akun karyawan secara terorganisir. Menggunakan teknik **Pagination** untuk efisiensi bandwidth dan mendukung pencarian dinamis (Filtering). Data ditarik langsung dari tabel `users`.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`

### Parameters :

Daftar filter pencarian dan pengaturan data melalui Query String.

| Key       | Tipe   | In    | Default | Deskripsi                                                           |
| --------- | ------ | ----- | ------- | ------------------------------------------------------------------- |
| page      | Int    | Query | 1       | Menentukan halaman data yang ingin diambil.                         |
| limit     | Int    | Query | 10      | Jumlah data yang ditampilkan per halaman.                           |
| name      | String | Query | -       | (Optional) Filter pencarian berdasarkan nama lengkap atau username. |
| role      | Enum   | Query | -       | (Optional) Filter: owner, cashier, staff, courier.                  |
| is_active | Int    | Query | -       | (Optional) Filter status. Isi 1 untuk aktif, 0 untuk non-aktif.     |

```
GET /api/users?page=1&limit=10&role=cashier&is_active=1
```

### Request Body :

```
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK

Daftar pengguna berhasil diambil secara sukses beserta informasi metadata halaman untuk kebutuhan navigasi di Frontend.

```json
{
  "success": true,
  "message": "Data retrieved successfully",
  "data": [
    {
      "id": 1,
      "full_name": "Farhan Rizki Maulana",
      "username": "farhanrizkimln",
      "email": "farhanrizki@gmail.com",
      "role": "owner",
      "phone_number": "081234567890",
      "is_active": 1,
      "last_login_at": "2026-01-12 14:00:00",
      "created_at": "2025-12-01 08:00:00",
      "updated_at": null
    },
    {
      "id": 2,
      "full_name": "Siti Aminah",
      "username": "sitiaminah",
      "email": "sitiaminah@gmail.com",
      "role": "cashier",
      "phone_number": "082345678901",
      "is_active": 1,
      "last_login_at": "2026-01-11 10:00:00",
      "created_at": "2025-12-28 07:24:03",
      "updated_at": null
    }
  ],
  "meta": {
    "current_page": 1,
    "per_page": 10,
    "total_items": 48,
    "total_pages": 5
  }
}
```

#### ⚠️ 400 Bad Request

Parameter query tidak valid (format salah).

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "page": "Must be a number",
      "limit": "Max limit is 100"
    }
  }
}
```

#### ⚠️ 401 Unauthorized

Token tidak valid atau tidak disertakan.

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

Akses ditolak karena pengguna bukan Owner.

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

Kegagalan sistem saat pengambilan data.

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

## Endpoint : `GET /users/{id}`

### Description :

Endpoint ini digunakan oleh **Owner** untuk mengambil data profil lengkap satu orang karyawan secara spesifik berdasarkan ID unik mereka. Backend akan memvalidasi apakah ID tersebut merupakan angka yang valid, memeriksa izin akses (Permissions), lalu melakukan query langsung ke baris data terkait di tabel `users`.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`

### Parameters :

Menargetkan ID unik melalui jalur URL (Path Parameter).

| Key | Tipe | In   | Default | Deskripsi                                                    |
| --- | ---- | ---- | ------- | ------------------------------------------------------------ |
| id  | Int  | Path | -       | ID Unik (Primary Key) pengguna yang ingin dilihat detailnya. |

```
GET /api/users/1
```

### Request Body :

```
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK (Success)

Data pengguna ditemukan. Perhatikan bahwa `data` berbentuk Object `{}`, bukan Array `[]`.

```json
{
  "success": true,
  "message": "Data retrieved successfully",
  "data": {
    "id": 1,
    "full_name": "Farhan Rizki Maulana",
    "username": "farhanrizkimln",
    "email": "farhanrizki@gmail.com",
    "role": "owner",
    "phone_number": "081234567890",
    "is_active": 1,
    "last_login_at": "2025-12-28 05:12:36",
    "created_at": "2025-12-28 03:12:36",
    "updated_at": null
  }
}
```

#### ⚠️ 400 Bad Request

Format ID salah (misal: mengirimkan string alih-alih angka).

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

Token Anda valid, tetapi Anda tidak memiliki hak akses (Role bukan owner) untuk melihat detail data karyawan lain.

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

ID berupa angka yang valid, namun datanya tidak ada di database.

```json
{
  "success": false,
  "message": "User not found",
  "data": {
    "error_code": "RESOURCE_NOT_FOUND",
    "errors": null
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

Kesalahan pada sistem atau query database.

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

## Endpoint : `PUT /users/{id}`

### Description :

Endpoint ini digunakan untuk memperbarui profil pengguna. Terdapat logika **Authorization** ketat:

- **Owner**: Akses penuh ke seluruh baris data di tabel `users`.
- **Non-Owner**: Hanya diizinkan mengubah data profil mereka sendiri. Field `role` dan `is_active` akan diabaikan oleh sistem demi keamanan jika dikirim oleh non-owner.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner, cashier, staff, courier`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`
- `Content-Type`: `application/json`

### Parameters :

Bagian ini menggunakan Path Parameter untuk menentukan target pengguna yang akan diperbarui.

| Key | Tipe | In   | Default | Deskripsi                       |
| --- | ---- | ---- | ------- | ------------------------------- |
| id  | Int  | Path | -       | ID Unik user yang ingin diedit. |

```
PUT /api/users/1
```

### Logic Guard :

| Skenario                   | Izin Akses              | Field yang Boleh Diubah            |
| -------------------------- | ----------------------- | ---------------------------------- |
| Owner akses ID siapa saja  | DIIZINKAN               | Semua field tanpa pengecualian     |
| User akses ID diri sendiri | DIIZINKAN               | Semua kecuali `role` & `is_active` |
| User akses ID orang lain   | DITOLAK (403 Forbidden) | -                                  |

### Request Body :

Kirimkan data yang ingin diubah. Field yang tidak dikirimkan akan tetap menggunakan nilai lama di database.

```json
{
  "full_name": "Farhan Rizki Maulana",
  "username": "farhanrizkimln",
  "email": "farhanrizki@gmail.com",
  "password": "barurahasia123",
  "phone_number": "081234567890",
  "role": "owner", // Opsional, hanya berlaku jika pengirim adalah Owner
  "is_active": 1 // Opsional, hanya berlaku jika pengirim adalah Owner
}
```

### Responses Body :

#### ✅ 200 OK

Data pengguna berhasil diperbarui secara sukses.

```json
{
  "success": true,
  "message": "User updated successfully",
  "data": {
    "id": 5,
    "full_name": "Farhan Rizki Maulana",
    "username": "farhanrizkimln",
    "email": "farhanrizki@gmail.com",
    "role": "owner",
    "phone_number": "081234567890",
    "is_active": 1,
    "last_login_at": "2026-01-12 08:12:36",
    "created_at": "2025-12-28 03:12:36",
    "updated_at": "2026-01-13 15:30:00"
  }
}
```

#### ⚠️ 400 Bad Request

Gagal validasi input atau format ID salah.

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "password": "Password must be at least 8 characters"
    }
  }
}
```

#### ⚠️ 401 Unauthorized

Token tidak valid atau tidak disertakan.

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

Role yang lain mencoba mengakses ID milik orang lain.

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

ID user yang dituju tidak ada di database.

```json
{
  "success": false,
  "message": "User not found",
  "data": {
    "error_code": "RESOURCE_NOT_FOUND",
    "errors": null
  }
}
```

#### 🚫 409 Conflict

Username atau Email baru yang dimasukkan sudah digunakan oleh orang lain.

```json
{
  "success": false,
  "message": "Data already exists",
  "data": {
    "error_code": "DUPLICATE_DATA",
    "errors": {
      "username": "Username already taken"
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

Kegagalan teknis pada server atau database.

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

## Endpoint : `DELETE /users/{id}`

### Description :

Endpoint ini digunakan untuk menghapus pengguna secara logika (**Soft Delete**). Sistem tidak akan menghapus baris data dari database, melainkan mengubah nilai `is_active` menjadi `0`. Hal ini bertujuan untuk menjaga integritas data pada riwayat transaksi laundry yang pernah ditangani oleh pengguna tersebut.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`

### Parameters :

Identifikasi pengguna yang akan dihapus dilakukan melalui Path Parameter.

| Key | Tipe    | In   | Default | Deskripsi                              |
| --- | ------- | ---- | ------- | -------------------------------------- |
| id  | Integer | Path | -       | ID Unik user yang ingin dinonaktifkan. |

```
DELETE /api/users/1
```

### Request Body :

```
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK

User berhasil dinonaktifkan. Kita tetap mengirimkan objek ID sebagai konfirmasi.

```json
{
  "success": true,
  "message": "User deleted successfully",
  "data": {
    "id": 1
  }
}
```

#### ⚠️ 400 Bad Request

Format ID pada URL tidak valid.

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

Terjadi jika pengirim bukan Owner, atau Owner mencoba menghapus ID-nya sendiri.

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

ID user tidak ditemukan di database.

```json
{
  "success": false,
  "message": "User not found",
  "data": {
    "error_code": "RESOURCE_NOT_FOUND",
    "errors": null
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

Kesalahan teknis saat proses update di database.

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
