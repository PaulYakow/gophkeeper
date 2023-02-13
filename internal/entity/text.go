package entity

import "time"

// TextDTO - объект текстового типа (заметка) для API
type TextDTO struct {
	ID       int
	Note     string
	Metadata string
}

// TextDAO - объект текстового типа (заметка) для БД
type TextDAO struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	Note      string    `db:"note"`
	Metadata  string    `db:"metadata,omitempty"`
	CreatedAt time.Time `db:"created_at,omitempty"`
}
