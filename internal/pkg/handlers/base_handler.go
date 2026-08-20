// internal/pkg/handlers/base_handler.go
package handlers

import (
	apierrors "github.com/DylanBergmann2502/go-maleficent/internal/pkg/api/errors"
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/api/responses"
	"github.com/labstack/echo/v5"
)

type BaseHandler struct{}

func (h *BaseHandler) HandleError(c *echo.Context, err error) error {
	if err == nil {
		return nil
	}

	apiErr := apierrors.FromError(err)
	return c.JSON(apiErr.Status, responses.Error(c, apiErr.Code, apiErr.Message, apiErr.Details))
}
