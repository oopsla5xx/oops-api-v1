# Wiring, shared packages, and the SQL pipeline

Module structure and dependency-direction rules live in `internal/modules/AGENTS.md` — this doc does not repeat them. This covers the application-level composition root and cross-module shared packages instead.

## Application-level wiring

| File | Responsibility |
|------|---------------|
| `internal/app/dependency.go` | construct all infra + module deps, bind cross-module ports |
| `internal/app/router.go` | register all routes |
| `internal/app/server.go` | HTTP server start / graceful shutdown |

Config is loaded once in `main.go` via `config.Load()` using `caarlos0/env/v11`. Every env var is `required` — no silent defaults. See `internal/config/config.go`.

## Key shared packages

| Package | Purpose |
|---------|---------|
| `internal/shared/errors` | `AppError`, sentinel errors (`ErrNotFound`, `ErrBadRequest`, `ErrUnauthorized`, `ErrForbidden`), `errors.Wrap(err, code, status)` |
| `internal/shared/response` | `response.OK/Created/Error/NoContent` — all handlers must use these |
| `internal/shared/constants` | API paths, error codes, timeouts, Redis keys, env names |
| `internal/shared/middleware` | CORS, recovery, request ID, request logger, timeout |
| `internal/shared/version` | version/commit/build-time injected via `ldflags` in `make build`, exposed at `GET /api/v1/health` |
| `internal/tests` | `tests.NewTestDB(t)`, `tests.Truncate(t, pool, "table")` — see `internal/tests/AGENTS.md` for the full testing guide |
| `internal/tests/factory` | gofakeit-backed test factories with functional override pattern |

## SQL data pipeline

```
database/queries/<name>.sql
    → make sqlc
    → sqlc-generated code (see the Makefile `sqlc` target for the current output path — do not edit generated code by hand)
    → used by internal/infrastructure/database/postgres.go and related infrastructure code
```
