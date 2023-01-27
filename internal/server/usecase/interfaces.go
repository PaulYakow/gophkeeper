package usecase

import "github.com/PaulYakow/gophkeeper/internal/entity"

//go:generate mockgen -source=interfaces.go -destination=../mocks/mocks_usecase.go -package=mocks

type (
	// IAuthorizationService абстракция сервиса авторизации.
	IAuthorizationService interface {
		// RegisterUser - регистрация нового пользователя с переданным логином и паролем.
		//
		// Возвращает токен или ошибку (например, если логин уже существует).
		RegisterUser(login, password string) (string, error)

		// LoginUser - авторизация существующего пользователя.
		//
		// Возвращает токен или ошибку (например, если логина не существует).
		LoginUser(login, password string) (string, error)
	}

	// IAuthorizationRepo абстракция взаимодействия с частью хранилища отвечающей за авторизацию пользователей.
	IAuthorizationRepo interface {
		// CreateUser - создание пользователя с заданными логином и хэшем пароля.
		//
		// Возвращает id пользователя или ошибку (если логин уже существует).
		CreateUser(login, passwordHash string) (int, error)

		// GetUser - находит пользователя в БД по логину.
		//
		// Возвращает объект пользователя или ошибку (при отсутствии логина).
		GetUser(login string) (entity.UserDAO, error)

		// CloseConnection - дожидается завершения запросов и закрывает все открытые соединения.
		CloseConnection() error
	}
)
