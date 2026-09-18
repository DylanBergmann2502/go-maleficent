// internal/app/auth/forms/update_user_form_test.go
package forms

import (
	stderrors "errors"
	"testing"

	apperrors "github.com/DylanBergmann2502/go-maleficent/internal/pkg/errors"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateUserFormValidate(t *testing.T) {
	tests := []struct {
		name     string
		email    *string
		password *string
		valid    bool
		field    string
	}{
		{name: "email only is normalized", email: stringPtr(" User@Example.COM "), valid: true},
		{name: "password only is valid", password: stringPtr("correct horse battery staple"), valid: true},
		{name: "both fields are valid", email: stringPtr("user@example.com"), password: stringPtr("correct horse battery staple"), valid: true},
		{name: "no fields", field: "base"},
		{name: "invalid email", email: stringPtr("not-an-email"), field: "email"},
		{name: "short password", password: stringPtr("short"), field: "password"},
		{name: "long password", password: stringPtr("123456789012345678901234567890123456789012345678901234567890123456789012345678901"), field: "password"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := NewUpdateUserForm(tt.email, tt.password)
			err := form.Validate(validator.New())

			if tt.valid {
				require.NoError(t, err)
				if tt.email != nil {
					require.NotNil(t, form.Email)
					assert.Equal(t, "user@example.com", *form.Email)
				}
				return
			}

			var appErr *apperrors.ApplicationError
			require.ErrorAs(t, err, &appErr)
			assert.True(t, stderrors.Is(err, apperrors.ErrValidationCategory))
			assert.Contains(t, appErr.Details, tt.field)
		})
	}
}

func stringPtr(value string) *string {
	return &value
}
