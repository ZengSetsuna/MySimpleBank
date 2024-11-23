package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)
	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	err = CheckPassword(hashedPassword, password)
	require.NoError(t, err)
	err = CheckPassword(hashedPassword, RandomString(6))
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	password2 := RandomString(6)
	hashedPassword2, err := HashPassword(password2)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)
	require.NotEqual(t, hashedPassword, hashedPassword2)
}
