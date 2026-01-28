package handlers

import (
	"laundry-backend/internal/dto"
	"laundry-backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// 1. Create User (POST /api/v1/users)
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Success: false,
			Message: "Input validation failed",
			Data: dto.ErrorResponseData{
				ErrorCode: "VALIDATION_ERROR",
				Errors:    err.Error(),
			},
		})
		return
	}

	resp, err := h.userService.RegisterUser(c.Request.Context(), req)

	if err != nil {
		errorCode := "INTERNAL_SERVER_ERROR"
		statusCode := http.StatusInternalServerError

		if err.Error() == "USERNAME_EXISTS" || err.Error() == "EMAIL_EXISTS" || err.Error() == "PHONE_EXISTS" {
			errorCode = "DUPLICATE_DATA"
			statusCode = http.StatusConflict
		}

		c.JSON(statusCode, dto.BaseResponse{
			Success: false,
			Message: err.Error(),
			Data: dto.ErrorResponseData{
				ErrorCode: errorCode,
				Errors:    nil,
			},
		})
		return
	}

	c.JSON(http.StatusCreated, dto.BaseResponse{
		Success: true,
		Message: "User created successfully",
		Data:    resp,
	})
}

// 2. Get List Users (GET /api/v1/users)
func (h *UserHandler) GetListUsers(c *gin.Context) {

	// Ambil Query Params (page, per_page, search, role, status)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	search := c.Query("search")
	role := c.Query("role")
	status, _ := strconv.Atoi(c.DefaultQuery("status", "-1"))

	resp, err := h.userService.RetrievedUserDirectory(c.Request.Context(), page, perPage, search, role, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{
			Success: false,
			Message: "Failed to retrieve users",
			Data: dto.ErrorResponseData{
				ErrorCode: "INTERNAL_SERVER_ERROR",
				Errors:    err.Error(),
			},
		})
		return
	}

	// Sukses (Perhatikan cara kita tempel Meta ke BaseResponse)
	c.JSON(http.StatusOK, dto.BaseResponse{
		Success: true,
		Message: "Users retrieved successfully",
		Data:    resp.Data,
		Meta:    resp.Meta,
	})
}

// 3. Get Detail User (GET /api/v1/users/:id)
func (h *UserHandler) GetDetailUser(c *gin.Context) {

	// Ambil ID dari URL
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Success: false,
			Message: "Invalid ID format",
			Data: dto.ErrorResponseData{
				ErrorCode: "VALIDATION_ERROR",
				Errors:    "ID must be a number",
			},
		})
		return
	}

	resp, err := h.userService.GetUserProfile(c.Request.Context(), id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "INTERNAL_SERVER_ERROR"

		if err.Error() == "User not found" {
			statusCode = http.StatusNotFound
			errorCode = "RESOURCE_NOT_FOUND"
		}

		c.JSON(statusCode, dto.BaseResponse{
			Success: false,
			Message: err.Error(),
			Data: dto.ErrorResponseData{
				ErrorCode: errorCode,
				Errors:    nil,
			},
		})
		return
	}

	c.JSON(http.StatusOK, dto.BaseResponse{
		Success: true,
		Message: "User detail retrieved successfully",
		Data:    resp,
	})
}

// 4. Update User (PUT /api/v1/users/:id)
func (h *UserHandler) UpdateUser(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Success: false,
			Message: "Invalid ID format",
		})
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Success: false,
			Message: "Invalid validation failed",
			Data: dto.ErrorResponseData{
				ErrorCode: "VALIDATION_ERROR",
				Errors:    err.Error(),
			},
		})
		return
	}

	resp, err := h.userService.ModifyUserData(c.Request.Context(), id, req)
	if err != nil {
		// Mapping Error Code
		statusCode := http.StatusInternalServerError
		errorCode := "INTERNAL_SERVER_ERROR"

		if err.Error() == "User not found" {
			statusCode = http.StatusNotFound
			errorCode = "RESOURCE_NOT_FOUND"
		} else if err.Error() == "USERNAME_EXISTS" || err.Error() == "EMAIL_EXISTS" || err.Error() == "PHONE_EXISTS" {
			statusCode = http.StatusConflict
			errorCode = "DUPLICATE_DATA"
		}

		c.JSON(statusCode, dto.BaseResponse{
			Success: false,
			Message: err.Error(),
			Data: dto.ErrorResponseData{
				ErrorCode: errorCode,
				Errors:    nil,
			},
		})
		return
	}

	c.JSON(http.StatusOK, dto.BaseResponse{
		Success: true,
		Message: "User updated successfully",
		Data:    resp,
	})
}

// 5. Delete User (DELETE /api/v1/users/:id)
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Success: false,
			Message: "Invalid ID Format",
		})
		return
	}

	if err := h.userService.DeactivateUserAccount(c.Request.Context(), id); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, dto.BaseResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.BaseResponse{
		Success: true,
		Message: "User deleted successfully",
		Data:    gin.H{"id": id}, // Kembalikan ID yang dihapus
	})
}
