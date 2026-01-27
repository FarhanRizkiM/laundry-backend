package handlers

import (
	"laundry-backend/internal/dto"
	"laundry-backend/internal/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthHandler menangani permintaan HTTP yang berkaitan dengan autentikasi.
type AuthHandler struct {
	authService services.AuthService
}

// NewAuthHandler membuat instance baru untuk AuthHandler.
func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login menangani permintaan masuk untuk autentikasi pengguna.
// @Summary User Login
// @Produce json
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	// Validasi input JSON berdasarkan tag di DTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Success: false,
			Message: "Input validation failed",
			Data: dto.ErrorResponseData{
				ErrorCode: "VALIDATION_ERROR",
				Errors:    err.Error(),
			},
		})
		return
	}

	// Memanggil Service untuk logika login
	res, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	// Mengirimkan response sukses
	c.JSON(http.StatusOK, dto.BaseResponse{
		Success: true,
		Message: "Login successfully",
		Data:    res,
	})
}

// Login menangani permintaan masuk untuk autentikasi pengguna.
// @Summary Refresh Token
// @Produce json
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest

	// Validasi Input
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Success: false,
			Message: "Refresh token is required",
			Data: dto.ErrorResponseData{
				ErrorCode: "VALIDATION_ERROR",
				Errors:    err.Error(),
			},
		})
		return
	}

	// Panggil Service
	res, err := h.authService.RefreshToken(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	// sukses
	c.JSON(http.StatusOK, dto.BaseResponse{
		Success: true,
		Message: "Access token refreshed successfully",
		Data:    res,
	})
}

// Logout menangani permintaan untuk keluar dari sistem.
// @Summary User Logout
// @Produce json
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// 1. Ambil Refresh Token dari request body atau header (sesuai kesepakatan DTO)
	// Untuk saat ini kita asumsikan dikirim via JSON body sederhana
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Success: false,
			Message: "Refresh token is required for logout",
			Data: dto.ErrorResponseData{
				ErrorCode: "VALIDATION_ERROR",
				Errors:    err.Error(),
			},
		})
		return
	}

	// 2. [BARU] Ambil Data dari Context (Hasil titipan Middleware)
	// Menggunakan MustGet karena kita yakin Middleware sudah menaruhnya di sana
	jti := c.MustGet("jti").(string)
	exp := c.MustGet("exp").(time.Time)

	// 3. Panggil Service dengan data lengkap
	err := h.authService.Logout(c.Request.Context(), req.RefreshToken, jti, exp)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.BaseResponse{
		Success: true,
		Message: "Logout successfully",
		Data: gin.H{
			"status": "access_terminated",
		},
	})
}

// GetMe menangani permintaan untuk melihat profil user sendiri (BARU).
// @Summary Get Me Profile
// @Produce json
// @Router /api/v1/auth/me [get]
func (h *AuthHandler) GetMe(c *gin.Context) {

	// 1. Ambil User ID dari Context (Titipan Middleware)
	userID := c.MustGet("user_id").(int64)

	// 2. Panggil Service
	res, err := h.authService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		// Handle khusus error ACCOUNT_INACTIVE agar return 403 Forbidden sesuai Spec
		if err.Error() == "ACCOUNT_INACTIVE" {
			c.JSON(http.StatusForbidden, dto.BaseResponse{
				Success: false,
				Message: "Access denied: your account is currently inactive",
				Data: dto.ErrorResponseData{
					ErrorCode: "ACCOUNT_INACTIVE",
					Errors:    nil,
				},
			})
			return
		}

		h.handleError(c, err)
		return
	}

	// 3. Response Sukses
	c.JSON(http.StatusOK, dto.BaseResponse{
		Success: true,
		Message: "User profile retrieved successfully",
		Data:    res,
	})
}

// handleError adalah helper untuk memetakan error bisnis ke response HTTP yang sesuai.
// Helper: Handle Error & HTTP Codes.
func (h *AuthHandler) handleError(c *gin.Context, err error) {
	switch err.Error() {

	case "INVALID_CREDENTIALS":
		c.JSON(http.StatusUnauthorized, dto.BaseResponse{
			Success: false,
			Message: "Invalid username or password",
			Data: dto.ErrorResponseData{
				ErrorCode: "INVALID_CREDENTIALS",
				Errors:    nil,
			},
		})

	case "ACCOUNT_INACTIVE":
		c.JSON(http.StatusForbidden, dto.BaseResponse{
			Success: false,
			Message: "Your account is inactive. Please contact the administrator.",
			Data: dto.ErrorResponseData{
				ErrorCode: "ACCOUNT_INACTIVE",
				Errors:    nil,
			},
		})

	case "TOKEN_EXPIRED":
		c.JSON(http.StatusUnauthorized, dto.BaseResponse{
			Success: false, Message: "Refresh token has expired",
			Data: dto.ErrorResponseData{ErrorCode: "TOKEN_EXPIRED", Errors: "Please login again"},
		})

	case "INVALID_TOKEN":
		c.JSON(http.StatusUnauthorized, dto.BaseResponse{
			Success: false, Message: "Invalid refresh token",
			Data: dto.ErrorResponseData{ErrorCode: "INVALID_TOKEN", Errors: "Token not found or revoked"},
		})

	case "sql: no rows in result set": // Handle jika user tidak ditemukan di DB
		c.JSON(http.StatusNotFound, dto.BaseResponse{
			Success: false,
			Message: "User not found",
			Data: dto.ErrorResponseData{
				ErrorCode: "USER_NOT_FOUND",
				Errors:    nil,
			},
		})

	default:
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{
			Success: false,
			Message: "An unexpected server error occurred during authentication",
			Data: dto.ErrorResponseData{
				ErrorCode: "INTERNAL_SERVER_ERROR",
				Errors:    nil,
			},
		})
	}
}
