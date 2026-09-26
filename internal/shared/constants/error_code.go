package constants

const (
	// Common error codes
	ErrInternalServer      = "INTERNAL_SERVER_ERROR"
	ErrNotFound            = "RESOURCE_NOT_FOUND"
	ErrBadRequest          = "BAD_REQUEST"
	ErrUnauthorized        = "UNAUTHORIZED"
	ErrForbidden           = "FORBIDDEN"
	ErrValidation          = "VALIDATION_ERROR"
	ErrConflict            = "CONFLICT"
	ErrGatewayTimeout      = "GATEWAY_TIMEOUT"
	ErrTooManyRequests     = "TOO_MANY_REQUESTS"
	ErrUnprocessableEntity = "UNPROCESSABLE_ENTITY"

	// Identity module specific error codes
	ErrEmailAlreadyTaken     = "EMAIL_ALREADY_TAKEN"
	ErrUsernameAlreadyTaken  = "USERNAME_ALREADY_TAKEN"
	ErrInvalidEmailFormat    = "INVALID_EMAIL_FORMAT"
	ErrInvalidUsernameFormat = "INVALID_USERNAME_FORMAT"
	ErrPasswordTooWeak       = "PASSWORD_TOO_WEAK"
)
