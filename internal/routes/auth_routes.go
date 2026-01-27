package routes

import (
	"laundry-backend/internal/handlers"
	middleware "laundry-backend/internal/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes adalah fungsi untuk mendaftarkan semua rute terkait Autentikasi
func SetupAuthRoutes(router *gin.RouterGroup, authHandler *handlers.AuthHandler) {

	auth := router.Group("/auth")
	{
		auth.POST("/login", middleware.RateLimiter(), authHandler.Login)
	}
}
