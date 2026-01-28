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
	// 1. Memuat konfigurasi dari file .env
	if err := godotenv.Load(); err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan variable sistem.")
	}

	// 2. Inisialisasi Koneksi Database
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

	// Konfigurasi Database Pool (Standar VIP)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatal("Database tidak merespon:", err)
	}

	// 3. Dependency Injection (Menyambungkan Kabel-Kabel)

	// A. Repositories (Layer Data)
	userRepo := repositories.NewUserRepository(db)
	authRepo := repositories.NewAuthRepository(db)

	// B. Services (Layer Bisnis)
	authService := services.NewAuthService(authRepo, userRepo)
	userService := services.NewUserService(userRepo) // --- BARU: Service User

	// C. Handlers (Layer HTTP)
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService) // --- BARU: Handler User

	// 4. Inisialisasi Framework Gin
	r := gin.Default()

	// Registrasi Rute
	v1 := r.Group("/api/v1")

	// Panggil setup route masing-masing modul
	routes.SetupAuthRoutes(v1, authHandler, authRepo)
	routes.SetupUserRoutes(v1, userHandler, authRepo) // --- BARU: Route User didaftarkan

	// 5. Menjalankan Server
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
