package repo

import (
	"context"
	"fmt"
	"time"

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
)

type Server struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) (*Server, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := pg.ExecContext(ctx, schema)
	if err != nil {
		return nil, fmt.Errorf("repo - create table failed: %w", err)
	}

	return &Server{pg}, nil
}

func (s *Server) CreateUser(login, passwordHash string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var id int
	err := s.GetContext(ctx, &id, createUser, login, passwordHash)
	if err != nil {
		return 0, err
	}

	return id, nil
}
