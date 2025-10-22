// cmd/server/main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/maleficent/go-maleficent/configs"
	"github.com/maleficent/go-maleficent/internal/pkg/database"
	"github.com/maleficent/go-maleficent/internal/pkg/logging"
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
	defer logger.Sync()

	db, err := database.NewDatabase(&config.Database, logger)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     config.CORS.AllowOrigins,
		AllowCredentials: config.CORS.AllowCredentials,
		AllowHeaders:     config.CORS.AllowHeaders,
	}))
	e.Use(logging.ZapLoggerMiddleware(logger))

	e.GET("/health", func(c echo.Context) error {
		logger := logging.GetLogger(c)
		logger.Info("Health check endpoint hit",
			zap.String("host", config.Server.Host),
			zap.Int("port", config.Server.Port),
			zap.String("log_level", config.Log.Level),
			zap.String("log_format", config.Log.Format),
		)

		return c.JSON(http.StatusOK, map[string]any{
			"status":   "ok",
			"message":  "Go Maleficent API is running",
			"database": db.GetStats(),
		})
	})

	address := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	logging.Info(logger, "Server listening", zap.String("address", address))
	e.Logger.Fatal(e.Start(address))
}
