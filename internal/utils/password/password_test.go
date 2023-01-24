package password_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"gophkeeper/internal/utils/password"
	"gophkeeper/internal/utils/test"
)

func TestPassword(t *testing.T) {
	pass := test.RandomPassword()
	hashedPassword1, err := password.Hash(pass)

	t.Run("correct hashing", func(t *testing.T) {
		require.NoError(t, err)
		require.NotEmpty(t, hashedPassword1)
	})

	t.Run("correct check", func(t *testing.T) {
		err = password.Check(pass, hashedPassword1)
		require.NoError(t, err)
	})

	t.Run("wrong password check", func(t *testing.T) {
		wrongPassword := test.RandomPassword()
		err = password.Check(wrongPassword, hashedPassword1)
		require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
	})

	t.Run("check that hashes of right and wrong password not equal", func(t *testing.T) {
		hashedPassword2, err := password.Hash(pass)
		require.NoError(t, err)
		require.NotEmpty(t, hashedPassword1)
		require.NotEqual(t, hashedPassword1, hashedPassword2)
	})
}
