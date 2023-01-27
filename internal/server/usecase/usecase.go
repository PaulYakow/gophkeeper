// Package usecase содержит реализацию логики взаимодействия с сервисами/хранилищем.
// А также интерфейсы для взаимодействия с этим слоем.
package usecase

import (
	"time"

	"github.com/PaulYakow/gophkeeper/internal/utils/password"
	"github.com/PaulYakow/gophkeeper/internal/utils/token"
)

// Usecase обеспечивает логику взаимодействия с сервисами/хранилищем.
//
//goland:noinspection SpellCheckingInspection
type Usecase struct {
	repo           IAuthorizationRepo
	passwordHasher password.IPasswordHash
	tokenMaker     token.IMaker
}

// New создаёт объект Usecase.
func New(repo IAuthorizationRepo, hasher password.IPasswordHash, maker token.IMaker) (*Usecase, error) {
	return &Usecase{
		repo:           repo,
		passwordHasher: hasher,
		tokenMaker:     maker,
	}, nil
}

// RegisterUser - регистрация нового пользователя с переданным логином и паролем.
//
// Возвращает токен или ошибку (например, если логин уже существует).
func (uc *Usecase) RegisterUser(login, pass string) (string, error) {
	passwordHash, err := uc.passwordHasher.Hash(pass)
	if err != nil {
		return "", err
	}

	id, err := uc.repo.CreateUser(login, passwordHash)
	if err != nil {
		return "", err
	}

	// fixme: change magic time to parameter from config
	// todo: вынести в отдельный метод? (повтор в LoginUser)
	userToken, err := uc.tokenMaker.Create(id, 12*time.Hour)
	if err != nil {
		return "", err
	}

	return userToken, nil
}

// LoginUser - авторизация существующего пользователя.
//
// Возвращает токен или ошибку (например, если логина не существует).
func (uc *Usecase) LoginUser(login, pass string) (string, error) {
	user, err := uc.repo.GetUser(login)
	if err != nil {
		return "", ErrLoginNotExist
	}

	err = uc.passwordHasher.Check(pass, user.PasswordHash)
	if err != nil {
		return "", ErrMismatchPassword
	}

	userToken, err := uc.tokenMaker.Create(user.ID, 12*time.Hour)
	if err != nil {
		return "", err
	}

	return userToken, nil
}
