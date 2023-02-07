// Package controller реализация протокола взаимодействия с сервером gRPC (описанного в gophkeeper/proto/...)
package controller

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/PaulYakow/gophkeeper/cmd/server/config"
	"github.com/PaulYakow/gophkeeper/internal/server/usecase"
	"github.com/PaulYakow/gophkeeper/pkg/logger"
	pb "github.com/PaulYakow/gophkeeper/proto"
)

type ctxKey string

const userIDKey ctxKey = "user_id"

type Controller struct {
	service usecase.IService
	logger  logger.ILogger
	port    string
}

// New создаёт объект Controller.
func New(service usecase.IService, l logger.ILogger, cfg *config.Config) *Controller {
	return &Controller{
		service: service,
		logger:  l,
		port:    ":" + cfg.GRPC.Port,
	}
}

// Run - запуск gRPC-сервера.
func (c *Controller) Run() {
	go func() {
		listen, err := net.Listen("tcp", c.port)
		if err != nil {
			c.logger.Fatal(fmt.Errorf("gRPC - net.Listen: %w", err))
		}

		// создаём gRPC-сервер без зарегистрированной службы
		grpcSrv := grpc.NewServer(grpc.UnaryInterceptor(c.userIdentity))

		// регистрируем сервисы
		pb.RegisterUserServer(grpcSrv, NewUserServer(c.service))
		pb.RegisterPairServer(grpcSrv, NewPairsServer(c.service))

		c.logger.Info("gRPC run: %s", c.port)

		// получаем запрос gRPC
		if err = grpcSrv.Serve(listen); err != nil {
			c.logger.Fatal(fmt.Errorf("gRPC - Serve: %w", err))
		}
	}()
}

// Идентификация пользователя.
func (c *Controller) userIdentity(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	if _, ok := info.Server.(*UserServer); ok {
		return handler(ctx, req)
	}

	var token string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("token")
		if len(values) > 0 {
			token = values[0]
		}
	}

	if token == "" {
		c.logger.Error(fmt.Errorf("missing token in metadata"))
		return nil, status.Error(codes.FailedPrecondition, "missing token")
	}

	userID, err := c.service.ParseToken(token)
	if err != nil {
		c.logger.Error(fmt.Errorf("user identity: %w", err))
		return nil, status.Error(codes.Aborted, "user identity error")
	}

	ctx = context.WithValue(ctx, userIDKey, userID)

	return handler(ctx, req)
}
