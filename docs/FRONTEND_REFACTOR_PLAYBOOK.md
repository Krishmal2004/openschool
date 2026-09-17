# Frontend Refactor Playbook

Status: **Phases 0 to 3 complete (17 Sep 2026). Phase 4 (performance verification) not started.**
Decision: `src/services/` and `src/queries/` are removed. Code moves to feature folders (Section 2).
Scope: `frontend/` only. Backend changes are listed as dependencies, not done here.
Goal: a frontend that stays fast with ~7,000 users, has one obvious place for every
piece of code, no duplicated or dead code, and a smaller attack surface.
See also `SECURITY_AND_PERFORMANCE_PLAYBOOK.md` (security, pagination, optimisation
across both workspaces) and `UX_REVIEW_PLAYBOOK.md` (usability backlog).

Baseline (17 Sep 2026, branch `development`): `tsc -b` clean, `eslint` clean,
`pnpm build` OK. 270 files, ~30,300 lines under `src/`. No test runner exists.

---

## 0. Writing rules (apply to this document, code comments and reports)

- Plain English, short sentences. One idea per sentence.
- No em-dashes. Use a comma, a colon or a new sentence.
- Keep prose short. Use a table or a list when there is more than one item.
- Code comments: one line, professional, only where the code is not obvious.
  Say why, not what. Example: `// Backend returns 404 for unlinked guardians, treat as empty.`
- No banner comments at the top of files. No commented-out code.
- Reports state numbers and results, not effort.

---

## 1. Audit findings

Everything below was verified against the current source. File references are
the evidence; fix each item in the phase noted.

### 1.1 Scale and performance (7,000 users)

| # | Finding | Evidence | Phase |
|---|---------|----------|-------|
| P1 | Every logged-in session polls **two** notification endpoints every 60 s (`/me/notifications` and `/me/notifications/unread-count`) because `NotificationsBell` is in every layout header. 7,000 open sessions = ~230 requests/s of polling alone, before any real traffic. The full list is fetched even when the dropdown is closed. | `queries/notifications/useNotifications.ts:14-32`, `components/common/NotificationsBell.tsx:21-22` | 0 |
| P2 | Backend list endpoints (`/students`, `/teachers`, `/guardians`, `/non-academic-staff`) have **no pagination**. The frontend downloads the full table and paginates client-side. Filtering runs on every render with no `useMemo`, so each keystroke re-filters the entire array and re-renders the Carbon `DataTable`. | `pages/admin/students/Students.tsx:77-96`, `pages/admin/teachers/Teachers.tsx:50-58`, backend `people/students_routes.go:48` | 1, backend follow-up |
| P3 | CSS bundle is **850 KB** because `index.scss` imports the whole of Carbon (`@use "@carbon/react"`). Main JS is 680 KB plus a 503 KB vendor chunk. `vite.config.ts` has no `manualChunks`, so every route shares one vendor bundle and any dependency bump invalidates the user's whole cache. | `src/index.scss:1`, `frontend/dist/assets/*`, `vite.config.ts` | 0, 4 |
| P4 | IBM Plex is loaded **twice**: once via `<link>` in `index.html` and again via a render-blocking `@import url(...)` inside the SCSS. | `index.html:6-9`, `src/index.scss:4` | 0 |
| P5 | `useApi` registers the axios auth interceptor inside a `useEffect` in `App`. Child effects run before parent effects, so any query fired from a child on first mount (for example `useSetupStatus` in `SignIn`) can go out **without** the bearer token. StrictMode also ejects and re-adds it. There is no response interceptor, so an expired session shows generic errors instead of redirecting to sign-in. | `hooks/useApi.ts`, `lib/api.ts`, `pages/SignIn.tsx:9` | 0 |
| P6 | Zero `React.memo`, two `useCallback` in the whole tree. Large composite pages (`ClassDetail`, `ChildDetail`, `TimetableEditor`) re-render all tabs on any state change. | `grep memo(` = 0 | 3 |
| P7 | Query keys are hand-written arrays. `MY_NOTIFICATIONS_KEY = ["me","notifications"]` nests under the raw `["me"]` key used in `useAuth`/`useProvisionUser`, so a password change invalidates and refetches all notification queries. ~75 exported key constants are never imported outside their own file. | `queries/notifications/useNotifications.ts:5-7`, `queries/useAuth.ts:19,27`, `hooks/useProvisionUser.ts:22` | 0 |
| P8 | A single `ErrorBoundary` wraps the whole app. One render error in any page blanks every portal. | `main.tsx:21` | 1 |

### 1.2 Security

| # | Finding | Evidence | Phase |
|---|---------|----------|-------|
| S1 | `parseJwt` decodes the JWT payload with plain `atob`. JWT payloads are **base64url**; `atob` throws on `-` or `_`, `parseJwt` returns `null`, and the user lands on "Access restricted" with no error. Latent today, guaranteed to bite on some tokens. | `hooks/useRole.ts:7-13` | 0 |
| S2 | `react-router@8.0.1` has an open **high** advisory (GHSA, RSC-mode CSRF). Not exploitable in SPA/data mode, but it fails `pnpm audit --prod`. Upgrade. | `pnpm audit --prod` | 0 |
| S3 | Two `dompurify` advisories (moderate, low) via `@thunderid/react`. Transitive; fix with a `pnpm.overrides` pin until upstream bumps. | `pnpm audit --prod` | 0 |
| S4 | Password rules (min 8 chars) are re-implemented in three places with no shared validator, so a policy change needs three edits. | `pages/ResetPassword.tsx:18`, `components/common/ChangePasswordModal.tsx:25`, `pages/PasswordInterstitial.tsx` | 1 |
| S5 | `ErrorBoundary.componentDidCatch` logs the full error and component stack to the console in production builds. | `components/common/ErrorBoundary.tsx:22` | 0 |
| S6 | No `Content-Security-Policy` or other security headers are set for the SPA. This belongs on whatever serves `dist/`; noted as a deployment follow-up, not a frontend code change. | `index.html` | follow-up |

What is already correct and must stay: access token is never stored by app code
(read per request from the ThunderID SDK); role comes from the JWT `roles` claim and
the backend re-checks it on every route; idle logout after 30 min; dev-only routes are
gated by `import.meta.env.DEV`; no `dangerouslySetInnerHTML`, `eval`, or web storage use.

### 1.3 Structure

| # | Finding | Evidence | Phase |
|---|---------|----------|-------|
| A1 | Code for one feature is spread over up to five directories (`services/x.ts`, `queries/useX.ts`, `pages/admin/x/`, `pages/teacher/X.tsx`, `components/x/`). Fixing one student bug means opening `services/student.ts`, `queries/useStudents.ts`, `pages/admin/students/*`, `pages/student/*`, `services/studentSelf.ts`, `queries/useStudentSelf.ts`. | tree | 2 |
| A2 | No import alias. 700+ imports use `../../../` or `../../../../`; moving a file breaks every import path. | `grep "from \"../../../"` | 0 |
| A3 | `App.tsx` (277 lines) holds 70 lazy imports and four route trees in one ternary chain. | `App.tsx` | 2 |
| A4 | Four layouts each re-implement the same sidebar-group rendering with inline styles and their own `NAV_GROUPS` shape. | `layouts/*.tsx` (511 lines total) | 1 |
| A5 | 25 files exceed 250 lines; the largest are `Classes.tsx` (600), `ChildDetail.tsx` (533), `ClassDetail.tsx` (448), `TeacherMarks.tsx` (418), `TimetableEditor.tsx` (407), `Houses.tsx` (385). | `wc -l` | 3 |
| A6 | `index.scss` is one 1,451-line file mixing tokens, resets, Carbon overrides, page styles and auth styles. | `src/index.scss` | 1 |

### 1.4 Duplication

| # | Finding | Evidence | Phase |
|---|---------|----------|-------|
| D1 | **1,246** inline `style={{...}}` blocks across ~110 files, and **91** hard-coded hex colours in TSX, while a full `--os-*` token set already exists. Top offenders: `ChildDetail` (39), `AttendanceByClassSection` (32), `Classes` (31), `ChartPrimitives` (30), `AttendanceMark` (29). | `grep style={{` | 3 |
| D2 | The list-page shell (toolbar search + filter selects + "Active Filters" tag row + Clear All + skeleton/error/empty + pagination) is copy-pasted across `Students`, `Teachers`, `GuardiansDirectory`, `NonAcademicStaff` and partially in `Positions`, `Prefects`, `Societies`. | `Students.tsx:120-190` vs `Teachers.tsx:80-115` | 1, 3 |
| D3 | 26 files hand-roll `<table className="os-table">`; only 2 use Carbon `DataTable`. Two table systems to maintain. | `grep '<table className="os-table"'` | 1, 3 |
| D4 | Same constants declared in multiple files: `EMAIL_RE` x4, `STATUS_STYLES` x3, `SEVERITY_TAG` x3, `STATUS_TAG` x2, `PRIORITY_TAG` x2, `EMPLOYMENT_STATUSES` x2, `RANK_SECTION_HEAD` x2, `ACCENT` x2. | `grep -rn "const EMAIL_RE"` etc. | 1 |
| D5 | Every mutation hook repeats the same `useQueryClient` + `invalidateQueries` boilerplate (three invalidations per student mutation, copied four times). | `queries/useStudents.ts:29-84` | 2 |

### 1.5 Dead code

| # | Item | Evidence | Phase |
|---|------|----------|-------|
| X1 | `react-hook-form` and `@tanstack/react-query-devtools` are dependencies with zero imports. | `package.json`, `grep` | 0 |
| X2 | `src/App.css` is empty. | `wc -l` = 0 | 0 |
| X3 | Unused exports: `useUpdateNotificationDraft`, `useUpdateProgressReport`, `CHART_PINK`, `silentProgressStatus`, `TeacherProfileForm`, plus ~75 query-key constants never imported elsewhere. | `scratch unused-export scan` | 0, 2 |
| X4 | "This file defines the X component, which..." banner comments on `App.tsx` and every layout; several comments restate what the code does. | `App.tsx:1`, `layouts/*.tsx:1` | 3 |

### 1.6 Architecture score

Scored 1 to 10 against what an enterprise product with ~7,000 users needs.

| Area | Now | After refactor | What changes |
|------|-----|----------------|--------------|
| Structure and findability | 4 | 9 | One folder per feature, `@/` alias, layer rules enforced by ESLint |
| Scale (7,000 users) | 3 | 8 (9 with backend pagination) | Polling cut to 1 request per 120 s, memoised lists, `DataGrid` ready for server pages |
| Performance and bundle | 4 | 8 | Manual chunks, one font load, Carbon CSS trimmed, route-level error boundaries |
| Security | 6 | 8 | base64url JWT fix, 401 handling, patched dependencies, one password validator |
| Code quality and duplication | 4 | 9 | One table, one list shell, one layout shell, shared constants, no dead code |
| Consistency and maintainability | 5 | 9 | 250-line limit, tokens only, same page template everywhere |
| Guardrails (lint, tests, CI) | 3 | 8 | Layer and size rules in ESLint, unit tests for shared hooks and helpers |
| **Overall** | **4** | **8.5** | |

What "after" means in practice: a bug in students is fixed inside `features/students/`
only; a new list page is under 150 lines because it composes `DataGrid` and
`FilterBar`; a dependency bump does not invalidate the whole app cache; the
server sees one poll per session every two minutes; and CI fails if anyone
imports across layers or writes a 300-line page.

---

## 2. Target architecture

### 2.1 Folder layout (feature-first)

```
src/
  app/
    main.tsx                 providers only
    App.tsx                  role switch -> one of four route modules
    queryClient.ts
    routes/
      admin.routes.tsx       lazy pages + <Route> list for the admin portal
      teacher.routes.tsx
      student.routes.tsx
      parent.routes.tsx
      public.routes.tsx      signin / setup / forgot / reset
  shared/
    api/
      client.ts              axios instance, token provider, 401 handler (module-level)
      errors.ts              getErrorMessage, isNotFoundError
      keys.ts                query-key factory root (see 2.3)
    auth/                    useRole (base64url-safe), ProtectedRoute, useIdleLogout,
                             useProvisionUser, password.ts (one validator)
    ui/                      DataGrid, FilterBar, ActiveFilterTags, ListState,
                             FormModal, ConfirmDeleteModal, EmptyState, ErrorMessage,
                             StatusView, StatusTag, Avatar, InfoRow, SectionHeader,
                             RouteErrorBoundary, skeletons
    hooks/                   useDebounced, usePagination, useListFilters
    lib/                     date, name, phone, constants (statusStyles, severityTag,
                             priorityTag, timetableStatusTag, employmentStatuses)
    styles/
      index.scss             @use only
      _tokens.scss  _base.scss  _carbon-overrides.scss  _layout.scss
      _page.scss  _table.scss  _forms.scss  _auth.scss  _utilities.scss
  layouts/
    PortalShell.tsx          header + sidebar renderer (one implementation)
    nav/admin.ts teacher.ts student.ts parent.ts   nav configs only
  features/
    students/    api/  queries/  keys.ts  components/  pages/admin/  pages/self/
    teachers/    (same shape; teacher self-service pages live in pages/self/)
    guardians/
    staff/
    academics/   grades, classes, streams, promotion
    curriculum/  subjects, buckets, mediums, presets
    attendance/  student + staff attendance
    marks/       term marks entry (admin, teacher, student views)
    timetable/
    notifications/
    portfolio/   activities, awards, discipline, progress reports, prefects, societies
    positions/
    school/      settings, houses, setup wizard, academic years, terms
    reports/     reports + analytics + dashboard charts
    system/      audit log, automation, orphaned accounts, global search
    parent/      ParentDashboard, ChildDetail (composes other features' queries)
    auth/        sign-in, forgot / reset password, first-login interstitial
```

Rules:

- A feature owns `api/` (one axios wrapper file per entity, with its request and
  response types beside the calls), `queries/` (TanStack hooks, one file per entity),
  `keys.ts` (its query-key factory), `components/` and `pages/`. Nothing else may call
  that feature's `api`.
- `shared/` never imports from `features/`. A feature may import another feature's
  `queries`, `keys` and `components`, and its `api` types (type-only). It never imports
  another feature's `pages` or calls its `api`.
- Imports use the `@/` alias (`@/features/students/queries`). No `../../..`.
- One default export per page file. Page files are route targets and nothing else.
- Hard limit **250 lines** per file; a page over the limit is split into
  `components/` under its feature.

### 2.2 Layer contract

| Layer | May import | Must not |
|-------|-----------|----------|
| `features/*/api/` | `shared/api/client`, other api files' types | React, Carbon, TanStack, hooks, components |
| `features/*/queries/` | own `api`, own `keys`, other features' `keys` | components, pages |
| `features/*/components`, `pages` | own queries, `shared/*`, other features' queries and components | axios, any `api` value, other features' pages |
| `layouts/`, `app/` | `shared/*`, feature queries, feature pages (lazy) | any `api` value, axios |

Enforced by `eslint.config.js` (`@typescript-eslint/no-restricted-imports`, one block
per feature, generated from the folder list). A violation fails `pnpm lint`.

### 2.3 Query-key factory

`shared/api/keys.ts` holds only cross-cutting keys (`me`). Each feature owns
`features/<name>/keys.ts`:

```ts
export const studentKeys = {
  all: ["students"] as const,
  list: () => ["students", "list"] as const,
  detail: (id: string) => ["students", "detail", id] as const,
  byClass: (classId: string) => ["students", "by-class", classId] as const,
};
```

Rules: every key is a function on the factory, keys are `as const`, no raw array
literals in hooks, one root segment per feature so `invalidate(studentKeys.all)` is
safe and never crosses features. Mutations call `useInvalidate()` from
`shared/api/useInvalidate.ts` and pass the key prefixes to refresh; most pass the
feature root. Data another feature needs to batch (rosters, sessions) is exposed as
`queryOptions` (`studentsByClassOptions`) so it can be composed with `useQueries`
without touching the owning feature's `api`.

### 2.4 Data-fetching rules for 7,000 users

- Global defaults stay `staleTime: 5 min`, `retry: 1`. Add `gcTime: 10 min`.
- Polling: only `keys.notifications.unread()` polls, every **120 s**, and only while
  the tab is visible. The notification list is `enabled: open` (fetched on click).
- Any list that can exceed ~500 rows (students, guardians, teachers, audit log) must
  use `useMemo` for filtering, `useDebounced(300)` for search, and `placeholderData:
  keepPreviousData` once server pagination exists.
- Backend follow-up (required to fully meet the 7,000-user goal, tracked separately):
  `limit`/`offset`/`search` params on `/students`, `/teachers`, `/guardians`,
  `/non-academic-staff`, `/audit-logs`; keep `/me/notifications/unread-count` an
  indexed `COUNT(*)`. The frontend `DataGrid` is built to accept either client-side
  arrays or server pages so this switch is a one-line change per list.

### 2.5 Styling rules

- Zero hex literals in TSX. Colours come from `--os-*` tokens only; add missing tokens
  to `_tokens.scss` (Sass var + CSS var mirror). Header stays `#007D9C`.
- Zero `style={{}}` for layout and typography. Use existing `.os-*` classes or add a
  small utility set in `_utilities.scss` (`.os-stack`, `.os-row`, `.os-gap-2`,
  `.os-text-sm`, `.os-muted`, `.os-mono`). Inline style is allowed only for values
  that are truly dynamic (chart widths, computed colours).
- Listing sections use `var(--os-layer)` background and `var(--os-border-subtle)`.
- Fonts load once, from `index.html`.

### 2.6 Shared UI to build once (Phase 1)

| Component | Replaces |
|-----------|----------|
| `DataGrid` (Carbon `DataTable` + `Pagination` + skeleton/empty/error states, no select checkboxes, optional server-mode props) | 26 hand-rolled `<table>` blocks and the two `DataTable` pages |
| `FilterBar` + `ActiveFilterTags` + `useListFilters()` | the copy-pasted toolbar / tag row / Clear All in 7 list pages |
| `ListState` (`loading` / `error` / `empty` / children) | the same three-way ternary in every page |
| `PortalShell` + nav config type | four layouts |
| `RouteErrorBoundary` | nothing today; wraps `<Outlet>` in `PortalShell` |
| `PasswordFields` + `shared/auth/password.ts` | three separate password forms |
| `StatusTagMap` helpers in `shared/lib/constants` | D4 duplicates |

---

## 3. Execution plan

Each phase is one PR (Phase 2 and 3 may be one PR per feature). Every commit keeps
`tsc`, `eslint` and `pnpm build` green. No behaviour change is bundled with a
structural move; moves are `git mv` so history survives.

### Phase 0 - Foundation and quick wins (no visual change)  [DONE]

1. Add `@/` alias (`tsconfig.app.json` paths + `vite.config.ts` resolve.alias).
2. Add `manualChunks` (react, carbon, query, router, thunderid) and `build.target`.
3. Remove the SCSS `@import url(...)` font line (P4).
4. Rewrite `shared/api/client.ts`: module-level axios instance, `setAccessTokenProvider()`
   called once from `main.tsx`, request interceptor, response interceptor that signs
   out on 401 (P5).
5. Fix `parseJwt` base64url (S1). Strip prod console logging in `ErrorBoundary` (S5).
6. Notifications: unread-count only, 120 s, list on open (P1).
7. Introduce `keys.ts` for `me` and notifications first; remove the `["me"]` overlap (P7).
8. Remove dead deps, `App.css`, dead exports (X1-X3). Upgrade `react-router`;
   add `pnpm.overrides` for `dompurify` (S2, S3). Re-run `pnpm audit --prod`.
9. Add ESLint: `no-restricted-imports` (no `../../..`), `max-lines: 250` (warn until
   Phase 3, then error).
10. Add `vitest` + Testing Library (required, not optional) with tests for
    `parseJwt`, `getErrorMessage`, `usePagination`, `useListFilters`.

Exit: build green, audit shows only upstream items, bundle report saved to
`docs/perf/bundle-phase0.txt`.

### Phase 1 - Shared UI and styles  [DONE]

1. Split `index.scss` into `shared/styles/` partials; add `_utilities.scss`.
2. Build the components in 2.6. Migrate **one** list page (`Students`) onto
   `DataGrid` + `FilterBar` + `useListFilters` as the reference implementation, with
   memoised filtering and debounced search (P2).
3. `PortalShell` + four nav configs; wrap `<Outlet>` in `RouteErrorBoundary` (A4, P8).
4. Consolidate D4 constants into `shared/lib/constants/`.
5. `shared/auth/password.ts` + `PasswordFields` (S4).

Exit: `Students` page has 0 inline styles, 0 hex, <200 lines; layouts total <150 lines.

### Phase 0 and 1 results

| Metric | Before | After |
|--------|--------|-------|
| App entry JS | 680 KB (+503 KB vendor) | 46 KB app + vendor chunks: carbon 444, thunderid 401, react 232, query 84, icons 56 |
| Polling per session | 2 requests / 60 s | 1 request / 120 s, tab visible only |
| `pnpm audit --prod` | 1 high, 1 moderate, 1 low | 0 |
| Relative `../` imports | 700+ | 0 (ESLint blocks them) |
| Layout files (4) | 511 lines | 43 lines + `PortalShell` |
| `Students.tsx` | 334 lines, 2 hex, 25 inline styles | 193 lines, 0 hex, 0 inline styles |
| Unit tests | none | 12 (jwt, errors, pagination, list filters) |
| Inline styles / hex in TSX | 1,246 / 91 | 1,210 / 79 (rest is Phase 3) |
| CSS bundle | 850 KB | 853 KB (Carbon trim is Phase 4) |

Carried into Phase 2: `shared/ui/NotificationsBell`, `GlobalSearch`,
`AgentFindingsBanner`, `ChangePasswordModal` and `shared/lib/classroom.ts`
still import feature queries; they move into their features. Feature
`constants.ts` files re-export `EMPLOYMENT_STATUSES` and `STATUS_STYLES`
from `shared/lib/constants` until their pages are moved.

### Phase 2 - Feature-first move  [DONE]

Per feature, in this order (smallest blast radius first): `notifications`, `students`,
`teachers`, `guardians`, `staff`, `positions`, `portfolio`, `school`, `curriculum`,
`academics`, `attendance`, `timetable`, `reports`, `system`, `parent`.

For each: `git mv` service -> `api.ts`, query hook -> `queries.ts`, types ->
`types.ts`, pages/components into the feature; convert keys to the factory; replace
mutation boilerplate with `useInvalidate`; drop unused exports. Split `App.tsx`
into `app/routes/*.routes.tsx` at the end (A3).

Exit: `src/services`, `src/queries`, `src/pages` no longer exist; ESLint layer rules
are `error`.

### Phase 2 results

| Metric | Before | After |
|--------|--------|-------|
| Folders a student bug can live in | 5 (`services`, `queries`, `pages/admin`, `pages/student`, `components`) | 1 (`features/students/`) |
| `src/services`, `src/queries`, `src/pages`, `src/components` | 247 files | removed |
| Query-key factories | 0 (hand-written arrays) | 16 `keys.ts`, no raw keys in hooks |
| `useQueryClient` boilerplate in mutations | every mutation | 0, all via `useInvalidate` |
| Route table | one 277-line `App.tsx` | `App.tsx` 55 lines + 5 route modules |
| Layer rules in ESLint | 0 | 5 rule blocks, 0 violations |
| Polling | jobs panel 30 s | 60 s, tab visible only |

Decisions made during the move: types stay beside their API calls (`api/<entity>.ts`)
instead of a separate `types.ts`, one file per entity keeps every file well under the
limit; another feature's components may be composed (detail pages are built from
tabs owned by several features), only pages and api values are off limits; marks got
its own feature; tab components that were route-less "pages" (`ClassMarks`,
`AuditLog`, `OrphanedAccounts`) moved to `components/`.

### Phase 3 - Page-by-page cleanup  [DONE]

Apply Section 4 to every page, feature by feature, same order as Phase 2. Priority
targets: the six files over 380 lines (A5) and the five worst inline-style offenders
(D1). Add `memo` only where the profiler shows a re-render problem (P6).

Exit: `grep style={{` under 100 (dynamic values only); `grep '#[0-9a-f]{6}' src --include=*.tsx` = 0;
no file over 250 lines; `max-lines` rule is `error`.

### Phase 3 results

| Metric | Before | After |
|--------|--------|-------|
| Inline `style={{}}` in TSX | 1,210 | 31, every one a runtime value (data colours, computed widths) |
| Hex literals in TSX/TS | 79 | 6, all in one data constant (house colour palette saved to the backend) |
| Files over 250 lines | 22 (largest 600) | 0 (largest 248); `max-lines` is now an ESLint error |
| Hand-rolled `<table>` | 29 | 9, all editing grids with form controls in every row |
| Pages on `DataGrid` | 2 | 23 |
| Multi-line comment blocks | 30 | 0; comments are one line and say why |
| Source lines | 30,300 | 25,700 across 395 files |
| `pnpm lint` | 16 warnings | 0 errors, 0 warnings |

How the styling was done: a generated utility class set on a 0.25rem grid in
`shared/styles/_utilities.scss` (spacing, flex/grid, typography, token colours,
borders, sizes), applied by a converter that mapped every static declaration
to a class and left only runtime values inline. State-driven looks became
toggle classes (`os-pill.is-active`, `os-subject-card.is-hover`, `os-select-row`).

Editing grids stay as plain tables on purpose: `MarksEntryTable`, `ClassMarks`,
`AttendanceMark`, `StaffAttendance`, `PeriodsEditor`, `TimetableGrid`,
`SubjectRequirements`, `PromotionGroup`, `TeacherSubjects`. `DataGrid` is
read-only by design; putting inputs through a render callback would make
those pages harder to follow, not easier.

Read-only tables shared across portals now live in their owning feature:
`AttendanceHistoryTable`, `TermMarksTable`, `EnrollmentsTable`,
`ProgressReportsTable`, `GuardiansTable`, `TimetableByDay` and the
`PortfolioTables` set. The student, parent and admin views render the same
component, so a column change happens once.

### Phase 4 - Performance verification

1. Load `Students`, `Guardians`, `AuditLog` with a 7,000-row seeded backend; record
   time-to-interactive and search keystroke latency with React Profiler.
2. Bundle report vs Phase 0 baseline; target: CSS < 400 KB, largest JS chunk < 350 KB.
   If the Carbon CSS is still the bulk, switch `index.scss` to per-component Carbon
   `@use` imports.
3. Confirm polling budget: one request / 120 s per session.
4. File the backend pagination tickets with the exact param contract `DataGrid` expects.

---

## 4. Per-page standard (applied in Phase 3, and to every new page after)

Read every file in the named feature first, plus any `shared/ui` component it uses.
Then enforce, file by file:

1. **Tokens only.** No hex in TSX/SCSS outside `_tokens.scss`.
2. **No layout/typography inline styles.** Use `.os-*` classes or utilities.
3. **Tables** use `DataGrid` (`pagination={false}` for short lists inside tabs). A grid
   whose rows hold form controls stays a plain `<table className="os-table">`. No
   `TableSelectAll` / batch actions unless the page has a real bulk operation.
4. **Lists** use `FilterBar` + `ActiveFilterTags` + `useListFilters`; search debounced;
   filtering memoised.
5. **Mutations** that delete or reset go through `ConfirmDeleteModal`; create/edit go
   through `FormModal`. Direct `mutate()` from a click is allowed only for reversible
   toggles (mark read, archive, set current).
6. **Queries** use the key factory; mutations invalidate via `useInvalidate`.
7. **Size** under 250 lines; extract to `components/` inside the feature.
8. **Comments**: delete any that restate the code. Keep one-liners for non-obvious
   invariants or backend quirks only.
9. **Dead code**: no unused exports, imports or props.
10. **Error and loading states** via `ListState`; never a bare `<div>` spinner.

---

## 5. Verification (every commit)

```bash
cd frontend
pnpm exec tsc -b --force
pnpm lint
pnpm build
pnpm audit --prod        # only upstream (@thunderid) items may remain
```

Plus, per phase, the exit criteria listed in Section 3. A phase is not complete
until its criteria are met and reported with the actual numbers.

---

## 6. Out of scope / follow-ups outside `frontend/`

- Server-side pagination and search on people list endpoints (backend).
- Security headers (CSP, HSTS, X-Frame-Options) on the static host / reverse proxy.
- `@thunderid/react` bumping its `dompurify` dependency (upstream).
- Any change to ThunderID auth flow or role model.

---

## 7. Completion report format

For each phase or feature, report:

1. Files moved / created / deleted (counts) and the largest file after the change.
2. Inline-style and hex counts before -> after.
3. Bundle sizes before -> after (Phase 0 and 4).
4. Polling requests per session per minute (Phase 0).
5. Verification output: `tsc`, `lint`, `build`, `audit` status.
6. Anything deliberately left, with the reason.
