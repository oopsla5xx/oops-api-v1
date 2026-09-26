package factory

import (
	"strings"

	"github.com/brianvoe/gofakeit/v6"
)

type UserInput struct {
	FirstName string
	LastName  string
	Username  string
	Email     string
	Password  string
}

func NewUser(overrides ...func(*UserInput)) UserInput {
	u := UserInput{
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Username:  strings.ToLower(gofakeit.Username()),
		Email:     gofakeit.Email(),
		Password:  gofakeit.Password(true, true, true, true, false, 12),
	}
	for _, o := range overrides {
		o(&u)
	}
	return u
}
