package repositories

import (
	"context"
	"database/sql"
	"laundry-backend/internal/models"
)

// UserRepository mendefinisikan kontrak untuk interaksi database terkait data pengguna.
type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (*models.User, error)
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
