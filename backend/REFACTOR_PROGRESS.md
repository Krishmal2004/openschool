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
| DONE | Student route boundary | Migrated student administration HTTP ownership into the people module | `internal/modules/people` |
| DONE | Students | Migrated student profile CRUD, enrollment status, house assignment, identity provisioning, rollback, and lifecycle operations | `internal/modules/people` |
| DONE | Teacher read slice | Migrated teacher profile, subject, workload, and subject-qualified teacher reads into People-owned repository contracts | `internal/modules/people` |
| DONE | Teacher route boundary | Migrated all teacher administration HTTP route registration into the People module | `internal/modules/people` |
| DONE | Guardian read slice | Migrated guardian directory, child, and student-linked guardian reads into People-owned repository contracts | `internal/modules/people` |
| DONE | Guardian route boundary | Migrated guardian administration write route registration into the People module | `internal/modules/people` |
| DONE | Guardian parent-access contract | Migrated parent child listing and guardian ownership checks to a narrow People port | `internal/ports`, `internal/modules/people` |
| DONE | Guardian authentication contract | Migrated guardian NIC credential verification behind a People-owned authentication port | `internal/ports`, `internal/modules/people` |
| DONE | Guardian notification access | Migrated guardian notification history reads into a People-owned repository adapter | `internal/modules/people` |
| DONE | Non-academic staff | Migrated CRUD, employee numbering, employment status, and audited house assignment | `internal/modules/people` |
| DONE | Student portfolio service | Migrated progress reports, activities, leadership, awards, discipline, validation, and persistence | `internal/modules/people` |
| DONE | Student portfolio routes | Migrated all portfolio HTTP handlers and endpoint registration into People | `internal/modules/people` |
| DONE | Student attendance | Migrated sessions, records, authorization, 24-hour locking, correction auditing, absence notifications, self-service reads, and report reads | `internal/modules/attendance` |
| DONE | Staff attendance | Migrated daily marking, directories, monthly summaries, staff histories, and teacher self-service attendance | `internal/modules/attendance` |
| DONE | Notifications | Migrated composer authorization, recipient resolution, drafts, delivery, inbox, unread/archive state, and system delivery | `internal/modules/notifications` |
| DONE | Audit | Migrated append-only audit recording, admin audit-log reads, and the shared recorder port | `internal/modules/audit`, `internal/ports` |
| DONE | Leadership positions | Migrated Principal and Vice Principal appointments, scoped grants, rank resolution, auditing, and leadership overview | `internal/modules/leadership` |
| DONE | Section heads | Migrated grade/stream TIC assignment, listing, deletion, and leadership-scope resolution | `internal/modules/leadership` |
| DONE | Prefects | Migrated academic-year appointments, rank updates, archive years, student history, and routes | `internal/modules/studentleadership` |
| DONE | Societies | Migrated society CRUD, archives, rosters, TIC authorization, student memberships, and teacher self-service reads | `internal/modules/studentleadership` |
| DONE | School setup | Migrated setup status, first-admin registration, ThunderID provisioning, role assignment, and compensating rollbacks | `internal/modules/setup` |
| DONE | Report export | Migrated attendance and marks PDF exports, query adapters, column selection, and HTTP delivery | `internal/modules/reports` |
| DONE | Dead repository cleanup | Removed fourteen superseded School, Academics, Curriculum, People, and Timetable repository adapters | Module repository adapters |
| DONE | Repository debt guard | Added an exact allowlist for the nine active horizontal repository files so no new compatibility repository can be introduced | `internal/architecture/dependencies_test.go` |
| DONE | Module tests | Added focused unit tests for migrated business rules and adapters | Module `*_test.go` files |
| DONE | Verification | `go test ./...`, `go vet ./...`, `go build ./...`, architecture checks, and `git diff --check` pass | Backend repository |

## 4. Migration progress

The architecture guard originally tracked 35 legacy service files importing
sqlc. All thirty-five have now been migrated out of the legacy service layer.

| Measure | Current status |
|---|---:|
| Original legacy sqlc service files | 35 |
| Migrated legacy sqlc service files | 35 |
| Remaining legacy sqlc service files | 0 |
| Foundation and composition work | DONE |
| Feature migration estimate | Approximately 93% |

> Note: the exact service-file debt is the authoritative metric. Run
> `rg -l 'db/sqlc' internal/services | sort` from `backend/` to inspect it.

## 5. Remaining feature migrations

| Status | Domain | Remaining work | Notes |
|---|---|---|---|
| TODO | Authentication | Login/setup lifecycle, password lifecycle, and related identity operations | ThunderID remains the identity provider |
| TODO | Identity reconciliation | Admin reconciliation and identity-provider/local-user consistency operations | Keep provider calls behind `internal/identity` |
| DONE | School setup | None | Setup status, one-time admin provisioning, role assignment, and rollbacks are module-owned |
| DONE | Curriculum | None | Levels, groups, subjects, tree, and mediums are module-owned |
| DONE | Curriculum presets | None | Preview and transactional, idempotent preset seeding are module-owned |
| DONE | Classes | None | Classes, assignments, subject-teacher qualification, and enrollment are module-owned |
| DONE | Enrollments | None | Admin, public, and student self-service enrollment workflows are owned by Academics |
| DONE | Promotions | None | Preview and transactional assignment are module-owned |
| DONE | Term marks | None | Mark workflows and report-export reads are module-owned |
| DONE | Students | None for student profile CRUD and lifecycle | Identity provisioning, rollback, house assignment, and deletion are migrated |
| DONE | Teachers | None for teacher profile and qualification workflows | Identity provisioning, rollback, lifecycle, house, and subject operations are migrated |
| DONE | Guardians | None | CRUD, relationships, provisioning and rollback, parent access, authentication, notifications, reads, and routes are migrated |
| DONE | Non-academic staff | None | CRUD, employee numbering, employment status, and audited house assignment are migrated |
| DONE | Student portfolio | None | Business logic, persistence, handlers, and route ownership are migrated |
| TODO | Student self-service | Move the remaining profile resolver and HTTP ownership into a module | Attendance, marks, and enrollments already use module contracts |
| TODO | Parent self-service | Move the remaining HTTP ownership into a module | Guardian, attendance, marks, and timetable reads already use module contracts |
| TODO | Teacher self-service | Move the remaining profile resolver and HTTP ownership into a module | Leadership, society, and timetable reads already use module contracts |
| DONE | Student attendance | None | Sessions, records, scoped authorization, lock/correction rules, notifications, self/parent reads, and report reads are module-owned |
| DONE | Staff attendance | None | Admin marking, date views, monthly summaries, histories, and teacher self-service are module-owned |
| DONE | Notifications | None | Composer, scoped recipients, drafts, delivery, inbox state, and system-triggered delivery are module-owned |
| DONE | Positions | None | Appointments, scoped grants, rank resolution, audit recording, overview, and routes are module-owned |
| DONE | Section heads | None | Grade/stream assignments, reads, deletion, and leadership scope are module-owned |
| DONE | Prefects | None | Appointments, rank updates, archive years, student history, and routes are module-owned |
| DONE | Societies | None | CRUD, roster authorization, membership history, and teacher self-service are module-owned |
| DONE | Reports | None for attendance and marks PDF exports | Query adapters, templates, column selection, and routes are module-owned |
| DONE | Audit | None | Recording is exposed through a narrow shared port; admin reads and persistence are module-owned |
| TODO | Dashboard | Dashboard aggregates and house distributions | Depends on people and attendance data |
| TODO | Search | Cross-entity search | Should be isolated behind a search capability |
| TODO | Automation | Jobs, scheduler, and admin job endpoints | Keep process lifecycle ownership in `internal/app` |

## 6. Temporary compatibility pieces

These are intentionally retained until their active consumers are migrated:

- Nine horizontal repository files used by Authentication, Identity
  Reconciliation, Dashboard, Search, Automation, and self-service.
- Legacy handlers and services for those same remaining capabilities.
- Legacy route registration grouped behind `internal/routes/modules.go`.

Fourteen superseded repositories for School, Academics, Curriculum, People,
and Timetable have been deleted. The architecture guard rejects any new file
outside the exact nine-file compatibility allowlist.

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
6. Student and Staff Attendance are complete.
7. Notifications, Audit, and Report Export are complete; migrate Dashboard,
   Search, and Automation.
8. Positions, Section Heads, Prefects, and Societies are complete.
9. School Setup and Report Export are complete; remove compatibility
   repositories and the legacy route bridge.
10. Run full integration/API tests and update this document before deleting it.

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
