package routes

import (
	"laundry-backend/internal/handlers"
	middleware "laundry-backend/internal/middlewares" // Pastikan nama package-nya benar (middlewares atau middleware?)
	"laundry-backend/internal/repositories"

	"github.com/gin-gonic/gin"
)

// SetupCategoryRoutes mengatur semua endpoint untuk modul Category
func SetupCategoryRoutes(router *gin.RouterGroup, categoryHandler *handlers.CategoryHandler, authRepo repositories.AuthRepository) {

	// Grouping URL: /api/v1/categories
	categories := router.Group("/categories")

	// 1. Pasang Global Auth Middleware untuk module ini
	// Artinya: Siapapun yang mau akses URL di bawah ini HARUS Login dulu
	categories.Use(middleware.AuthMiddleware(authRepo))

	// --- ENDPOINT YANG BISA DIAKSES SEMUA USER LOGIN (Owner, Kasir, Staff) ---

	// GET /api/v1/categories
	categories.GET("", categoryHandler.HandleGetCategoryList)

	// GET /api/v1/categories/:id
	categories.GET("/:id", categoryHandler.HandleGetCategoryDetail)

	// --- ENDPOINT KHUSUS OWNER ---
	// Kita buat sub-group lagi untuk proteksi Role
	protected := categories.Group("/")
	protected.Use(middleware.RoleMiddleware("owner"))
	{
		// POST /api/v1/categories
		protected.POST("", categoryHandler.HandleCreateCategory)

		// PUT /api/v1/categories/:id
		protected.PUT("/:id", categoryHandler.HandleUpdateCategory)

		// DELETE /api/v1/categories/:id
		protected.DELETE("/:id", categoryHandler.HandleDeleteCategory)
	}
}
