package usecase

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	IAuthorizationService interface {
		RegisterUser(login, password string) (string, error)
	}

	IAuthorizationRepo interface {
		CreateUser(login, passwordHash string) (int, error)
	}
)
