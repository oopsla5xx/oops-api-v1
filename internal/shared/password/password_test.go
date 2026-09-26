package password_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oopsla5xx/oops-api-v1/internal/shared/password"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		plain   string
		wantErr error
	}{
		{name: "accepts a 10-character password", plain: "xk7qz2mdrp"},
		{name: "accepts a strong long password", plain: "Tr0ub4dor&3xtra"},
		{name: "rejects a 9-character password", plain: "xk7qz2mdr", wantErr: password.ErrTooShort},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := password.Validate(tt.plain)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestHashAndVerify(t *testing.T) {
	hash, err := password.Hash("a-strong-password-1")
	require.NoError(t, err)
	assert.NotEqual(t, "a-strong-password-1", hash)

	assert.NoError(t, password.Verify(hash, "a-strong-password-1"))
	assert.Error(t, password.Verify(hash, "wrong-password"))
}
