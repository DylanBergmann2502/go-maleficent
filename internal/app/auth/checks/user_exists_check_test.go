// internal/app/auth/checks/user_exists_check_test.go
package checks

import (
	"testing"

	"github.com/DylanBergmann2502/go-maleficent/internal/app/auth/models"
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnsureEmailAvailableFindsExistingUser(t *testing.T) {
	db := testutil.NewDatabase(t)
	tx := testutil.Begin(t, db)

	require.NoError(t, tx.Create(&models.User{
		Email:        "registered@example.com",
		PasswordHash: "argon2-hash",
	}).Error)

	err := EnsureEmailAvailable(tx, "registered@example.com")

	assert.ErrorIs(t, err, ErrEmailAlreadyRegistered)
}

func TestEnsureEmailAvailableAcceptsUnusedEmail(t *testing.T) {
	db := testutil.NewDatabase(t)
	tx := testutil.Begin(t, db)

	assert.NoError(t, EnsureEmailAvailable(tx, "available@example.com"))
}
