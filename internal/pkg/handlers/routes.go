// internal/pkg/handlers/routes.go
package handlers

import (
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/database"
	"github.com/labstack/echo/v5"
)

// RegisterRoutes registers application-wide HTTP routes.
func RegisterRoutes(e *echo.Echo, db *database.Database) {
	NewHealthHandler(db).RegisterRoutes(e)
}
