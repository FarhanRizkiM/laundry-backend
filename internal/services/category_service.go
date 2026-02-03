package services

import (
	"context"
	"errors"
	"laundry-backend/internal/dto"
	"laundry-backend/internal/models"
	"laundry-backend/internal/repositories"
	"laundry-backend/pkg/response"
	"time"
)

// CategoryService defines the contract for business logic related to service categories.
type CategoryService interface {
	CreateCategory(ctx context.Context, req dto.CreateCategoryRequest) (*dto.CategoryDetailResponse, error)
	GetCategoryList(ctx context.Context, page, perPage int, search, status, sortBy, sortOrder string) (*dto.CategoryListResponse, error)
	GetCategoryDetail(ctx context.Context, id int64) (*dto.CategoryDetailResponse, error)

	// ModifyCategoryData updates category information with validation logic.
	ModifyCategory(ctx context.Context, targetID int64, req dto.UpdateCategoryRequest) (*dto.CategoryDetailResponse, error)

	// DeleteCategory handles soft deletion of a category.
	DeactivateCategory(ctx context.Context, targetID int64) error
}

type categoryService struct {
	categoryRepo repositories.CategoryRepository
}

// NewCategoryService creates a new instance of CategoryService.
func NewCategoryService(categoryRepo repositories.CategoryRepository) CategoryService {
	return &categoryService{categoryRepo: categoryRepo}
}

// CreateCategory handles the creation of a new service classification.
func (s *categoryService) CreateCategory(ctx context.Context, req dto.CreateCategoryRequest) (*dto.CategoryDetailResponse, error) {

	// 1. Check for Duplicate Data (Unique Name)
	existingCategory, _ := s.categoryRepo.FindCategoryByName(ctx, req.CategoryName)
	if existingCategory != nil {
		return nil, errors.New(response.ErrDuplicate)
	}

	// 2. Prepare Model
	categoryModel := &models.ServiceCategory{
		CategoryName: req.CategoryName,
		Description:  req.Description,
		IsActive:     true, // Default active upon creation
		CreatedAt:    time.Now(),
		UpdatedAt:    nil, // Explicitly nil
	}

	// 3. Insert into DB
	createdCategory, err := s.categoryRepo.InsertCategory(ctx, categoryModel)
	if err != nil {
		return nil, err
	}

	// 4. Return Response
	return &dto.CategoryDetailResponse{
		ID:           createdCategory.ID,
		CategoryName: createdCategory.CategoryName,
		Description:  createdCategory.Description,
		IsActive:     createdCategory.IsActive,
		CreatedAt:    createdCategory.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    nil,
	}, nil
}

// GetCategoryList fetches a list of categories with pagination, filters, and sorting.
func (s *categoryService) GetCategoryList(ctx context.Context, page, perPage int, search, status, sortBy, sortOrder string) (*dto.CategoryListResponse, error) {

	// 1. Calculate Page limits (Safety check)
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	// 2. Call Repository
	categories, totalItems, err := s.categoryRepo.FetchCategories(ctx, page, perPage, search, status, sortBy, sortOrder)
	if err != nil {
		return nil, err
	}

	// 3. Map to DTO Summary
	var categoryResponses []dto.CategorySummaryResponse
	for _, c := range categories {
		categoryResponses = append(categoryResponses, dto.CategorySummaryResponse{
			ID:           c.ID,
			CategoryName: c.CategoryName,
			Description:  c.Description,
			IsActive:     c.IsActive,
		})
	}

	// Handle empty slice to avoid "null" in JSON
	if categoryResponses == nil {
		categoryResponses = []dto.CategorySummaryResponse{}
	}

	// 4. Calculate Total Pages
	// Using integer math to avoid float precision issues
	var totalPages int
	if perPage > 0 {
		totalPages = int((totalItems + int64(perPage) - 1) / int64(perPage))
	}

	// 5. Return Final Response
	return &dto.CategoryListResponse{
		Data: categoryResponses,
		Meta: dto.MetaData{
			CurrentPage: page,
			PerPage:     perPage,
			TotalItems:  totalItems, // FIX: Langsung pakai int64 (sesuai DTO MetaData)
			TotalPages:  totalPages,
		},
	}, nil
}

// GetCategoryDetail retrieves detailed category information by ID.
func (s *categoryService) GetCategoryDetail(ctx context.Context, id int64) (*dto.CategoryDetailResponse, error) {

	// 1. Fetch Data
	category, err := s.categoryRepo.FindCategoryByID(ctx, id)
	if err != nil {
		return nil, err // Repo returns "RESOURCE_NOT_FOUND" string error
	}

	// 2. Format Dates safely
	createdAtStr := category.CreatedAt.Format("2006-01-02 15:04:05")
	var updatedAtPtr *string // DTO expects pointer string

	if category.UpdatedAt != nil {
		formatted := category.UpdatedAt.Format("2006-01-02 15:04:05")
		updatedAtPtr = &formatted
	}

	// 3. Return Response
	return &dto.CategoryDetailResponse{
		ID:           category.ID,
		CategoryName: category.CategoryName,
		Description:  category.Description,
		IsActive:     category.IsActive,
		CreatedAt:    createdAtStr,
		UpdatedAt:    updatedAtPtr,
	}, nil
}

// ModifyCategoryData updates category profile with validation logic.
func (s *categoryService) ModifyCategory(ctx context.Context, targetID int64, req dto.UpdateCategoryRequest) (*dto.CategoryDetailResponse, error) {

	// 1. Retrieve Existing Category
	existingCategory, err := s.categoryRepo.FindCategoryByID(ctx, targetID)
	if err != nil {
		return nil, err
	}

	// 2. Update Fields (Partial Update Logic)

	// Name update with duplicate check
	if req.CategoryName != "" && req.CategoryName != existingCategory.CategoryName {
		duplicateCheck, _ := s.categoryRepo.FindCategoryByName(ctx, req.CategoryName)
		// If exists AND it's not the same record we are currently editing
		if duplicateCheck != nil && duplicateCheck.ID != targetID {
			return nil, errors.New(response.ErrDuplicate)
		}
		existingCategory.CategoryName = req.CategoryName
	}

	if req.Description != "" {
		existingCategory.Description = req.Description
	}

	// Boolean Check (Using pointer check to distinguish false vs nil)
	if req.IsActive != nil {
		existingCategory.IsActive = *req.IsActive
	}

	// 3. Update Timestamp
	now := time.Now()
	existingCategory.UpdatedAt = &now

	// 4. Save Changes to DB
	updatedCategory, err := s.categoryRepo.UpdateCategory(ctx, existingCategory)
	if err != nil {
		return nil, err
	}

	// 5. Return Updated Detail
	// Manually construct response to avoid re-fetching
	createdAtStr := updatedCategory.CreatedAt.Format("2006-01-02 15:04:05")
	var updatedAtPtr *string
	if updatedCategory.UpdatedAt != nil {
		formatted := updatedCategory.UpdatedAt.Format("2006-01-02 15:04:05")
		updatedAtPtr = &formatted
	}

	return &dto.CategoryDetailResponse{
		ID:           updatedCategory.ID,
		CategoryName: updatedCategory.CategoryName,
		Description:  updatedCategory.Description,
		IsActive:     updatedCategory.IsActive,
		CreatedAt:    createdAtStr,
		UpdatedAt:    updatedAtPtr,
	}, nil
}

// DeleteCategory handles soft deletion of a category.
func (s *categoryService) DeactivateCategory(ctx context.Context, targetID int64) error {

	// 1. Check if category exists
	_, err := s.categoryRepo.FindCategoryByID(ctx, targetID)
	if err != nil {
		return err
	}

	// 2. Execute Soft Delete
	return s.categoryRepo.DeleteCategory(ctx, targetID)
}
