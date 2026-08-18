---
date: 2026-08-18
status: accepted
---

# Migrate golangci-lint from v1 to v2

## Decision

Bumped `golangci-lint` from `v1.64.8` to `v2.12.2` (Makefile `GOLANGCI_LINT_VERSION`, `ci.yml` lint step), migrated `.golangci.yml` to the v2 config schema via the official `golangci-lint migrate` command, and accepted the dependabot bump of `golangci-lint-action` to `v9` alongside `actions/checkout`, `actions/setup-go`, `actions/upload-artifact` to `v7`.

## Rationale

- `golangci-lint-action@v7+` only supports golangci-lint v2 — the open dependabot PR bumping the action to `v9` was failing CI because the repo was still pinned to `v1.64.8` with a v1-schema config.
- golangci-lint v1 no longer receives new features; staying on it would permanently block that action from updating and force manually overriding every future dependabot bump for it.
- Verified locally: `golangci-lint linters` confirms all 12 previously-enabled linters (including the ones the v1→v2 config migration silently dropped from the explicit `enable` list because they are on by default in v2: `errcheck`, `gosimple`, `staticcheck`, `govet`, `ineffassign`, `unused`) are still active. `make lint` and `make build` pass with no new issues.

## Limitations

- `run.timeout` (previously `5m`) has no v2 equivalent — v2 disables the internal timeout by default. The GitHub Actions job-level `timeout-minutes: 10` remains as the outer bound.
- Did not verify `make test` on this change — the checked-out repo's `database/migrations/` has no actual migration files yet (unrelated pre-existing state), so the test job's migration step cannot run locally.

## When to revisit

- If a future golangci-lint v2.x release changes default-enabled linters again, re-run `golangci-lint linters` to confirm the effective set still matches intent.
