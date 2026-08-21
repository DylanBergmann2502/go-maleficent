// internal/pkg/logging/middleware.go
package logging

import (
	"log/slog"

	"github.com/labstack/echo/v5"
)

// SlogLoggerMiddleware adds a slog logger to the Echo context.
func SlogLoggerMiddleware(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.SetLogger(logger)
			return next(c)
		}
	}
}

// GetLogger returns the slog logger from the Echo context.
func GetLogger(c *echo.Context) *slog.Logger {
	return c.Logger()
}
