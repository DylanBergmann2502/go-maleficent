// internal/pkg/logging/middleware.go
package logging

import (
	"log/slog"

	"github.com/labstack/echo/v5"
)

const LoggerKey = "slog_logger"

// SlogLoggerMiddleware adds a slog logger to the Echo context.
func SlogLoggerMiddleware(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Set(LoggerKey, logger)
			return next(c)
		}
	}
}

// GetLogger returns the slog logger from the Echo context.
func GetLogger(c *echo.Context) *slog.Logger {
	logger, ok := c.Get(LoggerKey).(*slog.Logger)
	if !ok {
		panic("slog logger not found in context - make sure SlogLoggerMiddleware is registered")
	}
	return logger
}
