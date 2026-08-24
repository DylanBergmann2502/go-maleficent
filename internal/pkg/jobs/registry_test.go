// internal/pkg/jobs/registry_test.go
package jobs

import (
	"io"
	"log/slog"
	"testing"

	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/jobs/tasks"
	"github.com/hibiken/asynq"
)

func TestRegisterHandlers(t *testing.T) {
	mux := asynq.NewServeMux()
	RegisterHandlers(mux, slog.New(slog.NewTextHandler(io.Discard, nil)))

	handler, pattern := mux.Handler(tasks.NewExampleTask("test"))
	if handler == nil {
		t.Fatal("expected example task handler to be registered")
	}
	if pattern != tasks.ExampleTaskType {
		t.Fatalf("expected pattern %q, got %q", tasks.ExampleTaskType, pattern)
	}
}
