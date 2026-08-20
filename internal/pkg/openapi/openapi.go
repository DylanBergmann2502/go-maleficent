// internal/pkg/openapi/openapi.go
package openapi

import (
	"context"
	"net/http"

	apierrors "github.com/DylanBergmann2502/go-maleficent/internal/pkg/api/errors"
	apiresponses "github.com/DylanBergmann2502/go-maleficent/internal/pkg/api/responses"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/labstack/echo/v5"
)

// New creates and mounts the Huma API and its documentation routes.
func New(e *echo.Echo) huma.API {
	huma.NewError = func(status int, message string, errs ...error) huma.StatusError {
		return apierrors.NewHumaError(context.TODO(), status, "", message, apierrors.HumaDetails(errs...))
	}
	huma.NewErrorWithContext = func(ctx huma.Context, status int, message string, errs ...error) huma.StatusError {
		return apierrors.NewHumaError(ctx.Context(), status, "", message, apierrors.HumaDetails(errs...))
	}

	api := humaecho.New(e, huma.Config{
		OpenAPI: &huma.OpenAPI{
			Info: &huma.Info{
				Title:       "Go Maleficent API",
				Version:     "1.0.0",
				Description: "Go Maleficent HTTP API",
			},
		},
		OpenAPIPath:   "/api/openapi",
		DocsPath:      "/swaggerui",
		DocsRenderer:  huma.DocsRendererSwaggerUI,
		SchemasPath:   "/api/schemas",
		Formats:       huma.DefaultFormats,
		DefaultFormat: "application/json",
	})

	api.UseMiddleware(func(ctx huma.Context, next func(huma.Context)) {
		echoContext := humaecho.Unwrap(ctx)
		requestID := echoContext.Response().Header().Get(echo.HeaderXRequestID)
		next(huma.WithContext(ctx, apiresponses.WithRequestID(ctx.Context(), requestID)))
	})

	e.GET("/api/openapi", func(c *echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/api/openapi.json")
	})
	e.GET("/swagger", func(c *echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/swaggerui")
	})

	return api
}
