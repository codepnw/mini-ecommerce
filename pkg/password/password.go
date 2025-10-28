package password

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashedPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password failed: %w", err)
	}
	return string(hashed), nil
}

func ComparePassword(hashed, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	if err != nil {
		return fmt.Errorf("compare password failed: %w", err)
	}
	return nil
}
