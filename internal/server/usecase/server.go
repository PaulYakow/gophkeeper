package usecase

import (
	"fmt"
	"time"

	"gophkeeper/internal/utils/password"
	"gophkeeper/internal/utils/token"
)

type Server struct {
	repo       IAuthorizationRepo
	tokenMaker token.IMaker
}

func New(repo IAuthorizationRepo) (*Server, error) {
	//todo: move key to config
	tokenMaker, err := token.NewPasetoMaker("_0987654321zyxwvutsrq1234567890_")
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	return &Server{
		repo:       repo,
		tokenMaker: tokenMaker,
	}, nil
}

func (s *Server) RegisterUser(login string, pass string) (string, error) {
	passwordHash, err := password.Hash(pass)
	if err != nil {
		return "", err
	}

	id, err := s.repo.CreateUser(login, passwordHash)
	if err != nil {
		return "", err
	}

	userToken, err := s.tokenMaker.CreateToken(id, 12*time.Hour)
	if err != nil {
		return "", err
	}

	return userToken, nil
}
