package handlers

import (
	"laundry-backend/internal/dto"
	"laundry-backend/internal/services"
	"laundry-backend/pkg/response"
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

// CreateUser handles POST /api/v1/users.
// Access: Owner only.
func (h *UserHandler) CreateUser(c *gin.Context) {

	var req dto.CreateUserRequest

	// 1. Validate Input
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, response.ErrValidation, "Input validation failed", err.Error())
		return
	}

	// 2. Call Service
	res, err := h.userService.RegisterUser(c.Request.Context(), req)
	if err != nil {
		// Map Service Errors to HTTP Responses
		if err.Error() == "USERNAME_EXISTS" || err.Error() == "EMAIL_EXISTS" || err.Error() == "PHONE_EXISTS" {
			response.ErrorResponse(c, http.StatusConflict, response.ErrDuplicate, "User data already exists", err.Error())
			return
		}

		// Default Internal Error
		response.ErrorResponse(c, http.StatusInternalServerError, response.ErrInternalServer, "An unexpected server error occurred", nil)
		return
	}

	// 3. Success Response
	response.SuccessCreated(c, "User account created successfully", res)
}

// GetListUsers handles GET /api/v1/users.
// Access: Owner only.
func (h *UserHandler) GetListUsers(c *gin.Context) {

	// 1. Parse Query Parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	search := c.Query("search")
	role := c.Query("role")
	status, _ := strconv.Atoi(c.DefaultQuery("status", "-1"))

	// 2. Call Service
	res, err := h.userService.RetrievedUserDirectory(c.Request.Context(), page, perPage, search, role, status)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, response.ErrInternalServer, "An unexpected server error occurred", nil)
		return
	}

	// 3. Success Response with Meta
	response.SuccessMeta(c, "User retrieved successfully", res.Data, res.Meta)
}

// GetDetailUser handles GET /api/v1/users/:id.
// Access: Owner only.
func (h *UserHandler) GetDetailUser(c *gin.Context) {

	// 1. Parse ID
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, response.ErrValidation, "Invalid ID format", nil)
		return
	}

	// 2. Call Service
	res, err := h.userService.GetUserProfile(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "RESOURCE_NOT_FOUND" {
			response.ErrorResponse(c, http.StatusNotFound, response.ErrNotFound, "User not found", nil)
			return
		}
		response.ErrorResponse(c, http.StatusInternalServerError, response.ErrInternalServer, "An unexpected server error occurred", nil)
		return
	}

	// 3. Success Response
	response.SuccessOK(c, "User detail retrieved successfully", res)
}

// UpdateUser handles PUT /api/v1/users/:id.
// Access: Owner (All), Others (Self only).
func (h *UserHandler) UpdateUser(c *gin.Context) {

	// 1. Parse Target ID
	targetID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, response.ErrValidation, "Invalid ID format", nil)
		return
	}

	// 2. Extract Requester Info (From Auth Middleware)
	// Note: We assume AuthMiddleware sets "user_id" (float64/int64) and "role" (string).
	// We use Type Assertion to ensure safety.
	requesterIDRaw, okID := c.Get("user_id")
	requesterRole, okRole := c.Get("role")

	if !okID || !okRole {
		response.ErrorResponse(c, http.StatusUnauthorized, response.ErrUnauthorized, "Invalid authentication context", nil)
		return
	}

	// Convert float64 (JWT default) to int64
	// --- FIX START: ROBUST TYPE ASSERTION ---
	// Mendeteksi apakah middleware mengirim int64 atau float64
	requesterID, okAssert := requesterIDRaw.(int64)
	if !okAssert {
		// Jika ternyata bukan int64 (misal middleware berubah), return error 500 (Jangan Panic)
		response.ErrorResponse(c, http.StatusInternalServerError, response.ErrInternalServer, "Invalid user ID type (expected int64)", nil)
		return
	}
	// --- FIX END ---

	// 3. Bind JSON Input
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, response.ErrValidation, "Input validation failed", err.Error())
		return
	}

	// 4. Call Service with Requester Context
	res, err := h.userService.ModifyUserData(c.Request.Context(), targetID, req, requesterID, requesterRole.(string))
	if err != nil {
		// Map Specific Errors
		if err.Error() == "FORBIDDEN_ACCESS" {
			response.ErrorResponse(c, http.StatusForbidden, response.ErrForbidden, "You do not have permission to modify this profile", nil)
			return
		}
		if err.Error() == "RESOURCE_NOT_FOUND" {
			response.ErrorResponse(c, http.StatusNotFound, response.ErrNotFound, "User not found", nil)
			return
		}
		if err.Error() == "USERNAME_EXISTS" || err.Error() == "EMAIL_EXISTS" || err.Error() == "PHONE_EXISTS" {
			response.ErrorResponse(c, http.StatusConflict, response.ErrDuplicate, "Data conflict detected (duplicate entry)", err.Error())
			return
		}

		response.ErrorResponse(c, http.StatusInternalServerError, response.ErrInternalServer, "An unexpected server error occurred", nil)
		return
	}

	// 5. Success Response
	response.SuccessOK(c, "User updated successfully", res)
}

// DeleteUser handles DELETE /api/v1/users/:id.
// Access: Owner only (Anti-self delete logic in Service).
func (h *UserHandler) DeleteUser(c *gin.Context) {
	// 1. Parse Target ID
	targetID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, response.ErrValidation, "Invalid ID format", nil)
		return
	}

	// 2. Extract Requester ID (From Auth Middleware)
	requesterIDRaw, ok := c.Get("user_id")
	if !ok {
		response.ErrorResponse(c, http.StatusUnauthorized, response.ErrUnauthorized, "Invalid authentication context", nil)
		return
	}

	// --- STRICT TYPE ASSERTION (INT64 ONLY) ---
	requesterID, okAssert := requesterIDRaw.(int64)
	if !okAssert {
		response.ErrorResponse(c, http.StatusInternalServerError, response.ErrInternalServer, "Invalid user ID type (expected int64)", nil)
		return
	}
	// ------------------------------------------

	// 3. Call Service
	if err := h.userService.DeactivateUserAccount(c.Request.Context(), targetID, requesterID); err != nil {
		if err.Error() == "FORBIDDEN_ACCESS" {
			response.ErrorResponse(c, http.StatusForbidden, response.ErrForbidden, "Action not permitted (cannot delete self)", nil)
			return
		}
		if err.Error() == "RESOURCE_NOT_FOUND" {
			response.ErrorResponse(c, http.StatusNotFound, response.ErrNotFound, "User not found", nil)
			return
		}

		response.ErrorResponse(c, http.StatusInternalServerError, response.ErrInternalServer, "An unexpected server error occurred", nil)
		return
	}

	// 4. Success Response
	// We return the ID of the deleted user as data.
	response.SuccessOK(c, "User account deactivated successfully", gin.H{"id": targetID})
}
