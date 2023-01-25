package password

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Password struct{}

func New() Password {
	return Password{}
}

func (p Password) Hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 7)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedPassword), nil
}

func (p Password) Check(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
