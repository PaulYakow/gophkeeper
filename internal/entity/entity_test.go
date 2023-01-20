package entity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gophkeeper/internal/entity"
)

const (
	id       = 1
	login    = "user"
	pass     = "pass"
	passHash = "|0@$$)(@W"
)

func TestUserDTO(t *testing.T) {
	got := entity.UserDTO{Login: login, Password: pass}
	assert.NotNil(t, got)
	assert.Equal(t, got.Login, login)
	assert.Equal(t, got.Password, pass)
}

func TestUserDAO(t *testing.T) {
	got := entity.UserDAO{ID: id, Login: login, PasswordHash: passHash}
	assert.NotNil(t, got)
	assert.Equal(t, got.ID, id)
	assert.Equal(t, got.Login, login)
	assert.Equal(t, got.PasswordHash, passHash)
}
