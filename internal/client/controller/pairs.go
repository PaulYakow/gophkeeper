package controller

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/PaulYakow/gophkeeper/internal/entity"
	pb "github.com/PaulYakow/gophkeeper/proto"
)

// PairsClient обеспечивает обмен данными о сохранённых парах логин/пароль пользователя.
type PairsClient struct {
	conn *grpc.ClientConn
}

// NewPairsClient создаёт объект PairsClient.
func NewPairsClient(conn *grpc.ClientConn) *PairsClient {
	return &PairsClient{
		conn: conn,
	}
}

// ViewAllPairs запрашивает информацию обо всех имеющихся парах логин/пароль пользователя.
func (c *PairsClient) ViewAllPairs(ctx context.Context, token string) ([]entity.PairDTO, error) {
	client := pb.NewPairClient(c.conn)
	req := &pb.GetAllPairsRequest{
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

	out := make([]entity.PairDTO, len(resp.Pairs))
	for i, pair := range resp.GetPairs() {
		out[i] = entity.PairDTO{
			ID:       int(pair.GetId()),
			Login:    pair.GetLogin(),
			Password: pair.GetPassword(),
			Metadata: pair.GetMetadata(),
		}
	}

	return out, nil
}
