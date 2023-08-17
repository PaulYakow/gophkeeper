package controller

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/PaulYakow/gophkeeper/internal/server/usecase"
	pb "github.com/PaulYakow/gophkeeper/proto"
)

// BankServer реализация интерфейса proto.BankServer (описание - gophkeeper/proto/bank.proto)
type BankServer struct {
	pb.UnimplementedBankServer
	cards usecase.IBankService
}

// NewBankServer создаёт объект BankServer.
func NewBankServer(cards usecase.IBankService) *BankServer {
	return &BankServer{
		cards: cards,
	}
}

// GetAll - получение всех значений банковских карт.
func (s *BankServer) GetAll(ctx context.Context, req *pb.GetAllCardsRequest) (*pb.GetAllCardsResponse, error) {
	var resp pb.GetAllCardsResponse

	userID, ok := ctx.Value(userIDKey).(int)
	if !ok {
		return nil, status.Error(codes.Aborted, "missing user_id")
	}

	cards, err := s.cards.ViewAllCards(userID)
	if err != nil {
		return nil, err
	}

	for _, card := range cards {
		resp.Cards = append(resp.Cards, &pb.CardMsg{
			Id:             int64(card.ID),
			CardHolder:     card.CardHolder,
			Number:         card.Number,
			ExpirationDate: card.ExpirationDate,
			Metadata:       card.Metadata,
		})
	}

	return &resp, nil
}
