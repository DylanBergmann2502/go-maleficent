// internal/pkg/database/migrations.go
package database

import (
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/DylanBergmann2502/go-maleficent/configs"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type Migrator struct {
	migrate *migrate.Migrate
}

// NewMigrator creates a new migrator instance with embedded migrations
func NewMigrator(config *configs.DatabaseConfig) (*Migrator, error) {
	// Create iofs source from embedded files
	sourceDriver, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to create migration source: %w", err)
	}

	// Build database connection string for migrations
	dbURL := buildMigrationDSN(config)

	// Create migrate instance
	m, err := migrate.NewWithSourceInstance("iofs", sourceDriver, dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	return &Migrator{migrate: m}, nil
}

// Up runs all pending migrations
func (m *Migrator) Up() error {
	if err := m.migrate.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration up failed: %w", err)
	}
	return nil
}

// UpSteps runs N migrations up
func (m *Migrator) UpSteps(steps int) error {
	if err := m.migrate.Steps(steps); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration up failed: %w", err)
	}
	return nil
}

// Down rolls back N migrations (default: 1)
func (m *Migrator) Down(steps int) error {
	if err := m.migrate.Steps(-steps); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration down failed: %w", err)
	}
	return nil
}

// Close closes the migrator
func (m *Migrator) Close() error {
	srcErr, dbErr := m.migrate.Close()
	if srcErr != nil {
		return srcErr
	}
	return dbErr
}

// buildMigrationDSN builds a postgres connection string for migrations
// This format is required by golang-migrate (different from GORM format)
func buildMigrationDSN(config *configs.DatabaseConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		config.SSLMode,
	)
}
