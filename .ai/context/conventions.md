# Conventions

## Environment variables

**Rule: no silent defaults.**

Every env var must be explicitly set in the `.env.*` file. Config uses `github.com/caarlos0/env/v11` with struct tags — variables tagged `required` cause `env.Parse` to return an error immediately if missing or malformed. There are no fallback default values in code.

**Correct:**
```go
type RedisConfig struct {
    Addr string `env:"REDIS_ADDR,required"`
    DB   int    `env:"REDIS_DB,required"`
    Password string `env:"REDIS_PASSWORD"` // intentionally optional — may be empty
}
```

**Incorrect:**
```go
v := os.Getenv("REDIS_ADDR")
if v == "" {
    v = "localhost:6379" // silent fallback — hides misconfiguration
}
```

**Exception:** `REDIS_PASSWORD` may be empty (Redis with no auth configured) — omit `required` tag.

All variables and their expected values are documented in `.env.example`. When adding a new env var:
1. Add it to `.env.example` with a placeholder or example value
2. Add it to `.env.development` and `.env.test`
3. Add it to the CI workflow env block in `.github/workflows/ci.yml`
4. Add a `required` struct tag in `internal/config/config.go` — never add a default value

## Error handling

- All errors go through `response.Error(c, err)` in handlers — never write raw JSON error responses
- Use `errors.Wrap(err, code, status)` to attach an error code and HTTP status to domain errors
- Use sentinel errors (`ErrNotFound`, `ErrBadRequest`, etc.) from `internal/shared/errors` for known failure cases

## Naming

- File names: `snake_case.go`
- Package names: lowercase singular noun (`handler`, `service`, `repository`) — no `utils`, `helpers`, `common`
- No hardcoded values: use `internal/shared/constants` for API paths, error codes, timeouts, Redis keys

## SQL

- No ORM — raw SQL via sqlc only
- Seed files must be idempotent (`INSERT ... ON CONFLICT DO NOTHING`)
- Never edit a migration file that has already run in production — always add a new migration
