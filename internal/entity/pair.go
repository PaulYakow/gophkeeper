package entity

import "time"

// PairDTO - объект типа логин/пароль для API
type PairDTO struct {
	ID       int
	Login    string
	Password string
	Metadata string
}

// PairDAO - объект типа логин/пароль для БД
type PairDAO struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	Login     string    `db:"login"`
	Password  string    `db:"password"`
	Metadata  string    `db:"metadata,omitempty"`
	CreatedAt time.Time `db:"created_at,omitempty"`
}
