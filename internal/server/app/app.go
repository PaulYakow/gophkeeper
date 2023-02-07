// Package app точка входа сервера.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/PaulYakow/gophkeeper/cmd/server/config"
	"github.com/PaulYakow/gophkeeper/internal/server/controller"
	"github.com/PaulYakow/gophkeeper/internal/server/usecase"
	"github.com/PaulYakow/gophkeeper/internal/server/usecase/repo"
	"github.com/PaulYakow/gophkeeper/internal/utils/password"
	"github.com/PaulYakow/gophkeeper/internal/utils/token"
	"github.com/PaulYakow/gophkeeper/pkg/logger"
	"github.com/PaulYakow/gophkeeper/pkg/postgres"
)

// App основная структура приложения (сервера).
type App struct {
	config  *config.Config
	logger  *logger.Logger
	repo    *repo.Repo
	service usecase.IService

	// todo: по сути это вспомогательные утилиты - можно сделать отдельную структуру
	// и потом, например a.utils = a.createUtils
	passwordHasher password.IPasswordHash
	tokenMaker     token.IMaker

	grpcSrv *controller.Controller
}

// New собирает сервер из слоёв (хранилище, сервисы, логика, контроллер).
func New(cfg *config.Config) *App {
	var err error

	// Config + Logger + Password hasher
	a := &App{
		config:         cfg,
		logger:         logger.New(cfg.App.Name),
		passwordHasher: password.New(),
	}

	// Repo
	a.repo = a.createPostgresRepo()

	// Token
	a.tokenMaker, err = token.NewPasetoMaker(a.config.Token.Key)
	if err != nil {
		a.logger.Fatal(fmt.Errorf("create token maker: %w", err))
	}

	// Usecase
	auth := usecase.NewAuthService(a.repo, a.passwordHasher, a.tokenMaker)
	pairs := usecase.NewPairsService(a.repo)
	a.service, err = usecase.New(auth, pairs)
	if err != nil {
		a.logger.Fatal(fmt.Errorf("create service: %w", err))
	}

	// Controller
	a.grpcSrv = controller.New(a.service, a.logger, cfg)

	return a
}

// Run - запуск сервера.
func (a *App) Run() {
	defer a.logger.Exit()

	if a.grpcSrv != nil {
		a.grpcSrv.Run()
	}

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	sig := <-interrupt
	a.logger.Info("Run - signal: %v", sig.String())

	// Shutdown
	if err := a.repo.CloseConnection(); err != nil {
		a.logger.Error(fmt.Errorf("run - close connection to repo: %w", err))
	}
}

// Создание структуры взаимодействия с хранилищем данных.
func (a *App) createPostgresRepo() (r *repo.Repo) {
	pg, err := postgres.New(a.config.PG.URL,
		postgres.ConnAttempts(a.config.PG.ConnAttempts), postgres.MaxOpenConn(a.config.PG.MaxOpen))
	if err != nil {
		a.logger.Fatal(fmt.Errorf("create DB conn: %w", err))
	}

	a.logger.Info("PostgreSQL connection ok")

	auth := repo.NewAuthPostgres(pg)
	pairs := repo.NewPairPostgres(pg)

	r, err = repo.New(pg, auth, pairs)
	if err != nil {
		a.logger.Fatal(fmt.Errorf("Run - repo.New: %w", err))
	}
	a.logger.Info("PostgreSQL in use")

	return
}
