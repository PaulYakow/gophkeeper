// Package controller реализация протокола запросов к gRPC-серверу (описанного в gophkeeper/proto/...).
package controller

import (
	"google.golang.org/grpc"
)

// Controller обеспечивает обмен клиента данными с gRPC-сервером.
type Controller struct {
	Auth  *UserClient
	Pairs *PairsClient
	Cards *BankClient
	Notes *TextClient
	Token string
}

// New создаёт объект Controller.
func New(conn *grpc.ClientConn) *Controller {
	return &Controller{
		Auth:  NewUserClient(conn),
		Pairs: NewPairsClient(conn),
		Cards: NewBankClient(conn),
		Notes: NewTextClient(conn),
	}
}
