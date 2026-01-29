package routes

import (
	"laundry-backend/internal/handlers"
	middleware "laundry-backend/internal/middlewares"
	"laundry-backend/internal/repositories"

	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes registers all authentication-related routes.
func SetupAuthRoutes(router *gin.RouterGroup, authHandler *handlers.AuthHandler, authRepo repositories.AuthRepository) {

	// Grouping URL: /api/v1/auth
	auth := router.Group("/auth")
	{
		// --- PUBLIC ENDPOINTS (No Login Required) ---
		// Login: Menerima username/password -> Balas Token
		auth.POST("/login", middleware.RateLimiter(), authHandler.Login)

		// Refresh: Menerima refresh_token -> Balas Access Token Baru
		// (Saya ubah jadi /refresh-token agar sesuai dengan Swagger di Handler)
		auth.POST("/refresh-token", middleware.RateLimiter(), authHandler.RefreshToken)

		// --- PRIVATE ENDPOINTS (Token Required) ---
		// Middleware Auth dipasang di sini untuk:
		// 1. Validasi Token Signature
		// 2. Cek apakah Token Expired
		// 3. Cek apakah Token ada di Blacklist (Logout)

		// Logout: Revoke refresh token & Blacklist access token
		auth.POST("/logout", middleware.RateLimiter(), middleware.AuthMiddleware(authRepo), authHandler.Logout)

		// Me: Cek profil sendiri
		auth.GET("/me", middleware.RateLimiter(), middleware.AuthMiddleware(authRepo), authHandler.GetMe)
	}
}
