package routes

import (
	"laundry-backend/internal/handlers"
	middleware "laundry-backend/internal/middlewares"
	"laundry-backend/internal/repositories"

	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes adalah fungsi untuk mendaftarkan semua rute terkait Autentikasi
func SetupAuthRoutes(router *gin.RouterGroup, authHandler *handlers.AuthHandler, authRepo repositories.AuthRepository) {

	auth := router.Group("/auth")
	{
		// Endpoint publik (Tidak perlu login)
		auth.POST("/login", middleware.RateLimiter(), authHandler.Login)

		// Endpoint privat (WAJIB login / membawa token)
		auth.POST("/logout", middleware.AuthMiddleware(authRepo), authHandler.Logout)
	}
}
