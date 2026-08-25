// cmd/server/commands/api.go
package commands

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/DylanBergmann2502/go-maleficent/config"
	application "github.com/DylanBergmann2502/go-maleficent/internal/app"
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/database"
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/jobs"
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/logging"
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/storage"
	"github.com/spf13/cobra"
)

// NewAPICommand creates the HTTP API server command.
func NewAPICommand() *cobra.Command {
	return &cobra.Command{
		Use:   "api",
		Short: "Start the API server",
		Long:  "Start the Go Maleficent API server with all HTTP endpoints",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAPIServer()
		},
	}
}

func runAPIServer() error {
	config, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	logger, err := logging.NewLogger(&config.Log)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	db, err := database.NewDatabase(&config.Database, logger)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			logger.Error("Failed to close database", slog.Any("error", closeErr))
		}
	}()

	jobClient, err := jobs.NewClient(config)
	if err != nil {
		return fmt.Errorf("failed to initialize background job client: %w", err)
	}
	defer func() {
		if closeErr := jobClient.Close(); closeErr != nil {
			logger.Error("Failed to close background job client", slog.Any("error", closeErr))
		}
	}()

	objectStorage, err := storage.NewS3Storage(context.Background(), &config.Storage)
	if err != nil {
		return fmt.Errorf("failed to initialize object storage: %w", err)
	}

	application := application.New(config, logger, db, jobClient, objectStorage)
	address := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	logging.Info(logger, "Server listening", slog.String("address", address))

	if err := application.Echo.Start(address); err != nil {
		logger.Error("Server stopped", slog.Any("error", err))
		return err
	}

	return nil
}
