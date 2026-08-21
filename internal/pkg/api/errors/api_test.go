// internal/pkg/api/errors/api_test.go
package errors

import (
	"context"
	stderrors "errors"
	"net/http"
	"testing"

	apperrors "github.com/DylanBergmann2502/go-maleficent/internal/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromErrorTranslatesConflict(t *testing.T) {
	err := apperrors.NewApplicationError(
		apperrors.ErrConflictCategory,
		"email_already_registered",
		"user with this email already exists",
		nil,
	)

	apiErr := FromError(err)

	assert.Equal(t, http.StatusConflict, apiErr.Status)
	assert.Equal(t, "conflict", apiErr.Code)
	assert.Equal(t, "user with this email already exists", apiErr.Message)
}

func TestFromHumaErrorsPreservesApplicationError(t *testing.T) {
	err := apperrors.NewApplicationError(
		apperrors.ErrValidationCategory,
		"validation_failed",
		"validation failed",
		map[string][]string{"email": {"must be a valid email address"}},
	)

	humaErr := FromHumaErrors(context.Background(), http.StatusInternalServerError, "unexpected error", err)

	var statusErr interface{ GetStatus() int }
	require.ErrorAs(t, humaErr, &statusErr)
	assert.Equal(t, http.StatusUnprocessableEntity, statusErr.GetStatus())
	assert.True(t, stderrors.Is(err, apperrors.ErrValidationCategory))
}
