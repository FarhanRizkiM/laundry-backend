package dto

// LoginRequest digunakan untuk menangkap input dari user saat login.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse adalah struktur data yang dikembalikan setelah login sukses.
type LoginResponse struct {
	Token TokenResponse `json:"token"`
	User  UserResponse  `json:"user"`
}

// TokenResponse menyimpan detail access token, refresh token, dan durasi berlakunya.
type TokenResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

// UserResponse menyimpan informasi profil singkat user untuk kebutuhan di Frontend.
type UserResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// UserProfileResponse adalah format spesifik untuk endpoint profil
type UserProfileResponse struct {
	ID        int64  `json:"id"`
	FullName  string `json:"full_name"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	IsActive  int    `json:"is_active"`
	CreatedAt string `json:"created_at"`
}

// RefreshTokenRequest digunakan saat frontend minta access token baru.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshTokenResponse hanya mengembalikan Access Token baru (Refresh Token pakai yang lama).
type RefreshTokenResponse struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
