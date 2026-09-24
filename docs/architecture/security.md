# Security notes — oops-api-v1

For the generic pre-ship checklist, see `../../../docs/security-checklist.md` (workspace root). This file adds only what's specific to this stack — it does not repeat the root checklist.

## Go / this repo specifics

- **SQL**: no ORM, `sqlc`-generated queries only (see `docs/architecture/wiring.md`) — there is no raw string concatenation path for query input, so the `data-access` root checklist item is satisfied by construction as long as no handler builds SQL by hand.
- **Input validation**: request DTOs are validated in `interface/`, before reaching `application` — see `internal/modules/AGENTS.md` §3 for the layer boundary (`interface` must map to a DTO, never pass raw request data down).
- **Error exposure**: all handler errors go through `response.Error(c, err)` (see root `CLAUDE.md` hard constraints) — never construct a raw JSON error body, which is how internal details leak to clients.
- **Auth**: `oops-api` issues its own JWT (RS256 access token, Redis-backed revocable refresh token); `oops-agent` authenticates via a separate service-to-service token, not a user JWT. Full model: workspace root `docs/architecture/system-design/system-design.md` §4.
