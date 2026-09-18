// cmd/server/commands/jobs.go
package commands

import (
	"fmt"
	"log/slog"

	"github.com/DylanBergmann2502/go-maleficent/config"
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/jobs"
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/logging"
	"github.com/spf13/cobra"
)

func NewWorkerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "worker",
		Short: "Start the background task worker",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWorker()
		},
	}
}

func NewCronCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "cron",
		Short: "Start the background task scheduler",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCron()
		},
	}
}

func loadJobsDependencies() (*config.Config, *slog.Logger, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load config: %w", err)
	}

	logger, err := logging.NewLogger(&cfg.Log)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	return cfg, logger, nil
}

func runWorker() error {
	config, logger, err := loadJobsDependencies()
	if err != nil {
		return err
	}

	worker, err := jobs.NewWorker(config, logger)
	if err != nil {
		return fmt.Errorf("failed to initialize worker: %w", err)
	}
	defer func() {
		if closeErr := worker.Close(); closeErr != nil {
			logger.Error("Failed to close worker", slog.Any("error", closeErr))
		}
	}()

	logger.Info("Background worker started")
	return worker.Run()
}

func runCron() error {
	config, logger, err := loadJobsDependencies()
	if err != nil {
		return err
	}

	scheduler, err := jobs.NewScheduler(config, logger)
	if err != nil {
		return fmt.Errorf("failed to initialize scheduler: %w", err)
	}
	defer scheduler.Close()

	logger.Info("Background scheduler started", "schedule", jobs.ExampleSchedule)
	if err := scheduler.Run(); err != nil {
		return err
	}

	return nil
}
