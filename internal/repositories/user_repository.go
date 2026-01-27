package repositories

import (
	"context"
	"database/sql"
	"laundry-backend/internal/models"
)

// UserRepository mendefinisikan kontrak untuk interaksi database terkait data pengguna.
type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByID(ctx context.Context, id int64) (*models.User, error)
}

// userRepository adalah implementasi konkrit dari interface UserRepository  yang bertanggung jawab atas query SQL ke tabel users.
type userRepository struct {
	db *sql.DB
}

// NewUserRepository membuat instance baru dari UserRepository.
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// GetByUsername mengambil data user berdasarkan username. Biasanya digunakan untuk proses autentikasi (Login).
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User

	query := "SELECT id, username, password_hash, role, is_active FROM users WHERE username = ?"

	err := r.db.QueryRowContext(ctx, query, username).Scan(

		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Role,
		&user.IsActive,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByID mengambil data user lengkap berdasarkan ID (Untuk Profile).
func (r *userRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	var user models.User

	// Kita ambil kolom created_at juga.
	// Catatan: Jika kolom full_name belum ada di DB, kita tidak select dulu agar tidak error.
	query := "SELECT id, full_name, username, role, is_active, created_at FROM users WHERE id = ?"

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.FullName,
		&user.Username,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt, // Pastikan model User Anda punya field CreatedAt (time.Time)
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
