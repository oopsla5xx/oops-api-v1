APP_NAME   := oops-api-v1
CMD_DIR    := ./cmd/api
BIN        := ./bin/$(APP_NAME)
VERSION    := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT     := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS    := -s -w \
  -X github.com/oopsla5xx/oops-api-v1/internal/shared/version.Version=$(VERSION) \
  -X github.com/oopsla5xx/oops-api-v1/internal/shared/version.Commit=$(COMMIT) \
  -X github.com/oopsla5xx/oops-api-v1/internal/shared/version.BuildTime=$(BUILD_TIME)

ENV           ?= development
COVER_OUT     := coverage.out
COVER_HTML    := coverage.html
COVER_EXCLUDE := github.com/oopsla5xx/oops-api-v1/docs,\
                 github.com/oopsla5xx/oops-api-v1/internal/infrastructure/database/sqlc,\
                 github.com/oopsla5xx/oops-api-v1/cmd/api
COVER_MIN     ?= 0

# Tool versions — keep in sync with .github/workflows/ci.yml
GOLANGCI_LINT_VERSION := v2.12.2
MOCKERY_VERSION       := v2.53.6
GOOSE_VERSION         := v3.27.3
LEFTHOOK_VERSION      := v2.1.12

GOBIN := $(shell go env GOBIN)

INFRA_TF_DIR := ../oops-infra-v1/terraform/local
FLOCI_HEALTH_URL := http://localhost:4566/_localstack/health

.DEFAULT_GOAL := help

.PHONY: build run dev test test-cover cover-func lint fmt tidy clean \
        docker-up docker-down docker-logs \
        infra-apply infra-destroy dev-up \
        test-up test-down \
        migrate-up migrate-down migrate-status \
        seed sqlc docs generate setup help

## build: compile binary with version ldflags
build:
	go build -ldflags "$(LDFLAGS)" -o $(BIN) $(CMD_DIR)

## run: build and run the server (loads .env.ENV outside Docker)
run: build
	@set -a && . ./.env.$(ENV) && set +a && $(BIN)

## dev: hot reload via air (loads .env.ENV outside Docker)
dev:
	@test -f $(GOBIN)/air || go install github.com/air-verse/air@latest
	@set -a && . ./.env.$(ENV) && set +a && $(GOBIN)/air

## test: run all tests with race detector
test:
	go test -race -count=1 -covermode=atomic ./...

## test-cover: run tests, generate coverage report, enforce thresholds
## Usage: make test-cover              (report + warn only)
##        make test-cover COVER_MIN=90 (fail if total < 90%)
##        make test-cover COVER_PKG_WARN=80 (warn if any package < 80%)
COVER_PKG_WARN ?= 80
test-cover:
	go test -race -count=1 -covermode=atomic \
	  -coverprofile=$(COVER_OUT) \
	  -coverpkg=$$(go list ./... | grep -v -E 'docs|sqlc|cmd/api' | tr '\n' ',') \
	  ./...
	go tool cover -html=$(COVER_OUT) -o $(COVER_HTML)
	@go tool cover -func=$(COVER_OUT) | tail -1
	@if [ "$(COVER_MIN)" -gt 0 ]; then \
	  total=$$(go tool cover -func=$(COVER_OUT) | tail -1 | awk '{print $$3}' | tr -d '%'); \
	  echo "Coverage total: $${total}% (min: $(COVER_MIN)%)"; \
	  awk "BEGIN { exit ($${total} < $(COVER_MIN)) }"; \
	fi
	@echo "--- Per-package coverage (warn if < $(COVER_PKG_WARN)%) ---"; \
	go tool cover -func=$(COVER_OUT) | grep -v "^total:" | awk -F'\t' \
	  '{gsub(/%/,"",$$3); split($$1,a,"/"); pkg=a[1]; for(i=2;i<=length(a)-1;i++) pkg=pkg"/"a[i]; sum[pkg]+=$$3; cnt[pkg]++} \
	  END {for(p in sum){avg=sum[p]/cnt[p]; if(avg<$(COVER_PKG_WARN)) printf "  WARNING: %s %.1f%% (< $(COVER_PKG_WARN)%%)\n",p,avg}}'

## cover-func: show per-function coverage breakdown
cover-func:
	@go tool cover -func=$(COVER_OUT)

## lint: run golangci-lint
lint:
	@$(GOBIN)/golangci-lint version 2>/dev/null | grep -q "version $(GOLANGCI_LINT_VERSION:v%=%) " || \
	  go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	$(GOBIN)/golangci-lint run ./...

## fmt: format all Go files (gofmt + goimports)
fmt:
	@test -f $(GOBIN)/goimports || go install golang.org/x/tools/cmd/goimports@latest
	@find . -name "*.go" -not -path "*/vendor/*" | xargs $(GOBIN)/goimports -w -local github.com/oopsla5xx/oops-api-v1
	go fmt ./...

## tidy: tidy go.mod
tidy:
	go mod tidy

## clean: remove build artifacts
clean:
	rm -rf ./bin $(COVER_OUT) $(COVER_HTML)

## docker-up: start Floci (AWS emulator, backs RDS/ElastiCache/S3) — run infra-apply after
docker-up:
	docker compose -f ../oops-infra-v1/docker/dev/compose.yaml up -d --wait

## docker-down: stop Floci (graceful — safe on its own; infra-destroy first only matters if you force-kill instead)
docker-down:
	docker compose -f ../oops-infra-v1/docker/dev/compose.yaml down

## docker-logs: follow container logs
docker-logs:
	docker compose -f ../oops-infra-v1/docker/dev/compose.yaml logs -f

## infra-apply: provision RDS/ElastiCache/S3 on Floci (endpoints are fixed — see .env.example)
infra-apply:
	@test -d $(INFRA_TF_DIR) || \
	  (echo "$(INFRA_TF_DIR) not found — clone via oops-wiki-v1 and run 'git submodule update --init'" >&2 && exit 1)
	@curl -sf $(FLOCI_HEALTH_URL) >/dev/null || \
	  (echo "Floci is not reachable — run 'make docker-up' first" >&2 && exit 1)
	terraform -chdir=$(INFRA_TF_DIR) init -input=false
	terraform -chdir=$(INFRA_TF_DIR) apply -auto-approve
	terraform -chdir=$(INFRA_TF_DIR) output

## infra-destroy: tear down RDS/ElastiCache/S3
infra-destroy:
	@curl -sf $(FLOCI_HEALTH_URL) >/dev/null || \
	  (echo "Floci is not reachable — run 'make docker-up' first" >&2 && exit 1)
	terraform -chdir=$(INFRA_TF_DIR) init -input=false
	terraform -chdir=$(INFRA_TF_DIR) destroy -auto-approve

## dev-up: one-command local setup — docker-up then infra-apply
dev-up: docker-up infra-apply

## test-up: start isolated test containers (mirrors CI — Postgres :5433, Redis :6380)
test-up:
	docker compose -f ../oops-infra-v1/docker/test/compose.yaml up -d --wait

## test-down: stop test containers
test-down:
	docker compose -f ../oops-infra-v1/docker/test/compose.yaml down

## migrate-up: run all pending migrations (ENV=development|test|production)
migrate-up:
	@test -f $(GOBIN)/goose || go install github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION)
	@set -a && . ./.env.$(ENV) && set +a && \
	  $(GOBIN)/goose -dir database/migrations postgres "$$DATABASE_DSN" up

## migrate-down: rollback last migration (ENV=development|test|production)
migrate-down:
	@test -f $(GOBIN)/goose || go install github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION)
	@set -a && . ./.env.$(ENV) && set +a && \
	  $(GOBIN)/goose -dir database/migrations postgres "$$DATABASE_DSN" down

## migrate-status: show migration status
migrate-status:
	@test -f $(GOBIN)/goose || go install github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION)
	@set -a && . ./.env.$(ENV) && set +a && \
	  $(GOBIN)/goose -dir database/migrations postgres "$$DATABASE_DSN" status

## seed: run seed files for ENV (default: development)
seed:
	@bash scripts/seed.sh $(ENV)

## sqlc: generate type-safe Go code from SQL queries
sqlc:
	@test -f $(GOBIN)/sqlc || go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	$(GOBIN)/sqlc generate

## docs: generate Swagger docs via swag
docs:
	@test -f $(GOBIN)/swag || go install github.com/swaggo/swag/cmd/swag@latest
	$(GOBIN)/swag init -g cmd/api/main.go -o docs

## generate: generate mocks via mockery
generate:
	@test -f $(GOBIN)/mockery || go install github.com/vektra/mockery/v2@$(MOCKERY_VERSION)
	$(GOBIN)/mockery

## setup: install Lefthook git hooks (run once after clone)
setup:
	@test -f $(GOBIN)/lefthook || go install github.com/evilmartians/lefthook/v2@$(LEFTHOOK_VERSION)
	$(GOBIN)/lefthook install

## help: show available targets
help:
	@echo "Usage: make <target> [ENV=development|test|production]"
	@echo ""
	@grep -E '^## ' Makefile | sed 's/## /  /'
