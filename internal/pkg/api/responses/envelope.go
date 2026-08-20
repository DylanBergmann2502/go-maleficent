// internal/pkg/api/responses/envelope.go
package responses

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type requestIDContextKey struct{}

// WithRequestID stores the request ID in a request context for Huma handlers.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDContextKey{}, requestID)
}

// NewMetaDataFromContext builds response metadata for a Huma request.
func NewMetaDataFromContext(ctx context.Context) MetaData {
	requestID, _ := ctx.Value(requestIDContextKey{}).(string)
	if requestID == "" {
		requestID = uuid.NewString()
	}

	return MetaData{
		RequestID: requestID,
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
	}
}

// Success builds the standard envelope for a successful response.
func Success(c *echo.Context, data any) map[string]any {
	return map[string]any{
		"data": data,
		"meta": requestMeta(c),
	}
}

// NewMetaData builds response metadata for an Echo request.
func NewMetaData(c *echo.Context) MetaData {
	return requestMeta(c)
}

// Error builds the standard envelope for an error response.
// Details is omitted when nil.
func Error(c *echo.Context, code string, message string, details any) ErrorResponse {
	return ErrorResponse{
		Error: ErrorPayload{
			Code:    code,
			Message: message,
			Details: details,
		},
		Meta: requestMeta(c),
	}
}

func requestMeta(c *echo.Context) MetaData {
	requestID, ok := c.Get("request_id").(string)
	if !ok || requestID == "" {
		requestID = c.Request().Header.Get(echo.HeaderXRequestID)
	}
	if requestID == "" {
		requestID = uuid.NewString()
	}

	return MetaData{
		RequestID: requestID,
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
	}
}
