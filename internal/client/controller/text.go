package controller

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/PaulYakow/gophkeeper/internal/entity"
	pb "github.com/PaulYakow/gophkeeper/proto"
)

// TextClient обеспечивает обмен данными о сохранённых заметках пользователя.
type TextClient struct {
	conn *grpc.ClientConn
}

// NewTextClient создаёт объект TextClient.
func NewTextClient(conn *grpc.ClientConn) *TextClient {
	return &TextClient{
		conn: conn,
	}
}

// ViewAllNotes запрашивает информацию обо всех имеющихся заметках пользователя.
func (c *TextClient) ViewAllNotes(ctx context.Context, token string) ([]entity.TextDTO, error) {
	client := pb.NewTextClient(c.conn)
	req := &pb.GetAllNotesRequest{
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

	out := make([]entity.TextDTO, len(resp.Notes))
	for i, note := range resp.GetNotes() {
		out[i] = entity.TextDTO{
			ID:       int(note.GetId()),
			Note:     note.GetNote(),
			Metadata: note.GetMetadata(),
		}
	}

	return out, nil
}
