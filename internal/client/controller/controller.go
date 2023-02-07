// Package controller реализация протокола запросов к gRPC-серверу (описанного в gophkeeper/proto/...).
package controller

import (
	"google.golang.org/grpc"
)

type Controller struct {
	Auth  *UserClient
	Pairs *PairsClient
	Token string
}

// New создаёт объект Controller.
func New(conn *grpc.ClientConn) *Controller {
	return &Controller{
		Auth:  NewUserClient(conn),
		Pairs: NewPairsClient(conn),
	}
}
