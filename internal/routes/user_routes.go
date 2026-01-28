package routes

import (
	"laundry-backend/internal/handlers"
	middleware "laundry-backend/internal/middlewares"
	"laundry-backend/internal/repositories"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.RouterGroup, userHandler *handlers.UserHandler, authRepo repositories.AuthRepository) {

	// Grouping URL /users
	users := router.Group("/users")

	// Pasang Satpam (Auth) di gerbang utama /users
	users.Use(middleware.AuthMiddleware(authRepo))

	// POST /api/v1/users (Hanya Owner)
	users.POST("", middleware.RoleMiddleware("owner"), userHandler.CreateUser)

	// GET /api/v1/users (Hanya Owner)
	users.GET("", middleware.RoleMiddleware("owner"), userHandler.GetListUsers)

	// GET /api/v1/users/:id (Hanya Owner)
	users.GET("/:id", middleware.RoleMiddleware("owner"), userHandler.GetDetailUser)

	// PUT /api/v1/users/:id (Owner, Cashier, Staff boleh update profil mereka sendiri - Logic validasi ID ada di Service)
	// Kita izinkan role-role ini masuk dulu
	users.PUT("/:id", middleware.RoleMiddleware("owner", "cashier", "staff", "courier"), userHandler.UpdateUser)

	// DELETE /api/v1/users/:id (Hanya Owner)
	users.DELETE("/:id", middleware.RoleMiddleware("owner"), userHandler.DeleteUser)
}
