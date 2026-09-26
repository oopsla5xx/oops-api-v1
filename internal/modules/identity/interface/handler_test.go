package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oopsla5xx/oops-api-v1/internal/modules/identity/application/command"
	"github.com/oopsla5xx/oops-api-v1/internal/modules/identity/infrastructure/postgres"
	user_handler "github.com/oopsla5xx/oops-api-v1/internal/modules/identity/interface"
	"github.com/oopsla5xx/oops-api-v1/internal/tests"
	"github.com/oopsla5xx/oops-api-v1/internal/tests/factory"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	m.Run()
}

func requestBody(t *testing.T, u factory.UserInput) *bytes.Reader {
	t.Helper()
	body, err := json.Marshal(map[string]string{
		"username": u.Username,
		"email":    u.Email,
		"password": u.Password,
	})
	require.NoError(t, err)
	return bytes.NewReader(body)
}

type errorEnvelope struct {
	Success bool `json:"success"`
	Error   struct {
		Code  string `json:"code"`
		Field string `json:"field"`
	} `json:"error"`
}

func decodeError(t *testing.T, w *httptest.ResponseRecorder) errorEnvelope {
	t.Helper()
	var body errorEnvelope
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	return body
}

func TestHandler_CreateUser(t *testing.T) {
	pool := tests.NewTestDB(t)
	t.Cleanup(func() { tests.Truncate(t, pool, "users") })

	repo := postgres.NewUserRepository(pool)
	h := user_handler.NewHandler(command.NewCreateUserCommand(repo))

	r := gin.New()
	h.Register(&r.RouterGroup)

	post := func(u factory.UserInput) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/auth/register", requestBody(t, u))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		return w
	}

	t.Run("creates user and returns 201 with no password/hash/token in the response", func(t *testing.T) {
		w := post(factory.NewUser())

		require.Equal(t, http.StatusCreated, w.Code)
		var body map[string]any
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
		data, ok := body["data"].(map[string]any)
		require.True(t, ok, "expected a data object in the response")
		for _, forbidden := range []string{"password", "password_hash", "access_token", "refresh_token", "token"} {
			_, present := data[forbidden]
			assert.Falsef(t, present, "response must not contain %q", forbidden)
		}
		assert.Contains(t, data, "id")
		assert.Contains(t, data, "email")
	})

	t.Run("returns 409 EMAIL_ALREADY_TAKEN when email already taken", func(t *testing.T) {
		existing := factory.NewUser()
		require.Equal(t, http.StatusCreated, post(existing).Code)

		dup := existing
		dup.Username = factory.NewUser().Username

		w := post(dup)
		require.Equal(t, http.StatusConflict, w.Code)
		body := decodeError(t, w)
		assert.Equal(t, "EMAIL_ALREADY_TAKEN", body.Error.Code)
		assert.Equal(t, "email", body.Error.Field)
	})

	t.Run("returns 409 USERNAME_ALREADY_TAKEN when username already taken", func(t *testing.T) {
		existing := factory.NewUser()
		require.Equal(t, http.StatusCreated, post(existing).Code)

		dup := existing
		dup.Email = factory.NewUser().Email

		w := post(dup)
		require.Equal(t, http.StatusConflict, w.Code)
		body := decodeError(t, w)
		assert.Equal(t, "USERNAME_ALREADY_TAKEN", body.Error.Code)
		assert.Equal(t, "username", body.Error.Field)
	})

	t.Run("returns 409 EMAIL_ALREADY_TAKEN for a duplicate differing only by case", func(t *testing.T) {
		existing := factory.NewUser()
		require.Equal(t, http.StatusCreated, post(existing).Code)

		dup := existing
		dup.Email = strings.ToUpper(existing.Email)
		dup.Username = factory.NewUser().Username

		w := post(dup)
		require.Equal(t, http.StatusConflict, w.Code)
		body := decodeError(t, w)
		assert.Equal(t, "EMAIL_ALREADY_TAKEN", body.Error.Code)
	})

	t.Run("returns 400 on missing required field", func(t *testing.T) {
		u := factory.NewUser()
		u.Email = ""

		assert.Equal(t, http.StatusBadRequest, post(u).Code)
	})

	t.Run("returns 400 INVALID_EMAIL_FORMAT for a malformed email", func(t *testing.T) {
		u := factory.NewUser()
		u.Email = "not-an-email"

		w := post(u)
		require.Equal(t, http.StatusBadRequest, w.Code)
		body := decodeError(t, w)
		assert.Equal(t, "INVALID_EMAIL_FORMAT", body.Error.Code)
		assert.Equal(t, "email", body.Error.Field)
	})

	t.Run("returns 400 INVALID_USERNAME_FORMAT for an invalid username", func(t *testing.T) {
		u := factory.NewUser()
		u.Username = "Bad Name!"

		w := post(u)
		require.Equal(t, http.StatusBadRequest, w.Code)
		body := decodeError(t, w)
		assert.Equal(t, "INVALID_USERNAME_FORMAT", body.Error.Code)
		assert.Equal(t, "username", body.Error.Field)
	})

	t.Run("returns 400 PASSWORD_TOO_WEAK for a too-short password", func(t *testing.T) {
		u := factory.NewUser()
		u.Password = "short1!"

		w := post(u)
		require.Equal(t, http.StatusBadRequest, w.Code)
		body := decodeError(t, w)
		assert.Equal(t, "PASSWORD_TOO_WEAK", body.Error.Code)
		assert.Equal(t, "password", body.Error.Field)
	})

	t.Run("old /users path no longer routed", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/users", requestBody(t, factory.NewUser()))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
