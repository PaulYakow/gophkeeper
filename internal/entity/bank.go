package entity

import "time"

// BankDTO - объект типа банковской карты для API
type BankDTO struct {
	ID             int
	CardHolder     string
	Number         string
	ExpirationDate string
	Metadata       string
}

// BankDAO - объект типа банковской карты для БД
type BankDAO struct {
	ID             int       `db:"id"`
	UserID         int       `db:"user_id"`
	CardHolder     string    `db:"card_holder"`
	Number         string    `db:"number"`
	ExpirationDate string    `db:"expiration_date"`
	Metadata       string    `db:"metadata,omitempty"`
	CreatedAt      time.Time `db:"created_at,omitempty"`
}
