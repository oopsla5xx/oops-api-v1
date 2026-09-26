package port

import (
	"context"

	"github.com/oopsla5xx/oops-api-v1/internal/modules/identity/domain"
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) (domain.User, error)
}
