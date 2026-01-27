package middlewares

import (
	"laundry-backend/internal/dto"
	"laundry-backend/internal/repositories"
	"laundry-backend/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware adalah penjaga gerbang untuk memvalidasi Access Token (JWT).
func AuthMiddleware(authRepo repositories.AuthRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil header Authorization
		authHeader := c.GetHeader("Authorization")

		// 2. Cek apakah header kosong atau tidak diawali dengan "Bearer "
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, dto.BaseResponse{
				Success: false,
				Message: "Unauthorized: Missing or invalid token format",
				Data: dto.ErrorResponseData{
					ErrorCode: "UNAUTHORIZED",
					Errors:    "Authorization header is required",
				},
			})
			c.Abort() // Stop request di sini
			return
		}

		// 3. Potong string "Bearer " untuk mendapatkan token murni
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 4. Validasi token menggunakan jwt_utils yang sudah kita buat
		claims, err := utils.ValidateAccessToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.BaseResponse{
				Success: false,
				Message: "Unauthorized: Invalid or expired token",
				Data: dto.ErrorResponseData{
					ErrorCode: "INVALID_TOKEN",
					Errors:    err.Error(),
				},
			})
			c.Abort()
			return
		}

		// 5. [BARU] Cek apakah JTI token ini ada di daftar Blacklist database
		isBlacklisted, err := authRepo.IsBlacklisted(c.Request.Context(), claims.JTI)
		if err != nil {
			// Jika database error, kita assume aman atau tolak (fail-safe). Di sini kita tolak untuk keamanan.
			c.JSON(http.StatusInternalServerError, dto.BaseResponse{
				Success: false,
				Message: "Internal server error during auth check",
			})
			c.Abort()
			return
		}

		if isBlacklisted {
			c.JSON(http.StatusUnauthorized, dto.BaseResponse{
				Success: false,
				Message: "Unauthorized: Token has been logged out",
				Data: dto.ErrorResponseData{
					ErrorCode: "TOKEN_BLACKLISTED",
					Errors:    "Please login again",
				},
			})
			c.Abort()
			return
		}

		// 6. [BARU] Simpan data penting ke Context agar bisa dipakai di Handler Logout
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("jti", claims.JTI)            // Penting untuk Logout
		c.Set("exp", claims.ExpiresAt.Time) // Penting untuk Logout

		// 7. Lanjutkan ke proses berikutnya
		c.Next()
	}
}
