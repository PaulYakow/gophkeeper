package controller

import (
	"context"
	"time"

	"google.golang.org/grpc"

	pb "github.com/PaulYakow/gophkeeper/proto"
)

type UserClient struct {
	conn *grpc.ClientConn
}

func NewUserClient(conn *grpc.ClientConn) *UserClient {
	return &UserClient{
		conn: conn,
	}
}

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
