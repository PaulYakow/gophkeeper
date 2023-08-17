package usecase

import (
	"context"

	"github.com/PaulYakow/gophkeeper/internal/entity"
)

// TextService сервис доступа к заметкам.
type TextService struct {
	repo ITextRepo
}

// NewTextService создаёт объект типа TextService.
func NewTextService(repo ITextRepo) *TextService {
	return &TextService{
		repo: repo,
	}
}

// ViewAllNotes получение всех значений заметок.
func (s *TextService) ViewAllNotes(userID int) ([]entity.TextDTO, error) {
	notesDAO, err := s.repo.GetAllNotes(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	notesDTO := make([]entity.TextDTO, len(notesDAO))
	for i, card := range notesDAO {
		notesDTO[i] = entity.TextDTO{
			ID:       card.ID,
			Note:     card.Note,
			Metadata: card.Metadata,
		}
	}

	return notesDTO, nil
}
