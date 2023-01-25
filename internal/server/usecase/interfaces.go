package usecase

//go:generate mockgen -source=interfaces.go -destination=../mocks/mocks_usecase.go -package=mocks

type (
	IAuthorizationService interface {
		RegisterUser(login, password string) (string, error)
	}

	IAuthorizationRepo interface {
		CreateUser(login, passwordHash string) (int, error)
	}
)
