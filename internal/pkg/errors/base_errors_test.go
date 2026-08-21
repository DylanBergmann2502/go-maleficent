// internal/pkg/errors/base_errors_test.go
package errors

import (
	stderrors "errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplicationErrorPreservesCategoryAndCause(t *testing.T) {
	cause := stderrors.New("database connection failed")
	err := WrapApplicationError(
		cause,
		ErrUnavailableCategory,
		"database_unavailable",
		"database is unavailable",
		nil,
	)

	assert.Equal(t, "database is unavailable", err.Error())
	assert.ErrorIs(t, err, ErrUnavailableCategory)
	assert.ErrorIs(t, err, cause)

	var applicationError *ApplicationError
	require.ErrorAs(t, err, &applicationError)
	assert.Equal(t, "database_unavailable", applicationError.Code)
}

func TestApplicationErrorWithoutCausePreservesCategory(t *testing.T) {
	err := NewApplicationError(
		ErrConflictCategory,
		"email_already_registered",
		"user already exists",
		map[string]string{"field": "email"},
	)

	assert.ErrorIs(t, err, ErrConflictCategory)
	assert.Equal(t, "email_already_registered", err.Code)
	assert.Equal(t, map[string]string{"field": "email"}, err.Details)
}
