// Package app точка входа клиентского приложения.
package app

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/PaulYakow/gophkeeper/cmd/client/config"
	"github.com/PaulYakow/gophkeeper/internal/client/controller"
	"github.com/PaulYakow/gophkeeper/internal/client/views"
	"github.com/PaulYakow/gophkeeper/pkg/logger"
)

// App основная структура приложения (сервера).
type App struct {
	config *config.Config
	logger *logger.Logger
	conn   *grpc.ClientConn
	ctrl   *controller.Controller
	view   *views.View
}

// New собирает клиентское приложение из слоёв (хранилище, сервисы, логика, контроллер, представление).
func New(cfg *config.Config) (a *App) {
	a = &App{
		config: cfg,
		logger: logger.New(cfg.App.Name),
	}

	var err error

	target := cfg.GRPC.Address + ":" + cfg.GRPC.Port
	a.conn, err = grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		a.logger.Fatal(err)
	}

	a.ctrl = controller.New(a.conn)
	a.view = views.New(a.ctrl, cfg)

	return
}

func (a *App) Run() {
	defer a.conn.Close()
	defer a.logger.Exit()

	a.view.Run()
}
