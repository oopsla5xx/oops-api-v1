package handler

import (
	"errors"
	"net/http"

	"github.com/oopsla5xx/oops-api-v1/internal/modules/identity/domain"
	"github.com/oopsla5xx/oops-api-v1/internal/shared/constants"
	app_errors "github.com/oopsla5xx/oops-api-v1/internal/shared/errors"
	"github.com/oopsla5xx/oops-api-v1/internal/shared/password"
)

func toResponse(u domain.User) userResponse {
	return userResponse{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: u.DeletedAt,
	}
}

func mapError(err error) error {
	switch {
	case errors.Is(err, domain.ErrEmailTaken):
		return app_errors.Wrap(err, constants.ErrEmailAlreadyTaken, "email already registered", http.StatusConflict).
			WithField("email")
	case errors.Is(err, domain.ErrUsernameTaken):
		return app_errors.Wrap(err, constants.ErrUsernameAlreadyTaken, "username already taken", http.StatusConflict).
			WithField("username")
	case errors.Is(err, password.ErrTooShort):
		return app_errors.Wrap(err, constants.ErrPasswordTooWeak, "password does not meet strength requirements", http.StatusBadRequest).
			WithField("password")
	default:
		return err
	}
}
