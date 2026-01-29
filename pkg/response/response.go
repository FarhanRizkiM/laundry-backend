package response

import (
	"laundry-backend/internal/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SuccessOK sends a standard success response with HTTP 200 status code.
// Use this for successful GET (Detail), PUT, and DELETE operations.
func SuccessOK(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, dto.BaseResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SuccessCreated sends a standard success response with HTTP 201 status code.
// Use this specifically for successful POST (Creation) operations.
func SuccessCreated(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, dto.BaseResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SuccessWithMeta sends a standard success response containing pagination metadata.
// Use this for list endpoints (GET) that support pagination.
func SuccessMeta(c *gin.Context, message string, data interface{}, meta interface{}) {
	c.JSON(http.StatusOK, dto.BaseResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// ErrorResponse sends a standardized error response.
// This function acts as a general error handler for various HTTP status codes (400, 401, 404, 500, etc.).
//
// Parameters:
//   - code: HTTP Status Code (e.g., http.StatusBadRequest).
//   - errorCode: A standardized string code (e.g., "VALIDATION_ERROR", "DUPLICATE_DATA").
//   - message: A human-readable error message in English.
//   - errors: Detailed error info (can be nil, string, or struct).
func ErrorResponse(c *gin.Context, code int, errorCode string, message string, errors interface{}) {
	c.JSON(code, dto.BaseResponse{
		Success: false,
		Message: message,
		Data: dto.ErrorResponseData{
			ErrorCode: errorCode,
			Errors:    errors,
		},
	})
	// Abort ensures that no further handlers are executed in the Gin middleware chain.
	c.Abort()
}
