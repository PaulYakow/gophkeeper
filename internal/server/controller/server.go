// Package controller реализация протокола взаимодействия с сервером gRPC (описанного в gophkeeper/proto/...)
package controller

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/PaulYakow/gophkeeper/cmd/server/config"
	"github.com/PaulYakow/gophkeeper/internal/server/usecase"
	"github.com/PaulYakow/gophkeeper/pkg/logger"
	pb "github.com/PaulYakow/gophkeeper/proto"
)

// todo: также нужно будет упаковать реализацию мелких сервисов в более общую структуру с методом запуска gRPC-сервера.
// Условно будет структура Controller содержащая logger, port и частные структуры (возможно embedded) типа UserServer.

// UserServer реализация интерфейса pb.UserServer (описание - gophkeeper/proto/user.proto)
type UserServer struct {
	pb.UnimplementedUserServer
	auth   usecase.IAuthorizationService
	logger logger.ILogger
	port   string
}

// New создаёт объект UserServer.
func New(authService usecase.IAuthorizationService, l logger.ILogger, cfg *config.Config) *UserServer {
	return &UserServer{
		auth:   authService,
		logger: l,
		port:   ":" + cfg.GRPC.Port,
	}
}

// Run - запуск gRPC-сервера.
func (s *UserServer) Run() {
	go func() {
		listen, err := net.Listen("tcp", s.port)
		if err != nil {
			s.logger.Fatal(fmt.Errorf("gRPC - net.Listen: %w", err))
		}

		// создаём gRPC-сервер без зарегистрированной службы
		grpcSrv := grpc.NewServer()
		// регистрируем сервис
		pb.RegisterUserServer(grpcSrv, s)

		s.logger.Info("gRPC run: %s", s.port)

		// получаем запрос gRPC
		if err = grpcSrv.Serve(listen); err != nil {
			s.logger.Fatal(fmt.Errorf("gRPC - Serve: %w", err))
		}
	}()
}

// Register - регистрация пользователя.
func (s *UserServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	var resp pb.RegisterResponse
	token, err := s.auth.RegisterUser(req.GetLogin(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	resp.Token = token
	return &resp, nil
}

// Login - авторизация пользователя.
func (s *UserServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var resp pb.LoginResponse
	token, err := s.auth.LoginUser(req.GetLogin(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	resp.Token = token
	return &resp, nil
}
