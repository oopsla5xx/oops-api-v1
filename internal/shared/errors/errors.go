package errors

import (
	"errors"
	"fmt"
	"net/http"
)

type AppError struct {
	Code       string
	Message    string
	StatusCode int
	Err        error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code, message string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

func Wrap(err error, code, message string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}

var (
	ErrNotFound     = New("RESOURCE_NOT_FOUND", "resource not found", http.StatusNotFound)
	ErrBadRequest   = New("BAD_REQUEST", "bad request", http.StatusBadRequest)
	ErrUnauthorized = New("UNAUTHORIZED", "unauthorized", http.StatusUnauthorized)
	ErrForbidden    = New("FORBIDDEN", "forbidden", http.StatusForbidden)
	ErrInternal     = New("INTERNAL_SERVER_ERROR", "internal server error", http.StatusInternalServerError)
	ErrConflict     = New("CONFLICT", "conflict", http.StatusConflict)
)

func IsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}
