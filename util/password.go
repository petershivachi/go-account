package util

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword returns bcrypt hash of the password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("failed to hash password: %v", err)
	}

	return string(bytes), nil
}

// CheckPasswordHash checks if the provided password is correct or not
func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
