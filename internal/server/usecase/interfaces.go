package usecase

import (
	"context"

	"github.com/PaulYakow/gophkeeper/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=../mocks/mocks_usecase.go -package=mocks

type (
	// IService общая абстракция для взаимодействия с сервисами.
	IService interface {
		IAuthorizationService
		IPairsService
	}

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

		// ParseToken - проверяет переданный токен.
		//
		// Возвращает id пользователя или ошибку.
		ParseToken(token string) (int, error)
	}

	// IPairsService абстракция сервиса доступа к парам логин/пароль.
	IPairsService interface {
		ViewAllPairs(userID int) ([]entity.PairDTO, error)
	}

	// IRepo общая абстракция для взаимодействия с хранилищем.
	IRepo interface {
		IAuthorizationRepo
		IPairsRepo
		CloseConnection() error
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
	}

	// IPairsRepo абстракция взаимодействия с частью хранилища отвечающей за хранение пар логин/пароль.
	IPairsRepo interface {
		GetAllPairs(ctx context.Context, userID int) ([]entity.PairDAO, error)
	}
)
