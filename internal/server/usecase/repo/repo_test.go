package repo_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/PaulYakow/gophkeeper/internal/entity"
	"github.com/PaulYakow/gophkeeper/internal/server/usecase/repo"
	"github.com/PaulYakow/gophkeeper/pkg/postgres"
)

const (
	qDropTables = `
DROP TABLE IF EXISTS public.pairs;
DROP TABLE IF EXISTS public.users;
`
	qCreateUser = `
INSERT INTO public.users (login, password_hash)
VALUES ($1, $2)
RETURNING id;
`
	qCreatePair = `
INSERT INTO public.pairs (user_id, login, password, metadata)
VALUES ($1, $2, $3, $4)
RETURNING id;
`
)

var (
	testDB   *postgres.Postgres
	testRepo *repo.Repo

	userDTO = entity.UserDTO{
		Login:    "userDTO",
		Password: "pass_hash",
	}

	userDAO = entity.UserDAO{
		Login:        "userDAO",
		PasswordHash: "p@$$_|-|@$|-|",
	}

	testPairs = []entity.PairDAO{
		{
			Login:    "pairTest1",
			Password: "pairPass",
			Metadata: `
tag #1: test-1_1;
tag #2: test-2_2;`,
		},
		{
			Login:    "pairTest2",
			Password: "pairPass",
			Metadata: `
tag #1: test-2_1;
tag #2: test-2_2;`,
		},
	}
)

// Подготовка тестовой БД
func setup() {
	var err error
	testDB, err = postgres.New("postgresql://test:test@localhost:54321/postgres", postgres.ConnAttempts(1))
	if err != nil {
		log.Println(fmt.Errorf("skip repo tests: %w", err))
		os.Exit(0)
	}
	_ = testDB

	log.Println("prepare test database")
}

// Сброс тестовой БД
func teardown() {
	_, err := testDB.Exec(qDropTables)
	if err != nil {
		log.Println(fmt.Errorf("fail drop repo: %w", err))
	}

	log.Println("clean repo")
}

func TestMain(m *testing.M) {
	setup()
	defer testRepo.CloseConnection()

	var err error
	var code int

	auth := repo.NewAuthPostgres(testDB)
	pairs := repo.NewPairPostgres(testDB)

	testRepo, err = repo.New(testDB, auth, pairs)
	if err != nil {
		log.Println(fmt.Errorf("repo tests - repo.New: %w", err))
	}

	err = testDB.Get(&userDAO.ID, qCreateUser, userDTO.Login, userDTO.Password)
	if err != nil {
		log.Println(fmt.Errorf("repo tests - fail create user: %w", err))
		code = 1
	}

	for i := range testPairs {
		err = testDB.Get(&testPairs[i].ID,
			qCreatePair,
			userDAO.ID, testPairs[i].Login, testPairs[i].Password, testPairs[i].Metadata)
		if err != nil {
			log.Println(fmt.Errorf("repo tests - fail create pair: %w", err))
			code = 1
		}
	}

	if err == nil {
		code = m.Run()
	}
	teardown()

	os.Exit(code)
}

func TestAuthorization_CreateUser(t *testing.T) {
	t.Run("create new user", func(t *testing.T) {
		userID, err := testRepo.CreateUser("new_user", userDTO.Password)
		require.NoError(t, err)
		assert.NotEmpty(t, userID)
		assert.Greater(t, userID, 0)
	})

	t.Run("duplicate user", func(t *testing.T) {
		userID, err := testRepo.CreateUser(userDTO.Login, userDTO.Password)
		require.Error(t, err)
		assert.Empty(t, userID)
	})
}

func TestAuthorization_GetUser(t *testing.T) {
	t.Run("get exist user", func(t *testing.T) {
		user, err := testRepo.GetUser(userDTO.Login)
		require.NoError(t, err)
		require.IsType(t, entity.UserDAO{}, user)
		assert.Equal(t, userDTO.Login, user.Login)
		assert.Equal(t, userDTO.Password, user.PasswordHash)
	})

	t.Run("get not exist user", func(t *testing.T) {
		notExistUser := entity.UserDTO{
			Login:    "userNotExist",
			Password: "no_pass",
		}
		user, err := testRepo.GetUser(notExistUser.Login)
		require.Error(t, err)
		require.Empty(t, user)
	})
}

func TestPairs_GetAll(t *testing.T) {
	t.Run("get exist pairs", func(t *testing.T) {
		pairs, err := testRepo.GetAllPairs(context.Background(), userDAO.ID)
		require.NoError(t, err)
		require.IsType(t, []entity.PairDAO{}, pairs)
		require.NotEmpty(t, pairs)
		for i := range pairs {
			require.Equal(t, testPairs[i].ID, pairs[i].ID)
			require.Equal(t, testPairs[i].Login, pairs[i].Login)
			require.Equal(t, testPairs[i].Password, pairs[i].Password)
			require.Equal(t, testPairs[i].Metadata, pairs[i].Metadata)
		}
	})

	t.Run("get not exist pairs (user_id not exist)", func(t *testing.T) {
		pairs, err := testRepo.GetAllPairs(context.Background(), 777)
		require.NoError(t, err)
		require.Empty(t, pairs)
	})
}
