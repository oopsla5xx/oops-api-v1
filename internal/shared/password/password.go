package password

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// bcryptCost is fixed above bcrypt.DefaultCost (10) per current OWASP guidance for
// commodity hardware — see design.md of the user-registration-api change.
const bcryptCost = 12

const minLength = 10

var ErrTooShort = errors.New("password must be at least 10 characters")

// Validate enforces the registration password policy: a minimum length.
func Validate(plain string) error {
	if len(plain) < minLength {
		return ErrTooShort
	}
	return nil
}

func Hash(plain string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcryptCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func Verify(hash, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}
