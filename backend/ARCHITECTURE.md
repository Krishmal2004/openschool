# Backend architecture

OpenSchool is a modular monolith: one Go process and one PostgreSQL database,
with business capabilities separated into compile-time modules.

## Dependency direction

```text
cmd/api -> internal/app -> internal/modules
                              |-> handlers/use cases
                              |-> module contracts
                              `-> repository adapters -> db/sqlc -> PostgreSQL
```

- `internal/app` is the composition root and process-level module registry.
- `internal/modules/<name>` owns the HTTP, application, and persistence code for
  one business capability.
- Module handlers depend on use cases, never repositories or generated sqlc.
- sqlc types stay inside `repository.go` adapters and are mapped to module types.
- Cross-module dependencies use narrow consumer-owned interfaces and are wired
  by `internal/app`.
- `internal/platform` is reserved for shared technical infrastructure; business
  behavior must not be placed there.

## Incremental migration

The existing horizontal packages remain a migration bridge. New behavior must
be added to a module, and existing behavior moves one vertical slice at a time.
The architecture test rejects new legacy dependencies, all service-layer sqlc
imports, and any new file in the horizontal repository package. Its explicit
repository debt list is reduced whenever a remaining capability migrates.
