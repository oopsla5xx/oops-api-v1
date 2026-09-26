package factory

import (
	"strings"

	"github.com/brianvoe/gofakeit/v6"
)

type UserInput struct {
	Username string
	Email    string
	Password string
}

func NewUser(overrides ...func(*UserInput)) UserInput {
	u := UserInput{
		Username: strings.ToLower(gofakeit.Username()),
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, true, false, 12),
	}
	for _, o := range overrides {
		o(&u)
	}
	return u
}
