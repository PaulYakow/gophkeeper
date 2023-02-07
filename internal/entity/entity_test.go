package entity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/PaulYakow/gophkeeper/internal/entity"
)

var testUser = struct {
	id       int
	login    string
	pass     string
	passHash string
}{
	id:       1,
	login:    "user",
	pass:     "pass",
	passHash: "|0@$$)(@W",
}

func TestUserDTO(t *testing.T) {
	userDTO := entity.UserDTO{
		Login:    testUser.login,
		Password: testUser.pass,
	}

	require.NotNil(t, userDTO)
	assert.Equal(t, userDTO.Login, testUser.login)
	assert.Equal(t, userDTO.Password, testUser.pass)
}

func TestUserDAO(t *testing.T) {
	userDAO := entity.UserDAO{
		ID:           testUser.id,
		Login:        testUser.login,
		PasswordHash: testUser.passHash,
	}

	require.NotNil(t, userDAO)
	assert.Equal(t, userDAO.ID, testUser.id)
	assert.Equal(t, userDAO.Login, testUser.login)
	assert.Equal(t, userDAO.PasswordHash, testUser.passHash)
}

var testPair = struct {
	id     int
	userID int
	login  string
	pass   string
	meta   string
}{
	id:     10,
	userID: testUser.id,
	login:  "pairTest",
	pass:   "pairPass",
	meta: `
tag #1: test-1;
tag #2: test-2;
`,
}

func TestPairDTO(t *testing.T) {
	pairDTO := entity.PairDTO{
		ID:       testPair.id,
		Login:    testPair.login,
		Password: testPair.pass,
		Metadata: testPair.meta,
	}

	require.NotNil(t, pairDTO)
	assert.Equal(t, pairDTO.ID, testPair.id)
	assert.Equal(t, pairDTO.Login, testPair.login)
	assert.Equal(t, pairDTO.Password, testPair.pass)
	assert.Equal(t, pairDTO.Metadata, testPair.meta)
}

func TestPairDAO(t *testing.T) {
	pairDAO := entity.PairDAO{
		ID:       testPair.id,
		UserID:   testPair.userID,
		Login:    testPair.login,
		Password: testPair.pass,
		Metadata: testPair.meta,
	}

	require.NotNil(t, pairDAO)
	assert.Equal(t, pairDAO.ID, testPair.id)
	assert.Equal(t, pairDAO.UserID, testPair.userID)
	assert.Equal(t, pairDAO.Login, testPair.login)
	assert.Equal(t, pairDAO.Password, testPair.pass)
	assert.Equal(t, pairDAO.Metadata, testPair.meta)
}
