package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/PaulYakow/gophkeeper/internal/entity"
	"github.com/PaulYakow/gophkeeper/pkg/postgres"
)

const (
	getPairsByUserID = `
SELECT * FROM resources.pairs_data
WHERE user_id = $1
ORDER BY id;
`
)

// PairPostgres реализация интерфейса usecase.IPairsRepo
type PairPostgres struct {
	db *postgres.Postgres
}

// NewPairPostgres создаёт объект типа PairPostgres.
func NewPairPostgres(pg *postgres.Postgres) *PairPostgres {
	return &PairPostgres{pg}
}

// GetAllPairs находит в БД все записи типа логин/пароль принадлежащие конкретному пользователю (userID).
func (p *PairPostgres) GetAllPairs(ctx context.Context, userID int) ([]entity.PairDAO, error) {
	var result []entity.PairDAO

	ctxInner, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	if err := p.db.SelectContext(ctxInner, &result, getPairsByUserID, userID); err != nil {
		return nil, fmt.Errorf("repo - get all pairs by user: %w", err)
	}

	return result, nil
}
