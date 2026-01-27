package middlewares

import (
	"laundry-backend/internal/dto"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// client menyimpan informasi limiter untuk setiap IP
type client struct {
	limiter *rate.Limiter
}

var (
	clients = make(map[string]*client)
	mu      sync.Mutex
)

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Mengambil IP user sebagai pembeda
		ip := c.ClientIP()

		mu.Lock()
		if _, found := clients[ip]; !found {
			// Setelan: 5 request per detik, dengan maksimal ledakan (burst) 10 request
			clients[ip] = &client{limiter: rate.NewLimiter(5, 10)}
		}

		limiter := clients[ip].limiter
		mu.Unlock()

		// Cek apakah user sudah melewati batas
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, dto.BaseResponse{
				Success: false,
				Message: "Too many attempts. Please slow down.",
				Data: dto.ErrorResponseData{
					ErrorCode: "TOO_MANY_REQUESTS",
					Errors:    "Rate limit exceeded",
				},
			})
			c.Abort() // Menghentikan request agar tidak lanjut ke Handler
			return
		}

		c.Next()
	}
}
