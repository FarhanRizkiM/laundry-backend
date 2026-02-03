package handlers

import (
	"fmt"
	"laundry-backend/internal/dto"
	"laundry-backend/internal/services"
	"laundry-backend/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service services.CategoryService
}

func NewCategoryHandler(service services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) HandleCreateCategory(c *gin.Context) {

	var req dto.CreateCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, response.ErrValidation, "Invalid request payload", err.Error())
		return
	}

	result, err := h.service.CreateCategory(c.Request.Context(), req)
	if err != nil {
		if err.Error() == response.ErrDuplicate {
			response.ErrorResponse(c, http.StatusConflict, response.ErrDuplicate, "Category name already exists", nil)
			return
		}

		fmt.Printf("[ERROR] CreateCategory: %v\n", err)

		response.ErrorResponse(c, http.StatusInternalServerError, response.ErrInternalServer, "Failed to create category", nil)
		return
	}

	response.SuccessCreated(c, "Category created successfully", result)
}

func (h *CategoryHandler) HandleGetCategoryList(c *gin.Context) {

	pageStr := c.DefaultQuery("page", "1")
	perPageStr := c.DefaultQuery("per_page", "10")
	search := c.Query("search")
	status := c.Query("status")
	sortBy := c.Query("sort_by")
	sortOrder := c.Query("order")

	page, _ := strconv.Atoi(pageStr)
	perPage, _ := strconv.Atoi(perPageStr)

	result, err := h.service.GetCategoryList(c.Request.Context(), page, perPage, search, status, sortBy, sortOrder)
	if err != nil {
		fmt.Printf("[ERROR] GetCategoryList: %v\n", err)

		response.ErrorResponse(c, http.StatusInternalServerError, response.ErrInternalServer, "Failed to retrieve categories", nil)
		return
	}

	response.SuccessMeta(c, "Categories retrieved successfully", result.Data, result.Meta)
}

func (h *CategoryHandler) HandleGetCategoryDetail(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, response.ErrValidation, "Invalid ID Format", "ID must be a number")
		return
	}

	result, err := h.service.GetCategoryDetail(c.Request.Context(), id)
	if err != nil {
		if err.Error() == response.ErrNotFound {
			response.ErrorResponse(c, http.StatusNotFound, response.ErrNotFound, "Category not found", nil)
			return
		}
		fmt.Printf("[ERROR] GetCategoryDetail: %v\n", err)

		response.ErrorResponse(c, http.StatusInternalServerError, response.ErrInternalServer, "Failed to retrieve category detail", nil)
		return
	}

	response.SuccessOK(c, "Category detail retrieved successfully", result)
}

func (h *CategoryHandler) HandleUpdateCategory(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, response.ErrValidation, "Invalid ID Format", "ID must be a number")
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, response.ErrValidation, "Invalid request payload", err.Error())
		return
	}

	result, err := h.service.ModifyCategory(c.Request.Context(), id, req)
	if err != nil {
		if err.Error() == response.ErrNotFound {
			response.ErrorResponse(c, http.StatusNotFound, response.ErrNotFound, "Category not found", nil)
			return
		}
		if err.Error() == response.ErrDuplicate {
			response.ErrorResponse(c, http.StatusConflict, response.ErrDuplicate, "Category name already taken", nil)
			return
		}

		fmt.Printf("[ERROR] ModifyCategory: %v\n", err)

		response.ErrorResponse(c, http.StatusInternalServerError, response.ErrInternalServer, "Failed to update category", nil)
		return
	}

	response.SuccessOK(c, "Category updated successfully", result)
}

func (h *CategoryHandler) HandleDeleteCategory(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, response.ErrValidation, "Invalid ID format", "ID must be a number")
		return
	}

	err = h.service.DeactivateCategory(c.Request.Context(), id)
	if err != nil {
		if err.Error() == response.ErrNotFound {
			response.ErrorResponse(c, http.StatusNotFound, response.ErrNotFound, "Category not found", nil)
			return
		}

		fmt.Printf("[ERROR] DeactivateCategory: %v\n", err)

		response.ErrorResponse(c, http.StatusInternalServerError, response.ErrInternalServer, "Failed to delete category", nil)
		return
	}

	response.SuccessOK(c, "Category deleted successfully", map[string]int64{"id": id})
}
