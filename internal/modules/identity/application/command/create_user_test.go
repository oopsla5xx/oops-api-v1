package command_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/oopsla5xx/oops-api-v1/internal/modules/identity/application/command"
	"github.com/oopsla5xx/oops-api-v1/internal/modules/identity/domain"
	"github.com/oopsla5xx/oops-api-v1/internal/modules/identity/domain/port/mocks"
	"github.com/oopsla5xx/oops-api-v1/internal/tests/factory"
)

func TestCreateUserCommand_Execute(t *testing.T) {
	in := factory.NewUser()
	input := command.CreateUserInput{
		FirstName: in.FirstName,
		LastName:  in.LastName,
		Username:  in.Username,
		Email:     in.Email,
		Password:  in.Password,
	}

	tests := []struct {
		name    string
		mockFn  func(repo *mocks.MockUserRepository)
		wantErr error
	}{
		{
			name: "creates user with a generated id and hashed password",
			mockFn: func(repo *mocks.MockUserRepository) {
				repo.EXPECT().
					Create(mock.Anything, mock.MatchedBy(func(u domain.User) bool {
						return u.ID != uuid.Nil &&
							u.PasswordHash != "" &&
							u.PasswordHash != input.Password &&
							u.Email == input.Email
					})).
					Return(domain.User{ID: uuid.New(), Email: input.Email}, nil)
			},
		},
		{
			name: "returns domain error when email already taken",
			mockFn: func(repo *mocks.MockUserRepository) {
				repo.EXPECT().Create(mock.Anything, mock.Anything).Return(domain.User{}, domain.ErrEmailTaken)
			},
			wantErr: domain.ErrEmailTaken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockUserRepository(t)
			tt.mockFn(repo)

			cmd := command.NewCreateUserCommand(repo)
			_, err := cmd.Execute(context.Background(), input)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}
