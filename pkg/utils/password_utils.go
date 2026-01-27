package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword menghasilkan hash Bcrypt dari password teks biasa.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ComparePassword membandingkan password teks biasa dengan hash yang tersimpan.
func ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
