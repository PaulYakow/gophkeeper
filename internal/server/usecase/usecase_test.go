package usecase_test

import (
	"context"
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
	repo   *mocks.MockIRepo
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

	serverMock.repo = mocks.NewMockIRepo(serverMock.ctrl)
	serverMock.hasher = mocks.NewMockIPasswordHash(serverMock.ctrl)
	serverMock.maker = mocks.NewMockIMaker(serverMock.ctrl)

	var err error

	auth := usecase.NewAuthService(serverMock.repo, serverMock.hasher, serverMock.maker)
	pairs := usecase.NewPairsService(serverMock.repo)
	cards := usecase.NewBankService(serverMock.repo)
	notes := usecase.NewTextService(serverMock.repo)

	serverMock.uc, err = usecase.New(auth, pairs, cards, notes)

	t.Run("proper usecase create", func(t *testing.T) {
		require.NoError(t, err)
		require.NotEmpty(t, serverMock.uc)
	})
}

const (
	login, password = "user", "password"
	hashedPassword  = "hashed_password"
)

func TestAuthorization_RegisterUser(t *testing.T) {
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

func TestAuthorization_LoginUser(t *testing.T) {
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

func TestPairs_GetAll(t *testing.T) {
	userID := 1
	testPairs := []entity.PairDAO{
		{
			UserID:   userID,
			Login:    "pairTest1",
			Password: "pairPass",
			Metadata: `
tag #1: test-1_1;
tag #2: test-2_2;`,
		},
		{
			UserID:   userID,
			Login:    "pairTest2",
			Password: "pairPass",
			Metadata: `
tag #1: test-2_1;
tag #2: test-2_2;`,
		},
	}

	t.Run("get pairs from exist user", func(t *testing.T) {
		serverMock.repo.EXPECT().GetAllPairs(context.Background(), userID).Return(testPairs, nil)
		pairs, err := serverMock.uc.ViewAllPairs(userID)
		require.NoError(t, err)
		require.NotEmpty(t, pairs)
		require.IsType(t, []entity.PairDTO{}, pairs)
	})

	t.Run("get pairs from not exist user", func(t *testing.T) {
		serverMock.repo.EXPECT().GetAllPairs(context.Background(), userID).Return(nil, errors.New("user_id not exist"))
		pairs, err := serverMock.uc.ViewAllPairs(userID)
		require.Error(t, err)
		require.Empty(t, pairs)
	})
}

func TestBank_GetAll(t *testing.T) {
	userID := 1
	testCards := []entity.BankDAO{
		{
			CardHolder:     "Ivanov Ivan",
			Number:         "1234 0987 5678 6543",
			ExpirationDate: "09/99",
			Metadata: `
tag #1: test bank 1-1;
tag #2: test bank 1-2;`,
		},
		{
			CardHolder:     "Petrov Petr",
			Number:         "0987 1234 6543 5678",
			ExpirationDate: "01/11",
			Metadata: `
tag #1: test bank 2-1;
tag #2: test bank 2-2;`,
		},
	}

	t.Run("get cards from exist user", func(t *testing.T) {
		serverMock.repo.EXPECT().GetAllCards(context.Background(), userID).Return(testCards, nil)
		cards, err := serverMock.uc.ViewAllCards(userID)
		require.NoError(t, err)
		require.NotEmpty(t, cards)
		require.IsType(t, []entity.BankDTO{}, cards)
	})

	t.Run("get cards from not exist user", func(t *testing.T) {
		serverMock.repo.EXPECT().GetAllCards(context.Background(), userID).Return(nil, errors.New("user_id not exist"))
		cards, err := serverMock.uc.ViewAllCards(userID)
		require.Error(t, err)
		require.Empty(t, cards)
	})
}
