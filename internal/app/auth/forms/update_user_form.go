// internal/app/auth/forms/update_user_form.go
package forms

import (
	"strings"

	apperrors "github.com/DylanBergmann2502/go-maleficent/internal/pkg/errors"
	"github.com/go-playground/validator/v10"
)

// UpdateUserForm contains workflow-aware input for a partial user update.
type UpdateUserForm struct {
	Email    *string `validate:"omitempty,email"`
	Password *string `validate:"omitempty,min=8,max=72"`
}

// NewUpdateUserForm cleans values that are safe to normalize at this layer.
func NewUpdateUserForm(email, password *string) *UpdateUserForm {
	form := &UpdateUserForm{Email: email, Password: password}
	if form.Email != nil {
		normalized := strings.ToLower(strings.TrimSpace(*form.Email))
		form.Email = &normalized
	}

	return form
}

// Validate validates the fields supplied for the update workflow.
func (f *UpdateUserForm) Validate(v *validator.Validate) error {
	if f.Email == nil && f.Password == nil {
		return apperrors.NewApplicationError(
			apperrors.ErrValidationCategory,
			"validation_failed",
			"at least one field must be provided",
			map[string][]string{"base": {"at least one field must be provided"}},
		)
	}

	if err := v.Struct(f); err != nil {
		return apperrors.NewApplicationError(
			apperrors.ErrValidationCategory,
			"validation_failed",
			"validation failed",
			validationDetails(err),
		)
	}

	return nil
}
