// internal/app/auth/handlers/user_handler.go
package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/maleficent/go-maleficent/internal/app/auth/requests"
	"github.com/maleficent/go-maleficent/internal/app/auth/responses"
	"github.com/maleficent/go-maleficent/internal/app/auth/services"
)

type UserHandler struct {
	service   *services.UserService
	validator *validator.Validate
}

func NewUserHandler(service *services.UserService, v *validator.Validate) *UserHandler {
	return &UserHandler{
		service:   service,
		validator: v,
	}
}

// Create handles user registration
func (h *UserHandler) Create(c echo.Context) error {
	var req requests.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := h.validator.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	user, err := h.service.CreateUser(req.Email, req.Password)
	if err != nil {
		if err == services.ErrUserAlreadyExists {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user")
	}

	return c.JSON(http.StatusCreated, responses.FromModel(user))
}
