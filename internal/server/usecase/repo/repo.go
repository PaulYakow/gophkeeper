// Package repo содержит структуру и методы взаимодействия с базой данных.
package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/PaulYakow/gophkeeper/internal/server/usecase"
	"github.com/PaulYakow/gophkeeper/pkg/postgres"
)

const (
	schema = `
CREATE SCHEMA IF NOT EXISTS resources;

CREATE TABLE IF NOT EXISTS public.users
(
    id            SERIAL PRIMARY KEY,
    login         VARCHAR NOT NULL UNIQUE,
    password_hash VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS resources.pairs_data
(
    id         SERIAL PRIMARY KEY,
    user_id	   INT REFERENCES public.users (id) ON DELETE CASCADE,
    login      VARCHAR NOT NULL,
    password   VARCHAR NOT NULL,
    metadata   TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS resources.bank_data
(
    id         SERIAL PRIMARY KEY,
    user_id	   INT REFERENCES public.users (id) ON DELETE CASCADE,
	card_holder VARCHAR NOT NULL,
	number VARCHAR NOT NULL,
	expiration_date VARCHAR NOT NULL,
    metadata   TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS resources.text_data
(
    id         SERIAL PRIMARY KEY,
    user_id	   INT REFERENCES public.users (id) ON DELETE CASCADE,
	note       TEXT,
    metadata   TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
`
)

// Repo реализация хранилища. Хранение в БД Postgres (драйвер - sqlx).
//
// Содержит реализации необходимых интерфейсов (IAuthorizationRepo, ...).
type Repo struct {
	db *postgres.Postgres
	usecase.IAuthorizationRepo
	usecase.IPairsRepo
	usecase.IBankRepo
	usecase.ITextRepo
}

// New создаёт объект Repo.
func New(db *postgres.Postgres,
	auth usecase.IAuthorizationRepo,
	pairs usecase.IPairsRepo,
	cards usecase.IBankRepo,
	notes usecase.ITextRepo,
) (*Repo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, schema)
	if err != nil {
		return nil, fmt.Errorf("repo - create schema failed: %w", err)
	}

	return &Repo{
		db,
		auth,
		pairs,
		cards,
		notes,
	}, nil
}

// CloseConnection - дожидается завершения запросов и закрывает все открытые соединения.
func (s *Repo) CloseConnection() error {
	return s.db.Shutdown()
}
