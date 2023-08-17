package entity

// UserDTO - объект для  API
type UserDTO struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserDAO - объект для БД
type UserDAO struct {
	ID           int    `db:"id"`
	Login        string `db:"login"`
	PasswordHash string `db:"password_hash"`
}
