# LAUNDRY MANAGEMENT SYSTEM — API SPECIFICATION

## AUTH MODULE SPECIFICATION

---

## Endpoint : `POST /auth/login`

#### Description :

Endpoint ini digunakan untuk memverifikasi identitas user (Login). Client mengirimkan kredensial, dan jika valid, Server akan mengembalikan Access Token (untuk otorisasi) serta data profil user tersebut.

### Request Body :

Mengirimkan username dan password user yang sudah terdaftar.

```json
{
  "username": "farhanrizkimln", // Required, Max 180 chars
  "password": "rahasia123" // Required, Min 8 chars
}
```

### Responses Body :

#### ✅ 200 OK (Success)

Login berhasil. Token dikembalikan untuk digunakan pada request selanjutnya.

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
      "role": "owner"
    }
  }
}
```

#### ⚠️ 400 Bad Request (Validation Error)

Format input tidak sesuai (misal: password kurang dari 8 karakter).

```json
{
  "success": false,
  "message": "Input Validation Failed",
  "data": {
    "username": "Username is required",
    "password": "Password too short"
  }
}
```

#### 🚫 401 Unauthorized (Invalid Credentials)

Username tidak ditemukan atau password salah.

```json
{
  "success": false,
  "message": "Username or password is incorrect",
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

## Endpoint : `POST /auth/refresh`

#### Description :

Endpoint ini digunakan untuk memperbarui sesi user tanpa harus login ulang. Client menukarkan Refresh Token lama (yang masih valid) untuk mendapatkan Access Token baru dan Refresh Token pengganti (Token Rotation).

### Request Body :

Mengirimkan refresh token yang didapat saat login sebelumnya.

```json
{
  "refresh_token": "d9b2d63d-a233-4123-847a-..." // Required
}
```

### Responses Body :

#### ✅ 200 OK (Success)

Token berhasil diperbarui. Token lama hangus, dan server memberikan pasangan token baru.

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

Terjadi jika client lupa mengirimkan field refresh_token.

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

Token sudah kadaluarsa, format salah, atau token sudah pernah dipakai sebelumnya (Reuse Detection).

```json
{
  "success": false,
  "message": "Invalid or expired refresh token",
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

## Endpoint : `POST /auth/logout`

### Headers :

- `Authorization: Bearer <access_token>` (Required)

#### Description :

Endpoint ini digunakan untuk mengakhiri sesi user. Client wajib mengirimkan Access Token (di Header) untuk otorisasi, dan Refresh Token (di Body) untuk menghapus sesi spesifik tersebut dari database.

### Request Body :

Mengirimkan refresh token yang ingin dihapus/di-invalidate.

```json
{
  "refresh_token": "d9b2d63d-a233-4123-847a-..." // Required
}
```

### Responses Body :

#### ✅ 200 OK (Success)

Logout berhasil. Refresh token telah dihapus dari database, sesi berakhir.

```json
{
  "success": true,
  "message": "Logout successful",
  "data": null
}
```

#### ⚠️ 400 Bad Request (Validation Error)

Terjadi jika client lupa mengirimkan field refresh_token di body.

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

Terjadi jika Header Authorization tidak ada, format salah, atau Access Token sudah expired.

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

## Endpoint : `GET /auth/me`

### Headers :

- `Authorization: Bearer <access_token>` (Required)

#### Description :

Endpoint ini digunakan untuk melihat "Siapa saya?". Backend akan membaca Access Token dari Header, mencari pemiliknya di database, dan mengembalikan profil user tersebut (tanpa password).

### Request Body :

```json
None (Kosong).
```

### Responses Body :

#### ✅ 200 OK (Success)

Profil ditemukan. Data user dikembalikan lengkap (kecuali password hash).

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
    "last_login_at": "2025-01-01T10:00:00Z", // <--- Tambahkan ini
    "created_at": "2024-01-01T10:00:00Z", // <--- Tambahkan ini
    "updated_at": "2024-01-01T10:00:00Z" // <--- Tambahkan ini
  }
}
```

#### 🚫 401 Unauthorized (Invalid Credentials)

Terjadi jika Token tidak valid, Expired, atau User tersebut sudah dihapus/dinonaktifkan oleh admin.

```json
{
  "success": false,
  "message": "Unauthorized (User not found or deactivated)",
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
