// internal/pkg/handlers/health_handler.go
package handlers

import (
	"net/http"

	apiresponses "github.com/DylanBergmann2502/go-maleficent/internal/pkg/api/responses"
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/database"
	"github.com/labstack/echo/v5"
)

// HealthHandler exposes application-wide liveness and readiness endpoints.
type HealthHandler struct {
	*BaseHandler
	db *database.Database
}

// NewHealthHandler creates a health handler backed by the application database.
func NewHealthHandler(db *database.Database) *HealthHandler {
	return &HealthHandler{
		BaseHandler: &BaseHandler{},
		db:          db,
	}
}

// RegisterRoutes mounts the application-wide health endpoints.
func (h *HealthHandler) RegisterRoutes(e *echo.Echo) {
	e.GET("/ping", h.Ping)
	e.GET("/livez", h.Ping)
	e.GET("/health", h.Check)
	e.GET("/healthz", h.Check)
	e.GET("/readyz", h.Check)
}

// Ping is a minimal liveness check that does not inspect dependencies.
func (h *HealthHandler) Ping(c *echo.Context) error {
	return c.JSON(http.StatusOK, apiresponses.Success(c, map[string]any{
		"status": "ok",
	}))
}

// Check verifies that the application can reach its database.
func (h *HealthHandler) Check(c *echo.Context) error {
	if err := h.db.Exec("SELECT 1").Error; err != nil {
		return c.JSON(http.StatusServiceUnavailable, apiresponses.Error(
			c,
			"service_unavailable",
			"database is not healthy",
			nil,
		))
	}

	return c.JSON(http.StatusOK, apiresponses.Success(c, map[string]any{
		"status":   "ok",
		"message":  "Go Maleficent API is running",
		"database": h.db.GetStats(),
	}))
}
