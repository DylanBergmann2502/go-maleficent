// internal/app/auth/forms/create_user_form_test.go
package forms

import (
	stderrors "errors"
	"testing"

	apperrors "github.com/DylanBergmann2502/go-maleficent/internal/pkg/errors"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUserFormValidate(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		password string
		valid    bool
		field    string
	}{
		{name: "valid input is normalized", email: " User@Example.COM ", password: "correct horse battery staple", valid: true},
		{name: "invalid email", email: "not-an-email", password: "correct horse battery staple", field: "email"},
		{name: "short password", email: "user@example.com", password: "short", field: "password"},
		{name: "long password", email: "user@example.com", password: "123456789012345678901234567890123456789012345678901234567890123456789012345678901", field: "password"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := NewCreateUserForm(tt.email, tt.password)
			err := form.Validate(validator.New())

			if tt.valid {
				require.NoError(t, err)
				assert.Equal(t, "user@example.com", form.Email)
				return
			}

			var appErr *apperrors.ApplicationError
			require.ErrorAs(t, err, &appErr)
			require.Error(t, err)
			assert.True(t, stderrors.Is(err, apperrors.ErrValidationCategory))
			assert.Contains(t, appErr.Details, tt.field)
		})
	}
}
