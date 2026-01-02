# LAUNDRY MANAGEMENT SYSTEM — API SPECIFICATION

## AUTH MODULE SPECIFICATION

---

## Endpoint : `POST /auth/login`

### Headers :

- `Content-Type`: `application/json` (Required)

#### Description :

Endpoint ini digunakan untuk memverifikasi identitas pengguna. Klien mengirimkan kredensial (username dan password), dan jika valid, server akan mengembalikan Access Token (JWT) untuk otorisasi serta informasi profil pengguna.

### Request Body :

Objek JSON yang berisi kredensial pengguna. Password dikirim dalam bentuk teks biasa (plain text) melalui koneksi aman (HTTPS).

```json
{
  "username": "farhanrizkimln", // Required, Max 100 chars (Sesuai DB)
  "password": "rahasia123" // Required, Min 8 chars
}
```

### Responses Body :

#### ✅ 200 OK (Success)

Otentikasi berhasil. Mengembalikan token akses, token refresh, dan data profil pengguna yang sedang login.

```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": {
      "token_type": "Bearer",
      "access_token": "eyJhbGciOiJIUzI1NiIsInR...",
      "refresh_token": "d9b2d63d-a233-4123-847a...",
      "expires_in": 3600
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

#### ⚠️ 400 Bad Request (Validation Error)

Validasi input gagal. Terjadi kesalahan format pada data yang dikirim (misalnya field kosong atau password terlalu pendek).

```json
{
  "success": false,
  "message": "Input validation failed",
  "data": {
    "username": "Username is required",
    "password": "Password must be at least 8 characters"
  }
}
```

#### 🚫 401 Unauthorized (Invalid Credentials)

Otentikasi gagal. Username tidak ditemukan atau password yang dimasukkan salah.

```json
{
  "success": false,
  "message": "Username or password is incorrect",
  "data": null
}
```

#### ⛔ 403 Forbidden (Account Suspended)

Kredensial valid, tetapi akun pengguna telah dinonaktifkan oleh administrator (is_active = 0).

```json
{
  "success": false,
  "message": "Your account has been deactivated. Please contact the administrator.",
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

## Endpoint : `POST /auth/refresh`

### Headers :

- `Content-Type`: `application/json` (Required)

#### Description :

Endpoint ini digunakan untuk memperbarui Access Token yang sudah kadaluarsa tanpa mengharuskan pengguna login ulang. Endpoint ini menerapkan mekanisme Token Rotation: saat request berhasil, Refresh Token yang lama akan hangus dan digantikan dengan yang baru.

### Request Body :

Objek JSON yang berisi Refresh Token yang valid (didapatkan saat login sebelumnya).

```json
{
  "refresh_token": "d9b2d63d-a233-4123-847a-..." // Required
}
```

### Responses Body :

#### ✅ 200 OK (Success)

Token berhasil diperbarui. Server mengembalikan pasangan Access Token dan Refresh Token yang baru. Token lama dianggap tidak valid lagi.

```json
{
  "success": true,
  "message": "Token refreshed successfully",
  "data": {
    "token": {
      "token_type": "Bearer",
      "access_token": "NEW_ACCESS_TOKEN_ABC123XYZ",
      "refresh_token": "NEW_REFRESH_TOKEN_DEF456UVW",
      "expires_in": 3600
    }
  }
}
```

#### ⚠️ 400 Bad Request (Validation Error)

Validasi input gagal. Field refresh_token kosong atau format JSON salah.

```json
{
  "success": false,
  "message": "Input Validation Failed",
  "data": {
    "refresh_token": "Refresh token is required"
  }
}
```

#### 🚫 401 Unauthorized (Invalid Credentials)

Refresh token tidak valid, sudah kadaluarsa, atau sudah pernah digunakan sebelumnya (terdeteksi replay attack). Pengguna harus login ulang.

```json
{
  "success": false,
  "message": "Invalid or expired refresh token",
  "data": null
}
```

#### ⛔ 403 Forbidden (Account Suspended)

Token valid, namun akun pengguna terkait telah dinonaktifkan oleh administrator (is_active = 0).

```json
{
  "success": false,
  "message": "Your account has been deactivated. Please contact the administrator.",
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

## Endpoint : `POST /auth/logout`

### Headers :

- `Authorization: Bearer <access_token>` (Required)
- `Content-Type`: `application/json` (Required)

#### Description :

Endpoint ini digunakan untuk mengakhiri sesi pengguna. Sistem akan memvalidasi Access Token dan menghapus/mem-blacklist Refresh Token yang dikirimkan agar tidak bisa digunakan lagi untuk memperbarui sesi. Ini memastikan pengguna benar-benar keluar dari sistem.

### Request Body :

Objek JSON berisi Refresh Token yang ingin dihapus dari sistem (whitelist database).

```json
{
  "refresh_token": "d9b2d63d-a233-4123-847a-..." // Required
}
```

### Responses Body :

#### ✅ 200 OK (Success)

Logout berhasil. Access Token dan Refresh Token telah dinonaktifkan.

```json
{
  "success": true,
  "message": "Logout successful",
  "data": null
}
```

#### ⚠️ 400 Bad Request (Validation Error)

Validasi gagal. Refresh token tidak dikirimkan dalam body request.

```json
{
  "success": false,
  "message": "Input Validation Failed",
  "data": {
    "refresh_token": "Refresh token is required"
  }
}
```

#### 🚫 401 Unauthorized (Invalid Credentials)

Access Token pada header tidak valid, kadaluarsa, atau tidak disertakan.

```json
{
  "success": false,
  "message": "Invalid or missing access token",
  "data": null
}
```

#### 🔥 500 Internal Server Error

Terjadi kesalahan saat menghapus sesi di database.

```json
{
  "success": false,
  "message": "An unexpected error occurred",
  "data": null
}
```

---

## Endpoint : `GET /auth/me`

### Headers :

- `Authorization: Bearer <access_token>` (Required)

#### Description :

Endpoint ini digunakan untuk mendapatkan profil lengkap dari pengguna yang sedang login berdasarkan Access Token yang dikirimkan. Server akan mendekripsi token untuk mendapatkan ID pengguna, lalu mengambil data terbaru dari database. Endpoint ini sering digunakan untuk inisialisasi data user di Frontend (misalnya: menampilkan nama & foto profil di dashboard).

### Request Body :

```json
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK (Success)

Profil pengguna berhasil diambil.

```json
{
  "success": true,
  "message": "User profile retrieved",
  "data": {
    "id": 1,
    "full_name": "Farhan Rizki Maulana",
    "username": "owner01",
    "email": "owner@laundry.com",
    "role": "owner",
    "phone_number": "081234567890",
    "is_active": true,
    "last_login_at": "2025-01-01T10:00:00Z",
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
}
```

#### 🚫 401 Unauthorized (Invalid Token)

Token tidak valid, kadaluarsa, atau header Authorization tidak dikirim.

```json
{
  "success": false,
  "message": "Unauthorized access",
  "data": null
}
```

#### ⛔ 403 Forbidden (Account Suspended)

Token valid, namun saat dicek ke database, akun pengguna tersebut telah dinonaktifkan (is_active = 0).

```json
{
  "success": false,
  "message": "Your account has been deactivated",
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
