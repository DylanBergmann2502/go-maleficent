// internal/pkg/jobs/tasks/example_task.go
package tasks

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/hibiken/asynq"
)

const ExampleTaskType = "example:print"

type ExampleTaskPayload struct {
	Message string `json:"message"`
}

func NewExampleTask(message string) *asynq.Task {
	payload, _ := json.Marshal(ExampleTaskPayload{Message: message})
	return asynq.NewTask(ExampleTaskType, payload)
}

type ExampleHandler struct {
	logger *slog.Logger
}

func NewExampleHandler(logger *slog.Logger) *ExampleHandler {
	return &ExampleHandler{logger: logger}
}

func (h *ExampleHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	var payload ExampleTaskPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}

	h.logger.InfoContext(ctx, "Example background task executed", "message", payload.Message)
	return nil
}
