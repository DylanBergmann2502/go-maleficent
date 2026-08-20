// internal/app/auth/handlers/user_handler.go
package handlers

import (
	"context"

	"github.com/DylanBergmann2502/go-maleficent/internal/app/auth/inputs"
	"github.com/DylanBergmann2502/go-maleficent/internal/app/auth/outputs"
	"github.com/DylanBergmann2502/go-maleficent/internal/app/auth/services"
	apiresponses "github.com/DylanBergmann2502/go-maleficent/internal/pkg/api/responses"
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/handlers"
)

type UserHandler struct {
	*handlers.BaseHandler
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{
		BaseHandler: &handlers.BaseHandler{},
		service:     service,
	}
}

func (h *UserHandler) Create(ctx context.Context, input *inputs.CreateUserInput) (*outputs.CreateUserOutput, error) {
	user, err := h.service.CreateUser(input.Body.Email, input.Body.Password)
	if err != nil {
		return nil, err
	}

	output := &outputs.CreateUserOutput{}
	output.Body.Data = outputs.FromModel(user)
	output.Body.Meta = apiresponses.NewMetaDataFromContext(ctx)
	return output, nil
}
