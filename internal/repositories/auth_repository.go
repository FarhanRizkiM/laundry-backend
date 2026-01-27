package repositories

import (
	"context"
	"database/sql"
	"laundry-backend/internal/models"
)

// AuthRepository mendefinisikan kontrak untuk interaksi database terkait autentikasi.
type AuthRepository interface {
	CreateRefreshToken(ctx context.Context, refreshToken *models.RefreshToken) error
	DeleteRefreshToken(ctx context.Context, token string) error
	GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error)
	AddToBlacklist(ctx context.Context, blacklist *models.TokenBlacklist) error
	IsBlacklisted(ctx context.Context, jti string) (bool, error)
}

// authRepository adalah implementasi konkrit dari interface AuthRepository, yang menggunakan sql.DB sebagai database engine-nya.
type authRepository struct {
	db *sql.DB
}

// NewAuthRepository membuat instance baru dari AuthRepository.
func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{db: db}
}

// CreateRefreshToken menyimpan refresh token baru ke database.
func (r *authRepository) CreateRefreshToken(ctx context.Context, rt *models.RefreshToken) error {
	query := "INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES (?, ?, ?)"
	_, err := r.db.ExecContext(ctx, query, rt.UserID, rt.Token, rt.ExpiresAt)
	return err
}

// AddToBlacklist memasukkan JTI token yang sudah logout ke dalam daftar hitam.
func (r *authRepository) AddToBlacklist(ctx context.Context, tb *models.TokenBlacklist) error {
	query := "INSERT INTO token_blacklist (jti, expires_at) VALUES (?, ?)"
	_, err := r.db.ExecContext(ctx, query, tb.JTI, tb.ExpiresAt)

	return err
}

// IsBlacklisted mengecek apakah sebuah JTI terdaftar dalam blacklist.
func (r *authRepository) IsBlacklisted(ctx context.Context, jti string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM token_blacklist WHERE jti = ?)"
	err := r.db.QueryRowContext(ctx, query, jti).Scan(&exists)
	return exists, err
}

// DeleteRefreshToken menghapus refresh token dari database (digunakan saat logout).
func (r *authRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	query := "DELETE FROM refresh_tokens WHERE token = ?"
	_, err := r.db.ExecContext(ctx, query, token)
	return err
}

// GetRefreshToken mengambil data lengkap refresh token berdasarkan string tokennya.
func (r *authRepository) GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	query := "SELECT id, user_id, token, expires_at, created_at FROM refresh_tokens WHERE token = ?"
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&rt.ID,
		&rt.UserID,
		&rt.Token,
		&rt.ExpiresAt,
		&rt.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &rt, nil
}
