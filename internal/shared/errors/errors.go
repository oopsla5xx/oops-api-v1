package errors

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/oopsla5xx/oops-api-v1/internal/shared/constants"
)

type AppError struct {
	Status  int
	Code    string
	Message string
	Err     error
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

func New(code, message string, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

func Wrap(err error, code, message string, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
		Err:     err,
	}
}

var (
	ErrNotFound            = New(constants.ErrNotFound, "resource not found", http.StatusNotFound)
	ErrBadRequest          = New(constants.ErrBadRequest, "bad request", http.StatusBadRequest)
	ErrUnauthorized        = New(constants.ErrUnauthorized, "unauthorized", http.StatusUnauthorized)
	ErrForbidden           = New(constants.ErrForbidden, "forbidden", http.StatusForbidden)
	ErrInternal            = New(constants.ErrInternalServer, "internal server error", http.StatusInternalServerError)
	ErrConflict            = New(constants.ErrConflict, "conflict", http.StatusConflict)
	ErrUnprocessableEntity = New(constants.ErrUnprocessableEntity, "unprocessable entity", http.StatusUnprocessableEntity)
	ErrValidation          = New(constants.ErrValidation, "validation error", http.StatusBadRequest)
	ErrGatewayTimeout      = New(constants.ErrGatewayTimeout, "gateway timeout", http.StatusGatewayTimeout)
	ErrTooManyRequests     = New(constants.ErrTooManyRequests, "too many requests", http.StatusTooManyRequests)
)

func IsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}
