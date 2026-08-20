// internal/app/auth/handlers/user_handler.go
package handlers

import (
	"net/http"

	"github.com/DylanBergmann2502/go-maleficent/internal/app/auth/requests"
	userresponses "github.com/DylanBergmann2502/go-maleficent/internal/app/auth/responses"
	"github.com/DylanBergmann2502/go-maleficent/internal/app/auth/services"
	apiresponses "github.com/DylanBergmann2502/go-maleficent/internal/pkg/api/responses"
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/errors"
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/handlers"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
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

func (h *UserHandler) Create(c *echo.Context) error {
	var req requests.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return h.HandleError(c, errors.ErrBadRequest)
	}

	if err := h.validator.Struct(req); err != nil {
		return h.HandleError(c, errors.NewApplicationError(errors.ErrValidationCategory, "validation_failed", "validation failed", err.Error()))
	}

	user, err := h.service.CreateUser(req.Email, req.Password)
	if err != nil {
		return h.HandleError(c, err)
	}

	return c.JSON(http.StatusCreated, apiresponses.Success(c, userresponses.FromModel(user)))
}
