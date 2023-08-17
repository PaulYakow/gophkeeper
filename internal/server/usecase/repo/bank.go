package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/PaulYakow/gophkeeper/internal/entity"
	"github.com/PaulYakow/gophkeeper/pkg/postgres"
)

const (
	getCardsByUserID = `
SELECT * FROM resources.bank_data
WHERE user_id = $1
ORDER BY id;
`
)

// BankPostgres реализация интерфейса usecase.IBankRepo
type BankPostgres struct {
	db *postgres.Postgres
}

// NewBankPostgres создаёт объект типа BankPostgres.
func NewBankPostgres(pg *postgres.Postgres) *BankPostgres {
	return &BankPostgres{pg}
}

// GetAllCards находит в БД все записи банковских карт принадлежащие конкретному пользователю (userID).
func (p *BankPostgres) GetAllCards(ctx context.Context, userID int) ([]entity.BankDAO, error) {
	var result []entity.BankDAO

	ctxInner, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	if err := p.db.SelectContext(ctxInner, &result, getCardsByUserID, userID); err != nil {
		return nil, fmt.Errorf("repo - get all cards by user: %w", err)
	}

	return result, nil
}
