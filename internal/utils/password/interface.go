package password

//go:generate mockgen -source=interface.go -destination=../../server/mocks/mocks_password.go -package=mocks

// IPasswordHash абстракция утилиты хеширования и проверки паролей.
type IPasswordHash interface {
	// Hash - хеширование пароля.
	Hash(password string) (string, error)

	// Check - проверка переданного пароля и оригинального хеша на соответствие.
	Check(password, hashedPassword string) error
}
