// internal/pkg/errors/errors.go
package errors

type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

var (
	ErrBadRequest     = NewAppError(400, "bad request")
	ErrUnauthorized   = NewAppError(401, "unauthorized")
	ErrForbidden      = NewAppError(403, "forbidden")
	ErrNotFound       = NewAppError(404, "not found")
	ErrConflict       = NewAppError(409, "conflict")
	ErrUnprocessable  = NewAppError(422, "unprocessable entity")
	ErrServiceUnavail = NewAppError(503, "service unavailable")
)
