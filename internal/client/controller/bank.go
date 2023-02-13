package controller

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/PaulYakow/gophkeeper/internal/entity"
	pb "github.com/PaulYakow/gophkeeper/proto"
)

// BankClient обеспечивает обмен данными о банковских картах пользователя.
type BankClient struct {
	conn *grpc.ClientConn
}

// NewBankClient создаёт объект BankClient.
func NewBankClient(conn *grpc.ClientConn) *BankClient {
	return &BankClient{
		conn: conn,
	}
}

// ViewAllCards запрашивает информацию обо всех имеющихся картах текущего пользователя.
func (c *BankClient) ViewAllCards(ctx context.Context, token string) ([]entity.BankDTO, error) {
	client := pb.NewBankClient(c.conn)
	req := &pb.GetAllCardsRequest{
		Token: "",
	}

	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second))
	defer cancel()

	md := metadata.New(map[string]string{"token": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := client.GetAll(ctx, req)
	if err != nil {
		return nil, err
	}

	out := make([]entity.BankDTO, len(resp.Cards))
	for i, card := range resp.GetCards() {
		out[i] = entity.BankDTO{
			ID:             int(card.GetId()),
			CardHolder:     card.GetCardHolder(),
			Number:         card.GetNumber(),
			ExpirationDate: card.GetExpirationDate(),
			Metadata:       card.GetMetadata(),
		}
	}

	return out, nil
}
