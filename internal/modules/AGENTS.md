# AGENTS.md — internal/modules

## Purpose

Architecture contract for every module under `internal/modules/<n>/`. AI agents and humans
generating code here MUST follow these rules as hard constraints, not style suggestions.
Rationale, tutorials, and full code examples live in the wiki — this file only states decisions.

---

## 1. Architecture Rules (Dependency Direction)

- MUST: `domain` has zero internal package dependencies.
- SHOULD: `domain` use only the Go standard library.
- MAY: `domain` use a minimal, dependency-light value library (e.g. `uuid`) when the domain model
  requires it.
- MUST NOT: `domain` import framework, database, driver, transport, logging, configuration, or
  any infrastructure package.
- MUST: `application` depend on `domain` and approved shared technical packages only (e.g. clock,
  ID generator) — MUST NOT depend on `infrastructure` or `interface`.
- MUST: `application` reach external capabilities only through `domain` ports or its own
  `application/port` interfaces — never a concrete infrastructure type.
- MUST: `infrastructure` depend on `domain` and approved shared technical packages only; MUST NOT
  depend on `application` or `interface`. It implements `domain/port` interfaces only.
- MUST: `interface` depend on `application`; it MAY import `domain` types only for error matching
  or DTO mapping — importing a domain type is not the same as calling domain behavior (see §3).
- MUST: `module.go` MAY import all layers belonging to its own module; MUST NOT import another
  module.
- MUST: `internal/app` MAY import each module's `module.go` output and cross-module application
  contracts (`Service`, `application/port`); it is the only place cross-module dependencies are
  bound.
- `domain/port/` contains ports owned by the domain: persistence and any domain-required
  technical capability.
- `application/port/` contains consumer-owned ports for synchronous cross-module application
  calls only. MUST NOT represent persistence, caching, messaging, or infrastructure access.
- `application/port` interfaces MUST be implemented/adapted at the composition root
  (`internal/app`); the consuming module MUST NOT import the providing module to satisfy its own
  port.

### Dependency Matrix

| From \ To        | domain | application | infrastructure | interface |
|-------------------|:---:|:---:|:---:|:---:|
| domain             | — | ❌ | ❌ | ❌ |
| application        | ✅ | — | ❌ | ❌ |
| infrastructure     | ✅ | ❌ | — | ❌ |
| interface          | ⚠️ types only | ✅ | ❌ | — |
| module.go          | ✅ | ✅ | ✅ | ✅ (own module only) |

✅ allowed · ❌ forbidden · ⚠️ allowed only for error matching / DTO mapping

---

## 2. Module Structure

```
<module>/
  domain/
    entity.go, errors.go
    port/            persistence / domain-required technical ports
  application/
    query/           read use cases, no side effects
    command/         write use cases, side effects
    port/            consumer-owned ports for synchronous calls to other modules
    service.go        optional public application API for other modules (see §4)
  infrastructure/     implementations of domain/port (postgres/, mongo/, redis/)
  interface/           handler.go, request.go, response.go, mapper.go, handler_test.go
  module.go            module-level composition root (wires this module's own layers only)
```

- MUST create only the layers a module actually needs. A folder that would be empty MUST NOT
  exist — no placeholder `doc.go` just to keep four folders. `health` legitimately has no
  `domain/`.

---

## 3. Layer Responsibilities

### domain
- MUST own all entities, value objects, invariants, domain errors, and `domain/port` interfaces,
  named in business language (`FindByID`, `Save` — not `SelectByID`, `Insert`).

### application (CQRS convention)
- MUST split use cases into `query/` (read) and `command/` (write, owns all business operations
  that mutate state).
- `query` MUST NOT mutate domain state or perform business writes. Read-only caching/telemetry is
  allowed only when it does not alter business state.
- MUST accept stdlib `context.Context` in every exported method — MUST NOT accept `*gin.Context`
  or any framework type.
- MUST call only `domain` entities, `domain/port` interfaces, and its own `application/port`
  interfaces — MUST NOT call `infrastructure` or `interface` directly.
- SHOULD expose a small `Service` in `application/service.go` as a cross-module API only if
  another module genuinely needs synchronous calls into this one (§4). `Service` MUST expose
  business/application operations only (e.g. `GetUser(ctx, id) (UserView, error)`), never
  repositories, entities, or infrastructure implementations.

### infrastructure
- MUST implement `domain/port` interfaces only, using pgx/sqlc, mongo-driver, or go-redis.
- MUST translate driver-specific errors into domain errors before returning them.
- MUST NOT contain business rules.

### interface
- MUST call `application/query` and `application/command` only — MUST NOT call domain behavior,
  repository ports, or `infrastructure` directly.
- MUST NOT let `*gin.Context` cross into `application` or `domain`; convert to
  `c.Request.Context()` first.
- MUST NOT return domain entities directly as JSON — always map to a response DTO.

---

## 4. Cross-module Rules

- MUST NOT import another module's `domain`, `application`, or `infrastructure` package.
- `internal/shared` is reserved for genuinely generic technical primitives. MUST NOT contain
  module-specific business logic, entities, repository/port interfaces, or DTOs.
- The consuming module MUST declare the interface it needs from another module in its own
  `application/port/`.
- The composition root (`internal/app`) is the only place allowed to bind a consumer's
  `application/port` to another module's exposed `Service`. Modules MUST NOT wire each other
  directly.

---

## 5. Error Handling

- `domain` MUST define errors as sentinels or typed errors, not ad-hoc strings.
- `infrastructure` MUST wrap technical errors (`%w`) and translate driver errors
  (`pgx.ErrNoRows`, `mongo.ErrNoDocuments`) into domain errors.
- `interface` is the only layer allowed to map a domain error to an HTTP status code.
- `application` and `domain` MUST NOT reference HTTP status codes.

---

## 6. Wiring

- MUST: expose a `New(...)` constructor in `module.go` that constructs and returns the module's
  public entrypoint(s) (`Handler`, and `application.Service` only if another module depends on it).
- `module.go` is the **module-level composition root** — it wires this module's own internal
  layers only (domain, application, infrastructure, interface) and MUST NOT import another
  module.
- `internal/app` is the **application-level composition root** — it builds every module and wires
  cross-module dependencies (binding a consumer's `application/port` to another module's
  `Service`), then calls `.Register(...)` on each handler.
- MUST NOT wire cross-module dependencies inside a module itself — only in `internal/app`.

---

## 7. Testing

- Handler tests MUST use `httptest`, in an external `_test` package (black-box).
- Application (`query`/`command`) tests SHOULD use domain fakes/mocks, not a real DB.
- Infrastructure tests MAY use a real or containerized DB; MUST NOT run as part of the default
  unit test suite.

---

## 8. New Module Checklist

1. Decide which layers this module needs — create only those.
2. If the module has domain logic, define entities + `domain/port` interfaces in `domain/`.
3. Write `query`/`command` in `application/`, depending only on `domain` and its own
   `application/port`.
4. If persistence is needed, implement `domain/port` interfaces in `infrastructure/`.
5. Write `Handler`, DTOs, and `mapper.go` in `interface/`; handler calls `application` only.
6. Write `module.go`; expose `application.Service` only if another module genuinely needs it.
7. Wire any cross-module ports in `internal/app`, never inside the module.
8. Test the handler via `httptest`.
9. Do not create placeholder folders for unused layers.