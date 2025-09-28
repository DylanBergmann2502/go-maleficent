// internal/pkg/logging/logger.go
package logging

import (
	"github.com/maleficent/go-maleficent/configs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func customColorLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	switch level {
	case zapcore.DebugLevel:
		enc.AppendString("\033[37mDEBUG\033[0m") // White/Gray (default)
	case zapcore.InfoLevel:
		enc.AppendString("\033[32mINFO\033[0m") // Green (custom)
	case zapcore.WarnLevel:
		enc.AppendString("\033[33mWARN\033[0m") // Yellow
	case zapcore.ErrorLevel:
		enc.AppendString("\033[31mERROR\033[0m") // Red
	case zapcore.FatalLevel:
		enc.AppendString("\033[35mFATAL\033[0m") // Magenta
	default:
		enc.AppendString(level.CapitalString())
	}
}

func NewLogger(logConfig *configs.LogConfig) (*zap.Logger, error) {
	var config zap.Config

	if logConfig.Format == "console" {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = customColorLevelEncoder
	} else {
		config = zap.NewProductionConfig()
	}

	config.Level = zap.NewAtomicLevelAt(getLogLevel(logConfig.Level))
	config.OutputPaths = logConfig.OutputPath
	config.ErrorOutputPaths = []string{"stderr"}

	logger, err := config.Build(zap.AddCallerSkip(0))
	if err != nil {
		return nil, err
	}

	return logger.WithOptions(
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.WithCaller(false),
	), nil
}

func NewLoggerWithCaller(logConfig *configs.LogConfig) (*zap.Logger, error) {
	logger, err := NewLogger(logConfig)
	if err != nil {
		return nil, err
	}

	return logger.WithOptions(zap.AddCaller()), nil
}

func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func Info(logger *zap.Logger, msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Error(logger *zap.Logger, msg string, fields ...zap.Field) {
	loggerWithCaller := logger.WithOptions(zap.AddCaller())
	loggerWithCaller.Error(msg, fields...)
}

func Debug(logger *zap.Logger, msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Warn(logger *zap.Logger, msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func Fatal(logger *zap.Logger, msg string, fields ...zap.Field) {
	loggerWithCaller := logger.WithOptions(zap.AddCaller())
	loggerWithCaller.Fatal(msg, fields...)
}
