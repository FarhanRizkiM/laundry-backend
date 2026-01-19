# LAUNDRY MANAGEMENT SYSTEM — API SPECIFICATION

## AUTH MODULE SPECIFICATION

---

## Endpoint : `POST /auth/login`

### Description :

Endpoint ini digunakan untuk memverifikasi identitas pengguna (Otentikasi). Backend akan melakukan query ke tabel `users` untuk mencocokkan `username` dan melakukan verifikasi hash password. Jika valid, server akan menghasilkan Access Token (JWT).

### Role Based Access Control (RBAC) :

- `Permissions`: `owner, cashier, staff, courier`

### Headers :

- `Accept`: `application/json`
- `Content-Type`: `application/json`

### Parameters :

Bagian ini mendefinisikan field yang wajib dikirimkan dalam Request Body untuk proses otentikasi. Password dikirim dalam bentuk teks biasa (plain text) melalui koneksi aman (HTTPS).

| Key      | Tipe   | Location | Deskripsi                              |
| -------- | ------ | -------- | -------------------------------------- |
| username | String | Body     | Username unik pengguna.                |
| password | String | Body     | Kata sandi pengguna (Min. 8 karakter). |

```
{
  "username": "farhanrizkimln",
  "password": "rahasia123"
}
```

### Request Body :

Objek JSON yang berisi kredensial pengguna. Password dikirim dalam bentuk teks biasa (plain text) melalui koneksi aman (HTTPS).

```json
{
  "username": "farhanrizkimln",
  "password": "rahasia123"
}
```

### Responses Body :

#### ✅ 200 OK

Otentikasi berhasil. Mengembalikan sepasang token dan profil lengkap pengguna.

```json
{
  "success": true,
  "message": "Login successfully",
  "data": {
    "token": {
      "token_type": "Bearer",
      "access_token": "eyJhbGciOiJIUzI1NiIsInR...",
      "expires_in": 86400 // 24 jam dalam detik (24 * 3600)
    },
    "user": {
      "id": 1,
      "full_name": "Farhan Rizki Maulana",
      "username": "farhanrizkimln",
      "email": "farhanrizkimln@gmail.com",
      "phone_number": "081234567890",
      "role": "owner"
    }
  }
}
```

#### ⚠️ 400 Bad Request

Terjadi jika format input tidak valid atau ada field yang kosong.

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "error_code": "VALIDATION_ERROR",
    "errors": {
      "username": "Username is required",
      "password": "Password must be at least 8 characters"
    }
  }
}
```

#### ⚠️ 401 Unauthorized

Username tidak ditemukan atau password salah.

```json
{
  "success": false,
  "message": "Authentication failed: invalid credentials or expired token",
  "data": {
    "error_code": "INVALID_CREDENTIALS",
    "errors": null
  }
}
```

#### 🚫 403 Forbidden

Kredensial benar, tetapi akun diblokir atau dinonaktifkan (is_active = 0).

```json
{
  "success": false,
  "message": "Access denied: account inactive or insufficient permission",
  "data": {
    "error_code": "ACCOUNT_INACTIVE",
    "errors": null
  }
}
```

#### 🚫 429 Too Many Requests

Terlalu banyak percobaan login dalam waktu singkat, memicu mekanisme rate limiting.

```json
{
  "success": false,
  "message": "Too many login attempts, please try again later.",
  "data": {
    "error_code": "RATE_LIMIT_EXCEEDED",
    "errors": null
  }
}
```

#### 🔥 500 Internal Server Error

Kegagalan teknis pada server, seperti database timeout atau gagal melakukan signing pada JWT.

```json
{
  "success": false,
  "message": "An unexpected server error occurred during authentication",
  "data": {
    "error_code": "TOKEN_GENERATION_FAILED",
    "errors": null
  }
}
```

---

## Endpoint : `POST /auth/logout`

### Description :

Endpoint ini digunakan untuk mengakhiri sesi pengguna secara formal. Karena sistem menggunakan `Stateless JWT`, backend memvalidasi keabsahan token. Proses penghapusan sesi yang sebenarnya dilakukan oleh Frontend dengan menghapus Access Token, namun endpoint ini memastikan server mengakui akhir sesi tersebut.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner, cashier, staff, courier`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`
- `Content-Type`: `application/json`

### Parameters :

Bagian ini tidak memerlukan data tambahan karena identifikasi sesi sepenuhnya diambil dari Authorization Header.

| Key  | Tipe | In  | Deskripsi                                    |
| ---- | ---- | --- | -------------------------------------------- |
| None | -    | -   | Tidak ada parameter yang diperlukan di body. |

```
None (Identifikasi sesi dilakukan melalui Authorization Header).
```

### Request Body :

```json
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK

Proses logout berhasil dilakukan secara sukses. Token yang dikirimkan telah divalidasi dan klien diinstruksikan untuk membersihkan sesi.

```json
{
  "success": true,
  "message": "Logout successfully",
  "data": {
    "status": "session_terminated"
  }
}
```

#### ⚠️ 401 Unauthorized

Terjadi jika Access Token tidak valid, sudah kadaluarsa, atau tidak disertakan.

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

#### 🔥 500 Internal Server Error

Terjadi kegagalan sistem internal saat memproses permintaan logout.

```json
{
  "success": false,
  "message": "An unexpected server error occurred during logout",
  "data": {
    "error_code": "LOGOUT_PROCESS_FAILED",
    "errors": null
  }
}
```

---

## Endpoint : `GET /auth/me`

### Description :

Endpoint ini digunakan untuk mengambil data profil lengkap pengguna yang sedang aktif (pemilik token). Backend akan mendekripsi token dari header untuk mendapatkan `user_id`, lalu melakukan query ke tabel `users` untuk mengambil data terbaru. Sangat berguna untuk sinkronisasi state aplikasi saat pertama kali dimuat.

### Role Based Access Control (RBAC) :

- `Permissions`: `owner, cashier, staff, courier`

### Headers :

- `Authorization`: `Bearer <access_token>` (Required)
- `Accept`: `application/json`
- `Content-Type`: `application/json`

### Parameters :

Identifikasi dilakukan sepenuhnya melalui Authorization Header. Tidak ada parameter tambahan.

| Key  | Tipe | In  | Deskripsi                                                 |
| ---- | ---- | --- | --------------------------------------------------------- |
| None | -    | -   | Identifikasi sesi dilakukan melalui Authorization Header. |

```
None (Identifikasi sesi dilakukan melalui Authorization Header).
```

### Request Body :

```json
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK

Profil pengguna berhasil diambil secara sukses.

```json
{
  "success": true,
  "message": "User profile retrieved successfully",
  "data": {
    "id": 1,
    "full_name": "Farhan Rizki Maulana",
    "username": "farhanrizkimln",
    "email": "farhanrizki@gmail.com",
    "role": "owner",
    "phone_number": "081234567890",
    "is_active": 1,
    "last_login_at": "2026-01-12 08:12:36",
    "created_at": "2026-01-12 08:12:36",
    "updated_at": null
  }
}
```

#### ⚠️ 401 Unauthorized

Terjadi jika Access Token tidak valid, kadaluarsa, atau header tidak disertakan.

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

Token valid, namun akun pengguna terkait dalam status non-aktif (is_active = 0).

```json
{
  "success": false,
  "message": "Access denied: account inactive or insufficient permission",
  "data": {
    "error_code": "ACCOUNT_INACTIVE",
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

Terjadi kegagalan saat proses query ke database atau dekripsi token.

```json
{
  "success": false,
  "message": "An unexpected server error occurred during profile retrieval",
  "data": {
    "error_code": "PROFILE_FETCH_FAILED",
    "errors": null
  }
}
```
