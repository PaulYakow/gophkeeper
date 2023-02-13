package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/PaulYakow/gophkeeper/internal/entity"
	"github.com/PaulYakow/gophkeeper/pkg/postgres"
)

const (
	getNotesByUserID = `
SELECT * FROM resources.text_data
WHERE user_id = $1
ORDER BY id;
`
)

// TextPostgres реализация интерфейса usecase.ITextRepo
type TextPostgres struct {
	db *postgres.Postgres
}

// NewTextPostgres создаёт объект типа TextPostgres.
func NewTextPostgres(pg *postgres.Postgres) *TextPostgres {
	return &TextPostgres{pg}
}

// GetAllNotes находит в БД все записи заметок принадлежащие конкретному пользователю (userID).
func (p *TextPostgres) GetAllNotes(ctx context.Context, userID int) ([]entity.TextDAO, error) {
	var result []entity.TextDAO

	ctxInner, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	if err := p.db.SelectContext(ctxInner, &result, getNotesByUserID, userID); err != nil {
		return nil, fmt.Errorf("repo - get all cards by user: %w", err)
	}

	return result, nil
}
