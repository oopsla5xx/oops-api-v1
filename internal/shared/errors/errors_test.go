package errors_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	app_errors "github.com/oopsla5xx/oops-api-v1/internal/shared/errors"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name       string
		code       string
		message    string
		statusCode int
	}{
		{
			name:       "creates error with all fields set",
			code:       "TEST_CODE",
			message:    "test message",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "creates error with 500 status",
			code:       "INTERNAL_SERVER_ERROR",
			message:    "internal server error",
			statusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := app_errors.New(tt.code, tt.message, tt.statusCode)
			assert.Equal(t, tt.code, err.Code)
			assert.Equal(t, tt.message, err.Message)
			assert.Equal(t, tt.statusCode, err.StatusCode)
			assert.Nil(t, err.Err)
		})
	}
}

func TestWrap(t *testing.T) {
	original := errors.New("original error")

	tests := []struct {
		name        string
		wrapped     error
		wantWrapped bool
	}{
		{
			name:        "wraps original error and exposes it via Unwrap",
			wrapped:     original,
			wantWrapped: true,
		},
		{
			name:        "wraps nil error",
			wrapped:     nil,
			wantWrapped: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := app_errors.Wrap(tt.wrapped, "CODE", "message", http.StatusBadRequest)
			require.NotNil(t, err)
			if tt.wantWrapped {
				assert.True(t, errors.Is(err, original))
			} else {
				assert.Nil(t, err.Unwrap())
			}
		})
	}
}

func TestAppError_Error(t *testing.T) {
	tests := []struct {
		name    string
		err     *app_errors.AppError
		wantMsg string
	}{
		{
			name:    "returns message when no wrapped error",
			err:     app_errors.New("CODE", "something went wrong", http.StatusBadRequest),
			wantMsg: "something went wrong",
		},
		{
			name:    "includes wrapped error in message",
			err:     app_errors.Wrap(errors.New("db timeout"), "CODE", "database error", http.StatusInternalServerError),
			wantMsg: "database error: db timeout",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantMsg, tt.err.Error())
		})
	}
}

func TestIsAppError(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		wantMatch bool
	}{
		{
			name:      "detects AppError",
			err:       app_errors.New("CODE", "msg", http.StatusBadRequest),
			wantMatch: true,
		},
		{
			name:      "detects wrapped AppError",
			err:       app_errors.Wrap(app_errors.ErrNotFound, "CODE", "msg", http.StatusBadRequest),
			wantMatch: true,
		},
		{
			name:      "returns false for stdlib error",
			err:       errors.New("plain error"),
			wantMatch: false,
		},
		{
			name:      "returns false for nil",
			err:       nil,
			wantMatch: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appErr, ok := app_errors.IsAppError(tt.err)
			assert.Equal(t, tt.wantMatch, ok)
			if tt.wantMatch {
				assert.NotNil(t, appErr)
			} else {
				assert.Nil(t, appErr)
			}
		})
	}
}

func TestSentinelErrors(t *testing.T) {
	tests := []struct {
		name           string
		err            *app_errors.AppError
		wantCode       string
		wantStatusCode int
	}{
		{
			name:           "ErrNotFound has 404 status",
			err:            app_errors.ErrNotFound,
			wantCode:       "RESOURCE_NOT_FOUND",
			wantStatusCode: http.StatusNotFound,
		},
		{
			name:           "ErrBadRequest has 400 status",
			err:            app_errors.ErrBadRequest,
			wantCode:       "BAD_REQUEST",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "ErrUnauthorized has 401 status",
			err:            app_errors.ErrUnauthorized,
			wantCode:       "UNAUTHORIZED",
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name:           "ErrForbidden has 403 status",
			err:            app_errors.ErrForbidden,
			wantCode:       "FORBIDDEN",
			wantStatusCode: http.StatusForbidden,
		},
		{
			name:           "ErrInternal has 500 status",
			err:            app_errors.ErrInternal,
			wantCode:       "INTERNAL_SERVER_ERROR",
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:           "ErrConflict has 409 status",
			err:            app_errors.ErrConflict,
			wantCode:       "CONFLICT",
			wantStatusCode: http.StatusConflict,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantCode, tt.err.Code)
			assert.Equal(t, tt.wantStatusCode, tt.err.StatusCode)
		})
	}
}
