package usecase_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/PaulYakow/gophkeeper/internal/server/mocks"
	"github.com/PaulYakow/gophkeeper/internal/server/usecase"
	"github.com/PaulYakow/gophkeeper/internal/server/usecase/repo"
)

func TestRegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	login, password := "user", "password"
	hashedPassword := "hashed_password"

	authRepo := mocks.NewMockIAuthorizationRepo(ctrl)
	hasher := mocks.NewMockIPasswordHash(ctrl)

	uc, err := usecase.New(authRepo, hasher)

	t.Run("usecase created without errors & not nil", func(t *testing.T) {
		require.NoError(t, err)
		require.NotEmpty(t, uc)
	})

	t.Run("create new user", func(t *testing.T) {
		hasher.EXPECT().Hash(password).Return(hashedPassword, nil)
		authRepo.EXPECT().CreateUser(login, hashedPassword).Return(1, nil)
		token, err := uc.RegisterUser(login, password)
		require.NoError(t, err)
		require.NotEmpty(t, token)
	})

	t.Run("duplicate user", func(t *testing.T) {
		hasher.EXPECT().Hash(password).Return(hashedPassword, nil)
		authRepo.EXPECT().CreateUser(login, hashedPassword).Return(0, repo.ErrUserExist)
		token, err := uc.RegisterUser(login, password)
		require.ErrorIs(t, err, repo.ErrUserExist)
		require.Empty(t, token)
	})

	t.Run("invalid password hash", func(t *testing.T) {
		hasher.EXPECT().Hash(password).Return("", errors.New("invalid password hash"))
		_, err := uc.RegisterUser(login, password)
		require.Error(t, err)
	})
}
