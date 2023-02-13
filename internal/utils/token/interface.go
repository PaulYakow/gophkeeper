// Package token содержит реализацию генерации и проверки паролей.
// А также интерфейс для взаимодействия с этим модулем.
package token

//go:generate mockgen -source=interface.go -destination=../../server/mocks/mocks_token.go -package=mocks

import "time"

// IMaker абстракция для управления токенами.
type IMaker interface {
	// Create создаёт токен для переданных id пользователя и продолжительности.
	Create(userID int, duration time.Duration) (string, error)

	// Verify проверяет, является ли токен действительным.
	Verify(in string) (*Payload, error)
}
