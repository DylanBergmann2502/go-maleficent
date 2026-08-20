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
		return apierrors.FromHumaErrors(ctx.Context(), status, message, errs...)
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
	e.GET("/redoc", func(c *echo.Context) error {
		return c.HTML(http.StatusOK, redocHTML)
	})

	return api
}

const redocHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Go Maleficent API Reference</title>
</head>
<body>
  <redoc spec-url="/api/openapi.json"></redoc>
  <script src="https://cdn.jsdelivr.net/npm/redoc@2.5.0/bundles/redoc.standalone.js"></script>
</body>
</html>`
