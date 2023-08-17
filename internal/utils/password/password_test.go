package password_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/PaulYakow/gophkeeper/internal/utils/password"
	"github.com/PaulYakow/gophkeeper/internal/utils/test"
)

func TestPassword(t *testing.T) {
	pass := password.New()
	t.Run("object pass created", func(t *testing.T) {
		require.NotNil(t, pass)
	})

	rndPassword := test.RandomPassword()
	hashedPassword1, err := pass.Hash(rndPassword)

	t.Run("correct hashing", func(t *testing.T) {
		require.NoError(t, err)
		require.NotEmpty(t, hashedPassword1)
	})

	t.Run("correct check", func(t *testing.T) {
		err = pass.Check(rndPassword, hashedPassword1)
		require.NoError(t, err)
	})

	t.Run("wrong password check", func(t *testing.T) {
		wrongPassword := test.RandomPassword()
		err = pass.Check(wrongPassword, hashedPassword1)
		require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
	})

	t.Run("check that hashes of right and wrong password not equal", func(t *testing.T) {
		hashedPassword2, err := pass.Hash(rndPassword)
		require.NoError(t, err)
		require.NotEmpty(t, hashedPassword1)
		require.NotEqual(t, hashedPassword1, hashedPassword2)
	})
}
