package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"laundry-backend/internal/models"
	"strings"
	"time"
)

type CategoryRepository interface {
	InsertCategory(ctx context.Context, category *models.ServiceCategory) (*models.ServiceCategory, error)
	FetchCategories(ctx context.Context, page, perPage int, search, status, sortBy, sortOrder string) ([]models.ServiceCategory, int64, error)
	FindCategoryByID(ctx context.Context, id int64) (*models.ServiceCategory, error)
	FindCategoryByName(ctx context.Context, categoryName string) (*models.ServiceCategory, error)
	UpdateCategory(ctx context.Context, category *models.ServiceCategory) (*models.ServiceCategory, error)
	DeleteCategory(ctx context.Context, id int64) error
}

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) InsertCategory(ctx context.Context, category *models.ServiceCategory) (*models.ServiceCategory, error) {

	query := `
			INSERT INTO service_categories (category_name, description, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?)
	`

	res, err := r.db.ExecContext(ctx, query,
		category.CategoryName,
		category.Description,
		category.IsActive,
		category.CreatedAt,
		category.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	category.ID = id
	return category, nil
}

func (r *categoryRepository) FetchCategories(ctx context.Context, page, perPage int, search, status, sortBy, sortOrder string) ([]models.ServiceCategory, int64, error) {

	baseQuery := " FROM service_categories WHERE 1=1"
	var args []interface{}

	if search != "" {
		baseQuery += " AND LOWER(category_name) LIKE ?"
		args = append(args, "%"+strings.ToLower(search)+"%")
	}

	if status != "" {
		if status == "1" {
			baseQuery += " AND is_active = 1"
		} else if status == "0" {
			baseQuery += " AND is_active = 0"
		}
	}

	countQuery := "SELECT COUNT(*)" + baseQuery
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	validSortColumns := map[string]bool{
		"category_name": true,
		"created_at":    true,
		"id":            true,
	}
	if !validSortColumns[sortBy] {
		sortBy = "category_name"
	}

	sortOrder = strings.ToUpper(sortOrder)
	if sortOrder != "ASC" && sortOrder != "DESC" {
		sortOrder = "ASC"
	}

	query := "SELECT id, category_name, description, is_active, created_at, updated_at" + baseQuery

	query += fmt.Sprintf(" ORDER BY %s %s LIMIT ? OFFSET ?", sortBy, sortOrder)

	offset := (page - 1) * perPage
	args = append(args, perPage, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var categories []models.ServiceCategory
	for rows.Next() {
		var c models.ServiceCategory
		err := rows.Scan(
			&c.ID,
			&c.CategoryName,
			&c.Description,
			&c.IsActive,
			&c.CreatedAt,
			&c.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		categories = append(categories, c)
	}

	return categories, total, nil
}

func (r *categoryRepository) FindCategoryByID(ctx context.Context, id int64) (*models.ServiceCategory, error) {

	query := `
		SELECT id, category_name, description, is_active, created_at, updated_at FROM service_categories WHERE id = ?
	`

	var c models.ServiceCategory
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID,
		&c.CategoryName,
		&c.Description,
		&c.IsActive,
		&c.CreatedAt,
		&c.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("RESOURCE_NOT_FOUND")
		}
		return nil, err
	}
	return &c, nil
}

func (r *categoryRepository) FindCategoryByName(ctx context.Context, categoryName string) (*models.ServiceCategory, error) {

	query := `
		SELECT id, category_name, description, is_active, created_at, updated_at FROM service_categories WHERE category_name = ?
	`

	var c models.ServiceCategory
	err := r.db.QueryRowContext(ctx, query, categoryName).Scan(
		&c.ID,
		&c.CategoryName,
		&c.Description,
		&c.IsActive,
		&c.CreatedAt,
		&c.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("RESOURCE_NOT_FOUND")
		}
		return nil, err
	}

	return &c, nil
}

func (r *categoryRepository) UpdateCategory(ctx context.Context, category *models.ServiceCategory) (*models.ServiceCategory, error) {

	query := `
		UPDATE service_categories SET category_name = ?, description = ?, is_active = ?, updated_at = ?
		WHERE id = ? 
	`

	now := time.Now()
	category.UpdatedAt = &now

	_, err := r.db.ExecContext(ctx, query,
		category.CategoryName,
		category.Description,
		category.IsActive,
		category.UpdatedAt,
		category.ID,
	)

	if err != nil {
		return nil, err
	}

	return category, nil
}

func (r *categoryRepository) DeleteCategory(ctx context.Context, id int64) error {

	query := "UPDATE service_categories SET is_active = 0 WHERE id = ?"

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
