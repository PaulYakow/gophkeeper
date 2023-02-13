package usecase

import (
	"context"

	"github.com/PaulYakow/gophkeeper/internal/entity"
)

// BankService сервис доступа к банковским картам.
type BankService struct {
	repo IBankRepo
}

// NewBankService создаёт объект типа BankService.
func NewBankService(repo IBankRepo) *BankService {
	return &BankService{
		repo: repo,
	}
}

// ViewAllCards получение всех значений банковских карт.
func (s *BankService) ViewAllCards(userID int) ([]entity.BankDTO, error) {
	cardsDAO, err := s.repo.GetAllCards(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	cardsDTO := make([]entity.BankDTO, len(cardsDAO))
	for i, card := range cardsDAO {
		cardsDTO[i] = entity.BankDTO{
			ID:             card.ID,
			CardHolder:     card.CardHolder,
			Number:         card.Number,
			ExpirationDate: card.ExpirationDate,
			Metadata:       card.Metadata,
		}
	}

	return cardsDTO, nil
}
