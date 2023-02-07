package controller

import (
	"google.golang.org/grpc"
)

type Controller struct {
	Auth  *UserClient
	Pairs *PairsClient
	Token string
}

func New(conn *grpc.ClientConn) *Controller {
	return &Controller{
		Auth:  NewUserClient(conn),
		Pairs: NewPairsClient(conn),
	}
}
