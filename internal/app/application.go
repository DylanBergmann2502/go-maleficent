// internal/app/application.go
package app

import (
	"log/slog"

	"github.com/DylanBergmann2502/go-maleficent/configs"
	"github.com/DylanBergmann2502/go-maleficent/internal/app/auth"
	authhandlers "github.com/DylanBergmann2502/go-maleficent/internal/app/auth/handlers"
	"github.com/DylanBergmann2502/go-maleficent/internal/app/auth/services"
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/database"
	debug "github.com/DylanBergmann2502/go-maleficent/internal/pkg/debug"
	apphandlers "github.com/DylanBergmann2502/go-maleficent/internal/pkg/handlers"
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/jobs"
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/logging"
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/openapi"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

// Application contains the configured HTTP application.
type Application struct {
	Echo *echo.Echo
}

// New builds the HTTP application and wires its dependencies.
func New(config *configs.Config, logger *slog.Logger, db *database.Database, jobClient *jobs.Client) *Application {
	e := echo.New()
	e.Logger = logger

	e.Use(middleware.RequestID())
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     config.CORS.AllowOrigins,
		AllowCredentials: config.CORS.AllowCredentials,
		AllowHeaders:     config.CORS.AllowHeaders,
	}))
	e.Use(logging.SlogLoggerMiddleware(logger))

	validate := validator.New()
	userService := services.NewUserService(db.DB, validate, jobClient)
	userHandler := authhandlers.NewUserHandler(userService)

	humaAPI := openapi.New(e)
	auth.RegisterRoutes(humaAPI, userHandler)
	apphandlers.RegisterRoutes(e, db)
	debug.RegisterRoutes(e, config.Debug.PprofEndpointsEnabled)

	return &Application{Echo: e}
}
