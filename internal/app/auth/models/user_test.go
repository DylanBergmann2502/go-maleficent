// internal/app/auth/models/user_test.go
package models

import (
	stderrors "errors"
	"testing"

	apperrors "github.com/DylanBergmann2502/go-maleficent/internal/pkg/errors"
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestUserBeforeSaveValidatesModelInvariants(t *testing.T) {
	user := &User{
		Email:        "invalid-email",
		PasswordHash: "",
	}

	err := user.BeforeSave(nil)

	var appErr *apperrors.ApplicationError
	require.ErrorAs(t, err, &appErr)
	assert.True(t, stderrors.Is(err, apperrors.ErrValidationCategory))
	assert.Equal(t, "model_validation_failed", appErr.Code)
	assert.Contains(t, appErr.Message, "user model validation failed")
}

func TestUserBeforeSaveAcceptsValidModel(t *testing.T) {
	user := &User{
		Email:        "user@example.com",
		PasswordHash: "argon2-hash",
	}

	require.NoError(t, user.BeforeSave(nil))
}

func TestUserPersistsThroughGORM(t *testing.T) {
	db := testutil.NewDatabase(t)
	tx := testutil.Begin(t, db)

	user := &User{
		Email:        "gorm-user@example.com",
		PasswordHash: "argon2-hash",
	}

	require.NoError(t, tx.Create(user).Error)
	assert.NotEqual(t, uuid.Nil, user.ID)
	assert.False(t, user.CreatedAt.IsZero())

	var persisted User
	require.NoError(t, tx.First(&persisted, "id = ?", user.ID).Error)
	assert.Equal(t, user.Email, persisted.Email)
}

func TestUserRejectsInvalidDataThroughGORM(t *testing.T) {
	db := testutil.NewDatabase(t)
	tx := testutil.Begin(t, db)

	err := tx.Create(&User{
		Email:        "not-an-email",
		PasswordHash: "",
	}).Error

	var appErr *apperrors.ApplicationError
	require.ErrorAs(t, err, &appErr)
	assert.True(t, stderrors.Is(err, apperrors.ErrValidationCategory))
	assert.Equal(t, "model_validation_failed", appErr.Code)
}

func TestUserDatabaseRejectsDuplicateEmail(t *testing.T) {
	db := testutil.NewDatabase(t)
	tx := testutil.Begin(t, db)

	first := &User{Email: "duplicate@example.com", PasswordHash: "argon2-hash"}
	second := &User{Email: first.Email, PasswordHash: "another-argon2-hash"}

	require.NoError(t, tx.Create(first).Error)
	err := tx.Create(second).Error

	assert.ErrorIs(t, err, gorm.ErrDuplicatedKey)
}
