package repo

import (
	"context"
	"time"

	"github.com/PaulYakow/gophkeeper/internal/entity"
	"github.com/PaulYakow/gophkeeper/pkg/postgres"
)

const (
	createUser = `
INSERT INTO public.users (login, password_hash)
VALUES ($1, $2)
RETURNING id;
`
	getUser = `
SELECT *
FROM users
WHERE login=$1;
`
)

// AuthPostgres реализация интерфейса usecase.IAuthorizationRepo
type AuthPostgres struct {
	db *postgres.Postgres
}

// NewAuthPostgres создаёт объект типа AuthPostgres.
func NewAuthPostgres(db *postgres.Postgres) *AuthPostgres {
	return &AuthPostgres{db}
}

// CreateUser - создание пользователя с заданными логином и хэшем пароля.
//
// Возвращает id пользователя или ошибку (например, если логин уже существует).
func (a *AuthPostgres) CreateUser(login, passwordHash string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var id int
	err := a.db.GetContext(ctx, &id, createUser, login, passwordHash)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetUser - находит пользователя в БД по логину.
//
// Возвращает объект пользователя или ошибку (например, при отсутствии логина).
func (a *AuthPostgres) GetUser(login string) (entity.UserDAO, error) {
	var user entity.UserDAO
	err := a.db.Get(&user, getUser, login)
	return user, err
}
