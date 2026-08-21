// internal/app/auth/utils/password_test.go
package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPasswordHashRoundTrip(t *testing.T) {
	password := "correct horse battery staple"
	hash, err := HashPassword(password)
	require.NoError(t, err)

	valid, err := CheckPasswordHash(password, hash)
	require.NoError(t, err)
	assert.True(t, valid)

	valid, err = CheckPasswordHash("wrong password", hash)
	require.NoError(t, err)
	assert.False(t, valid)
}

func TestCheckPasswordHashRejectsMalformedHash(t *testing.T) {
	valid, err := CheckPasswordHash("password", "not-a-password-hash")

	assert.False(t, valid)
	assert.Error(t, err)
}
