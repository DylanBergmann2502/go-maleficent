// internal/pkg/api/responses/envelope.go
package responses

import (
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

// Success builds the standard envelope for a successful response.
func Success(c *echo.Context, data any) map[string]any {
	return map[string]any{
		"data": data,
		"meta": requestMeta(c),
	}
}

// Error builds the standard envelope for an error response.
// Details is omitted when nil.
func Error(c *echo.Context, code string, message string, details any) map[string]any {
	errorPayload := map[string]any{
		"code":    code,
		"message": message,
	}

	if details != nil {
		errorPayload["details"] = details
	}

	return map[string]any{
		"error": errorPayload,
		"meta":  requestMeta(c),
	}
}

func requestMeta(c *echo.Context) map[string]any {
	requestID, ok := c.Get("request_id").(string)
	if !ok || requestID == "" {
		requestID = c.Request().Header.Get(echo.HeaderXRequestID)
	}
	if requestID == "" {
		requestID = uuid.NewString()
	}

	return map[string]any{
		"request_id": requestID,
		"timestamp":  time.Now().UTC().Format(time.RFC3339Nano),
	}
}
