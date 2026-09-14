# Backend Modular Monolith Refactor Progress

> Temporary tracking document for the incremental backend architecture refactor.
> This file records what has been migrated, what remains, and the rules being
> followed while the migration is in progress.

## 1. Refactor objective

Move the backend from a horizontally organized structure:

```text
routes -> handlers -> services -> repositories -> sqlc
```

to a pragmatic modular monolith where each business capability owns its:

- HTTP route registration
- Request and response types
- Application/use-case logic
- Repository adapter
- Tests

The application remains one Go process with one PostgreSQL database. This is
not a microservice split.

## 2. Current architecture

```text
cmd/api
  -> internal/app              composition root
      -> internal/modules      migrated business modules
          -> module handlers
          -> module use cases
          -> repository.go     sqlc adapter boundary
              -> db/sqlc
      -> internal/routes        temporary migration bridge
          -> legacy handlers/services/repositories
```

### Boundary rules

1. Module handlers must not import `internal/repositories` or `db/sqlc`.
2. Module use cases must not import `internal/repositories` or `db/sqlc`.
3. A module may import `db/sqlc` only from a file named `repository.go`.
4. sqlc rows must be mapped to module-owned types at the repository boundary.
5. New legacy service-layer sqlc imports are prohibited.
6. Cross-module capabilities must use narrow interfaces from `internal/ports`.
7. The composition root in `internal/app` wires modules and compatibility
   dependencies together.
8. Existing API paths, request payloads, response payloads, business rules,
   and authorization behavior must remain compatible unless explicitly
   changed.

## 3. Completed work

| Status | Area | Completed work | Current location |
|---|---|---|---|
| DONE | Composition root | Centralized router groups and module wiring | `internal/app/app.go` |
| DONE | Route migration bridge | Grouped remaining legacy route registration by business area | `internal/routes/modules.go` |
| DONE | Architecture guard | Added automated dependency-boundary and sqlc-debt checks | `internal/architecture/dependencies_test.go` |
| DONE | Architecture documentation | Documented the modular-monolith direction and dependency rules | `ARCHITECTURE.md` |
| DONE | HTTP binding | Added shared strict JSON binding infrastructure | `internal/platform/httpx` |
| DONE | Identity role seam | Added provider-neutral role resolution | `internal/identity/roles.go` |
| DONE | Identity `/me` | Migrated the current-user endpoint | `internal/modules/identity` |
| DONE | School profile | Migrated school profile create, read, and update | `internal/modules/school` |
| DONE | Academic years | Migrated academic-year CRUD and current-year operations | `internal/modules/school` |
| DONE | Grades | Migrated grade CRUD | `internal/modules/school` |
| DONE | Terms | Migrated term CRUD and current-term operations | `internal/modules/school` |
| DONE | Houses | Migrated house CRUD, balancing, reassignment, and audited member changes | `internal/modules/school` |
| DONE | Subjects | Migrated subject CRUD | `internal/modules/academics` |
| DONE | Streams | Migrated streams and stream groups | `internal/modules/academics` |
| DONE | Timetable settings | Migrated timetable settings upsert and read | `internal/modules/timetable` |
| DONE | Classrooms | Migrated classroom CRUD and lab-subject validation | `internal/modules/timetable` |
| DONE | Subject-period requirements | Migrated requirement upsert, listing, deletion, and scheduling validation | `internal/modules/timetable` |
| DONE | Teacher availability | Migrated teacher availability create, list, and delete | `internal/modules/timetable` |
| DONE | Grade sections | Migrated grade-section CRUD, grade assignment, section heads, and period grids | `internal/modules/timetable` |
| DONE | Period generation | Migrated settings-based period generation, fallback generation, and regeneration | `internal/modules/timetable` |
| DONE | Timetable entries | Migrated timetable entry listing, draft-only save, and entry clearing | `internal/modules/timetable` |
| DONE | Timetable validation | Migrated conflict, availability, assignment, requirement, and reviewer validation | `internal/modules/timetable` |
| DONE | Timetable CRUD | Migrated timetable create, copy, revision, read, listing, deletion, and archive endpoints | `internal/modules/timetable` |
| DONE | Timetable status history | Migrated timetable status-history reads | `internal/modules/timetable` |
| DONE | Timetable review workflow | Migrated submit, approve, reject, publish, and reviewer queue endpoints | `internal/modules/timetable` |
| DONE | Timetable portal views | Migrated teacher schedule, student timetable, and published class timetable endpoints | `internal/modules/timetable` |
| DONE | Timetable generation | Migrated grade-section automatic generation and draft entry persistence | `internal/modules/timetable` |
| DONE | Curriculum mediums | Migrated medium CRUD endpoints and reference-safe deletion | `internal/modules/curriculum` |
| DONE | Curriculum levels and groups | Migrated level, selection-group, group-subject, duplicate, delete, and tree endpoints | `internal/modules/curriculum` |
| DONE | Curriculum presets | Migrated preview and transactional, idempotent preset seeding | `internal/modules/curriculum` |
| DONE | Classes | Migrated class CRUD, assignments, subject-teacher qualification, and enrollment endpoints | `internal/modules/academics` |
| DONE | Enrollment routes | Migrated validation, submission, locking, deletion, and enrollment reporting routes | `internal/modules/academics` |
| DONE | Promotion route boundary | Migrated promotion preview and commit HTTP ownership into the academics module | `internal/modules/academics` |
| DONE | Promotions | Migrated preview, ranking, target validation, and transactional bulk assignment | `internal/modules/academics` |
| DONE | Term-mark route boundary | Migrated mark-entry and mark-report HTTP ownership into the academics module | `internal/modules/academics` |
| DONE | Term marks | Migrated mark entry, authorization, enrollment checks, reporting, and deletion into the academics module | `internal/modules/academics` |
| DONE | Module tests | Added focused unit tests for migrated business rules and adapters | Module `*_test.go` files |
| DONE | Verification | `go test ./...`, `go vet ./...`, `go build ./...`, architecture checks, and `git diff --check` pass | Backend repository |

## 4. Migration progress

The architecture guard originally tracked 35 legacy service files importing
sqlc. Nineteen have now been migrated out of the legacy service layer.

| Measure | Current status |
|---|---:|
| Original legacy sqlc service files | 35 |
| Migrated legacy sqlc service files | 19 |
| Remaining legacy sqlc service files | 16 |
| Foundation and composition work | DONE |
| Feature migration estimate | Approximately 50% |

> Note: the exact service-file debt is the authoritative metric. Run
> `rg -l 'db/sqlc' internal/services | sort` from `backend/` to inspect it.

## 5. Remaining feature migrations

| Status | Domain | Remaining work | Notes |
|---|---|---|---|
| TODO | Authentication | Login/setup lifecycle, password lifecycle, and related identity operations | ThunderID remains the identity provider |
| TODO | Identity reconciliation | Admin reconciliation and identity-provider/local-user consistency operations | Keep provider calls behind `internal/identity` |
| TODO | School setup | Initial setup workflow and setup status | Spans school and academic configuration |
| TODO | Curriculum | None for the core curriculum configuration endpoints | Levels, groups, subjects, tree, and mediums are migrated |
| TODO | Curriculum presets | None | Preview and transactional, idempotent preset seeding are migrated |
| TODO | Classes | None for the class configuration endpoints | Classes, assignments, subject-teacher qualification, and enrollment are migrated |
| TODO | Enrollments | Migrate the remaining student self-service dependency | Public enrollment routes are migrated; student self-service still consumes the compatibility service |
| TODO | Promotions | None | Preview and transactional assignment are migrated |
| TODO | Term marks | None for the mark workflows | Reports still use the compatibility repository until Reports is migrated |
| TODO | Students | Student profile CRUD and lifecycle | Must preserve ThunderID provisioning and rollback behavior |
| TODO | Teachers | Teacher profile CRUD and subject assignment | Must preserve ThunderID provisioning and rollback behavior |
| TODO | Guardians | Guardian profiles and student relationships | Includes parent access behavior |
| TODO | Non-academic staff | Staff profile and house operations | House reassignment currently remains in the legacy service |
| TODO | Student portfolio | Portfolio records and student access | Requires narrow student-access contracts |
| TODO | Student self-service | Student-facing profile, timetable, and academic endpoints | Some endpoints currently depend on legacy timetable repositories |
| TODO | Parent self-service | Parent-facing student and timetable endpoints | Some endpoints currently depend on legacy repositories |
| TODO | Teacher self-service | Teacher-facing timetable and workload endpoints | Depends on timetable engine and teacher module |
| TODO | Student attendance | Attendance sessions and student attendance records | Includes lock and correction rules |
| TODO | Staff attendance | Staff attendance records and reporting | Includes date-based operations |
| TODO | Notifications | Notifications, unread counts, and delivery behavior | Includes timetable and people dependencies |
| TODO | Positions | Teacher positions and scoped notifications | Check existing ADR before changing semantics |
| TODO | Section heads | Section-head assignments and access rules | Related to grade sections and teacher roles |
| TODO | Prefects | Prefect assignments and student leadership | Academic-year scoped |
| TODO | Societies | Society management and membership | Includes student access |
| TODO | Reports | Report export and report-facing queries | Preserve existing output formats |
| TODO | Audit | Audit-log service and endpoints | Cross-cutting capability; should become a narrow port |
| TODO | Dashboard | Dashboard aggregates and house distributions | Depends on people and attendance data |
| TODO | Search | Cross-entity search | Should be isolated behind a search capability |
| TODO | Automation | Jobs, scheduler, and admin job endpoints | Keep process lifecycle ownership in `internal/app` |

## 6. Temporary compatibility pieces

These are intentionally retained until their consumers are migrated:

- Legacy timetable repositories used by attendance, notifications, parent and
  student views, teacher self-service, and jobs.
- Legacy grade-section repository used by attendance, notifications, parent and
  student views, teacher self-service, and jobs.
- Legacy classroom repository used for subject-specific lab-room lookup.
- Legacy subject and grade repositories used by curriculum preset seeding and
  other unmigrated services.
- Legacy school repository used by current-year and school-type lookups.
- Legacy term repository used by report export.
- Legacy people services that still handle profile provisioning and lifecycle
  operations.
- Legacy route registration grouped behind `internal/routes/modules.go`.

These are migration debt, not new architecture targets. No new feature should
be added to them unless it is required to keep an unmigrated consumer working.

## 7. Verification checklist

Run from `backend/`:

```bash
GOCACHE=/tmp/openschool-go-cache go test ./...
GOCACHE=/tmp/openschool-go-cache go vet ./...
GOCACHE=/tmp/openschool-go-cache go build ./...
git diff --check
```

The architecture test also verifies:

- No handler imports repositories or sqlc.
- No new service-layer sqlc dependency is introduced.
- Modules do not import legacy handlers, services, or repositories.
- Module sqlc access is limited to `repository.go`.

## 8. Recommended migration order

1. Complete the remaining timetable engine dependencies and move timetable
   generation, review, and publication into the Timetable module.
2. Migrate Curriculum and Curriculum Presets as one related aggregate.
3. Migrate Classes, Enrollments, Promotions, and Term Marks as an academic
   workflow group.
4. Migrate Students, Teachers, Guardians, Non-academic Staff, and portfolios
   as the People module.
5. Migrate Student, Parent, and Teacher self-service endpoints.
6. Migrate Student and Staff Attendance.
7. Migrate Notifications, Audit, Reports, Dashboard, Search, and Automation.
8. Remove compatibility repositories and the legacy route bridge.
9. Run full integration/API tests and update this document before deleting it.

## 9. Definition of complete

The refactor is complete when:

- Every feature endpoint is registered by a business module.
- No legacy handler or service owns active feature behavior.
- No service-layer file imports `db/sqlc`.
- sqlc is accessed only by repository adapters.
- Cross-module dependencies use narrow interfaces.
- Compatibility repositories and route bridges are removed.
- Unit, integration, API, vet, build, and architecture checks pass.
- This temporary progress document is no longer needed and can be deleted.
