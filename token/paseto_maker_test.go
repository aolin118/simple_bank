package token

import (
	"testing"
	"time"

	"github.com/aolin118/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	pasetoMaker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, pasetoMaker)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := pasetoMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	payload, err := pasetoMaker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, payload.Username, username)
	require.WithinDuration(t, payload.IssuedAt, issuedAt, time.Second)
	require.WithinDuration(t, payload.ExpiredAt, expiredAt, time.Second)

}

func TestExpiredPasetoToken(t *testing.T) {
	pasetoMaker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, pasetoMaker)

	username := util.RandomOwner()
	duration := time.Minute

	token, err := pasetoMaker.CreateToken(username, -duration)
	require.NotEmpty(t, token)
	require.NoError(t, err)
	payload, err := pasetoMaker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidPasetoToken(t *testing.T) {
	pasetoMaker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, pasetoMaker)

	username := util.RandomOwner()
	duration := time.Minute

	token, err := pasetoMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	token = token + " "
	_, err = pasetoMaker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
}
