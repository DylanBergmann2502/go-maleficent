// internal/pkg/errors/base_errors.go
package errors

import stderrors "errors"

// Application categories are surface-independent error values. They are the
// Go equivalent of tagged application reasons such as :conflict and
// :not_found; API, MCP, and CLI surfaces can translate them independently.
var (
	ErrBadRequestCategory   = stderrors.New("bad_request")
	ErrValidationCategory   = stderrors.New("validation")
	ErrUnauthorizedCategory = stderrors.New("unauthorized")
	ErrForbiddenCategory    = stderrors.New("forbidden")
	ErrNotFoundCategory     = stderrors.New("not_found")
	ErrConflictCategory     = stderrors.New("conflict")
	ErrUnavailableCategory  = stderrors.New("unavailable")
	ErrInternalCategory     = stderrors.New("internal")
)

// ApplicationError contains application semantics without HTTP or other
// presentation-layer details.
type ApplicationError struct {
	category error
	Code     string
	Message  string
	Details  any
	Cause    error
}

func (e *ApplicationError) Error() string {
	return e.Message
}

// Unwrap exposes both the application category and underlying cause so
// errors.Is/errors.As work for either semantic or technical inspection.
func (e *ApplicationError) Unwrap() error {
	if e.Cause != nil {
		return stderrors.Join(e.category, e.Cause)
	}
	return e.category
}

// NewApplicationError creates an application error with optional details.
func NewApplicationError(category error, code string, message string, details any) *ApplicationError {
	return &ApplicationError{
		category: category,
		Code:     code,
		Message:  message,
		Details:  details,
	}
}

// WrapApplicationError creates an application error while retaining the
// underlying cause for logging and errors.Is/errors.As checks.
func WrapApplicationError(cause error, category error, code string, message string, details any) *ApplicationError {
	return &ApplicationError{
		category: category,
		Code:     code,
		Message:  message,
		Details:  details,
		Cause:    cause,
	}
}

var (
	ErrBadRequest   = NewApplicationError(ErrBadRequestCategory, "bad_request", "bad request", nil)
	ErrUnauthorized = NewApplicationError(ErrUnauthorizedCategory, "unauthorized", "unauthorized", nil)
	ErrForbidden    = NewApplicationError(ErrForbiddenCategory, "forbidden", "forbidden", nil)
	ErrNotFound     = NewApplicationError(ErrNotFoundCategory, "not_found", "not found", nil)
)
