package models

import "time"

// User merepresentasikan entitas pengguna dalam sistem VIP Laundry.
type User struct {
	ID           int64      `json:"id"`
	FullName     string     `json:"full_name"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"` // Password tidak akan pernah muncul di JSON
	Role         string     `json:"role"`
	PhoneNumber  string     `json:"phone_number"`
	IsActive     bool       `json:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at"` // Pointer untuk kolom nullable
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}
