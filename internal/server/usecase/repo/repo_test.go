package repo_test

import (
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
	dropUsers = `
DROP TABLE IF EXISTS users;
`
)

var (
	testDB   *postgres.Postgres
	testRepo *repo.Repo
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
	_, err := testDB.Exec(dropUsers)
	if err != nil {
		log.Println(fmt.Errorf("fail drop repo: %w", err))
	}

	log.Println("clean repo")
}

func TestMain(m *testing.M) {
	setup()
	defer testDB.Close()

	var err error
	testRepo, err = repo.New(testDB)
	if err != nil {
		log.Println(fmt.Errorf("repo tests - repo.New: %w", err))
	}
	code := m.Run()
	teardown()

	os.Exit(code)
}

var user = entity.UserDTO{
	Login:    "user",
	Password: "pass_hash",
}

func TestCreateUser(t *testing.T) {
	t.Run("create new user", func(t *testing.T) {
		userID, err := testRepo.CreateUser(user.Login, user.Password)
		require.NoError(t, err)
		assert.NotEmpty(t, userID)
	})

	t.Run("duplicate user", func(t *testing.T) {
		userID, err := testRepo.CreateUser(user.Login, user.Password)
		require.Error(t, err)
		assert.Empty(t, userID)
	})
}

func TestGetUser(t *testing.T) {
	t.Run("get exist user", func(t *testing.T) {
		userDAO, err := testRepo.GetUser(user.Login)
		require.NoError(t, err)
		require.IsType(t, entity.UserDAO{}, userDAO)
		assert.Equal(t, user.Login, userDAO.Login)
		assert.Equal(t, user.Password, userDAO.PasswordHash)
	})

	t.Run("get not exist user", func(t *testing.T) {
		notExistUser := entity.UserDTO{
			Login:    "userNotExist",
			Password: "no_pass",
		}
		userDAO, err := testRepo.GetUser(notExistUser.Login)
		require.Error(t, err)
		require.Empty(t, userDAO)
	})
}
