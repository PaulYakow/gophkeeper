package token

import (
	"errors"
	"time"
)

var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalid")
)

// todo: выделить в интерфейс, чтобы можно было создавать структуры под конкретные задачи (?)

// Payload содержит полезную нагрузку для токена.
type Payload struct {
	UserID    int       `json:"user_id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// NewPayload создаёт объект Payload.
func NewPayload(userID int, duration time.Duration) (*Payload, error) {
	return &Payload{
		UserID:    userID,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}, nil
}

// Valid - проверяет валидность токена.
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
