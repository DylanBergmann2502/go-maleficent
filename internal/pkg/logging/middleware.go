// internal/pkg/logging/middleware.go
package logging

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

const LoggerKey = "zap_logger"

// Logger provides easy access to Zap logger from Echo context
type Logger struct {
	*zap.Logger
}

// Info logs an info level message
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.Logger.Info(msg, fields...)
}

// Error logs an error level message with caller info
func (l *Logger) Error(msg string, fields ...zap.Field) {
	loggerWithCaller := l.Logger.WithOptions(zap.AddCaller())
	loggerWithCaller.Error(msg, fields...)
}

// Debug logs a debug level message
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.Logger.Debug(msg, fields...)
}

// Warn logs a warn level message
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.Logger.Warn(msg, fields...)
}

// Fatal logs a fatal level message with caller info
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	loggerWithCaller := l.Logger.WithOptions(zap.AddCaller())
	loggerWithCaller.Fatal(msg, fields...)
}

// ZapLoggerMiddleware adds Zap logger to Echo context
func ZapLoggerMiddleware(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			zapLogger := &Logger{Logger: logger}
			c.Set(LoggerKey, zapLogger)
			return next(c)
		}
	}
}

// GetLogger returns the logger from context as a variable for cleaner syntax
// Usage: logger := logging.GetLogger(c); logger.Info("message")
func GetLogger(c echo.Context) *Logger {
	logger, ok := c.Get(LoggerKey).(*Logger)
	if !ok {
		panic("zap logger not found in context - make sure ZapLoggerMiddleware is registered")
	}
	return logger
}
