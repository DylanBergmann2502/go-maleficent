// cmd/server/main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/DylanBergmann2502/go-maleficent/configs"
	"github.com/DylanBergmann2502/go-maleficent/internal/app/auth/handlers"
	"github.com/DylanBergmann2502/go-maleficent/internal/app/auth/services"
	apiresponses "github.com/DylanBergmann2502/go-maleficent/internal/pkg/api/responses"
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/database"
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/logging"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "server",
		Short: "Go Maleficent - API Server and Management CLI",
		Long:  "A comprehensive Go API server with database migrations, background workers, and more.",
	}

	rootCmd.AddCommand(apiCmd())
	rootCmd.AddCommand(migrateCmd())

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Failed to execute command: %v", err)
	}
}

func apiCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "api",
		Short: "Start the API server",
		Long:  "Start the Go Maleficent API server with all HTTP endpoints",
		Run: func(cmd *cobra.Command, args []string) {
			runAPIServer()
		},
	}
}

func runAPIServer() {
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger, err := logging.NewLogger(&config.Log)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Printf("Failed to sync logger: %v", err)
		}
	}()

	db, err := database.NewDatabase(&config.Database, logger)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Failed to close database: %v", err)
		}
	}()

	e := echo.New()

	e.Use(middleware.RequestID())
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     config.CORS.AllowOrigins,
		AllowCredentials: config.CORS.AllowCredentials,
		AllowHeaders:     config.CORS.AllowHeaders,
	}))
	e.Use(logging.ZapLoggerMiddleware(logger))

	validate := validator.New()

	// User Resource
	userService := services.NewUserService(db.DB)
	userHandler := handlers.NewUserHandler(userService, validate)

	// API Routes
	api := e.Group("/api")
	v1 := api.Group("/v1")
	users := v1.Group("/users")
	users.POST("/", userHandler.Create)

	e.GET("/health", func(c *echo.Context) error {
		logger := logging.GetLogger(c)
		logger.Info("Health check endpoint hit",
			zap.String("host", config.Server.Host),
			zap.Int("port", config.Server.Port),
			zap.String("log_level", config.Log.Level),
			zap.String("log_format", config.Log.Format),
		)

		return c.JSON(http.StatusOK, apiresponses.Success(c, map[string]any{
			"status":   "ok",
			"message":  "Go Maleficent API is running",
			"database": db.GetStats(),
		}))
	})

	address := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	logging.Info(logger, "Server listening", zap.String("address", address))
	if err := e.Start(address); err != nil {
		logging.Fatal(logger, "Server stopped", zap.Error(err))
	}
}

func migrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Database migration commands",
		Long:  "Run database migrations up or down using embedded migration files",
	}

	cmd.AddCommand(migrateUpCmd())
	cmd.AddCommand(migrateDownCmd())

	return cmd
}

func migrateUpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "up [steps]",
		Short: "Run migrations up",
		Long:  "Apply pending database migrations. Optionally specify number of steps.",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			config, err := configs.LoadConfig()
			if err != nil {
				log.Fatalf("Failed to load config: %v", err)
			}

			migrator, err := database.NewMigrator(&config.Database)
			if err != nil {
				log.Fatalf("Failed to create migrator: %v", err)
			}
			defer func() {
				if err := migrator.Close(); err != nil {
					log.Printf("Failed to close migrator: %v", err)
				}
			}()

			if len(args) == 0 {
				log.Println("Running all pending migrations...")
				if err := migrator.Up(); err != nil {
					log.Fatalf("Migration failed: %v", err)
				}
			} else {
				steps, parseErr := strconv.Atoi(args[0])
				if parseErr != nil {
					log.Fatalf("Invalid steps argument: %v", parseErr)
				}
				log.Printf("Running %d migration(s) up...\n", steps)
				if err := migrator.UpSteps(steps); err != nil {
					log.Fatalf("Migration failed: %v", err)
				}
			}

			log.Println("Migrations applied successfully")
		},
	}
}

func migrateDownCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "down [steps]",
		Short: "Run migrations down",
		Long:  "Rollback database migrations. Specify number of steps (default: 1).",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			config, err := configs.LoadConfig()
			if err != nil {
				log.Fatalf("Failed to load config: %v", err)
			}

			migrator, err := database.NewMigrator(&config.Database)
			if err != nil {
				log.Fatalf("Failed to create migrator: %v", err)
			}
			defer func() {
				if err := migrator.Close(); err != nil {
					log.Printf("Failed to close migrator: %v", err)
				}
			}()

			steps := 1
			if len(args) > 0 {
				steps, err = strconv.Atoi(args[0])
				if err != nil {
					log.Fatalf("Invalid steps argument: %v", err)
				}
			}

			log.Printf("Rolling back %d migration(s)...\n", steps)
			if err := migrator.Down(steps); err != nil {
				log.Fatalf("Migration rollback failed: %v", err)
			}

			log.Println("Migrations rolled back successfully")
		},
	}
}
