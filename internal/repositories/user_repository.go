package repositories

import (
	"context"
	"database/sql"
	"errors"
	"laundry-backend/internal/models"
)

// UserRepository mendefinisikan kontrak untuk interaksi database terkait data pengguna.
type UserRepository interface {
	// Method POST
	InsertUser(ctx context.Context, user *models.User) error

	// Method GET
	FetchAllUsers(ctx context.Context, limit, offset int, search, role string, status int) ([]models.User, int64, error)
	FindUserByUsername(ctx context.Context, username string) (*models.User, error)
	FindUserByID(ctx context.Context, id int64) (*models.User, error)

	// Method PUT
	UpdatedRowUser(ctx context.Context, user *models.User) error

	// Method DELETE
	SetUserInactive(ctx context.Context, id int64) error

	// --- Helper Validasi Unik ---
	IsEmailExists(ctx context.Context, email string, excludeID int64) (bool, error)
	IsPhoneExists(ctx context.Context, phone string, excludeID int64) (bool, error)
	IsUsernameExists(ctx context.Context, username string, excludeID int64) (bool, error)
}

// userRepository adalah implementasi konkrit dari interface UserRepository  yang bertanggung jawab atas query SQL ke tabel users.
type userRepository struct {
	db *sql.DB
}

// NewUserRepository membuat instance baru dari UserRepository.
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// --- IMPLEMENTASI ---

// 1. InsertUser (Insert Data Baru)
func (r *userRepository) InsertUser(ctx context.Context, user *models.User) error {
	// A. Query Dasar
	query := "INSERT INTO users (full_name, username, email, password_hash, role, phone_number, is_active, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"

	// Eksekusi Query
	res, err := r.db.ExecContext(ctx, query, user.FullName, user.Username, user.Email, user.PasswordHash, user.Role, user.PhoneNumber, user.IsActive, user.CreatedAt)
	if err != nil {
		return err
	}

	// Ambil ID yang baru saja digenerate (Auto Increment)
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = id
	return nil
}

// 2. FetchAllUsers (List dengan Pagination & Filter)
func (r *userRepository) FetchAllUsers(ctx context.Context, limit, offset int, search, role string, status int) ([]models.User, int64, error) {

	// A. Query Dasar
	query := "SELECT id, full_name, username, role, is_active FROM users WHERE 1=1"
	countQuery := "SELECT COUNT(*) FROM users WHERE 1=1"
	var args []interface{}

	// B. Filter Dinamis
	if search != "" {
		// Cari di nama ATAU username
		query += " AND (full_name LIKE ? OR username LIKE ?)"
		countQuery += " AND (full_name LIKE ? OR username LIKE ?)"
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm)
	}

	if role != "" {
		query += " AND role = ?"
		countQuery += " AND role = ?"
		args = append(args, role)
	}

	// Filter Status (Jika -1 berarti tampilkan semua, jika 0/1 filter sesuai nilai)
	if status != -1 {
		query += " AND is_active = ?"
		countQuery += " AND is_active = ?"
		args = append(args, status)
	}

	// C. Pagination & Sorting
	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"

	// D. Hitung Total Data (Untuk Meta Pagination)
	var totalItems int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalItems)
	if err != nil {
		return nil, 0, err
	}

	// E. Eksekusi Query Data
	// Tambahkan limit & offset ke args untuk query data
	argsData := append(args, limit, offset)
	rows, err := r.db.QueryContext(ctx, query, argsData...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(
			&u.ID,
			&u.FullName,
			&u.Username,
			&u.Role,
			&u.IsActive,
		); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}
	return users, totalItems, nil
}

// 3. FindUserByID (Detail User Lengkap)
func (r *userRepository) FindUserByID(ctx context.Context, id int64) (*models.User, error) {

	// A. Query Dasar
	query := "SELECT id, full_name, username, email, password_hash, role, phone_number, is_active, last_login_at, created_at, updated_at FROM users WHERE id = ?"

	// Eksekusi Query
	var u models.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID,
		&u.FullName,
		&u.Username,
		&u.Email,
		&u.PasswordHash,
		&u.Role,
		&u.PhoneNumber,
		&u.IsActive,
		&u.LastLoginAt,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("User not found")
		}
		return nil, err
	}
	return &u, nil
}

// 4. FindUserByUsername (Dipakai Login & Auth Middleware).
func (r *userRepository) FindUserByUsername(ctx context.Context, username string) (*models.User, error) {

	// A. Query Dasar
	query := "SELECT id, full_name, username, password_hash, role, is_active FROM users WHERE username = ?"

	// Eksekusi Query
	var user models.User
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.FullName,
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

// 5. UpdatedRowUser (Update Data Profil)
func (r *userRepository) UpdatedRowUser(ctx context.Context, user *models.User) error {

	// A. Query Dasar
	query := "UPDATE users SET full_name=?, username=?, email=?, password_hash=?, role=?, phone_number=?, is_active=?, updated_at=? WHERE id=?"

	// Eksekusi Query
	_, err := r.db.ExecContext(ctx, query,
		user.FullName,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.PhoneNumber,
		user.IsActive,
		user.UpdatedAt,
		user.ID,
	)
	return err
}

// 6. SetUserInactive (Soft Delete)
func (r *userRepository) SetUserInactive(ctx context.Context, id int64) error {

	// A. Query Dasar
	query := "UPDATE users SET is_active = 0 WHERE id = ?"

	// Eksekusi Query
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// --- HELPER VALIDATION ---

func (r *userRepository) IsEmailExists(ctx context.Context, email string, excludeID int64) (bool, error) {

	var exists bool

	// A. Query Dasar
	// Cek email, TAPI abaikan user dengan ID = excludeID (dipakai saat update diri sendiri)
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = ? AND id != ?)"

	// Eksekusi Query
	err := r.db.QueryRowContext(ctx, query, email, excludeID).Scan(&exists)
	return exists, err
}

func (r *userRepository) IsPhoneExists(ctx context.Context, phone string, excludeID int64) (bool, error) {

	var exists bool

	// A. Query Dasar
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE phone_number = ? AND id != ?)"

	// Eksekusi Query
	err := r.db.QueryRowContext(ctx, query, phone, excludeID).Scan(&exists)
	return exists, err
}

func (r *userRepository) IsUsernameExists(ctx context.Context, username string, excludeID int64) (bool, error) {

	var exists bool

	// Query Dasar
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE username = ? AND id != ?)"

	// Eksekusi Query
	err := r.db.QueryRowContext(ctx, query, username, excludeID).Scan(&exists)
	return exists, err
}
