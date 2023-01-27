// Package repo содержит структуру и методы взаимодействия с базой данных.
package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/PaulYakow/gophkeeper/internal/entity"
	"github.com/PaulYakow/gophkeeper/pkg/postgres"
)

const (
	schema = `
CREATE TABLE IF NOT EXISTS users
(
    id            SERIAL PRIMARY KEY,
    login         VARCHAR NOT NULL UNIQUE,
    password_hash VARCHAR NOT NULL
);
`
	createUser = `
INSERT INTO users (login, password_hash)
VALUES ($1, $2)
RETURNING id;
`
	getUser = `
SELECT *
FROM users
WHERE login=$1;
`
)

// Repo реализация хранилища. Хранение в БД Postgres (драйвер - sqlx).
//
// Содержит реализации необходимых интерфейсов (IAuthorizationRepo, ...).
type Repo struct {
	*postgres.Postgres
}

// New создаёт объект Repo.
func New(pg *postgres.Postgres) (*Repo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := pg.ExecContext(ctx, schema)
	if err != nil {
		return nil, fmt.Errorf("repo - create table failed: %w", err)
	}

	return &Repo{pg}, nil
}

// todo: переместится в 'user.go' и свой интерфейс в 'interfaces.go' ==>
// (аналогично остальные интерфейсы в свои файлы размещать)

// CreateUser - создание пользователя с заданными логином и хэшем пароля.
//
// Возвращает id пользователя или ошибку (если логин уже существует).
func (s *Repo) CreateUser(login, passwordHash string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var id int
	err := s.GetContext(ctx, &id, createUser, login, passwordHash)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetUser - находит пользователя в БД по логину.
//
// Возвращает объект пользователя или ошибку (при отсутствии логина).
func (s *Repo) GetUser(login string) (entity.UserDAO, error) {
	var user entity.UserDAO
	err := s.Get(&user, getUser, login)
	return user, err
}

// <==

// todo: должно остаться в данной структуре ==>

// CloseConnection - дожидается завершения запросов и закрывает все открытые соединения.
func (s *Repo) CloseConnection() error {
	return s.Shutdown()
}

// <==
