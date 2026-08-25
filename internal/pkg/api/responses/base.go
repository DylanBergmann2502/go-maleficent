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
	RequestID  string          `json:"request_id"`
	Timestamp  string          `json:"timestamp" format:"date-time"`
	Pagination *PaginationData `json:"pagination,omitempty"`
}

// PaginationData contains metadata for a paginated collection response.
type PaginationData struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalCount int64 `json:"total_count"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
	NextPage   *int  `json:"next_page"`
	PrevPage   *int  `json:"prev_page"`
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
