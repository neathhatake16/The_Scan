package apperrors

import "net/http"

// AppError is a typed, HTTP-aware error that can be returned from any layer
// and rendered directly by the response helper.
type AppError struct {
	Code    int    // HTTP status code
	Message string // user-facing message
	Err     error  // underlying cause (logged, never exposed to client)
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error { return e.Err }

// ── Constructors ──────────────────────────────────────────────

func BadRequest(msg string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: msg}
}

func Unauthorized(msg string) *AppError {
	return &AppError{Code: http.StatusUnauthorized, Message: msg}
}

func Forbidden(msg string) *AppError {
	return &AppError{Code: http.StatusForbidden, Message: msg}
}

func NotFound(msg string) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: msg}
}

func Conflict(msg string) *AppError {
	return &AppError{Code: http.StatusConflict, Message: msg}
}

func Internal(msg string, cause error) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: msg, Err: cause}
}
