package command

import (
	"context"
	"strings"

	"github.com/oopsla5xx/oops-api-v1/internal/modules/identity/domain"
	"github.com/oopsla5xx/oops-api-v1/internal/modules/identity/domain/port"
	"github.com/oopsla5xx/oops-api-v1/internal/shared/id"
	"github.com/oopsla5xx/oops-api-v1/internal/shared/password"
)

type CreateUserInput struct {
	FirstName string
	LastName  string
	Username  string
	Email     string
	Password  string
}

type CreateUserCommand struct {
	repo port.UserRepository
}

func NewCreateUserCommand(repo port.UserRepository) *CreateUserCommand {
	return &CreateUserCommand{
		repo: repo,
	}
}

func (c *CreateUserCommand) Execute(ctx context.Context, input CreateUserInput) (domain.User, error) {
	if err := password.Validate(input.Password); err != nil {
		return domain.User{}, err
	}

	hashedPassword, err := password.Hash(input.Password)
	if err != nil {
		return domain.User{}, err
	}

	user := domain.User{
		ID:           id.New(),
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		Username:     strings.ToLower(input.Username),
		Email:        strings.ToLower(input.Email),
		PasswordHash: hashedPassword,
	}

	return c.repo.Create(ctx, user)
}
