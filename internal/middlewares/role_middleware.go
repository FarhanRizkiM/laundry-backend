package middleware

import (
	"laundry-backend/pkg/response"
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
			// REFACTOR: Jika role tidak ketemu di context, berarti Auth bermasalah (Unauthorized)
			response.ErrorResponse(c, http.StatusUnauthorized, response.ErrUnauthorized, "Unauthorized", "User role not found in context")
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
			// REFACTOR: Pakai Helper response.ErrorResponse
			// Code: 403, ErrorCode: FORBIDDEN_ACCESS
			response.ErrorResponse(c, http.StatusForbidden, response.ErrForbidden, "Forbidden Access", "You do not have permission to access this resource")
			return
		}

		// 4. Lanjut
		c.Next()
	}
}
