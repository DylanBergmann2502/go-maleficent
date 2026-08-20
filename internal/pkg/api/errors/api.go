// internal/pkg/api/errors/api.go
package errors

import (
	"context"
	stderrors "errors"
	"net/http"
	"strings"

	apiresponses "github.com/DylanBergmann2502/go-maleficent/internal/pkg/api/responses"
	apperrors "github.com/DylanBergmann2502/go-maleficent/internal/pkg/errors"
	"github.com/danielgtaylor/huma/v2"
)

// Error is an API presentation error. It contains HTTP-specific information
// and must not be returned by application services.
type Error struct {
	Status  int
	Code    string
	Message string
	Details any
}

// HumaError is Huma's HTTP error representation backed by the public API
// error envelope.
type HumaError struct {
	status int
	apiresponses.ErrorResponse
}

// Error implements error.
func (e *HumaError) Error() string {
	return e.ErrorResponse.Error.Message
}

// GetStatus implements huma.StatusError.
func (e *HumaError) GetStatus() int {
	return e.status
}

// NewHumaError builds a Huma status error using the public API envelope.
func NewHumaError(ctx context.Context, status int, code string, message string, details any) huma.StatusError {
	if code == "" {
		code = codeForStatus(status)
	}

	return &HumaError{
		status: status,
		ErrorResponse: apiresponses.ErrorResponse{
			Error: apiresponses.ErrorPayload{
				Code:    code,
				Message: message,
				Details: detailsForStatus(status, details),
			},
			Meta: apiresponses.NewMetaDataFromContext(ctx),
		},
	}
}

// ToHumaError translates an application error into Huma's status-aware error
// type at the HTTP/OpenAPI boundary.
func ToHumaError(ctx context.Context, err error) error {
	apiErr := FromError(err)
	return NewHumaError(ctx, apiErr.Status, apiErr.Code, apiErr.Message, apiErr.Details)
}

func codeForStatus(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "bad_request"
	case http.StatusUnauthorized:
		return "unauthorized"
	case http.StatusForbidden:
		return "forbidden"
	case http.StatusNotFound:
		return "not_found"
	case http.StatusConflict:
		return "conflict"
	case http.StatusUnprocessableEntity:
		return "validation_failed"
	case http.StatusServiceUnavailable:
		return "service_unavailable"
	default:
		return "internal_server_error"
	}
}

func detailsForStatus(status int, details any) any {
	if status >= http.StatusInternalServerError {
		return nil
	}

	return details
}

// HumaDetails converts Huma's field-level errors into the API details map.
func HumaDetails(errs ...error) map[string][]string {
	details := make(map[string][]string)
	for _, err := range errs {
		if err == nil {
			continue
		}

		key := "base"
		message := err.Error()
		if detailer, ok := err.(huma.ErrorDetailer); ok {
			detail := detailer.ErrorDetail()
			key = strings.TrimPrefix(detail.Location, "body.")
			key = strings.TrimPrefix(key, "query.")
			key = strings.TrimPrefix(key, "path.")
			key = strings.TrimPrefix(key, "header.")
			if key == "" {
				key = "base"
			}
			message = detail.Message
		}
		details[key] = append(details[key], message)
	}

	if len(details) == 0 {
		return nil
	}

	return details
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
