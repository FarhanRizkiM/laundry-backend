package services

import (
	"context"
	"errors"
	"time"

	"laundry-backend/internal/dto"
	"laundry-backend/internal/models"
	"laundry-backend/internal/repositories"
	"laundry-backend/pkg/utils"
)

type UserService interface {
	RegisterUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserDetailResponse, error)
	RetrievedUserDirectory(ctx context.Context, page, perPage int, search, role string, status int) (*dto.UserListResponse, error)
	GetUserProfile(ctx context.Context, id int64) (*dto.UserDetailResponse, error)
	ModifyUserData(ctx context.Context, id int64, req dto.UpdateUserRequest) (*dto.UserDetailResponse, error)
	DeactivateUserAccount(ctx context.Context, id int64) error
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

// 1. CreateUser (Register Karyawan Baru)
func (s *userService) RegisterUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserDetailResponse, error) {
	// A. Cek Duplikasi Data (Username, Email, No HP)
	if exists, _ := s.userRepo.IsUsernameExists(ctx, req.Username, 0); exists {
		return nil, errors.New("USERNAME_EXISTS")
	}
	if exists, _ := s.userRepo.IsEmailExists(ctx, req.Email, 0); exists {
		return nil, errors.New("EMAIL_EXISTS")
	}
	if exists, _ := s.userRepo.IsPhoneExists(ctx, req.PhoneNumber, 0); exists {
		return nil, errors.New("PHONE_EXISTS")
	}

	// B. Hash Password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// C. Siapkan Model
	userModel := &models.User{
		FullName:     req.FullName,
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         req.Role,
		PhoneNumber:  req.PhoneNumber,
		IsActive:     true, // Default aktif saat dibuat
		CreatedAt:    time.Now(),
	}

	// D. Simpan ke DB
	if err := s.userRepo.InsertUser(ctx, userModel); err != nil {
		return nil, err
	}

	// E. Kembalikan Response
	return &dto.UserDetailResponse{
		ID:          userModel.ID,
		FullName:    userModel.FullName,
		Username:    userModel.Username,
		Email:       userModel.Email,
		Role:        userModel.Role,
		PhoneNumber: userModel.PhoneNumber,
		IsActive:    userModel.IsActive,
		CreatedAt:   userModel.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// 2. GetAllUsers (List dengan Pagination)
func (s *userService) RetrievedUserDirectory(ctx context.Context, page, perPage int, search, role string, status int) (*dto.UserListResponse, error) {
	// Hitung Offset
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}
	offset := (page - 1) * perPage

	// Panggil Repository
	users, totalItems, err := s.userRepo.FetchAllUsers(ctx, perPage, offset, search, role, status)
	if err != nil {
		return nil, err
	}

	// Mapping ke DTO Summary
	var userResponses []dto.UserSummaryResponse
	for _, u := range users {
		userResponses = append(userResponses, dto.UserSummaryResponse{
			ID:       u.ID,
			FullName: u.FullName,
			Username: u.Username,
			Role:     u.Role,
			IsActive: u.IsActive,
		})
	}

	// Hitung Total Pages
	totalPages := int((totalItems + int64(perPage) - 1) / int64(perPage))

	return &dto.UserListResponse{
		Data: userResponses,
		Meta: dto.MetaData{
			CurrentPage: page,
			PerPage:     perPage,
			TotalItems:  totalItems,
			TotalPages:  totalPages,
		},
	}, nil
}

// 3. GetUserByID (Detail Profile)
func (s *userService) GetUserProfile(ctx context.Context, id int64) (*dto.UserDetailResponse, error) {
	user, err := s.userRepo.FindUserByID(ctx, id)
	if err != nil {
		return nil, err // Error "user not found" akan diteruskan
	}

	// Format Tanggal
	createdAtStr := user.CreatedAt.Format("2006-01-02 15:04:05")
	updatedAtStr := ""
	if user.UpdatedAt != nil {
		updatedAtStr = user.UpdatedAt.Format("2006-01-02 15:04:05")
	}
	lastLoginStr := ""
	if user.LastLoginAt != nil {
		lastLoginStr = user.LastLoginAt.Format("2006-01-02 15:04:05")
	}

	return &dto.UserDetailResponse{
		ID:          user.ID,
		FullName:    user.FullName,
		Username:    user.Username,
		Email:       user.Email,
		Role:        user.Role,
		PhoneNumber: user.PhoneNumber,
		IsActive:    user.IsActive,
		LastLoginAt: lastLoginStr,
		CreatedAt:   createdAtStr,
		UpdatedAt:   updatedAtStr,
	}, nil
}

// 4. UpdateUser (Edit Profile)
func (s *userService) ModifyUserData(ctx context.Context, id int64, req dto.UpdateUserRequest) (*dto.UserDetailResponse, error) {
	// A. Cari User Lama dulu
	existingUser, err := s.userRepo.FindUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// B. Update field jika ada di request (Partial Update)
	// Kita cek pointer atau string kosong.

	if req.FullName != "" {
		existingUser.FullName = req.FullName
	}

	// Cek Username unik jika berubah
	if req.Username != "" && req.Username != existingUser.Username {
		if exists, _ := s.userRepo.IsUsernameExists(ctx, req.Username, id); exists {
			return nil, errors.New("USERNAME_EXISTS")
		}
		existingUser.Username = req.Username
	}

	// Cek Email unik jika berubah
	if req.Email != "" && req.Email != existingUser.Email {
		if exists, _ := s.userRepo.IsEmailExists(ctx, req.Email, id); exists {
			return nil, errors.New("EMAIL_EXISTS")
		}
		existingUser.Email = req.Email
	}

	// Cek No HP unik jika berubah
	if req.PhoneNumber != "" && req.PhoneNumber != existingUser.PhoneNumber {
		if exists, _ := s.userRepo.IsPhoneExists(ctx, req.PhoneNumber, id); exists {
			return nil, errors.New("PHONE_EXISTS")
		}
		existingUser.PhoneNumber = req.PhoneNumber
	}

	if req.Role != "" {
		existingUser.Role = req.Role
	}

	// Cek jika password mau diubah
	if req.Password != "" {
		hashedPwd, err := utils.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		existingUser.PasswordHash = hashedPwd
	}

	// Cek IsActive (Pointer check karena boolean bisa false)
	if req.IsActive != nil {
		existingUser.IsActive = *req.IsActive
	}

	// C. Update Timestamp
	now := time.Now()
	existingUser.UpdatedAt = &now

	// D. Simpan Perubahan
	if err := s.userRepo.UpdatedRowUser(ctx, existingUser); err != nil {
		return nil, err
	}

	// E. Kembalikan Data Terbaru (Rekursif panggil GetUserByID biar formatnya sama)
	return s.GetUserProfile(ctx, id)
}

// 5. DeleteUser (Soft Delete)
func (s *userService) DeactivateUserAccount(ctx context.Context, id int64) error {
	// Cek dulu user-nya ada gak
	_, err := s.userRepo.FindUserByID(ctx, id)
	if err != nil {
		return err
	}

	return s.userRepo.SetUserInactive(ctx, id)
}
