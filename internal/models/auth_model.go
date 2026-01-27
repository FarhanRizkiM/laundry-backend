package models

import "time"

// RefreshToken merepresentasikan data token untuk pembaruan sesi user.
type RefreshToken struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// TokenBlacklist menyimpan JTI dari token yang sudah tidak valid (logout/revoke).
type TokenBlacklist struct {
	ID        int64     `json:"id"`
	JTI       string    `json:"jti"` // Menggunakan JTI (all caps) sesuai idiom Go
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
