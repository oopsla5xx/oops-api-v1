# CLAUDE.md

This file provides guidance to Claude Code when working with code in this repository.

## Where to read

| Need to know | Read this |
|---|---|
| Commands (build, test, migrate, lint…) | `README.md` |
| Module structure, dependency-direction rules, layer responsibilities | `internal/modules/AGENTS.md` — hard constraints, not suggestions |
| Testing conventions, available test helpers | `internal/tests/AGENTS.md` |
| App-level wiring, shared packages, SQL pipeline | `docs/architecture/wiring.md` |
| Security checklist for this stack | `docs/architecture/security.md` |
| Business/domain terms, cross-repo architecture, which source wins on conflict | `../CLAUDE.md` (workspace root) |

**When unsure about anything — read the relevant file above first. Do not guess.**

---

## Agent rules

### Before writing code
- Read `internal/modules/AGENTS.md` before touching any module structure, layer, or cross-module wiring
- Read `internal/tests/AGENTS.md` before writing any test
- Read `docs/architecture/security.md` before writing authentication, authorization, input handling, secrets, or other security-sensitive code

### After making changes
| Changed | Run |
|---|---|
| Any repository interface | `make generate` |
| Any `.sql` query file | `make sqlc` |
| Any handler or route | `make docs` |
| Any `.go` file | `make fmt` then `make lint` |

### Hard constraints
- Never push to `main` directly
- Never edit a migration file that has already run in production — add a new one
- Never hardcode values — use `internal/shared/constants`
- Never write raw JSON error responses in handlers — use `response.Error(c, err)`
- Never add default values to env var config — every var must be `required`
- Never add a new dependency without recording the decision in the `design.md` of the OpenSpec change introducing it (there is no standalone ADR folder in this workspace — see `../docs/agent-context/authority.md`)

### When adding an env var
1. Add it to `.env.example` with a placeholder or example value
2. Add it to `.env.development` and `.env.test`
3. Add it to the CI workflow env block in `.github/workflows/ci.yml`
4. Add a `required` struct tag in `internal/config/config.go` — never add a default value

### Definition of done
A task is not done until:
- `make test` passes
- `make lint` passes
- `make build` succeeds
- All changed files went through `make fmt`
