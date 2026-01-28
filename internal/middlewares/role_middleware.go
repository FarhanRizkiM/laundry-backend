package middleware

import (
	"laundry-backend/internal/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RoleMiddleware memastikan user memiliki role yang sesuai untuk mengakses endpoint.
// Param `allowedRoles` bisa lebih dari satu, misal: "owner", "cashier".
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {

	return func(c *gin.Context) {
		// 1. Ambil role user dari Context (yang sudah dipasang oleh AuthMiddleware)
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, dto.BaseResponse{
				Success: false,
				Message: "Unauthorized",
				Data: dto.ErrorResponseData{
					ErrorCode: "MISSING_ROLE",
					Errors:    "User role not found in context",
				},
			})
			c.Abort()
			return
		}

		roleStr := userRole.(string)

		// 2. Cek apakah role user ada di dalam daftar allowedRoles
		roleAllowed := false
		for _, role := range allowedRoles {
			if role == roleStr {
				roleAllowed = true
				break
			}
		}

		// 3. Jika tidak cocok, tolak akses (403 Forbidden)
		if !roleAllowed {
			c.JSON(http.StatusForbidden, dto.BaseResponse{
				Success: false,
				Message: "Forbidden Access",
				Data: dto.ErrorResponseData{
					ErrorCode: "FORBIDDEN_ACCESS",
					Errors:    "You do not have permission to access this resource",
				},
			})
			c.Abort()
			return
		}

		// 4. Lanjut
		c.Next()
	}
}
