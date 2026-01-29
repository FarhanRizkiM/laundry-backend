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

// AuthService defines the business logic contract for authentication.
type AuthService interface {
	AuthenticateUser(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	RenewUserSession(ctx context.Context, req dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error)
	RevokeUserSession(ctx context.Context, refreshToken string, jti string, expiresAt time.Time, userID int64) error
	GetAccountProfile(ctx context.Context, userID int64) (*dto.AuthMeResponse, error)
}

// authService is the concrete implementation combining Auth and User repositories.
type authService struct {
	authRepo repositories.AuthRepository
	userRepo repositories.UserRepository
}

// NewAuthService creates a new instance of AuthService.
func NewAuthService(authRepo repositories.AuthRepository, userRepo repositories.UserRepository) AuthService {
	return &authService{
		authRepo: authRepo,
		userRepo: userRepo,
	}
}

// AuthenticateUser handles credential verification, token generation, and login timestamp update.
func (s *authService) AuthenticateUser(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {

	// 1. Find user by username
	user, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		// Return generic error for security (avoid username enumeration)
		return nil, errors.New("INVALID_CREDENTIALS")
	}

	// 2. Verify password using the new Utils
	err = utils.VerifyPassword(user.PasswordHash, req.Password)
	if err != nil {
		return nil, errors.New("INVALID_CREDENTIALS")
	}

	// 3. Guard: Check if account is active
	if !user.IsActive {
		return nil, errors.New("ACCOUNT_INACTIVE")
	}

	// 4. Generate Access Token (JWT)
	accessToken, _, err := utils.GenerateAccessToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}

	// 5. Generate Refresh Token (UUID)
	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// 6. Persist Refresh Token to Database
	rtModel := &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(24 * time.Hour), // Refresh token validity (e.g., 1 day)
	}

	if err := s.authRepo.CreateRefreshToken(ctx, rtModel); err != nil {
		return nil, err
	}

	// 7. [CRITICAL FIX] Update Last Login Timestamp
	// We execute this immediately. If this fails, strictly speaking, we might still want
	// to allow login, but for data consistency, let's catch the error.
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		// Option: Log error here if you have a logger, but don't stop the login process.
		// For now, we proceed even if timestamp update fails.
	}

	// 8. Return Success Response
	return &dto.LoginResponse{
		Token: dto.TokenResponse{
			TokenType:    "Bearer",
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    int(utils.GetExpiryDuration().Seconds()), // Dynamic from Env
		},
		User: dto.AuthUserSummary{
			ID:       user.ID,
			FullName: user.FullName,
			Username: user.Username,
			Role:     user.Role,
		},
	}, nil
}

// RenewUserSession generates a new access token using a valid refresh token.
func (s *authService) RenewUserSession(ctx context.Context, req dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error) {

	// 1. Check if Refresh Token exists in DB
	storedToken, err := s.authRepo.GetRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		// "RESOURCE_NOT_FOUND" from repo maps to invalid token here
		return nil, errors.New("INVALID_TOKEN")
	}

	// 2. Check Expiration
	if storedToken.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("TOKEN_EXPIRED")
	}

	// 3. Retrieve User Data (Ensure user still exists/active)
	user, err := s.userRepo.FindByID(ctx, storedToken.UserID)
	if err != nil {
		return nil, errors.New("USER_NOT_FOUND")
	}

	// 4. Generate NEW Access Token
	newAccessToken, _, err := utils.GenerateAccessToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}

	// 5. Return only the new Access Token
	return &dto.RefreshTokenResponse{
		TokenType:   "Bearer",
		AccessToken: newAccessToken,
		ExpiresIn:   int(utils.GetExpiryDuration().Seconds()),
	}, nil
}

// RevokeUserSession handles logout by deleting the refresh token and blacklisting the JTI.
func (s *authService) RevokeUserSession(ctx context.Context, refreshToken string, jti string, expiresAt time.Time, userID int64) error {

	storedToken, err := s.authRepo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil
	}

	if storedToken.UserID != userID {
		return errors.New("UNAUTHORIZED_ACCESS")
	}

	// 1. Delete Refresh Token from DB (Revoke future access)
	// We ignore error here because if it's already gone, that's fine.
	_ = s.authRepo.DeleteRefreshToken(ctx, refreshToken)

	// 2. Add Access Token (JTI) to Blacklist (Revoke current access)
	blacklistData := &models.TokenBlacklist{
		JTI:       jti,
		ExpiresAt: expiresAt,
	}

	return s.authRepo.AddToBlacklist(ctx, blacklistData)
}

// GetAccountProfile retrieves the currently authenticated user's profile.
func (s *authService) GetAccountProfile(ctx context.Context, userID int64) (*dto.AuthMeResponse, error) {

	// 1. Fetch user from Repo
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 2. Guard: Active Status
	if !user.IsActive {
		return nil, errors.New("ACCOUNT_INACTIVE")
	}

	// 3. Format Dates
	createdAtStr := user.CreatedAt.Format("2006-01-02 15:04:05")

	// 4. Map to Response DTO
	return &dto.AuthMeResponse{
		ID:          user.ID,
		FullName:    user.FullName,
		Username:    user.Username,
		Email:       user.Email, // Added: Now available in DTO
		Role:        user.Role,
		PhoneNumber: user.PhoneNumber, // Added: Now available in DTO
		IsActive:    user.IsActive,    // Direct Bool (No int conversion needed)
		CreatedAt:   createdAtStr,
	}, nil
}
