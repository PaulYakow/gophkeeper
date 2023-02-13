package controller

import (
	"context"

	"github.com/PaulYakow/gophkeeper/internal/server/usecase"
	pb "github.com/PaulYakow/gophkeeper/proto"
)

// UserServer реализация интерфейса proto.UserServer (описание - gophkeeper/proto/user.proto)
type UserServer struct {
	pb.UnimplementedUserServer
	auth usecase.IAuthorizationService
}

func NewUserServer(auth usecase.IAuthorizationService) *UserServer {
	return &UserServer{
		auth: auth,
	}
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
