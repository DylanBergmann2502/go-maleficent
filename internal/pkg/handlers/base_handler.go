// internal/pkg/handlers/base_handler.go
package handlers

import (
	"net/http"

	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/errors"
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
			return echo.NewHTTPError(http.StatusBadRequest, e.Message)
		case 401:
			return echo.NewHTTPError(http.StatusUnauthorized, e.Message)
		case 403:
			return echo.NewHTTPError(http.StatusForbidden, e.Message)
		case 404:
			return echo.NewHTTPError(http.StatusNotFound, e.Message)
		case 409:
			return echo.NewHTTPError(http.StatusConflict, e.Message)
		case 422:
			return echo.NewHTTPError(http.StatusUnprocessableEntity, e.Message)
		case 503:
			return echo.NewHTTPError(http.StatusServiceUnavailable, e.Message)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
}
