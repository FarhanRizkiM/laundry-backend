package dto

// --- 1. REQUEST DTO (Input) ---

type CreateUserRequest struct {
	FullName    string `json:"full_name" binding:"required,min=3,max=150"`
	Username    string `json:"username" binding:"required,alphanum,min=3,max=100"`
	Email       string `json:"email" binding:"required,email,max=150"`
	Password    string `json:"password" binding:"required,min=8"`
	PhoneNumber string `json:"phone_number" binding:"required,numeric,max=30"`
	Role        string `json:"role" binding:"required,oneof=owner cashier staff courier"`
}

type UpdateUserRequest struct {
	FullName    string `json:"full_name" binding:"omitempty,min=3,max=150"`
	Username    string `json:"username" binding:"omitempty,alphanum,min=3,max=100"`
	Email       string `json:"email" binding:"omitempty,email,max=150"`
	Password    string `json:"password" binding:"omitempty,min=8"`
	PhoneNumber string `json:"phone_number" binding:"omitempty,numeric,max=30"`
	Role        string `json:"role" binding:"omitempty,oneof=owner cashier staff courier"`
	IsActive    *bool  `json:"is_active" binding:"omitempty"`
}

// --- 2. RESPONSE DTO (Output) ---

// UserSummaryResponse: KHUSUS untuk Endpoint GET /api/v1/users (List) atau ringkasan response.
// Isinya ringkas, hanya info publik karyawan.
type UserSummaryResponse struct {
	ID       int64  `json:"id"`
	FullName string `json:"full_name"`
	Username string `json:"username"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}

// UserDetailResponse: KHUSUS untuk Endpoint GET /api/v1/users/{id} (Detail) atau detail response.
// Isinya lengkap termasuk Email, No HP, dan Time Logs.
type UserDetailResponse struct {
	ID          int64  `json:"id"`
	FullName    string `json:"full_name"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	PhoneNumber string `json:"phone_number"`
	IsActive    bool   `json:"is_active"`
	LastLoginAt string `json:"last_login_at"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

type UserListResponse struct {
	Data []UserSummaryResponse `json:"data"`
	Meta MetaData              `json:"meta"`
}
