# oops-api-v1

Backend API for **Oops** — an AI-native Software Development Workspace.

---

## Tech Stack

| | |
|---|---|
| Language | Go 1.26+ |
| Framework | Gin v1.12 |
| Database | PostgreSQL via pgx v5 (pgxpool) + sqlc |
| Cache | Redis via go-redis v9 |
| Migration | goose |
| Logging | zap |
| Validation | go-playground/validator v10 |
| API Docs | swaggo/swag + gin-swagger |

---

## Development Setup

**Prerequisites:** Go 1.26+, Docker, `make`

```bash
# 1. Install dependencies
go mod download

# 2. Start dev infrastructure (Postgres :5432, Redis :6379)
make docker-up

# 3. Set up environment
cp .env.example .env.development
# Edit .env.development — all variables are required, no defaults

# 4. Apply migrations
make migrate-up ENV=development

# 5. Seed data (optional)
make seed ENV=development

# 6. Run with hot reload
make dev
```

Server runs at `http://localhost:8080`.

---

## Commands

```bash
# Development
make dev              # hot reload via Air
make run              # build + run once
make build            # compile binary into ./bin/

# Quality
make test             # all tests with race detector
make test-cover       # tests + coverage report
make lint             # golangci-lint
make fmt              # goimports + gofmt

# Database
make migrate-up ENV=development
make migrate-down ENV=development
make migrate-status ENV=development
make seed ENV=development

# Code generation
make sqlc             # SQL queries → Go (run after editing database/queries/)
make generate         # mocks via mockery (run after changing a repository interface)
make docs             # Swagger docs

# Infrastructure
make docker-up        # start dev infra (Postgres :5432, Redis :6379)
make docker-down
make test-up          # start isolated test infra (Postgres :5433, Redis :6380)
make test-down
```

Run a single test:
```bash
go test -run TestFunctionName ./internal/modules/...
```

Integration tests (require real DB):
```bash
make test-up
DATABASE_DSN=postgres://oops:oops@localhost:5433/oops_test?sslmode=disable \
REDIS_ADDR=localhost:6380 REDIS_DB=0 make test
make test-down
```

---

## Environment

All variables in `.env.example` are required — `config.Load()` returns an error immediately if any variable is missing or malformed. Copy `.env.example` to `.env.development` and fill in the values. Never commit `.env.*` files.

---

## API

```
GET /api/v1/health     → {"status":"ok","service":"oops-api-v1","version":"<version>"}
GET /swagger/index.html
```

Version is injected at build time via `ldflags` — never hardcoded.

---

## CI

Three jobs run on every push and pull request:

| Job | What it does |
|-----|---|
| `test` | Starts isolated containers, runs migrations, `make test-cover COVER_MIN=90` |
| `lint` | `golangci-lint v1.64.8` |
| `build` | `make build` |

---

## Architecture & Conventions

- Architecture and module rules → `.ai/context/architecture.md` + `internal/modules/AGENTS.md`
- Coding conventions → `.ai/context/conventions.md`
- Testing conventions → `.ai/context/testing-conventions.md`
