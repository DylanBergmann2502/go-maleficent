// internal/pkg/debug/pprof_test.go
package debug

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
)

func TestRegisterRoutesDoesNothingWhenDisabled(t *testing.T) {
	e := echo.New()
	RegisterRoutes(e, false)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/debug/pprof/", nil)
	e.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestRegisterRoutesMountsPprofWhenEnabled(t *testing.T) {
	e := echo.New()
	RegisterRoutes(e, true)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/debug/pprof/", nil)
	e.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "profile")
}
