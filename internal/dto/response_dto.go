package dto

// BaseResponse defines the standard JSON structure for all API responses.
// This structure ensures consistency across success and error responses.
type BaseResponse struct {
	Success bool        `json:"success"`        // Success indicates whether the request was processed successfully.
	Message string      `json:"message"`        // Message provides a human-readable description of the operation result.
	Data    interface{} `json:"data"`           // Data contains the requested payload. It can be a single object, an array, or null.
	Meta    interface{} `json:"meta,omitempty"` // Meta contains pagination information. It is omitted if empty.
}

// MetaData defines the structure for pagination details.
// It is used within the 'Meta' field of BaseResponse.
type MetaData struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
}

// ErrorResponseData defines the structure for detailed error reporting.
// It is typically put inside the 'Data' field when Success is false.
type ErrorResponseData struct {
	ErrorCode string      `json:"error_code"` // ErrorCode is a standardized machine-readable code (e.g., "VALIDATION_ERROR").
	Errors    interface{} `json:"errors"`     // Errors contains specific validation messages or debugging info.
}
