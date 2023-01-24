package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"gophkeeper/internal/utils/test"
)

func TestPasetoMaker(t *testing.T) {
	userID := test.RandomInt(1, 100)
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := time.Now().Add(duration)

	var maker IMaker
	var err error
	t.Run("create token maker", func(t *testing.T) {
		maker, err = NewPasetoMaker(test.RandomString(32))
		require.NoError(t, err)
	})

	var token string
	t.Run("create token", func(t *testing.T) {
		token, err = maker.CreateToken(userID, duration)
		require.NoError(t, err)
		require.NotEmpty(t, token)
	})

	t.Run("verify token", func(t *testing.T) {
		payload, err := maker.VerifyToken(token)
		require.NoError(t, err)
		require.NotEmpty(t, payload)

		require.Equal(t, userID, payload.UserID)
		require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
		require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
	})
}

func TestExpiredPasetoToken(t *testing.T) {
	var err error

	maker, err := NewPasetoMaker(test.RandomString(32))
	require.NoError(t, err)

	var token string
	t.Run("create expired token", func(t *testing.T) {
		token, err = maker.CreateToken(test.RandomInt(1, 100), -time.Minute)
		require.NoError(t, err)
		require.NotEmpty(t, token)
	})

	t.Run("verify expired token", func(t *testing.T) {
		payload, err := maker.VerifyToken(token)
		require.Error(t, err)
		require.EqualError(t, err, ErrExpiredToken.Error())
		require.Nil(t, payload)
	})
}
