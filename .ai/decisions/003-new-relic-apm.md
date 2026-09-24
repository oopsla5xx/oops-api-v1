---
date: 2026-08-25
status: accepted
---

# New Relic APM: HTTP-tier only (nrgin)

## Decision

Added `github.com/newrelic/go-agent/v3` (v3.44.2) and its
`integrations/nrgin` package as direct dependencies. Instrumented only the
HTTP tier:

- `internal/infrastructure/newrelic` wraps `newrelic.NewApplication`,
  returning a nil `*Application` (no error) when `NEW_RELIC_ENABLED=false`.
- `nrgin.Middleware(app)` runs in `internal/app/router.go` right after
  `middleware.Recovery`, so the New Relic transaction covers the full
  middleware chain (request ID, CORS, logging, timeout, rate limiting) and
  is still caught by `Recovery` if it panics.
- `middleware.Recovery` reports recovered panics via
  `nrgin.Transaction(c).NoticeError(err)` so crashes show up in New Relic's
  Errors Inbox, not just as a logged 500.
- `Server.shutdown()` calls `newRelic.Shutdown(timeout)` last (after
  `db`/`redis` close) to flush any buffered data before the process exits.
- `NEW_RELIC_ENABLED` gates the whole thing per environment: `false` in
  development/test (no agent connection, no CI dependency on a real
  license key), `true` in staging/production. `NEW_RELIC_APP_NAME` is
  environment-suffixed (`oops-api-v1 (staging)`, `oops-api-v1
  (production)`) so each environment is a distinct APM entity.
- `Attributes.Exclude` explicitly lists `request.headers.authorization`
  and `request.headers.cookie` as defense-in-depth, even though the
  current agent version does not capture arbitrary headers by default
  (verified by reading `nrgin.go` and `attributes.go` in the vendored
  source — only a small fixed set of headers like content-type/host are
  ever turned into attributes).

## Rationale

- The project had zero APM/observability beyond structured `zap` logs —
  no visibility into response times, throughput, or error rates in
  production.
- `nrgin` is the officially maintained New Relic integration for Gin and
  needs no additional wrapping beyond what's implemented here.
- Scoping to HTTP-tier only keeps this change small and reviewable;
  Postgres (`nrpgx5`) and Redis (`nrredis-v9`) instrumentation are
  deliberately deferred (see Limitations) rather than bundled in, to avoid
  growing the dependency surface and blast radius of a single change.
- `caarlos0/env`'s `required` tag only checks that a variable is present
  in the environment, not that it's non-empty (confirmed by reading
  `env.go`) — combined with the go-agent's own `validate()` allowing an
  empty `License`/`AppName` when `Enabled` is false (confirmed by reading
  `config.go`), this let every `NEW_RELIC_*` var stay `required` (no
  silent defaults) while dev/test/CI carry explicit empty values instead
  of a fake license key.

## Limitations

- No database or cache instrumentation — slow queries and Redis latency
  are invisible to New Relic. Everything reported is HTTP-transaction
  level only.
- No custom segments/spans inside handlers or use cases; only the
  automatic transaction boundaries `nrgin.Middleware` creates.
- New Relic error reporting is wired only into the panic-recovery path
  (`middleware.Recovery`). Ordinary 4xx/5xx responses returned via
  `response.Error` are not separately reported to New Relic's error
  collector.
- The local `.env.development` already contained a real-looking license
  key before this change (untracked, gitignored — never committed). It
  was left as-is; `NEW_RELIC_ENABLED=false` still keeps it inert by
  default.

## When to revisit

- When query-level or cache-level latency needs to be visible in New
  Relic: add `nrpgx5` (Postgres) and/or `nrredis-v9` (Redis) as a
  follow-up decision.
- If handlers start doing enough non-trivial work that transaction-level
  timing isn't granular enough: add custom segments via
  `newrelic.FromContext`/`nrgin.Transaction` inside individual handlers.
- If `response.Error` should also report errors to New Relic beyond
  panics: revisit whether that belongs in `response.Error` itself or
  stays panic-only.
