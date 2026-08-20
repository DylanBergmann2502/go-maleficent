// internal/app/auth/forms/create_user_form.go
package forms

import (
	"errors"
	"strings"

	apperrors "github.com/DylanBergmann2502/go-maleficent/internal/pkg/errors"
	"github.com/go-playground/validator/v10"
)

// CreateUserForm contains workflow-aware input for creating a user.
// Huma validates the transport shape; the form cleans and validates the
// values again before they enter the application workflow.
type CreateUserForm struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=72"`
}

// NewCreateUserForm cleans values that are safe to normalize at this layer.
func NewCreateUserForm(email string, password string) *CreateUserForm {
	return &CreateUserForm{
		Email:    strings.ToLower(strings.TrimSpace(email)),
		Password: password,
	}
}

// Validate validates workflow-level input and returns an application error.
func (f *CreateUserForm) Validate(v *validator.Validate) error {
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

func validationDetails(err error) map[string][]string {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return map[string][]string{"base": {err.Error()}}
	}

	details := make(map[string][]string, len(validationErrors))
	for _, validationError := range validationErrors {
		field := strings.ToLower(validationError.Field())
		details[field] = append(details[field], validationMessage(validationError))
	}

	return details
}

func validationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "can't be blank"
	case "email":
		return "must be a valid email address"
	case "min":
		return "is too short"
	case "max":
		return "is too long"
	default:
		return "is invalid"
	}
}
