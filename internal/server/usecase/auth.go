package usecase

import (
	"time"

	"github.com/PaulYakow/gophkeeper/internal/utils/password"
	"github.com/PaulYakow/gophkeeper/internal/utils/token"
)

// AuthService сервис аутентификации пользователей.
type AuthService struct {
	repo           IAuthorizationRepo
	passwordHasher password.IPasswordHash
	tokenMaker     token.IMaker
}

// NewAuthService создаёт объект типа AuthService.
func NewAuthService(repo IAuthorizationRepo, hasher password.IPasswordHash, maker token.IMaker) *AuthService {
	return &AuthService{
		repo:           repo,
		passwordHasher: hasher,
		tokenMaker:     maker,
	}
}

// RegisterUser - регистрация нового пользователя с переданным логином и паролем.
//
// Возвращает токен или ошибку (например, если логин уже существует).
func (s *AuthService) RegisterUser(login, pass string) (string, error) {
	passwordHash, err := s.passwordHasher.Hash(pass)
	if err != nil {
		return "", err
	}

	id, err := s.repo.CreateUser(login, passwordHash)
	if err != nil {
		return "", err
	}

	userToken, err := s.generateToken(id)
	if err != nil {
		return "", err
	}

	return userToken, nil
}

// LoginUser - авторизация существующего пользователя.
//
// Возвращает токен или ошибку (например, если логина не существует).
func (s *AuthService) LoginUser(login, pass string) (string, error) {
	user, err := s.repo.GetUser(login)
	if err != nil {
		return "", ErrLoginNotExist
	}

	err = s.passwordHasher.Check(pass, user.PasswordHash)
	if err != nil {
		return "", ErrMismatchPassword
	}

	userToken, err := s.generateToken(user.ID)
	if err != nil {
		return "", err
	}

	return userToken, nil
}

// ParseToken - проверяет переданный токен.
//
// Возвращает id пользователя или ошибку.
func (s *AuthService) ParseToken(token string) (int, error) {
	payload, err := s.tokenMaker.Verify(token)
	if err != nil {
		return 0, err
	}

	return payload.UserID, nil
}

func (s *AuthService) generateToken(userID int) (string, error) {
	// fixme: change magic time to parameter from config
	authToken, err := s.tokenMaker.Create(userID, 12*time.Hour)
	if err != nil {
		return "", err
	}

	return authToken, nil
}
