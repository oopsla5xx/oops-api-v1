---
date: 2026-08-31
status: accepted
---

# Pre-push quality gate via Lefthook

## Decision

Added `github.com/evilmartians/lefthook` (v2.1.12) as a dev-only tool dependency.
`lefthook.yml` runs `make lint` then `make test` on `git push`. Devs enable it once via `make setup`.

## Rationale

- Before this, `make test`/`make lint` were only ever run manually or by CI — nothing stopped a broken push from reaching GitHub.
- Lefthook is a single static Go binary (no Node/Python runtime needed), fits a Go-only repo, and its config is committed to the repo like Husky's is in `oops-web-v1`.
- `go vet` was deliberately left out of the hook: `golangci-lint`'s default linter set already runs `govet`, so a separate `go vet` step would just re-run the same check (see `.golangci.yml` — no `default: none`, so the default linters stay enabled).
- `make test-cover` (Docker containers + migrations + 90% coverage gate) was deliberately left out of the hook — too slow for every push and requires Docker to be running locally. That full gate still runs in CI (`.github/workflows/ci.yml`) and is enforced there via required branch protection checks.

## Limitations

- `git push --no-verify` bypasses this hook entirely — it is a local convenience/early-warning layer, not a real gate. The real gate is GitHub branch protection requiring the CI jobs to pass before merge (see `scripts/setup-branch-protection.sh`).
- A dev who never runs `make setup` gets no local hook at all; nothing currently enforces that step.

## When to revisit

- If pre-push ever feels too slow, move `test` to a `pre-commit`-only lint pass and keep full `test` on push, or vice versa.
- If integration tests (real DB) are ever fast/reliable enough to run without Docker (e.g. via `sqlmock` or an in-memory Postgres), reconsider adding `make test-cover` back into the pre-push hook.
