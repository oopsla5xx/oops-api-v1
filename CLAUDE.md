# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Where to read

| Need to know | Read this |
|---|---|
| Commands (build, test, migrate, lint…) | `README.md` |
| Architecture, module layout, wiring, data flow | `.ai/context/architecture.md` |
| Coding conventions, error handling, env vars, SQL rules | `.ai/context/conventions.md` |
| Testing patterns, mocks, factory, coverage | `.ai/context/testing-conventions.md` |
| Current task status, what's in progress | `.ai/status.md` |

**When unsure about anything — read the relevant file above first. Do not guess.**

---

## Agent rules

### Before writing code
- Read `.ai/context/architecture.md` before touching any module structure or wiring
- Read `.ai/context/conventions.md` before writing handlers, error handling, or env vars
- Read `.ai/context/testing-conventions.md` before writing any test

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
- Never add a new dependency without recording it in `.ai/decisions/`

### When adding an env var
Follow the 4-step checklist in `.ai/context/conventions.md` — all 4 places must be updated together.

### Definition of done
A task is not done until:
- `make test` passes
- `make lint` passes
- `make build` succeeds
- All changed files went through `make fmt`
