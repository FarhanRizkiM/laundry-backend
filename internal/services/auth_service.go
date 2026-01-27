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
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
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
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
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
