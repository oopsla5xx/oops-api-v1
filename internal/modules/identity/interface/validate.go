package handler

import (
	"net/http"
	"net/mail"
	"regexp"

	"github.com/oopsla5xx/oops-api-v1/internal/shared/constants"
	app_errors "github.com/oopsla5xx/oops-api-v1/internal/shared/errors"
)

var usernamePattern = regexp.MustCompile(`^[a-z0-9_-]{3,30}$`)

func validateRequest(req createUserRequest) *app_errors.AppError {
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return app_errors.New(constants.ErrInvalidEmailFormat, "invalid email format", http.StatusBadRequest).
			WithField("email")
	}
	if !usernamePattern.MatchString(req.Username) {
		return app_errors.New(constants.ErrInvalidUsernameFormat, "invalid username format", http.StatusBadRequest).
			WithField("username")
	}
	return nil
}
