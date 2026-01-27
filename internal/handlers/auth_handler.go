package handlers

import (
	"laundry-backend/internal/dto"
	"laundry-backend/internal/services"
	"net/http"

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

// handleError adalah helper untuk memetakan error bisnis ke response HTTP yang sesuai.
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
