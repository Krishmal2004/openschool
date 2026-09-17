# Security and Performance Playbook

Status: **DRAFT for approval. No code changed yet.**
Scope: the whole project, `backend/` (Go, Gin, pgx, sqlc) and `frontend/` (React, Carbon).
Companion to `FRONTEND_REFACTOR_PLAYBOOK.md` (structure) and `UX_REVIEW_PLAYBOOK.md` (usability).

Writing rules from the refactor playbook apply here: short sentences, tables over
prose, no em-dashes, one-line code comments that say why.

---

## 1. Score

| Area | Now | After this playbook | What changes |
|------|-----|---------------------|--------------|
| Authentication and sessions | 7 | 9 | Default-password expiry, forgot-password throttling, shorter idle timeout on shared devices |
| Input handling | 7 | 9 | Search wildcard escaping, server-side length caps, generic error responses |
| Transport and headers | 5 | 9 | HSTS, CSP, Permissions-Policy, correct client IP behind the proxy |
| Data protection (PDPA) | 4 | 8 | Access logging, retention rules, erasure endpoint, encrypted backups |
| Scale (7,000 users) | 4 | 9 | Server pagination on every people list, search endpoints for pickers, trigram indexes |
| Delivery and caching | 4 | 8 | Compression, immutable asset caching, Carbon CSS trimmed, size budgets in CI |
| Observability | 3 | 8 | Structured logs with request ids, metrics, graceful shutdown, load test in CI |
| **Overall** | **5** | **8.5** | |

What is already right and must stay: sqlc parameterised queries everywhere (no
string-built SQL), RS256 pinned with issuer check, JWKS cached, bearer token never
stored by the app (no cookies, so no CSRF surface), no `dangerouslySetInnerHTML`,
body size limit, per-IP and per-account rate limiters, TLS 1.2+ for SMTP, reset
tokens hashed at rest with a 15-minute TTL, Swagger disabled in production,
nightly backup agent.

---

## 2. Threat model for a Sri Lankan school

| Who | How | What they want |
|-----|-----|----------------|
| Students | Guess a classmate's or teacher's password; the initial password is the NIC or index number printed on ID cards and report cards | Change marks or attendance, read others' records |
| Anyone with a phone | Phishing over WhatsApp or SMS with a fake reset link | Account takeover |
| Shared lab PCs | Previous user still signed in | Access to whoever signed in last |
| Insider (clerk, teacher) | Bulk export or edit outside their scope | Data leak, grade tampering |
| Bots | Credential stuffing, scraping the public forgot-password endpoint | Enumerate accounts, lock users out |
| Environment | Low-end Android phones, metered mobile data, power cuts, one Nginx in front of one binary | Slow pages and lost work count as failures |

Legal context: the Personal Data Protection Act No. 9 of 2022 (PDPA) applies. Student
records, NIC numbers, phone numbers and guardian details are personal data. The
system needs a lawful basis, retention limits, access logging and a way to erase.

---

## 3. Security findings and fixes

Severity: H = fix before the next release, M = next sprint, L = when convenient.

| # | Sev | Finding | Evidence | Fix |
|---|-----|---------|----------|-----|
| S1 | H | Initial passwords are the NIC number (teacher, parent) or index number (student). Both are printed on documents other students see. `must_change_password` exists but nothing forces or expires it. | `auth/service.go:67`, `teacher.ts` comment | Expire an unchanged default password after 7 days (account then needs the forgot-password flow). Reject a new password equal to the NIC or index number. Raise the minimum to 10 characters and block the top 10k breached passwords with `zxcvbn` on the client and a small deny list on the server. |
| S2 | H | `/auth/forgot-password` is public and only covered by the per-IP limiter. Behind Nginx every user shares one IP (see S3), so the limiter is either useless or locks the whole school out. | `groups.go`, `main.go:100` | Add a per-identifier limiter (5 per hour per email) and a per-IP limiter using the real client IP. Keep the response generic and constant-time. Log attempts. |
| S3 | H | `r.SetTrustedProxies(nil)` ignores `X-Forwarded-For`, so `ClientIP()` returns the proxy address for every request. Per-IP rate limits and audit IPs are wrong in production. | `main.go:88`, `ratelimit.go:75` | Read `TRUSTED_PROXIES` from the environment and pass it to `SetTrustedProxies`. Document the Nginx `proxy_set_header X-Forwarded-For` line. |
| S4 | H | 342 handlers return `err.Error()` to the client. pgx and validation errors expose table and column names and internal shapes. | `grep 'err.Error()' internal` | One error mapper: known domain errors become 4xx with a plain message, everything else becomes `500 {"error":"Something went wrong", "request_id":"..."}` and is logged server-side with the id. |
| S5 | M | The reset token travels in the URL (`/reset-password?token=`). It lands in browser history, proxy logs and the SPA host's access log. | `ResetPassword.tsx`, `auth/service.go:115` | Put the token in the URL fragment (`#token=`), read it on the client and POST it. Fragments are never sent to servers. Keep the 15-minute TTL. |
| S6 | M | Global search passes `q` straight into `ILIKE '%' || $1 || '%'`. `%` and `_` are wildcards, so `%` alone scans four tables. Minimum length is enforced only in the browser. | `search.sql`, `search/routes.go:11` | Escape `%`, `_` and `\` server-side, require 2 to 64 characters, and add `pg_trgm` GIN indexes on the searched columns. |
| S7 | M | Only three security headers on the API, none on the SPA. No HSTS, CSP or Permissions-Policy. | `security_headers.go`, `index.html` | At the proxy: HSTS (1 year, preload), `Permissions-Policy` (camera, microphone, geolocation off), `X-Content-Type-Options`. CSP for the SPA: `default-src 'self'; connect-src 'self' <thunderid>; img-src 'self' data:; font-src 'self'; style-src 'self' 'unsafe-inline'` (Carbon injects inline styles; test before removing `unsafe-inline`). Self-host IBM Plex so `fonts.googleapis.com` leaves the policy. |
| S8 | M | The school logo is a data URL stored in the `schools` row (up to 500 KB) and returned on every `GET /school`, which every admin page load calls. Only the prefix is validated, not the bytes. | `school/school.go:84` | Store uploads on disk or object storage under a random name, validate magic bytes and dimensions, serve with a long cache header, keep only the path in the row. |
| S9 | M | Audience is checked only when `THUNDERID_AUDIENCE` is set. Without it, any token from the same issuer is accepted. | `auth.go:156` | Make the audience required at startup, like the issuer. Add 30 seconds of leeway for clock skew. |
| S10 | M | `gin.Default()` logs every request line including query strings. Search terms, names and index numbers end up in plain-text logs. | `main.go:86` | Replace with `gin.New()` plus `Recovery` and a `slog` middleware that logs method, path template, status, latency, request id and user id only. |
| S11 | M | PDPA gaps: no access log for reads of student data, no retention policy, no erasure endpoint, backups not encrypted or copied off the host. | `audit` module covers 8 of 18 modules | Log reads of a student's profile and exports with actor and reason. Define retention (for example 7 years after leaving) and a nightly purge. Add an admin "erase person" flow that anonymises rather than deletes referenced rows. Encrypt backup dumps with `age` and copy them to a second location. Run a restore drill quarterly. |
| S12 | M | Idle logout is 30 minutes everywhere. Lab PCs and staffroom machines are shared. | `useIdleLogout.ts` | 15 minutes for student and teacher portals, 30 for admin. Show a 60-second warning before sign-out. Add "sign out of all devices" through ThunderID if the SDK supports it. |
| S13 | L | No graceful shutdown. A deploy during a mark save loses the request. | `main.go` | Handle `SIGTERM`, call `srv.Shutdown` with a 10-second deadline, then close the pool. |
| S14 | L | `govulncheck` and the dead-code scan are informational in backend CI. | `backend-ci.yml` | Make `govulncheck` blocking like `pnpm audit` now is. |
| S15 | L | No per-statement database timeout. A slow query holds a pooled connection under load. | `postgres.go` | Set `statement_timeout=10s` in the pool config and `context.WithTimeout` per request in the handlers. |
| S16 | L | Mailer credentials and database password live in `.env` on the host. | `backend/.env.example` | Docker or systemd secrets, file mode 600, rotate SMTP and DB passwords, never log them. |

Not a finding, keep as a rule: React escapes all text, the app has no
`dangerouslySetInnerHTML`, and notification bodies render as text. Add the ESLint
rule `react/no-danger` so it stays that way.

---

## 4. Pagination, front to back

Today `/students`, `/teachers`, `/guardians` and `/non-academic-staff` return whole
tables. 120 sqlc list queries exist and 32 have a `LIMIT`. Every picker
(`EntityCombobox`) downloads the full student list to search it in the browser.

### 4.1 API contract

```
GET /students?limit=25&offset=0&search=perera&sort=full_name&order=asc&grade_id=...
200 { "items": [...], "total": 6812, "limit": 25, "offset": 0 }
```

Rules: `limit` defaults to 25 and is capped at 100. `search` is trimmed, escaped
and capped at 64 characters. `sort` is validated against a per-endpoint allow list,
never interpolated. Filters are explicit query parameters, not a free-form
expression. Responses stay JSON objects, never bare arrays, so `total` can travel
with the page.

### 4.2 Database

```sql
-- name: ListStudentsPage :many
SELECT sp.*, c.name AS class_name, g.name AS grade_name, h.name AS house_name,
       COUNT(*) OVER () AS total
FROM student_profiles sp ...
WHERE ($3::text = '' OR sp.full_name ILIKE '%' || $3 || '%' OR sp.index_number ILIKE '%' || $3 || '%')
ORDER BY sp.full_name
LIMIT $1 OFFSET $2;
```

`COUNT(*) OVER ()` returns the total with the page in one round trip. Add
`CREATE EXTENSION pg_trgm` and GIN trigram indexes on `student_profiles.full_name`,
`student_profiles.index_number`, `teacher_profiles.full_name`,
`teacher_profiles.employee_number`, `guardians.full_name` and `guardians.phone` so
`ILIKE '%x%'` uses an index. Keep `ORDER BY` on an indexed column.

### 4.3 Frontend

`DataGrid` already accepts `server={{ page, pageSize, totalItems, onChange }}`.
Per list page: keep `useListFilters`, put `debouncedSearch`, filters, page and
page size into the query key, pass `placeholderData: keepPreviousData` so the old
page stays visible while the next loads, and mirror the filters into the URL so a
refresh or a shared link keeps them.

Pickers: `EntityCombobox` gets an `onSearch` prop backed by `?search=` with a
300 ms debounce and a 2-character minimum. `useStudents()` with no filters is
removed from every page that only needed a picker.

### 4.4 Order of work

1. `/students` and the Students page (reference implementation).
2. `/teachers`, `/guardians`, `/non-academic-staff`, `/audit-logs`.
3. Pickers: enrol student, guardians, assign teacher, prefects, societies.
4. Notifications sent list and any query still ending without `LIMIT` that a user can grow.

Exit: no list endpoint returns more than 100 rows; the Students page renders in
under 1 second against a 7,000-row seed on a mid-range Android over 3G.

---

## 5. Performance and optimisation checklist

### Backend

| Item | Why | How |
|------|-----|-----|
| Compression | JSON lists shrink 5 to 10 times | `gin-contrib/gzip` or gzip and Brotli at Nginx |
| Cache headers | Grades, subjects, mediums, houses change rarely | `Cache-Control: private, max-age=300` plus ETag on those endpoints |
| Select only needed columns | `sp.*` joins ship every column to list views | Dedicated list queries returning the columns the grid shows |
| Trigram indexes | Search on 7,000 rows without a scan | Section 4.2 |
| Pool sizing | 25 connections per binary must fit Postgres `max_connections` | Document the ratio; use PgBouncer if more than one replica |
| Statement timeout | One slow query must not hold the pool | S15 |
| Structured logs | Find a slow request in production | `slog` JSON with request id, user id, route, latency |
| Metrics | Know p95 latency and error rate | `/metrics` for Prometheus, protected by the proxy |
| Dashboard analytics | Check for N+1 patterns in `dashboard` and `reports` | `EXPLAIN ANALYZE` the top five queries, add covering indexes |
| Graceful shutdown | Zero-downtime deploys | S13 |

### Frontend

| Item | Why | How |
|------|-----|-----|
| Carbon CSS | 874 KB today, most unused | Replace `@use "@carbon/react"` with per-component `@use "@carbon/react/scss/components/..."` for the components in use; target under 400 KB |
| Fonts | Third-party request on every load, blocked by CSP | Self-host IBM Plex Sans and Mono as woff2 with `font-display: swap` |
| Per-entity stale time | Grades and subjects do not change during a session | `staleTime: 1 hour` for reference lists, 5 minutes elsewhere |
| Prefetch on hover | Sidebar links feel instant | `queryClient.prefetchQuery` on `onMouseEnter` for the next page's list |
| Virtualise long lists | Client-side pages over 200 rows | `react-window` inside `DataGrid` when `pagination={false}` and rows exceed 200 |
| Size budget | Stop regressions | `vite-bundle-visualizer` report in CI, fail if app chunk exceeds 100 KB or CSS exceeds 400 KB |
| Route error and loading shells | No blank white screens | Skeleton layout in `ProtectedRoute` and `App` fallback (also a UX item) |
| Images | Logo delivered as data URL | S8 |

### Infrastructure

| Item | How |
|------|-----|
| Static assets | Nginx serves `dist/` with `Cache-Control: public, max-age=31536000, immutable` for hashed files and `no-cache` for `index.html` |
| HTTP/2 and TLS | At Nginx, with HSTS |
| Docker | Multi-stage builds, non-root user, read-only filesystem, health checks on `/healthz` |
| Postgres | `shared_buffers` 25 percent of RAM, `work_mem` tuned for sorts, autovacuum on, weekly `REINDEX` of trigram indexes |
| Backups | Nightly encrypted dump, off-host copy, monthly restore test |
| Load test | `k6` script simulating 7,000 sessions: sign-in, dashboard, one list page, one poll; run in CI nightly with p95 under 500 ms |

---

## 6. Execution plan

| Phase | Scope | Exit criteria |
|-------|-------|---------------|
| A. Hardening | S1 to S10, S12 to S16 | `govulncheck` and `pnpm audit` blocking; error responses carry no internal text; OWASP ZAP baseline scan has no medium findings |
| B. Pagination | Section 4 | No endpoint returns more than 100 rows; Students page under 1 s on the 7,000-row seed |
| C. Delivery | Compression, caching, Carbon trim, fonts, size budget | CSS under 400 KB; Lighthouse performance 90+ on a mid-range Android profile |
| D. Data protection | S11 | Retention job running; erasure flow documented and tested; restore drill recorded |
| E. Observability | Logs, metrics, k6 | Dashboard with p95 latency, error rate, active sessions; nightly load test green |

One PR per phase. Each PR keeps `go test ./...`, `staticcheck`, `tsc`, `lint`,
`test` and `build` green.

---

## 7. Verification

```bash
cd backend && go build ./... && go vet ./... && go test ./... && staticcheck ./... && govulncheck ./...
cd frontend && pnpm exec tsc -b --force && pnpm lint && pnpm test && pnpm build && pnpm audit --prod
# after phase A
docker run --rm -t zaproxy/zap-stable zap-baseline.py -t https://<host>
# after phase B
k6 run scripts/load/7000-users.js
```
