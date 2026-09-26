package identity

import (
	"github.com/jackc/pgx/v5/pgxpool"

	identity_command "github.com/oopsla5xx/oops-api-v1/internal/modules/identity/application/command"
	identity_postgres "github.com/oopsla5xx/oops-api-v1/internal/modules/identity/infrastructure/postgres"
	identity_handler "github.com/oopsla5xx/oops-api-v1/internal/modules/identity/interface"
)

func New(pool *pgxpool.Pool) *identity_handler.Handler {
	repo := identity_postgres.NewUserRepository(pool)
	createUser := identity_command.NewCreateUserCommand(repo)
	return identity_handler.NewHandler(createUser)
}
