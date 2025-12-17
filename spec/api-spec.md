# API Specification — Laundry Management System

## Auth

### POST /auth/login

#### Deskripsi

Digunakan oleh user internal (owner, kasir, staff, courier)
untuk login ke sistem dan mendapatkan access token.

---

#### HTTP Request

```http
POST /api/v1/auth/login
Content-Type: application/json


#### Request Body:

{
  "username": "andi.kasir",
  "password": "password123"
}


Response — Success (200 OK):

{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "user": {
    "id": 3,
    "name": "Andi",
    "username": "andi.kasir",
    "role": "kasir"
  }
}

Response — Error (400 Bad Request):

{
  "error": "validation_error",
  "message": "Username and password are required"
}

Response — Error (401 Unauthorized):

{
  "error": "invalid_credentials",
  "message": "Username or password is incorrect"
}

Response — Error (403 Forbidden):

{
  "error": "user_inactive",
  "message": "User account is inactive"
}
```
