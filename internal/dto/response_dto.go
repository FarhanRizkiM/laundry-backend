package dto

// BaseResponse adalah struktur standar untuk semua respon API.
type BaseResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    interface{} `json:"meta,omitempty"`
}

// MetaData digunakan untuk informasi halaman (pagination).
type MetaData struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
}

// ErrorResponseData memberikan detail jika terjadi kegagalan sistem atau validasi.
type ErrorResponseData struct {
	ErrorCode string      `json:"error_code"`
	Errors    interface{} `json:"errors"`
}
