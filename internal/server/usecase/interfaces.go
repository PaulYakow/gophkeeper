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
		IBankService
		ITextService
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
		// ViewAllPairs получение всех значений типа логин/пароль.
		ViewAllPairs(userID int) ([]entity.PairDTO, error)
	}

	// IBankService абстракция сервиса доступа к банковским картам.
	IBankService interface {
		// ViewAllCards получение всех значений банковских карт.
		ViewAllCards(userID int) ([]entity.BankDTO, error)
	}

	// ITextService абстракция сервиса доступа к заметкам.
	ITextService interface {
		// ViewAllNotes получение всех значений заметок.
		ViewAllNotes(userID int) ([]entity.TextDTO, error)
	}

	// IRepo общая абстракция для взаимодействия с хранилищем.
	IRepo interface {
		IAuthorizationRepo
		IPairsRepo
		IBankRepo
		ITextRepo
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
		// GetAllPairs находит в БД все записи типа логин/пароль принадлежащие конкретному пользователю (userID).
		GetAllPairs(ctx context.Context, userID int) ([]entity.PairDAO, error)
	}

	// IBankRepo абстракция взаимодействия с частью хранилища отвечающей за хранение банковских данных о картах.
	IBankRepo interface {
		// GetAllCards находит в БД все записи банковских карт принадлежащие конкретному пользователю (userID).
		GetAllCards(ctx context.Context, userID int) ([]entity.BankDAO, error)
	}

	ITextRepo interface {
		// GetAllNotes абстракция взаимодействия с частью хранилища отвечающей за хранение заметок.
		GetAllNotes(ctx context.Context, userID int) ([]entity.TextDAO, error)
	}
)
