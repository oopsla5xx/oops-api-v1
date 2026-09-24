package newrelic_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oopsla5xx/oops-api-v1/internal/config"
	nrapp "github.com/oopsla5xx/oops-api-v1/internal/infrastructure/newrelic"
)

func TestNewApplication(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *config.NewRelicConfig
		wantNil bool
	}{
		{
			name: "returns nil application when disabled",
			cfg: &config.NewRelicConfig{
				Enabled: false,
			},
			wantNil: true,
		},
		{
			name: "returns a usable application when enabled with a valid license key",
			cfg: &config.NewRelicConfig{
				Enabled:                    true,
				LicenseKey:                 strings.Repeat("a", 40),
				AppName:                    "oops-api-v1 (test)",
				CodeLevelMetricsEnabled:    true,
				CodeLevelMetricsPathPrefix: "github.com/oopsla5xx/oops-api-v1",
				Labels:                     "Env:test;Service:oops-api-v1",
			},
			wantNil: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := nrapp.NewApplication(tt.cfg)
			require.NoError(t, err)
			if tt.wantNil {
				assert.Nil(t, app)
			} else {
				assert.NotNil(t, app)
			}
		})
	}
}

func TestNewApplication_InvalidLicenseKey(t *testing.T) {
	cfg := &config.NewRelicConfig{
		Enabled:    true,
		LicenseKey: "too-short",
		AppName:    "oops-api-v1 (test)",
	}

	app, err := nrapp.NewApplication(cfg)

	assert.Error(t, err)
	assert.Nil(t, app)
}
