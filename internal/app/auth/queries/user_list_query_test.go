// internal/app/auth/queries/user_list_query_test.go
package queries

import (
	stderrors "errors"
	"testing"
	"time"

	"github.com/DylanBergmann2502/go-maleficent/internal/app/auth/models"
	apperrors "github.com/DylanBergmann2502/go-maleficent/internal/pkg/errors"
	basemodels "github.com/DylanBergmann2502/go-maleficent/internal/pkg/models"
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUserListQueryTranslatesCanonicalParameters(t *testing.T) {
	params, err := NewUserListQuery(
		2,
		50,
		"-created_at,email",
		" *@example.com ",
		"one@example.com,two@example.com",
		"2026-01-01T00:00:00Z",
		"2026-02-01T00:00:00Z",
		"2026-01-15T00:00:00Z",
		"2026-02-15T00:00:00Z",
	)

	require.NoError(t, err)
	assert.Equal(t, 2, params.Page)
	assert.Equal(t, 50, params.PageSize)
	assert.Equal(t, "*@example.com", params.Email)
	assert.Equal(t, []string{"one@example.com", "two@example.com"}, params.EmailIn)
	assert.Len(t, params.Sort, 2)
	assert.False(t, params.Sort[0].Ascending)
	assert.True(t, params.Sort[1].Ascending)
	assert.Equal(t, "2026-01-01T00:00:00Z", params.CreatedAtGTE)
	assert.Equal(t, "2026-02-15T00:00:00Z", params.UpdatedAtLTE)
}

func TestNewUserListQueryRejectsUnsupportedSort(t *testing.T) {
	_, err := NewUserListQuery(1, 20, "-password_hash", "", "", "", "", "", "")

	var appErr *apperrors.ApplicationError
	require.ErrorAs(t, err, &appErr)
	assert.True(t, stderrors.Is(err, apperrors.ErrValidationCategory))
	assert.Contains(t, appErr.Details, "sort")
}

func TestUserListQueryAppliesFiltersSortsAndPaginates(t *testing.T) {
	db := testutil.NewDatabase(t)
	tx := testutil.Begin(t, db)
	createdAt := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	users := []models.User{
		{Email: "alpha@example.com", PasswordHash: "argon2-hash", BaseModel: baseModelAt(createdAt)},
		{Email: "beta@example.com", PasswordHash: "argon2-hash", BaseModel: baseModelAt(createdAt.Add(time.Hour))},
		{Email: "gamma@other.com", PasswordHash: "argon2-hash", BaseModel: baseModelAt(createdAt.Add(2 * time.Hour))},
	}
	for index := range users {
		require.NoError(t, tx.Create(&users[index]).Error)
	}

	params, err := NewUserListQuery(1, 1, "-email", "*@example.com", "", "", "", "", "")
	require.NoError(t, err)

	var result []models.User
	require.NoError(t, params.Apply(tx.Model(&models.User{})).Find(&result).Error)
	require.Len(t, result, 1)
	assert.Equal(t, "beta@example.com", result[0].Email)
}

func TestUserListQuerySupportsINAndDateFilters(t *testing.T) {
	db := testutil.NewDatabase(t)
	tx := testutil.Begin(t, db)
	createdAt := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	users := []models.User{
		{Email: "alpha@example.com", PasswordHash: "argon2-hash", BaseModel: baseModelAt(createdAt)},
		{Email: "beta@example.com", PasswordHash: "argon2-hash", BaseModel: baseModelAt(createdAt.Add(time.Hour))},
	}
	for index := range users {
		require.NoError(t, tx.Create(&users[index]).Error)
	}

	params, err := NewUserListQuery(
		1,
		20,
		"email",
		"",
		"beta@example.com",
		createdAt.Format(time.RFC3339),
		createdAt.Add(2*time.Hour).Format(time.RFC3339),
		"",
		"",
	)
	require.NoError(t, err)

	var result []models.User
	require.NoError(t, params.Apply(tx.Model(&models.User{})).Find(&result).Error)
	require.Len(t, result, 1)
	assert.Equal(t, "beta@example.com", result[0].Email)
}

func baseModelAt(createdAt time.Time) basemodels.BaseModel {
	return basemodels.BaseModel{CreatedAt: createdAt, UpdatedAt: createdAt}
}
