// Package v2 является обёрткой над библиотекой github.com/jmoiron/sqlx.
package postgres

import (
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

const (
	defaultMaxOpenConn     = 4
	defaultMaxIdleConn     = 4
	defaultMaxConnIdleTime = time.Second * 30
	defaultMaxConnLifetime = time.Minute * 2

	defaultConnAttempts = 10
	defaultConnTimeout  = time.Second
)

// Postgres структура с настройками подключения к БД и доступом к текущему соединению.
type Postgres struct {
	*sqlx.DB

	maxOpenConn     int
	maxIdleConn     int
	connAttempts    int
	maxConnIdleTime time.Duration
	maxConnLifetime time.Duration
	connTimeout     time.Duration
}

// New создаёт объект Postgres с заданными параметрами и подключается к БД.
func New(dsn string, opts ...Option) (*Postgres, error) {
	pg := &Postgres{
		maxOpenConn:     defaultMaxOpenConn,
		maxIdleConn:     defaultMaxIdleConn,
		maxConnIdleTime: defaultMaxConnIdleTime,
		maxConnLifetime: defaultMaxConnLifetime,
		connAttempts:    defaultConnAttempts,
		connTimeout:     defaultConnTimeout,
	}

	for _, opt := range opts {
		opt(pg)
	}

	var err error

	for pg.connAttempts > 0 {
		if pg.DB, err = sqlx.Connect("pgx", dsn); err == nil {
			break
		}

		log.Printf("Postgres is trying to connect, attempts left: %d", pg.connAttempts)

		time.Sleep(pg.connTimeout)

		pg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgre - connAttempts == 0: %w", err)
	}

	pg.SetMaxOpenConns(pg.maxOpenConn)
	pg.SetMaxIdleConns(pg.maxIdleConn)
	pg.SetConnMaxIdleTime(pg.maxConnIdleTime)
	pg.SetConnMaxLifetime(pg.maxConnLifetime)

	return pg, nil
}

func (pg *Postgres) Shutdown() error {
	return pg.Close()
}
