package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashedPassword ...
func HashedPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

// CheckPassword ...
func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
