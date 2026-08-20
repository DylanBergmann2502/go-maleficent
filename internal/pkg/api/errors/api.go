// internal/pkg/api/errors/api.go
package errors

import (
	stderrors "errors"
	"net/http"

	apperrors "github.com/DylanBergmann2502/go-maleficent/internal/pkg/errors"
)

// Error is an API presentation error. It contains HTTP-specific information
// and must not be returned by application services.
type Error struct {
	Status  int
	Code    string
	Message string
	Details any
}

// FromError translates an application error into an API error.
func FromError(err error) *Error {
	var appErr *apperrors.ApplicationError
	if !stderrors.As(err, &appErr) {
		return &Error{
			Status:  http.StatusInternalServerError,
			Code:    "internal_server_error",
			Message: "Internal server error",
		}
	}

	status, code := translate(err)
	return &Error{
		Status:  status,
		Code:    code,
		Message: appErr.Message,
		Details: appErr.Details,
	}
}

func translate(err error) (int, string) {
	switch {
	case stderrors.Is(err, apperrors.ErrBadRequestCategory):
		return http.StatusBadRequest, "bad_request"
	case stderrors.Is(err, apperrors.ErrValidationCategory):
		return http.StatusUnprocessableEntity, "validation_failed"
	case stderrors.Is(err, apperrors.ErrUnauthorizedCategory):
		return http.StatusUnauthorized, "unauthorized"
	case stderrors.Is(err, apperrors.ErrForbiddenCategory):
		return http.StatusForbidden, "forbidden"
	case stderrors.Is(err, apperrors.ErrNotFoundCategory):
		return http.StatusNotFound, "not_found"
	case stderrors.Is(err, apperrors.ErrConflictCategory):
		return http.StatusConflict, "conflict"
	case stderrors.Is(err, apperrors.ErrUnavailableCategory):
		return http.StatusServiceUnavailable, "service_unavailable"
	default:
		return http.StatusInternalServerError, "internal_server_error"
	}
}
