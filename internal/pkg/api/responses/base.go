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
