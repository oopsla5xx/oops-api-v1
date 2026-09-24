package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/oopsla5xx/oops-api-v1/internal/shared/middleware"
	"github.com/oopsla5xx/oops-api-v1/internal/shared/response"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	m.Run()
}

func TestRecovery(t *testing.T) {
	tests := []struct {
		name         string
		withNewRelic bool
		panicValue   any
	}{
		{
			name:         "recovers from a panic with no New Relic transaction present",
			withNewRelic: false,
			panicValue:   "boom",
		},
		{
			name:         "recovers from a panic and notices the error on the New Relic transaction",
			withNewRelic: true,
			panicValue:   "boom",
		},
		{
			name:         "recovers from a non-string panic value",
			withNewRelic: true,
			panicValue:   assertionError{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := zap.NewNop()

			r := gin.New()
			if tt.withNewRelic {
				app, err := newrelic.NewApplication(
					newrelic.ConfigAppName("recovery-test"),
					newrelic.ConfigLicense(strings.Repeat("a", 40)),
					newrelic.ConfigEnabled(false),
				)
				require.NoError(t, err)
				r.Use(nrgin.Middleware(app))
			}
			r.Use(middleware.Recovery(log))
			r.GET("/panic", func(c *gin.Context) {
				panic(tt.panicValue)
			})

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/panic", nil)
			require.NoError(t, err)

			assert.NotPanics(t, func() {
				r.ServeHTTP(w, req)
			})

			assert.Equal(t, http.StatusInternalServerError, w.Code)

			var body response.Response
			require.NoError(t, json.NewDecoder(w.Body).Decode(&body))
			assert.False(t, body.Success)
			require.NotNil(t, body.Error)
			assert.Equal(t, "INTERNAL_SERVER_ERROR", body.Error.Code)
		})
	}
}

type assertionError struct{}

func (assertionError) Error() string { return "assertion failed" }
