// internal/app/application.go
package app

import (
	"log/slog"

	"github.com/DylanBergmann2502/go-maleficent/config"
	"github.com/DylanBergmann2502/go-maleficent/internal/app/auth"
	authhandlers "github.com/DylanBergmann2502/go-maleficent/internal/app/auth/handlers"
	"github.com/DylanBergmann2502/go-maleficent/internal/app/auth/services"
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/database"
	debug "github.com/DylanBergmann2502/go-maleficent/internal/pkg/debug"
	apphandlers "github.com/DylanBergmann2502/go-maleficent/internal/pkg/handlers"
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/jobs"
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/logging"
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/mailer"
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/openapi"
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/storage"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

// Application contains the configured HTTP application.
type Application struct {
	Echo    *echo.Echo
	Mailer  *mailer.Mailer
	Storage *storage.S3Storage
}

// New builds the HTTP application and wires its dependencies.
func New(cfg *config.Config, logger *slog.Logger, db *database.Database, jobClient *jobs.Client, emailMailer *mailer.Mailer, objectStorage *storage.S3Storage) *Application {
	e := echo.New()
	e.Logger = logger

	e.Use(middleware.RequestID())
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     cfg.CORS.AllowOrigins,
		AllowCredentials: cfg.CORS.AllowCredentials,
		AllowHeaders:     cfg.CORS.AllowHeaders,
	}))
	e.Use(logging.SlogLoggerMiddleware(logger))

	validate := validator.New()
	userService := services.NewUserService(db.DB, validate, jobClient)
	userHandler := authhandlers.NewUserHandler(userService)

	humaAPI := openapi.New(e)
	auth.RegisterRoutes(humaAPI, userHandler)
	apphandlers.RegisterRoutes(e, db)
	debug.RegisterRoutes(e, cfg.Debug.PprofEndpointsEnabled)

	return &Application{Echo: e, Mailer: emailMailer, Storage: objectStorage}
}
