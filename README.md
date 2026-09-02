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

This repo is a submodule of the [`oops-wiki-v1`](https://github.com/oopsla5xx/oops-wiki-v1) workspace. Shared dev/test containers live in the sibling `oops-infra-v1` submodule, so clone through the workspace rather than cloning this repo alone:

```bash
git clone git@github.com:oopsla5xx/oops-wiki-v1.git
cd oops-wiki-v1
git submodule update --init
cd oops-api-v1
```

```bash
go mod download
make setup                        # install git hooks
cp .env.example .env.development  # values already work for local dev, nothing to edit
make dev-up                       # start Floci + provision RDS/ElastiCache/S3 (~1-2 min first time)
make migrate-up ENV=development
make dev                          # hot reload, http://localhost:8080
```

`make dev-up` is the only infra command you need day to day. It starts
[Floci](https://github.com/floci-io/floci) (a local AWS emulator, see `oops-infra-v1/README.md`) and
provisions RDS/ElastiCache/S3 against it — mirrors production instead of running plain
Postgres/Redis containers. `make docker-down` stops it when you're done for the day.

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
make dev-up            # start Floci + provision RDS/ElastiCache/S3 (see oops-infra-v1/README.md)
make docker-down       # stop it
make test-up           # start isolated test infra (Postgres :5433, Redis :6380 — unaffected by Floci)
make test-down
```

`dev-up` is `docker-up` (start Floci) + `infra-apply` (provision) chained — those two, plus
`infra-destroy` (tear down RDS/ElastiCache/S3 without stopping Floci itself), are also available
individually for finer control; see the Makefile.

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

All variables in `.env.example` are required — `config.Load()` returns an error immediately if any variable is missing or malformed. Its values already work as-is for local dev (`DATABASE_DSN`/`REDIS_ADDR` point at the Floci-provisioned RDS/ElastiCache instance — fixed values, safe to hardcode since `oops-infra-v1/terraform/local/` only ever provisions exactly one of each, see ADR-0003). Copy it to `.env.development`; only edit it if you're changing something on purpose. Never commit `.env.*` files.

---

## API

```
GET /api/v1/health     → {"status":"ok","service":"oops-api-v1","version":"<version>"}
GET /swagger/index.html
```

Version is injected at build time via `ldflags` — never hardcoded.

---

## CI

Four jobs run on every push and pull request:

| Job | What it does |
|-----|---|
| `test` | Postgres/Redis via GitHub Actions `services:`, runs migrations, `make test-cover COVER_MIN=90` |
| `lint` | `golangci-lint v2.12.2` |
| `generate` | Runs `make generate`, fails if mocks are stale |
| `build` | `make build` |

---

## Architecture & Conventions

- Architecture and module rules → `.ai/context/architecture.md` + `internal/modules/AGENTS.md`
- Coding conventions → `.ai/context/conventions.md`
- Testing conventions → `.ai/context/testing-conventions.md`
