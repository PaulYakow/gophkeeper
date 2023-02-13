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
DROP TABLE IF EXISTS resources.pairs_data;
DROP TABLE IF EXISTS resources.bank_data;
DROP TABLE IF EXISTS resources.text_data;
DROP TABLE IF EXISTS public.users;
`
	qCreateUser = `
INSERT INTO public.users (login, password_hash)
VALUES ($1, $2)
RETURNING id;
`
	qCreatePair = `
INSERT INTO resources.pairs_data (user_id, login, password, metadata)
VALUES ($1, $2, $3, $4)
RETURNING id;
`

	qCreateCard = `
INSERT INTO resources.bank_data (user_id, card_holder, number, expiration_date, metadata)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;
`

	qCreateNote = `
INSERT INTO resources.text_data (user_id, note, metadata)
VALUES ($1, $2, $3)
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

	testCards = []entity.BankDAO{
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

	testNotes = []entity.TextDAO{
		{
			Note: "some text from user",
			Metadata: `tag #1: test text 1-1;
tag #2: test text 1-2;`,
		},
		{
			Note: "another text from user",
			Metadata: `tag #1: test text 1-1;
tag #2: test text 1-2;`,
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
		return
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
	cards := repo.NewBankPostgres(testDB)
	notes := repo.NewTextPostgres(testDB)

	testRepo, err = repo.New(testDB, auth, pairs, cards, notes)
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

	for i := range testCards {
		err = testDB.Get(&testCards[i].ID,
			qCreateCard,
			userDAO.ID, testCards[i].CardHolder, testCards[i].Number, testCards[i].ExpirationDate, testCards[i].Metadata)
		if err != nil {
			log.Println(fmt.Errorf("repo tests - fail create card: %w", err))
			code = 1
		}
	}

	for i := range testNotes {
		err = testDB.Get(&testNotes[i].ID,
			qCreateNote,
			userDAO.ID, testNotes[i].Note, testNotes[i].Metadata)
		if err != nil {
			log.Println(fmt.Errorf("repo tests - fail create note: %w", err))
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

func TestCards_GetAll(t *testing.T) {
	t.Run("get exist cards", func(t *testing.T) {
		cards, err := testRepo.GetAllCards(context.Background(), userDAO.ID)
		require.NoError(t, err)
		require.IsType(t, []entity.BankDAO{}, cards)
		require.NotEmpty(t, cards)
		for i := range cards {
			require.Equal(t, testCards[i].ID, cards[i].ID)
			require.Equal(t, testCards[i].CardHolder, cards[i].CardHolder)
			require.Equal(t, testCards[i].Number, cards[i].Number)
			require.Equal(t, testCards[i].ExpirationDate, cards[i].ExpirationDate)
			require.Equal(t, testCards[i].Metadata, cards[i].Metadata)
		}
	})

	t.Run("get not exist cards (user_id not exist)", func(t *testing.T) {
		cards, err := testRepo.GetAllCards(context.Background(), 777)
		require.NoError(t, err)
		require.Empty(t, cards)
	})
}

func TestNotes_GetAll(t *testing.T) {
	t.Run("get exist notes", func(t *testing.T) {
		notes, err := testRepo.GetAllNotes(context.Background(), userDAO.ID)
		require.NoError(t, err)
		require.IsType(t, []entity.TextDAO{}, notes)
		require.NotEmpty(t, notes)
		for i := range notes {
			require.Equal(t, testNotes[i].ID, notes[i].ID)
			require.Equal(t, testNotes[i].Note, notes[i].Note)
			require.Equal(t, testNotes[i].Metadata, notes[i].Metadata)
		}
	})

	t.Run("get not exist notes (user_id not exist)", func(t *testing.T) {
		notes, err := testRepo.GetAllCards(context.Background(), 777)
		require.NoError(t, err)
		require.Empty(t, notes)
	})
}
