// cmd/server/commands/migrate.go
package commands

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/DylanBergmann2502/go-maleficent/config"
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/database"
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/logging"
	"github.com/spf13/cobra"
)

// NewMigrateCommand creates the database migration command tree.
func NewMigrateCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "migrate",
		Short: "Database migration commands",
		Long:  "Run database migrations up or down using embedded migration files",
	}

	command.AddCommand(newMigrateUpCommand())
	command.AddCommand(newMigrateDownCommand())

	return command
}

func newMigrateUpCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "up [steps]",
		Short: "Run migrations up",
		Long:  "Apply pending migrations. Optionally specify the number of steps.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMigrateUp(args)
		},
	}
}

func newMigrateDownCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "down [steps]",
		Short: "Rollback migrations",
		Long:  "Rollback database migrations. Specify the number of steps.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMigrateDown(args)
		},
	}
}

func runMigrateUp(args []string) (err error) {
	config, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	logger, err := logging.NewLogger(&config.Log)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	migrator, err := database.NewMigrator(&config.Database)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer func() {
		err = errors.Join(err, migrator.Close())
	}()

	if len(args) == 0 {
		logger.Info("Running all pending migrations")
		return migrator.Up()
	}

	steps, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid migration steps: %w", err)
	}
	logger.Info("Running migrations up", "steps", steps)
	return migrator.UpSteps(steps)
}

func runMigrateDown(args []string) (err error) {
	config, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	logger, err := logging.NewLogger(&config.Log)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	migrator, err := database.NewMigrator(&config.Database)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer func() {
		err = errors.Join(err, migrator.Close())
	}()

	steps := 1
	if len(args) > 0 {
		steps, err = strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid migration steps: %w", err)
		}
	}

	logger.Info("Rolling back migrations", "steps", steps)
	return migrator.Down(steps)
}
