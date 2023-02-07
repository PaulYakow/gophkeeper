package usecase

import (
	"context"

	"github.com/PaulYakow/gophkeeper/internal/entity"
)

// PairsService сервис доступа к типу логин/пароль.
type PairsService struct {
	repo IPairsRepo
}

// NewPairsService создаёт объект типа PairsService.
func NewPairsService(repo IPairsRepo) *PairsService {
	return &PairsService{
		repo: repo,
	}
}

// ViewAllPairs получение всех значений типа логин/пароль.
func (s *PairsService) ViewAllPairs(userID int) ([]entity.PairDTO, error) {
	pairsDAO, err := s.repo.GetAllPairs(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	pairsDTO := make([]entity.PairDTO, len(pairsDAO))
	for i, pair := range pairsDAO {
		pairsDTO[i] = entity.PairDTO{
			ID:       pair.ID,
			Login:    pair.Login,
			Password: pair.Password,
			Metadata: pair.Metadata,
		}
	}

	return pairsDTO, nil
}
