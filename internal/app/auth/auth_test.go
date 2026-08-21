// internal/app/auth/auth_test.go
package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DylanBergmann2502/go-maleficent/internal/app/auth"
	"github.com/DylanBergmann2502/go-maleficent/internal/app/auth/handlers"
	"github.com/DylanBergmann2502/go-maleficent/internal/app/auth/services"
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/openapi"
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/testutil"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type userResponse struct {
	Data struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	} `json:"data"`
	Meta responseMeta `json:"meta"`
}

type errorResponse struct {
	Error struct {
		Code    string         `json:"code"`
		Message string         `json:"message"`
		Details map[string]any `json:"details"`
	} `json:"error"`
	Meta responseMeta `json:"meta"`
}

type responseMeta struct {
	RequestID string `json:"request_id"`
	Timestamp string `json:"timestamp"`
}

func newAuthAPI(t *testing.T) (*echo.Echo, *gorm.DB) {
	t.Helper()

	db := testutil.NewDatabase(t)
	tx := testutil.Begin(t, db)
	service := services.NewUserService(tx, validator.New())
	handler := handlers.NewUserHandler(service)

	e := echo.New()
	e.Use(middleware.RequestID())
	humaAPI := openapi.New(e)
	auth.RegisterRoutes(humaAPI, handler)

	return e, tx
}

func TestCreateUserEndpointReturnsCreatedUser(t *testing.T) {
	e, _ := newAuthAPI(t)
	requestBody := `{"email":"user@example.com","password":"correct horse battery staple"}`
	request := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString(requestBody))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderXRequestID, "request-123")
	recorder := httptest.NewRecorder()

	e.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusCreated, recorder.Code)
	var response userResponse
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	assert.NotEmpty(t, response.Data.ID)
	assert.Equal(t, "user@example.com", response.Data.Email)
	assert.Equal(t, "request-123", response.Meta.RequestID)
	assert.NotEmpty(t, response.Meta.Timestamp)
}

func TestCreateUserEndpointReturnsBadRequestForMalformedJSON(t *testing.T) {
	e, _ := newAuthAPI(t)
	request := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString(`{"email":`))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()

	e.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	var response errorResponse
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	assert.Equal(t, "bad_request", response.Error.Code)
	assert.NotEmpty(t, response.Error.Message)
	assert.NotEmpty(t, response.Meta.RequestID)
}

func TestCreateUserEndpointReturnsConflictForRegisteredEmail(t *testing.T) {
	e, db := newAuthAPI(t)
	require.NoError(t, db.Exec(
		"INSERT INTO users (id, email, password_hash, created_at, updated_at) VALUES (?, ?, ?, NOW(), NOW())",
		uuid.New(),
		"registered@example.com",
		"argon2-hash",
	).Error)

	requestBody := `{"email":"registered@example.com","password":"correct horse battery staple"}`
	request := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString(requestBody))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()

	e.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusConflict, recorder.Code)
	var response errorResponse
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	assert.Equal(t, "conflict", response.Error.Code)
	assert.Equal(t, "user with this email already exists", response.Error.Message)
}
