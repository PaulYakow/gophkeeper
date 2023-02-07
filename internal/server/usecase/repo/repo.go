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
CREATE TABLE IF NOT EXISTS public.users
(
    id            SERIAL PRIMARY KEY,
    login         VARCHAR NOT NULL UNIQUE,
    password_hash VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS public.pairs
(
    id         SERIAL PRIMARY KEY,
    user_id	   INT REFERENCES public.users (id) ON DELETE CASCADE,
    login      VARCHAR NOT NULL,
    password   VARCHAR NOT NULL,
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
}

// New создаёт объект Repo.
func New(db *postgres.Postgres, auth usecase.IAuthorizationRepo, pairs usecase.IPairsRepo) (*Repo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, schema)
	if err != nil {
		return nil, fmt.Errorf("repo - create table failed: %w", err)
	}

	return &Repo{
		db,
		auth,
		pairs,
	}, nil
}

// CloseConnection - дожидается завершения запросов и закрывает все открытые соединения.
func (s *Repo) CloseConnection() error {
	return s.db.Shutdown()
}
