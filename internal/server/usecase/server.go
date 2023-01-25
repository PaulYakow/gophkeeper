package usecase

import (
	"fmt"
	"time"

	"github.com/PaulYakow/gophkeeper/internal/utils/password"
	"github.com/PaulYakow/gophkeeper/internal/utils/token"
)

type Server struct {
	repo       IAuthorizationRepo
	hasher     password.IPasswordHash
	tokenMaker token.IMaker
}

func New(repo IAuthorizationRepo, hasher password.IPasswordHash) (*Server, error) {
	//todo: move key to config
	//todo: move object to input args
	tokenMaker, err := token.NewPasetoMaker("_0987654321zyxwvutsrq1234567890_")
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	return &Server{
		repo:       repo,
		hasher:     hasher,
		tokenMaker: tokenMaker,
	}, nil
}

func (s *Server) RegisterUser(login string, pass string) (string, error) {
	passwordHash, err := s.hasher.Hash(pass)
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
