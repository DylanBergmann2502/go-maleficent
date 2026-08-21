// internal/pkg/debug/pprof.go
package debug

import (
	echopprof "github.com/labstack/echo-contrib/v5/pprof"
	"github.com/labstack/echo/v5"
)

// RegisterRoutes mounts Go profiling endpoints when explicitly enabled.
func RegisterRoutes(e *echo.Echo, enabled bool) {
	if !enabled {
		return
	}

	echopprof.Register(e)
}
