---
date: 2026-08-11
status: accepted
---

# Rate Limiter: in-memory per-IP (golang.org/x/time/rate)

## Decision

Added `golang.org/x/time/rate` (v0.15.0) as a direct dependency.
Implemented `middleware.RateLimit(rps, burst)` — a per-IP token-bucket limiter backed by `sync.Map`.

## Rationale

- The project had no rate limiting, exposing the API to brute-force and abuse.
- `golang.org/x/time/rate` is the standard library extension for rate limiting; no additional third-party dep required.
- Per-IP limits are more granular than a global limiter and protect against single-source abuse.

## Limitations

- `sync.Map` entries are never evicted: each unique IP adds one entry for the process lifetime. Acceptable for normal traffic; becomes a memory concern only under sustained DDoS (thousands of unique IPs/hour).
- State is not shared across replicas. In a multi-replica deployment, each instance enforces limits independently — the effective global limit is `rps × replicas`.

## When to revisit

- When the service runs more than one replica in production.
- Replace with a Redis sliding-window implementation (e.g., using `INCR` + `EXPIRE` per IP+window key) to share state across replicas.
