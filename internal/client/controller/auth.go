package controller

import (
	"context"
	"time"

	"google.golang.org/grpc"

	pb "github.com/PaulYakow/gophkeeper/proto"
)

// UserClient обеспечивает регистрацию/аутентификацию пользователя.
type UserClient struct {
	conn *grpc.ClientConn
}

// NewUserClient создаёт объект UserClient.
func NewUserClient(conn *grpc.ClientConn) *UserClient {
	return &UserClient{
		conn: conn,
	}
}

// Register регистрация пользователя с заданными логином и паролем.
// Возвращает ошибку если пользователь с таким логином уже существует.
func (c *UserClient) Register(ctx context.Context, login, password string) (string, error) {
	client := pb.NewUserClient(c.conn)
	req := &pb.RegisterRequest{
		Login:    login,
		Password: password,
	}

	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second))
	defer cancel()

	resp, err := client.Register(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.GetToken(), nil
}

// Login аутентификация пользователя по переданным логину и паролю.
func (c *UserClient) Login(ctx context.Context, login, password string) (string, error) {
	client := pb.NewUserClient(c.conn)
	req := &pb.LoginRequest{
		Login:    login,
		Password: password,
	}

	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second))
	defer cancel()

	resp, err := client.Login(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.GetToken(), nil
}
