package usecase_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"gophkeeper/internal/server/usecase"
	"gophkeeper/internal/server/usecase/repo"
)

func TestRegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	login, password := "user", "password"
	authRepo := NewMockIAuthorizationRepo(ctrl)
	uc, err := usecase.New(authRepo)

	t.Run("usecase created without errors", func(t *testing.T) {
		require.NoError(t, err)
	})

	t.Run("create new user", func(t *testing.T) {
		authRepo.EXPECT().CreateUser(login, password).Return(1, nil)
		token, err := uc.RegisterUser(login, password)
		require.NoError(t, err)
		require.NotEmpty(t, token)
	})

	t.Run("duplicate user", func(t *testing.T) {
		authRepo.EXPECT().CreateUser(login, password).Return(0, repo.ErrUserNotExist)
		token, err := uc.RegisterUser(login, password)
		require.ErrorIs(t, err, repo.ErrUserNotExist)
		require.Empty(t, token)
	})

}
