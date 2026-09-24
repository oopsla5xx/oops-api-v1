package newrelic

import (
	"fmt"
	"strings"

	"github.com/newrelic/go-agent/v3/newrelic"

	"github.com/oopsla5xx/oops-api-v1/internal/config"
)

// excludedAttributes are never sent to New Relic even if a future agent
// version starts capturing them by default.
var excludedAttributes = []string{
	"request.headers.authorization",
	"request.headers.cookie",
}

// NewApplication creates the New Relic agent application for HTTP-tier
// (nrgin) instrumentation. It returns a nil *Application, with no error,
// when New Relic is disabled — callers must treat a nil Application as a
// no-op; nrgin.Middleware, nrgin.Transaction, and Application.Shutdown all
// handle a nil Application/Transaction safely.
func NewApplication(cfg *config.NewRelicConfig) (*newrelic.Application, error) {
	if !cfg.Enabled {
		return nil, nil
	}

	app, err := newrelic.NewApplication(configOptions(cfg)...)
	if err != nil {
		return nil, fmt.Errorf("create new relic application: %w", err)
	}

	return app, nil
}

func configOptions(cfg *config.NewRelicConfig) []newrelic.ConfigOption {
	return []newrelic.ConfigOption{
		newrelic.ConfigAppName(cfg.AppName),
		newrelic.ConfigLicense(cfg.LicenseKey),
		newrelic.ConfigEnabled(cfg.Enabled),
		newrelic.ConfigCodeLevelMetricsEnabled(cfg.CodeLevelMetricsEnabled),
		newrelic.ConfigCodeLevelMetricsPathPrefixes(cfg.CodeLevelMetricsPathPrefix),
		withLabels(cfg.Labels),
		withExcludedAttributes(excludedAttributes...),
	}
}

func withLabels(raw string) newrelic.ConfigOption {
	return func(nrCfg *newrelic.Config) {
		nrCfg.Labels = parseLabels(raw)
	}
}

func withExcludedAttributes(attrs ...string) newrelic.ConfigOption {
	return func(nrCfg *newrelic.Config) {
		nrCfg.Attributes.Exclude = append(nrCfg.Attributes.Exclude, attrs...)
	}
}

// parseLabels parses the "Type1:value1;Type2:value2" format documented at
// https://docs.newrelic.com/docs/apm/agents/manage-apm-agents/app-naming/labels-categories-organize-your-apps-servers/.
func parseLabels(raw string) map[string]string {
	labels := make(map[string]string)
	for pair := range strings.SplitSeq(raw, ";") {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		k, v, ok := strings.Cut(pair, ":")
		if !ok {
			continue
		}
		labels[strings.TrimSpace(k)] = strings.TrimSpace(v)
	}
	return labels
}
