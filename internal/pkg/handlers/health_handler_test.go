// internal/pkg/handlers/health_handler_test.go
package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/testutil"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthHandlerPingDoesNotRequireDatabase(t *testing.T) {
	for _, path := range []string{"/ping", "/livez"} {
		t.Run(path, func(t *testing.T) {
			e := echo.New()
			e.GET(path, NewHealthHandler(nil).Ping)
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, path, nil)

			e.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response map[string]any
			require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
			assert.Equal(t, "ok", response["data"].(map[string]any)["status"])
		})
	}
}

func TestHealthHandlerCheckVerifiesDatabase(t *testing.T) {
	database := testutil.NewDatabaseClient(t)
	e := echo.New()
	e.GET("/readyz", NewHealthHandler(database).Check)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/readyz", nil)

	e.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	var response map[string]any
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	data := response["data"].(map[string]any)
	assert.Equal(t, "ok", data["status"])
	assert.NotNil(t, data["database"])
}

func TestHealthHandlerCheckReturnsDatabaseStats(t *testing.T) {
	database := testutil.NewDatabaseClient(t)
	e := echo.New()
	e.GET("/health", NewHealthHandler(database).Check)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/health", nil)

	e.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	var response map[string]any
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	data := response["data"].(map[string]any)
	assert.Equal(t, "ok", data["status"])
	assert.Equal(t, "Go Maleficent API is running", data["message"])
	assert.NotNil(t, data["database"])
}
