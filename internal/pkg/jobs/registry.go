// internal/pkg/jobs/registry.go
package jobs

import (
	"log/slog"

	"github.com/DylanBergmann2502/go-maleficent/config"
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/jobs/tasks"
	"github.com/hibiken/asynq"
)

const ExampleSchedule = "* * * * *"

func RegisterHandlers(mux *asynq.ServeMux, logger *slog.Logger) {
	mux.Handle(tasks.ExampleTaskType, tasks.NewExampleHandler(logger))
}

func RegisterSchedules(scheduler *asynq.Scheduler, config *config.Config) error {
	_, err := scheduler.Register(
		ExampleSchedule,
		tasks.NewExampleTask("scheduled example task"),
	)
	return err
}
