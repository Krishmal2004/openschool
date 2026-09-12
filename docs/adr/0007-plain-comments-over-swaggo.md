# 0007. Plain one-line doc comments in handlers instead of swaggo annotations

**Status:** Accepted

## Context

Every handler method carried a full `swaggo` annotation block
(`@Summary`/`@Description`/`@Tags`/`@Param`/`@Success`/`@Failure`/`@Router`/
etc.), 8–14 lines each, used to generate `backend/docs/swagger.json`
(served at `/swagger/index.html` in development). Much of that duplicated
what the handler's name, route registration, and request/response types
already convey.

## Decision

Replace every swaggo block with one plain doc-comment line per handler.
`backend/docs/swagger.json` (and the generated `docs.go`/`swagger.yaml`)
stay committed as a snapshot from before this change, but are no longer
regenerated from source - `swag init` requires at least `@Router` to
include a handler, so re-running it today would drop most endpoints.

## Consequences

- Source is shorter and easier to scan - a handler's comment is a
  sentence, not an annotation block.
- `/swagger/index.html` now serves a frozen, increasingly-stale snapshot,
  not live-generated docs. Anyone wanting current Swagger docs needs to
  reintroduce annotations (and regenerate) for the endpoints they care
  about.
- Don't add swaggo annotations back to individual handlers piecemeal - a
  half-annotated handler set is worse than none. Revisit this ADR first.
