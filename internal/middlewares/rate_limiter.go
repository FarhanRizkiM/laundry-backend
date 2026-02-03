package middleware

import (
	"laundry-backend/pkg/response"
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
			// REFACTOR: Pakai Helper response.ErrorResponse
			// Code: 429, ErrorCode: RATE_LIMIT_EXCEEDED
			response.ErrorResponse(c, http.StatusTooManyRequests, response.ErrRateLimit, "Too many attempts. Please slow down.", "Rate limit exceeded")
			return
		}

		c.Next()
	}
}
