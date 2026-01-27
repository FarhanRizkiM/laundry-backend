package utils

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims merepresentasikan payload data yang akan dimasukkan ke dalam JWT.
type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	JTI      string `json:"jti"`
	jwt.RegisteredClaims
}

// getSecretKey mengambil kunci rahasia dari environment variable untuk menandatangani token.
func getSecretKey() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// PENTING: Gunakan fallback hanya untuk development, jangan untuk production.
		return []byte("default_secret_key_fallback")
	}
	return []byte(secret)
}

// getExpiryDuration menentukan berapa lama access token berlaku (dalam menit).
func getExpiryDuration() time.Duration {
	expiryStr := os.Getenv("JWT_EXPIRY_MINUTE")
	expiry, err := strconv.Atoi(expiryStr)

	if err != nil || expiry <= 0 {
		return 15 * time.Minute // Default fallback: 15 menit
	}

	return time.Duration(expiry) * time.Minute
}

// GenerateAccessToken membuat token baru berjenis HS256 untuk akses autentikasi. Mengembalikan token string, ID unik token (JTI), dan error jika gagal.
func GenerateAccessToken(userID int64, username, role string) (string, string, error) {
	jti := uuid.New().String()

	expirationTime := time.Now().Add(getExpiryDuration())

	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		JTI:      jti,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(getSecretKey())

	return tokenString, jti, err
}

// GenerateRefreshToken menghasilkan string acak unik untuk kebutuhan pembaruan access token.
func GenerateRefreshToken() (string, error) {
	return uuid.New().String(), nil
}
