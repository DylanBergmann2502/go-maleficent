// internal/pkg/handlers/base_handler.go
package handlers

import (
	"net/http"

	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/errors"
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/responses"
	"github.com/labstack/echo/v5"
)

type BaseHandler struct{}

func (h *BaseHandler) HandleError(c *echo.Context, err error) error {
	if err == nil {
		return nil
	}

	switch e := err.(type) {
	case *errors.AppError:
		switch e.Code {
		case 400:
			return c.JSON(http.StatusBadRequest, responses.Error(c, "bad_request", e.Message, nil))
		case 401:
			return c.JSON(http.StatusUnauthorized, responses.Error(c, "unauthorized", e.Message, nil))
		case 403:
			return c.JSON(http.StatusForbidden, responses.Error(c, "forbidden", e.Message, nil))
		case 404:
			return c.JSON(http.StatusNotFound, responses.Error(c, "not_found", e.Message, nil))
		case 409:
			return c.JSON(http.StatusConflict, responses.Error(c, "conflict", e.Message, nil))
		case 422:
			return c.JSON(http.StatusUnprocessableEntity, responses.Error(c, "validation_failed", e.Message, nil))
		case 503:
			return c.JSON(http.StatusServiceUnavailable, responses.Error(c, "service_unavailable", e.Message, nil))
		default:
			return c.JSON(http.StatusInternalServerError, responses.Error(c, "internal_server_error", "Internal server error", nil))
		}
	default:
		return c.JSON(http.StatusInternalServerError, responses.Error(c, "internal_server_error", "Internal server error", nil))
	}
}
