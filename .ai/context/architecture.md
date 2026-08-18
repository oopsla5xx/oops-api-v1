# Architecture

## Pattern

**Clean Architecture + DDD + Modular Monolith.**

```
Interface Layer   (HTTP handler, Router, Middleware)
      ↓
Application Layer (Use Case: Command / Query)
      ↓
Domain Layer      (Entity, Value Object, Repository Interface, Domain Service)

Infrastructure Layer  ← implements Domain interfaces (Dependency Inversion)
```

Module structure, dependency rules, layer responsibilities, wiring, cross-module contracts, and the new module checklist are all defined as hard constraints in:

> **`internal/modules/AGENTS.md`** — read this before touching any module.

---

## Application-level wiring

| File | Responsibility |
|------|---------------|
| `internal/app/dependency.go` | construct all infra + module deps, bind cross-module ports |
| `internal/app/router.go` | register all routes |
| `internal/app/server.go` | HTTP server start / graceful shutdown |

Config is loaded once in `main.go` via `config.Load()` using `caarlos0/env/v11`. Every env var is `required` — no silent defaults. See `internal/config/config.go`.

---

## Key shared packages

| Package | Purpose |
|---------|---------|
| `internal/shared/errors` | `AppError`, sentinel errors (`ErrNotFound`, `ErrBadRequest`, `ErrUnauthorized`, `ErrForbidden`), `errors.Wrap(err, code, status)` |
| `internal/shared/response` | `response.OK/Created/Error/NoContent` — all handlers must use these |
| `internal/shared/constants` | API paths, error codes, timeouts, Redis keys, env names |
| `internal/shared/middleware` | CORS, recovery, request ID, request logger, timeout |
| `internal/tests` | `tests.NewTestDB(t)` (auto-skips if `DATABASE_DSN` unset), `tests.Truncate(t, pool, "table")` |
| `internal/tests/factory` | gofakeit-backed test factories with functional override pattern |

---

## SQL data pipeline

```
database/queries/<name>.sql
    → make sqlc
    → internal/infrastructure/database/sqlc/  (generated — do not edit)
    → used by infrastructure/postgres/<name>_repository.go
```

---

## Version injection

Version, commit, and build time are injected at build time via `ldflags` in `make build`. Exposed at `GET /api/v1/health` and sourced from `internal/shared/version/`.
