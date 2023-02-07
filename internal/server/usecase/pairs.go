package usecase

import (
	"context"

	"github.com/PaulYakow/gophkeeper/internal/entity"
)

type PairsService struct {
	repo IPairsRepo
}

func NewPairsService(repo IPairsRepo) *PairsService {
	return &PairsService{
		repo: repo,
	}
}

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
