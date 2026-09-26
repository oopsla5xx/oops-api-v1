package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	sqlcdb "github.com/oopsla5xx/oops-api-v1/internal/infrastructure/database/sqlc"
	"github.com/oopsla5xx/oops-api-v1/internal/modules/identity/domain"
)

const uniqueViolation = "23505"

type UserRepository struct {
	queries *sqlcdb.Queries
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		queries: sqlcdb.New(pool),
	}
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) (domain.User, error) {
	row, err := r.queries.CreateUser(ctx, sqlcdb.CreateUserParams{
		ID:           u.ID,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			switch pgErr.ConstraintName {
			case "users_username_unique":
				return domain.User{}, domain.ErrUsernameTaken
			case "users_email_unique":
				return domain.User{}, domain.ErrEmailTaken
			}
		}

		return domain.User{}, fmt.Errorf("create user: %w", err)
	}

	return toDomain(row), nil
}

func toDomain(u sqlcdb.User) domain.User {
	out := domain.User{
		ID:           u.ID,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt.Time,
		UpdatedAt:    u.UpdatedAt.Time,
	}
	if u.DeletedAt.Valid {
		out.DeletedAt = &u.DeletedAt.Time
	}
	return out
}
