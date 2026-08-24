// internal/pkg/jobs/asynq_logger.go
package jobs

import "log/slog"

type asynqLogger struct {
	logger *slog.Logger
}

func (l asynqLogger) Debug(args ...interface{}) { l.logger.Debug("Asynq", "args", args) }
func (l asynqLogger) Info(args ...interface{})  { l.logger.Info("Asynq", "args", args) }
func (l asynqLogger) Warn(args ...interface{})  { l.logger.Warn("Asynq", "args", args) }
func (l asynqLogger) Error(args ...interface{}) { l.logger.Error("Asynq", "args", args) }
func (l asynqLogger) Fatal(args ...interface{}) { l.logger.Error("Asynq fatal error", "args", args) }
