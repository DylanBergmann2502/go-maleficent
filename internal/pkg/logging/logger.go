// internal/pkg/logging/logger.go
package logging

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/DylanBergmann2502/go-maleficent/configs"
)

// NewLogger creates a structured application logger.
func NewLogger(logConfig *configs.LogConfig) (*slog.Logger, error) {
	options := &slog.HandlerOptions{Level: getLogLevel(logConfig.Level)}
	writer, err := outputWriter(logConfig.OutputPath)
	if err != nil {
		return nil, err
	}

	if logConfig.Format == "console" {
		return slog.New(slog.NewTextHandler(writer, options)), nil
	}
	return slog.New(slog.NewJSONHandler(writer, options)), nil
}

// NewLoggerWithCaller creates a logger that includes source locations.
func NewLoggerWithCaller(logConfig *configs.LogConfig) (*slog.Logger, error) {
	options := &slog.HandlerOptions{
		Level:     getLogLevel(logConfig.Level),
		AddSource: true,
	}
	writer, err := outputWriter(logConfig.OutputPath)
	if err != nil {
		return nil, err
	}

	if logConfig.Format == "console" {
		return slog.New(slog.NewTextHandler(writer, options)), nil
	}
	return slog.New(slog.NewJSONHandler(writer, options)), nil
}

func getLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error", "fatal":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func outputWriter(paths []string) (io.Writer, error) {
	if len(paths) == 0 {
		return os.Stdout, nil
	}

	writers := make([]io.Writer, 0, len(paths))
	for _, path := range paths {
		switch path {
		case "stdout":
			writers = append(writers, os.Stdout)
		case "stderr":
			writers = append(writers, os.Stderr)
		default:
			file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
			if err != nil {
				return nil, err
			}
			writers = append(writers, file)
		}
	}

	if len(writers) == 1 {
		return writers[0], nil
	}
	return io.MultiWriter(writers...), nil
}

func Info(logger *slog.Logger, message string, attrs ...slog.Attr) {
	logger.LogAttrs(context.Background(), slog.LevelInfo, message, attrs...)
}

func Error(logger *slog.Logger, message string, attrs ...slog.Attr) {
	logger.LogAttrs(context.Background(), slog.LevelError, message, attrs...)
}

func Debug(logger *slog.Logger, message string, attrs ...slog.Attr) {
	logger.LogAttrs(context.Background(), slog.LevelDebug, message, attrs...)
}

func Warn(logger *slog.Logger, message string, attrs ...slog.Attr) {
	logger.LogAttrs(context.Background(), slog.LevelWarn, message, attrs...)
}
