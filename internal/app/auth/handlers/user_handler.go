// internal/app/auth/handlers/user_handler.go
package handlers

import (
	"context"

	"github.com/DylanBergmann2502/go-maleficent/internal/app/auth/inputs"
	"github.com/DylanBergmann2502/go-maleficent/internal/app/auth/outputs"
	"github.com/DylanBergmann2502/go-maleficent/internal/app/auth/services"
	apierrors "github.com/DylanBergmann2502/go-maleficent/internal/pkg/api/errors"
	apiresponses "github.com/DylanBergmann2502/go-maleficent/internal/pkg/api/responses"
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/errors"
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/handlers"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	*handlers.BaseHandler
	service   *services.UserService
	validator *validator.Validate
}

func NewUserHandler(service *services.UserService, v *validator.Validate) *UserHandler {
	return &UserHandler{
		BaseHandler: &handlers.BaseHandler{},
		service:     service,
		validator:   v,
	}
}

func (h *UserHandler) Create(ctx context.Context, input *inputs.CreateUserInput) (*outputs.CreateUserOutput, error) {
	if err := h.validator.Struct(input.Body); err != nil {
		return nil, apierrors.ToHumaError(ctx,
			errors.NewApplicationError(errors.ErrValidationCategory, "validation_failed", "validation failed", err.Error()),
		)
	}

	user, err := h.service.CreateUser(input.Body.Email, input.Body.Password)
	if err != nil {
		return nil, apierrors.ToHumaError(ctx, err)
	}

	output := &outputs.CreateUserOutput{}
	output.Body.Data = outputs.FromModel(user)
	output.Body.Meta = apiresponses.NewMetaDataFromContext(ctx)
	return output, nil
}
