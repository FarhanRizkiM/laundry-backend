package response

// Standardized Error Codes for the Application.
// Frontend can use these strings to determine logic flow (e.g., show login modal if UNAUTHORIZED).
const (
	// ErrValidation: Input data is missing or invalid format. (400 Bad Request)
	ErrValidation = "VALIDATION_ERROR"

	// ErrUnauthorized: Missing or invalid authentication token. (401 Unauthorized)
	ErrUnauthorized = "UNAUTHORIZED_ACCESS"

	// ErrForbidden: Authenticated but insufficient permissions (e.g., User accessing Owner menu). (403 Forbidden)
	ErrForbidden = "FORBIDDEN_ACCESS"

	// ErrNotFound: The requested resource ID was not found in database. (404 Not Found)
	ErrNotFound = "RESOURCE_NOT_FOUND"

	// ErrDuplicate: Unique constraint violation (e.g., Email/Username already exists). (409 Conflict)
	ErrDuplicate = "DUPLICATE_DATA"

	// ErrRateLimit: Too many requests sent in short period. (429 Too Many Requests)
	ErrRateLimit = "RATE_LIMIT_EXCEEDED"

	// ErrInternal: Unexpected server or database failure. (500 Internal Server Error)
	ErrInternalServer = "INTERNAL_SERVER_ERROR"
)
