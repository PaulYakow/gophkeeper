package usecase_test

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/PaulYakow/gophkeeper/internal/entity"
	"github.com/PaulYakow/gophkeeper/internal/server/mocks"
	"github.com/PaulYakow/gophkeeper/internal/server/usecase"
	"github.com/PaulYakow/gophkeeper/internal/server/usecase/repo"
)

var serverMock = struct {
	ctrl   *gomock.Controller
	uc     *usecase.Usecase
	repo   *mocks.MockIAuthorizationRepo
	hasher *mocks.MockIPasswordHash
	maker  *mocks.MockIMaker
}{}

func setup() {
}

func teardown() {
	defer serverMock.ctrl.Finish()
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()

	os.Exit(code)
}

func TestNew(t *testing.T) {
	serverMock.ctrl = gomock.NewController(t)

	serverMock.repo = mocks.NewMockIAuthorizationRepo(serverMock.ctrl)
	serverMock.hasher = mocks.NewMockIPasswordHash(serverMock.ctrl)
	serverMock.maker = mocks.NewMockIMaker(serverMock.ctrl)

	var err error
	serverMock.uc, err = usecase.New(serverMock.repo, serverMock.hasher, serverMock.maker)

	t.Run("proper usecase create", func(t *testing.T) {
		require.NoError(t, err)
		require.NotEmpty(t, serverMock.uc)
	})
}

const (
	login, password = "user", "password"
	hashedPassword  = "hashed_password"
)

func TestRegisterUser(t *testing.T) {
	t.Run("proper create new user", func(t *testing.T) {
		userID := 1
		serverMock.hasher.EXPECT().Hash(password).Return(hashedPassword, nil)
		serverMock.repo.EXPECT().CreateUser(login, hashedPassword).Return(userID, nil)
		serverMock.maker.EXPECT().Create(userID, 12*time.Hour).Return("token", nil)
		token, err := serverMock.uc.RegisterUser(login, password)
		require.NoError(t, err)
		require.NotEmpty(t, token)
		assert.Equal(t, "token", token)
	})

	t.Run("duplicate user", func(t *testing.T) {
		serverMock.hasher.EXPECT().Hash(password).Return(hashedPassword, nil)
		serverMock.repo.EXPECT().CreateUser(login, hashedPassword).Return(0, repo.ErrUserExist)
		token, err := serverMock.uc.RegisterUser(login, password)
		require.ErrorIs(t, err, repo.ErrUserExist)
		require.Empty(t, token)
	})

	t.Run("invalid password hash", func(t *testing.T) {
		serverMock.hasher.EXPECT().Hash(password).Return("", errors.New("invalid password hash"))
		_, err := serverMock.uc.RegisterUser(login, password)
		require.Error(t, err)
	})

	t.Run("error from Create", func(t *testing.T) {
		userID := 1
		serverMock.hasher.EXPECT().Hash(password).Return(hashedPassword, nil)
		serverMock.repo.EXPECT().CreateUser(login, hashedPassword).Return(userID, nil)
		serverMock.maker.EXPECT().Create(userID, 12*time.Hour).Return("", errors.New("token create error"))
		_, err := serverMock.uc.RegisterUser(login, password)
		require.Error(t, err)
	})
}

func TestLoginUser(t *testing.T) {
	user := entity.UserDAO{
		ID:           1,
		Login:        login,
		PasswordHash: hashedPassword,
	}

	t.Run("proper login new user", func(t *testing.T) {
		testToken := "token"

		serverMock.repo.EXPECT().GetUser(login).Return(user, nil)
		serverMock.hasher.EXPECT().Check(password, hashedPassword).Return(nil)
		serverMock.maker.EXPECT().Create(user.ID, 12*time.Hour).Return(testToken, nil)
		token, err := serverMock.uc.LoginUser(login, password)
		require.NoError(t, err)
		require.Equal(t, testToken, token)
	})

	t.Run("login not exist", func(t *testing.T) {
		serverMock.repo.EXPECT().GetUser(login).Return(entity.UserDAO{}, errors.New("login not exist"))
		token, err := serverMock.uc.LoginUser(login, password)
		require.Error(t, err)
		require.Empty(t, token)
	})

	t.Run("invalid hash check", func(t *testing.T) {
		serverMock.repo.EXPECT().GetUser(login).Return(user, nil)
		serverMock.hasher.EXPECT().Check(password, hashedPassword).Return(errors.New("invalid hash"))
		token, err := serverMock.uc.LoginUser(login, password)
		require.Error(t, err)
		require.Empty(t, token)
	})

	t.Run("invalid token create", func(t *testing.T) {
		serverMock.repo.EXPECT().GetUser(login).Return(user, nil)
		serverMock.hasher.EXPECT().Check(password, hashedPassword).Return(nil)
		serverMock.maker.EXPECT().Create(user.ID, 12*time.Hour).Return("", errors.New("invalid token"))
		token, err := serverMock.uc.LoginUser(login, password)
		require.Error(t, err)
		require.Empty(t, token)
	})
}
