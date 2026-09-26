# oops-api-v1

Backend API for **Oops** — an AI-native Software Development Workspace.

---

## Tech Stack

|            |                                        |
| ---------- | -------------------------------------- |
| Language   | Go 1.26+                               |
| Framework  | Gin v1.12                              |
| Database   | PostgreSQL via pgx v5 (pgxpool) + sqlc |
| Cache      | Redis via go-redis v9                  |
| Migration  | goose                                  |
| Logging    | zap                                    |
| Validation | go-playground/validator v10            |
| API Docs   | swaggo/swag + gin-swagger              |

---

## Quick Start

This repo is a submodule of [`oops-wiki-v1`](https://github.com/oopsla5xx/oops-wiki-v1). Clone the workspace rather than this repo alone.

```bash
git clone git@github.com:oopsla5xx/oops-wiki-v1.git
cd oops-wiki-v1
git submodule update --init
cd oops-api-v1
```

Install dependencies and configure local environment:

```bash
go mod download
make setup
cp .env.example .env.development
```

Start local infrastructure:

```bash
make dev-up
```

Apply database migrations:

```bash
make migrate-up ENV=development
```

Start the API:

```bash
make dev
```

API: `http://localhost:8080`

Health check: `GET /api/v1/health`

Swagger: `GET /swagger/index.html`

---

## Commands

### Development

```bash
make dev                 # hot reload
make run                 # build and run
make build               # build ./bin/oops-api-v1
```

### Database

```bash
make migrate-create NAME=create_users_table
make migrate-up ENV=development
make migrate-down ENV=development
make migrate-status ENV=development
make migrate-reset ENV=development
make seed ENV=development
```

Create a new migration:

```bash
make migrate-create NAME=create_users_table
```

`migrate-reset` rolls back all migrations and applies them again. Use it only for development/test.

### Testing & Quality

```bash
make test
make test-cover
make lint
make fmt
make tidy
```

Run a specific test:

```bash
go test -run TestFunctionName ./internal/modules/...
```

Integration tests:

```bash
make test-up
make test
make test-down
```

### Code Generation

```bash
make query-create NAME=users   # scaffold database/queries/users.sql
make sqlc                      # database/queries → Go
make generate                  # generate mocks
make docs                      # generate Swagger docs
```

### Infrastructure

```bash
make dev-up
make docker-down

make test-up
make test-down
```

---

## Database Migration Rules

* Create migrations with `make migrate-create`.
* Never modify a migration already applied to production.
* Create a new migration for schema changes.
* Use `migrate-reset` only for development/test.

---

## Project Documentation

* Architecture → `.ai/context/architecture.md`
* Coding conventions → `.ai/context/conventions.md`
* Testing conventions → `.ai/context/testing-conventions.md`
* Module rules → `internal/modules/AGENTS.md`

Read the relevant documentation before making changes.
