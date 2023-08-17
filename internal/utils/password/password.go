// Package password содержит реализацию хеширования и проверки паролей.
// А также интерфейс для взаимодействия с этим модулем.
package password

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Password обеспечивает хеширование и проверку паролей.
type Password struct{}

// New создаёт объект Password.
func New() Password {
	return Password{}
}

// Hash - хеширование пароля.
func (p Password) Hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 7)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedPassword), nil
}

// Check - проверка переданного пароля и оригинального хеша на соответствие.
func (p Password) Check(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
