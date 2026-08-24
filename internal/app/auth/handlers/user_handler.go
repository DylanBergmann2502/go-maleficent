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

func (h *UserHandler) List(ctx context.Context, _ *inputs.ListUsersInput) (*outputs.ListUsersOutput, error) {
	users, err := h.service.ListUsers()
	if err != nil {
		return nil, err
	}

	data := make([]outputs.UserOutput, 0, len(users))
	for index := range users {
		data = append(data, outputs.FromModel(&users[index]))
	}

	output := &outputs.ListUsersOutput{}
	output.Body.Data = data
	output.Body.Meta = apiresponses.NewMetaDataFromContext(ctx)
	return output, nil
}

func (h *UserHandler) Get(ctx context.Context, input *inputs.GetUserInput) (*outputs.GetUserOutput, error) {
	user, err := h.service.GetUser(input.ID)
	if err != nil {
		return nil, err
	}

	output := &outputs.GetUserOutput{}
	output.Body.Data = outputs.FromModel(user)
	output.Body.Meta = apiresponses.NewMetaDataFromContext(ctx)
	return output, nil
}

func (h *UserHandler) Update(ctx context.Context, input *inputs.UpdateUserInput) (*outputs.UpdateUserOutput, error) {
	user, err := h.service.UpdateUser(input.ID, input.Body.Email, input.Body.Password)
	if err != nil {
		return nil, err
	}

	output := &outputs.UpdateUserOutput{}
	output.Body.Data = outputs.FromModel(user)
	output.Body.Meta = apiresponses.NewMetaDataFromContext(ctx)
	return output, nil
}

func (h *UserHandler) Delete(_ context.Context, input *inputs.DeleteUserInput) (*outputs.DeleteUserOutput, error) {
	if err := h.service.DeleteUser(input.ID); err != nil {
		return nil, err
	}

	return &outputs.DeleteUserOutput{}, nil
}
