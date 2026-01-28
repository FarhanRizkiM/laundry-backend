package services

import (
	"context"
	"errors"
	"time"

	"laundry-backend/internal/dto"
	"laundry-backend/internal/models"
	"laundry-backend/internal/repositories"
	"laundry-backend/pkg/utils"
)

// AuthService mendefinisikan kontrak logika bisnis untuk urusan autentikasi.
type AuthService interface {
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	Logout(ctx context.Context, refreshToken string, jti string, expiresAt time.Time) error
	GetProfile(ctx context.Context, userID int64) (*dto.UserProfileResponse, error)
	RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error)
}

// authService adalah implementasi dari AuthService yang menggabungkan repository user dan repository auth.
type authService struct {
	authRepo repositories.AuthRepository
	userRepo repositories.UserRepository
}

// NewAuthService membuat instance baru untuk AuthService.
func NewAuthService(authRepo repositories.AuthRepository, userRepo repositories.UserRepository) AuthService {
	return &authService{
		authRepo: authRepo,
		userRepo: userRepo,
	}
}

// Login menjalankan alur verifikasi kredensial hingga pembuatan token.
func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {

	// 1. Cari user berdasarkan username
	user, err := s.userRepo.FindUserByUsername(ctx, req.Username)
	if err != nil {
		// Kita kembalikan pesan yang sama untuk keamanan
		return nil, errors.New("INVALID_CREDENTIALS")
	}

	// 2. Bandingkan password asli dengan hash di DB
	err = utils.ComparePassword(user.PasswordHash, req.Password)
	if err != nil {
		return nil, errors.New("INVALID_CREDENTIALS")
	}

	// 3. Pastikan status akun aktif
	if !user.IsActive {
		return nil, errors.New("ACCOUNT_INACTIVE")
	}

	// 4. Generate Access Token
	accessToken, _, err := utils.GenerateAccessToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}

	// 5. Generate Refresh Token
	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// 6. Simpan Refresh Token ke database
	rtModel := &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.authRepo.CreateRefreshToken(ctx, rtModel); err != nil {
		return nil, err
	}

	// 7. Kembalikan response lengkap
	return &dto.LoginResponse{
		Token: dto.TokenResponse{
			TokenType:    "Bearer",
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    900, // Harus sinkron dengan setting di jwt_utils
		},
		User: dto.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Role:     user.Role,
		},
	}, nil
}

// RefreshToken mengganti access token yang baru berdasarkan refresh token yang dimiiliki.
func (s *authService) RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error) {

	// 1. Cek keberadaan token di database
	storedToken, err := s.authRepo.GetRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		// Jika error sql: no rows in result set, berarti token tidak valid/sudah logout
		return nil, errors.New("INVALID_TOKEN")
	}

	// 2. Cek apakah token sudah expired (Masa hidup 24 jam habis)
	if storedToken.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("TOKEN_EXPIRED")
	}

	// 3. Ambil data user terbaru untuk memastikan role & username update
	user, err := s.userRepo.FindUserByID(ctx, storedToken.UserID)
	if err != nil {
		return nil, errors.New("USER_NOT_FOUND")
	}

	// 4. Generate Access Token BARU (15 Menit)
	newAccessToken, _, err := utils.GenerateAccessToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}

	// 5. Kembalikan Access Token baru saja
	return &dto.RefreshTokenResponse{
		TokenType:   "Bearer",
		AccessToken: newAccessToken,
		ExpiresIn:   900,
	}, nil
}

// Logout menangani proses penghapusan sesi user
func (s *authService) Logout(ctx context.Context, refreshToken string, jti string, expiresAt time.Time) error {
	// 1. Hapus Refresh Token dari database agar tidak bisa ditukar lagi
	_ = s.authRepo.DeleteRefreshToken(ctx, refreshToken)

	// 2. [BARU] Masukkan Access Token (JTI) ke Blacklist
	blacklistData := &models.TokenBlacklist{
		JTI:       jti,
		ExpiresAt: expiresAt,
	}

	return s.authRepo.AddToBlacklist(ctx, blacklistData)
}

// GetProfile mengambil data berdasarkan id
func (s *authService) GetProfile(ctx context.Context, userID int64) (*dto.UserProfileResponse, error) {

	// 1. Ambil data user dari repository
	user, err := s.userRepo.FindUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 2. Cek Status Aktif (Logic untuk return 403 Forbidden)
	if !user.IsActive {
		return nil, errors.New("ACCOUNT_INACTIVE")
	}

	// 3. Konversi IsActive (Boolean -> Int 1/0)
	isActiveInt := 0
	if user.IsActive {
		isActiveInt = 1
	}

	// 4. Format Tanggal (Standard Go Time Layout -> "YYYY-MM-DD HH:MM:SS")
	createdAtStr := user.CreatedAt.Format("2006-01-02 15:04:05")

	// 5. Mapping ke DTO
	// Catatan: Karena kolom 'full_name' belum ada di DB, kita isi dengan username dulu
	// agar sesuai spesifikasi JSON Frontend.
	return &dto.UserProfileResponse{
		ID:        user.ID,
		FullName:  user.FullName,
		Username:  user.Username,
		Role:      user.Role,
		IsActive:  isActiveInt,
		CreatedAt: createdAtStr,
	}, nil
}
