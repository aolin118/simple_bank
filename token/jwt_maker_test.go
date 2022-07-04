package token

import (
	"testing"
	"time"

	"github.com/aolin118/simple_bank/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	jwtMaker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, jwtMaker)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := jwtMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	payload, err := jwtMaker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, payload.Username, username)
	require.WithinDuration(t, payload.IssuedAt, issuedAt, time.Second)
	require.WithinDuration(t, payload.ExpiredAt, expiredAt, time.Second)

}

func TestExpiredToken(t *testing.T) {
	jwtMaker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, jwtMaker)

	username := util.RandomOwner()
	duration := time.Minute

	token, err := jwtMaker.CreateToken(username, -duration)
	require.NotEmpty(t, token)
	require.NoError(t, err)
	payload, err := jwtMaker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidToken(t *testing.T) {
	jwtMaker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, jwtMaker)

	payload, err := NewPayload(util.RandomOwner(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	payload, err = jwtMaker.VerifyToken(token)
	require.Nil(t, payload)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
}
