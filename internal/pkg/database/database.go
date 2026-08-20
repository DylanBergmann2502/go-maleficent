// internal/pkg/database/database.go
package database

import (
	"fmt"

	"github.com/DylanBergmann2502/go-maleficent/configs"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
	config *configs.DatabaseConfig
	logger *zap.Logger
}

func NewDatabase(config *configs.DatabaseConfig, logger *zap.Logger) (*Database, error) {
	dsn := buildDSN(config)

	gormConfig := &gorm.Config{TranslateError: true}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB for connection pool configuration
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool settings
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database connection established successfully",
		zap.String("host", config.Host),
		zap.Int("port", config.Port),
		zap.String("database", config.Database),
		zap.String("user", config.User),
		zap.Int("max_open_conns", config.MaxOpenConns),
		zap.Int("max_idle_conns", config.MaxIdleConns),
		zap.Duration("conn_max_lifetime", config.ConnMaxLifetime),
		zap.Duration("conn_max_idle_time", config.ConnMaxIdleTime),
	)

	return &Database{
		DB:     db,
		config: config,
		logger: logger,
	}, nil
}

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	d.logger.Info("Closing database connection")
	return sqlDB.Close()
}

func (d *Database) GetStats() map[string]any {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return map[string]any{"error": err.Error()}
	}

	stats := sqlDB.Stats()
	return map[string]any{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration.String(),
	}
}

func buildDSN(config *configs.DatabaseConfig) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.Database,
		config.SSLMode,
		config.Timezone,
	)
}
