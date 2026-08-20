// internal/pkg/api/responses/base.go
package responses

import (
	"time"

	"github.com/google/uuid"
)

// BaseResponse contains common response fields.
type BaseResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// MetaData contains metadata included with every API response.
type MetaData struct {
	RequestID string `json:"request_id"`
	Timestamp string `json:"timestamp" format:"date-time"`
}

// ErrorPayload contains the public API error information.
type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// ErrorResponse is the standard public API error envelope.
type ErrorResponse struct {
	Error ErrorPayload `json:"error"`
	Meta  MetaData     `json:"meta"`
}
