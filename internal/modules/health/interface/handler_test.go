package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	health_query "github.com/oopsla5xx/oops-api-v1/internal/modules/health/application/query"
	health_handler "github.com/oopsla5xx/oops-api-v1/internal/modules/health/interface"
	"github.com/oopsla5xx/oops-api-v1/internal/shared/constants"
	"github.com/oopsla5xx/oops-api-v1/internal/shared/version"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	m.Run()
}

func TestHandler_Health(t *testing.T) {
	tests := []struct {
		name        string
		wantStatus  int
		wantStatus_ string
		wantService string
		wantVersion string
	}{
		{
			name:        "returns 200 with service liveness info",
			wantStatus:  http.StatusOK,
			wantStatus_: "ok",
			wantService: constants.AppName,
			wantVersion: version.Version,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			h := health_handler.NewHandler(health_query.NewHealthQuery())
			r.GET("/health", h.Health)

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/health", nil)
			require.NoError(t, err)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			var body map[string]string
			require.NoError(t, json.NewDecoder(w.Body).Decode(&body))
			assert.Equal(t, tt.wantStatus_, body["status"])
			assert.Equal(t, tt.wantService, body["service"])
			assert.Equal(t, tt.wantVersion, body["version"])
		})
	}
}

func TestHandler_Health_ResponseShape(t *testing.T) {
	t.Run("response contains exactly status, service, version fields", func(t *testing.T) {
		r := gin.New()
		h := health_handler.NewHandler(health_query.NewHealthQuery())
		r.GET("/health", h.Health)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/health", nil)
		require.NoError(t, err)
		r.ServeHTTP(w, req)

		var body map[string]interface{}
		require.NoError(t, json.NewDecoder(w.Body).Decode(&body))
		assert.Contains(t, body, "status")
		assert.Contains(t, body, "service")
		assert.Contains(t, body, "version")
		assert.Len(t, body, 3, "response must have exactly 3 fields")
	})
}
