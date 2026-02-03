package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"laundry-backend/internal/handlers"
	"laundry-backend/internal/repositories"
	"laundry-backend/internal/routes"
	"laundry-backend/internal/services"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {

	// 1. Load configuration from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan variable sistem.")
	}

	// 2. Initialize Database Connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Gagal menyambung ke database:", err)
	}
	defer db.Close()

	// Configure Database Connection Pool
	// Optimization for high-traffic scenarios
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatal("Database tidak merespon:", err)
	}

	// ==========================================
	// 3. DEPENDENCY INJECTION (WIRING)
	// ==========================================

	// We inject dependencies from the bottom up: DB -> Repo -> Service -> Handler

	// A. Repository Layer (Data Access)
	userRepo := repositories.NewUserRepository(db)
	authRepo := repositories.NewAuthRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)

	// B. Service Layer (Business Logic)
	authService := services.NewAuthService(authRepo, userRepo)
	userService := services.NewUserService(userRepo)
	categoryService := services.NewCategoryService(categoryRepo)

	// C. Handler Layer (HTTP Transport)
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// 4. Initialize Gin Framework
	r := gin.Default()

	// Setup Routes
	v1 := r.Group("/api/v1")

	// Register Module Routes
	routes.SetupAuthRoutes(v1, authHandler, authRepo)
	routes.SetupUserRoutes(v1, userHandler, authRepo)
	routes.SetupCategoryRoutes(v1, categoryHandler, authRepo)

	// 5. Start the Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("-------------------------------------------")
	fmt.Printf("🚀 VIP Laundry Backend menyala di port %s\n", port)
	fmt.Println("-------------------------------------------")

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Gagal menyalakan server:", err)
	}
}
