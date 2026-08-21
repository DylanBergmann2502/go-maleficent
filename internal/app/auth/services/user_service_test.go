// internal/app/auth/services/user_service_test.go
package services

import (
	stderrors "errors"
	"testing"

	"github.com/DylanBergmann2502/go-maleficent/internal/app/auth/checks"
	"github.com/DylanBergmann2502/go-maleficent/internal/app/auth/models"
	"github.com/DylanBergmann2502/go-maleficent/internal/app/auth/utils"
	apperrors "github.com/DylanBergmann2502/go-maleficent/internal/pkg/errors"
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/testutil"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserServiceCreateUserPersistsNormalizedUser(t *testing.T) {
	db := testutil.NewDatabase(t)
	tx := testutil.Begin(t, db)
	service := NewUserService(tx, validator.New())

	user, err := service.CreateUser("  USER@Example.COM  ", "correct horse battery staple")

	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, "user@example.com", user.Email)
	assert.NotEqual(t, "correct horse battery staple", user.PasswordHash)
	assert.NotEmpty(t, user.PasswordHash)

	passwordMatches, err := utils.CheckPasswordHash("correct horse battery staple", user.PasswordHash)
	require.NoError(t, err)
	assert.True(t, passwordMatches)

	var persisted models.User
	require.NoError(t, tx.First(&persisted, "id = ?", user.ID).Error)
	assert.Equal(t, user.Email, persisted.Email)
}

func TestUserServiceCreateUserRejectsInvalidForm(t *testing.T) {
	db := testutil.NewDatabase(t)
	tx := testutil.Begin(t, db)
	service := NewUserService(tx, validator.New())

	user, err := service.CreateUser("invalid-email", "short")

	assert.Nil(t, user)
	var appErr *apperrors.ApplicationError
	require.ErrorAs(t, err, &appErr)
	assert.True(t, stderrors.Is(err, apperrors.ErrValidationCategory))
	assert.Equal(t, "validation_failed", appErr.Code)
}

func TestUserServiceCreateUserRejectsRegisteredEmail(t *testing.T) {
	db := testutil.NewDatabase(t)
	tx := testutil.Begin(t, db)
	service := NewUserService(tx, validator.New())

	require.NoError(t, tx.Create(&models.User{
		Email:        "registered@example.com",
		PasswordHash: "argon2-hash",
	}).Error)

	user, err := service.CreateUser("REGISTERED@example.com", "correct horse battery staple")

	assert.Nil(t, user)
	assert.ErrorIs(t, err, checks.ErrEmailAlreadyRegistered)
	assert.True(t, stderrors.Is(err, apperrors.ErrConflictCategory))
}
