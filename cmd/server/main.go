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

	// SetMaxOpenConns membatasi jumlah koneksi yang terbuka ke database.
	db.SetMaxOpenConns(25)
	// SetMaxIdleConns menjaga jumlah koneksi standby agar tidak dibuat ulang terus menerus.
	db.SetMaxIdleConns(25)
	// SetConnMaxLifetime memastikan koneksi yang sudah tua ditutup dan diganti yang baru.
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatal("Database tidak merespon:", err)
	}

	// 3. Dependency Injection (Menyambungkan Kabel-Kabel)
	userRepo := repositories.NewUserRepository(db)
	authRepo := repositories.NewAuthRepository(db)

	authService := services.NewAuthService(authRepo, userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	// 4. Inisialisasi Framework Gin
	r := gin.Default()

	// Registrasi Rute
	v1 := r.Group("/api/v1")
	routes.SetupAuthRoutes(v1, authHandler)

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
