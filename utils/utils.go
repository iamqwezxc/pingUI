package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword хеширует пароль с использованием bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
