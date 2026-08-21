// internal/pkg/api/responses/envelope_test.go
package responses

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMetaDataFromContextUsesRequestID(t *testing.T) {
	ctx := WithRequestID(context.Background(), "request-123")
	meta := NewMetaDataFromContext(ctx)

	assert.Equal(t, "request-123", meta.RequestID)
	_, err := time.Parse(time.RFC3339Nano, meta.Timestamp)
	require.NoError(t, err)
}

func TestNewMetaDataFromContextGeneratesRequestIDWhenMissing(t *testing.T) {
	meta := NewMetaDataFromContext(context.Background())

	assert.NotEmpty(t, meta.RequestID)
	_, err := time.Parse(time.RFC3339Nano, meta.Timestamp)
	require.NoError(t, err)
}

func TestSuccessAndErrorUseEchoRequestID(t *testing.T) {
	e := echo.New()
	request := httptest.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()
	ctx := e.NewContext(request, recorder)
	ctx.Set("request_id", "request-456")

	success := Success(ctx, map[string]string{"status": "ok"})
	assert.Equal(t, map[string]string{"status": "ok"}, success["data"])
	successMeta, ok := success["meta"].(MetaData)
	require.True(t, ok)
	assert.Equal(t, "request-456", successMeta.RequestID)

	err := Error(ctx, "conflict", "resource already exists", nil)
	assert.Equal(t, "conflict", err.Error.Code)
	assert.Equal(t, "resource already exists", err.Error.Message)
	assert.Equal(t, "request-456", err.Meta.RequestID)
}
