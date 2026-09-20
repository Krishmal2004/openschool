# UX Review Playbook

Status: **"Now" backlog underway (started 20 Sep 2026).** Every finding below has a
Status column: ✅ Done · 🚧 Started (partly done) · ⬜ Not started. "Next"/"Later" items
are all ⬜ — this pass only touched the "Now" backlog.
Reviewer stance: senior product designer and frontend engineer, reviewing what a
school admin, a teacher on a phone in class, a student and a parent actually see.
Evidence comes from the current `frontend/src` after the refactor.

---

## 1. Score

| Area | Now | After | What changes |
|------|-----|-------|--------------|
| Navigation and information architecture | 4 | 8 | 27 admin items in 8 groups become 5 hubs; timetable and settings get one home each |
| Task efficiency (top tasks) | 5 | 8 | Attendance marking on a phone, bulk enrolment, sorting and export on lists |
| Feedback and system status | 4 | 8 | One toast system, idle-logout warning, no blank loading screens, page titles |
| Forms and input | 6 | 8 | Every control labelled, human error messages, unsaved-change guard |
| Content and copy | 4 | 8 | One spelling, one date format, sentence case, plain error text |
| Accessibility | 5 | 9 | Labels, focus styles, skip link, 12 px minimum text, axe in CI |
| Mobile | 3 | 8 | Responsive shell for all portals, touch-first attendance and marks |
| Localisation | 1 | 7 | Sinhala and Tamil for student and parent portals, en-LK dates everywhere |
| Onboarding | 6 | 8 | Setup checklist after the wizard; empty states that lead somewhere |
| Consistency | 6 | 9 | Design tokens and shared components already in place; copy guide added |
| **Overall** | **4.5** | **8** | |

---

## 2. Who uses it and what they need first

| Persona | Device | Top three tasks | Where it hurts today |
|---------|--------|-----------------|----------------------|
| Admin clerk | Desktop, all day | Enrol students, keep classes and teachers current, run reports | 27 sidebar items, 10-row pages on 7,000 students, no bulk actions, no export |
| Teacher | Phone in class, laptop in staffroom | Mark attendance, enter marks, see today's timetable | Sidebar is fixed-width on phones, roster rows are small targets, save bar scrolls away |
| Student | Low-end Android, mobile data | Today's timetable, my marks, notices | English only, heavy first load, no offline timetable |
| Parent | Phone, occasional | Child's attendance and marks, notices | English only, tabs hide the key numbers, no summary on the first screen |
| Principal / Vice Principal | Laptop | Attendance today, who has not marked, timetable status | Data is spread across Analytics, Attendance and Timetables |

---

## 3. Findings and fixes

Impact: H / M / L. Effort: S (under a day), M (days), L (a week or more).

### 3.1 Navigation and structure

| # | Status | Impact | Finding | Evidence | Fix | Effort |
|---|--------|--------|---------|----------|-----|--------|
| N1 | ✅ Done | H | The admin sidebar has 27 links in 8 groups. Timetable alone has 6 entries. New admins scan the whole list to find anything. | `layouts/nav/admin.ts` | Five hubs: Overview, People, Academics, Timetable, Operations. Timetable becomes one page with tabs (Timetables, Generate, Classrooms, Requirements, Settings). Audit Log and Orphaned Accounts were already under Settings; Automation now joins them. Ctrl+K search extended to actions ("Add student", "Mark attendance", ...). | M |
| N2 | ✅ Done | M | Two routes for the same page: `/subjects` and `/curriculum` both render `SubjectsCurriculum`. `/grades` redirects. | `admin.routes.tsx` | One canonical route per page; redirect the others. | S |
| N3 | ⬜ Not started | M | Breadcrumbs appear on one page only (`AddClass`). Detail pages rely on a Back button. | `AddClass.tsx`, detail pages | Breadcrumb in `os-page__header` for every nested page, generated from the route. | S |
| N4 | ⬜ Not started | M | Students, teachers and parents cannot collapse the sidebar; it takes 15 rem of a phone screen. | `PortalShell` `collapsible` only for admin | Responsive shell for all portals: sidebar becomes a drawer under 66 rem, bottom tab bar for the four most-used student and parent pages. | M |
| N5 | ⬜ Not started | L | The page title in the browser tab is always "OpenSchool". Tabs and history are unreadable. | `index.html` | `document.title = "<Page> · OpenSchool"` from a small `usePageTitle` hook in every page. | S |

### 3.2 Dashboards and first screens

| # | Status | Impact | Finding | Fix | Effort |
|---|--------|--------|---------|-----|--------|
| D1 | ⬜ Not started | H | Parent and student first screens do not answer "what matters today". Key numbers sit behind tabs. | A summary card row on the first screen: attendance this month, latest marks, next class, unread notices. Tabs stay for detail. | M |
| D2 | ⬜ Not started | M | After the setup wizard the admin lands on a dashboard with no next step. | A setup checklist card (academic year current, grades, subjects, teachers, first class, first students) with progress and links; hides itself when complete. | M |
| D3 | ⬜ Not started | M | Principal view is split across three pages. | One leadership overview: classes not marked today, timetables awaiting review, absence spikes, with links into each. Most data already exists in `LeadershipOverviewPanel`. | M |

### 3.3 Top tasks

| # | Status | Impact | Finding | Evidence | Fix | Effort |
|---|--------|--------|---------|----------|-----|--------|
| T1 | ✅ Done | H | Marking attendance on a phone: status buttons are small, the roster table scrolls horizontally, the Save button is below 40 rows. | `AttendanceMark.tsx`, `StudentAttendanceRow.tsx` | Card rows under 42rem with `data-label` stacking, 44px-min status buttons, sticky bottom bar showing unmarked count and Save, "Mark all Present" styled as the primary action so present-then-fix-exceptions is the visual default. Tap-to-toggle; swipe not implemented (kept in scope as a later enhancement, not core to the fix). | M |
| T2 | ⬜ Not started | H | Marks entry: one number input per row, no keyboard flow, absent checkbox separate. | `MarksEntryTable.tsx` | Enter moves to the next row, arrow keys move between rows, "AB" typed in the field marks absent, sticky header with subject and max marks, unsaved-change guard on navigation. | M |
| T3 | 🚧 Started | H | Lists default to 10 rows with 7,000 students. No sorting, no bulk actions, no export. | `DataGrid` default `pageSize=10` | Default is now 25 and the user's choice is remembered (`usePersistedPageSize`, wired into Students/Teachers/Guardians/Non-Academic Staff). Click-to-sort not done: the backend list endpoints have no `sort`/`order` param, so a client-side sort on a server-paginated grid would only reorder the current page — misleading. Needs a backend change first. Checkbox selection + "Enrol selected"/"Export CSV" not started. | M |
| T4 | ✅ Done | M | Pickers load every student to filter in the browser. | `EntityCombobox` | Server search with 2-character minimum and recent choices at the top. Already existed before this pass. | M |
| T5 | 🚧 Started | M | Idle logout after 30 minutes with no warning. A teacher mid-marking loses work. | `useIdleLogout.ts` | The 60-second warning modal with "Stay signed in" already existed before this pass. `sessionStorage` draft persistence + restore-after-sign-in does not exist anywhere in the codebase yet — still open. | S |
| T6 | ⬜ Not started | L | Filters are lost on refresh and cannot be shared. | `useListFilters` | Mirror filters into the URL query string. | S |

### 3.4 Feedback and system status

| # | Status | Impact | Finding | Evidence | Fix | Effort |
|---|--------|--------|---------|----------|-----|--------|
| F1 | ✅ Done | H | Loading states between routes and during auth are blank divs. Users see a white screen for up to a second. | `App.tsx` `Loading`, `ProtectedRoute.tsx` | `SkeletonShell` (header, sidebar, three grey blocks) for the initial/auth load; a nested `Suspense` around `PortalShell`'s `Outlet` now keeps header/sidebar mounted across route changes so only the content area shows `ContentSkeleton`. | S |
| F2 | 🚧 Started | H | Success feedback is inconsistent: "✓ Saved" text that vanishes after 2.5 seconds in two places, inline notifications elsewhere, nothing after most deletes. No toast system exists. | `GeneralSettingsTab.tsx`, `AttendanceMark.tsx`, 0 `ToastNotification` uses | `useToast()` + a Carbon `ToastNotification`/`ActionableNotification` stack now exists (`shared/ui/toast/`), rendered from `PortalShell`. Wired into the two sites named in the evidence column, plus the new bulk attendance-session button. Not yet wired into every mutation across the app, and no undo action is wired anywhere yet — both are follow-ups. | M |
| F3 | ⬜ Not started | M | 62 distinct "Failed to ..." messages, many showing the raw backend error string. | `MutationErrorNotification` fallback strings | Error copy pattern: what happened, what to do, in one sentence. Map backend error codes to copy in one file. Never show internal text. | M |
| F4 | ✅ Done | M | Delete confirmations are good, but "Saved, redirecting" waits 1.2 seconds on a timer. | `AttendanceMark.tsx` | Redirect immediately and show the toast on the destination page. | S |
| F5 | ⬜ Not started | L | Notification bell shows eight items with no "mark all read". | `NotificationsBell.tsx` | Add "Mark all read" and group by day. | S |

### 3.5 Forms and input

| # | Status | Impact | Finding | Evidence | Fix | Effort |
|---|--------|--------|---------|----------|-----|--------|
| I1 | 🚧 Started | H | 11 selects render with `labelText=""`. Sighted users lose context in filter bars; screen readers announce nothing. | `grep labelText=""` | Real labels, or `hideLabel` with a visible placeholder and an `aria-label`. All 9 flagged files fixed; `EntityCombobox` gained an `ariaLabel` prop since Carbon's `ComboBox` has no `hideLabel`. `FilterBar` should require a label per control — not done, would need a type-level change touching every call site. | S |
| I2 | ⬜ Not started | M | Long helper texts under selects (home classroom, medium) read like documentation. | `AddClass.tsx` | One short line; move the explanation to a tooltip icon. | S |
| I3 | ⬜ Not started | M | Validation appears only after blur. Submitting a long wizard step with three errors scrolls nowhere. | `SchoolSetup.tsx`, `FirstRunSetup.tsx` | On submit, focus the first invalid field and show an error summary at the top of the step. | S |
| I4 | ⬜ Not started | M | No unsaved-change guard when leaving marks, attendance or a profile in edit mode. | `useMarksDraft`, `useStudentProfileEditor` | `useBlocker` from React Router with a confirm dialog when `hasUnsaved` is true. | S |
| I5 | ⬜ Not started | L | Date pickers accept free text and the format hint is `YYYY-MM-DD` while the UI shows `en-LK` dates elsewhere. | `TermForm.tsx` | Carbon date picker in calendar mode with locale `en-LK`; keep ISO only in the API. | S |

### 3.6 Content and copy

| # | Status | Impact | Finding | Evidence | Fix | Effort |
|---|--------|--------|---------|----------|-----|--------|
| C1 | 🚧 Started | H | Spelling is mixed: "Enrol" 17 times, "Enroll" 61 times. | `grep` | British spelling throughout (Enrol, Enrolment, Organisation) in every user-facing string (nav label, tab, heading, button, toast). Internal identifiers/types (`EnrollmentStatus`, `useEnrollmentPicker`, ...) intentionally left alone — renaming those is a much larger, purely internal refactor with no user-visible payoff. An ESLint word-list rule is not set up yet. | S |
| C2 | ✅ Done | M | Five date formats in use (`toLocaleString()`, `en-LK` short, weekday long, and more). | `grep toLocale` | `shared/lib/date.ts` now has `formatDate`, `formatDateTime`, `formatMonth`, `formatLongDate`, `formatShortDayMonth`, `formatDayMonthYear`, all `en-LK`, all timezone-safe for bare `YYYY-MM-DD` strings. Every `toLocale*` call site in `src` (18 of them, including one that was hardcoded `en-US`) now goes through these. `no-restricted-properties` in `eslint.config.js` bans calling `toLocale*` directly outside `date.ts`. | S |
| C3 | 🚧 Started | M | Jargon and long labels: "Non-Academic Staff", "Mediums", "Leadership Positions", "Subject Requirements". | `nav/admin.ts` | Nav labels changed to "Staff", "Languages of instruction" (kept "Medium" in data/page titles), "Principal and VPs" (also updated the page's own `h1` to match). "Subject Requirements" no longer has a standalone nav entry — it's the "Requirements" tab inside the new Timetable hub (N1), which already shortens it. Not tested with clerks. | S |
| C4 | ⬜ Not started | M | Em-dashes as separators in 32 UI strings (between grade and class, for example). They render inconsistently on Android fonts. | grep for the character | Use a middle dot or a hyphen with spaces; add to the copy guide. | S |
| C5 | ⬜ Not started | M | Title Case and sentence case are mixed across headings and buttons. | Various | Sentence case for everything except proper nouns and nav labels. | S |
| C6 | ⬜ Not started | L | Subtitles under page titles are two-line explanations on some pages. | `Classes.tsx`, `Positions.tsx` | One line, under 80 characters, or none. | S |

### 3.7 Accessibility

| # | Status | Impact | Finding | Evidence | Fix | Effort |
|---|--------|--------|---------|----------|-----|--------|
| A1 | ✅ Done | H | Focus styles are defined in three places only; custom pills and select-rows rely on browser defaults. | `grep :focus` | A shared `:focus-visible` ring (`.os-pill`, `.os-select-card`, `.os-select-row`, `.os-focus-ring`) in `_utilities.scss`. The three existing context-specific `:focus` rules (dark header, search inputs) were left as-is — they're deliberate, not duplicates of this. `DataGrid`'s clickable rows also gained `tabIndex`/`onKeyDown`/the ring class; they were not keyboard-reachable at all before. | S |
| A2 | ⬜ Not started | M | 11 px text (`os-text-2xs`) used nine times for labels and metadata. | `grep os-text-2xs` | Minimum 12 px; reserve 11 px for uppercase stat labels only. | S |
| A3 | ⬜ Not started | M | No skip-to-content link; the sidebar is tabbed through on every page. | `PortalShell.tsx` | "Skip to main content" link as the first focusable element. | S |
| A4 | ⬜ Not started | M | Colour is the only signal in a few places (house dot, progress bars). | `Houses`, `AttendanceByClassSection` | Add text or a pattern next to the colour. | S |
| A5 | ⬜ Not started | M | No automated accessibility check. | CI | `axe-core` via Playwright on the six main pages; fail under 95 Lighthouse accessibility. | M |
| A6 | ⬜ Not started | L | `lang="en"` only. | `index.html` | Set per locale once localisation lands. | S |

### 3.8 Localisation

| # | Status | Impact | Finding | Fix | Effort |
|---|--------|--------|---------|-----|--------|
| L1 | ⬜ Not started | H | The UI is English only. Sinhala and Tamil exist as data (mediums) but not for users. Parents are the group most affected. | `react-i18next` with `en`, `si`, `ta`. Extract student and parent portal strings first (about 300), then teacher, then admin. Language switch in the header, remembered per user. Fonts: Noto Sans Sinhala and Noto Sans Tamil, self-hosted. | L |
| L2 | 🚧 Started | M | Numbers and dates are not locale-aware. | `Intl.NumberFormat("en-LK")` and the date helpers from C2. The date half is done (C2); number formatting is not. | S |

### 3.9 Mobile and low bandwidth

| # | Status | Impact | Finding | Evidence | Fix | Effort |
|---|--------|--------|---------|----------|-----|--------|
| M1 | 🚧 Started | H | Eight media queries in the whole stylesheet. Tables and the fixed sidebar do not adapt. | `grep @media` | Breakpoints at 42 rem and 66 rem; drawer sidebar; `DataGrid` switches to stacked cards under 42 rem when a `mobileRender` is provided. The 42rem breakpoint and the stacked-card pattern now exist and are proven on the attendance table (T1), but `DataGrid` itself has no generic `mobileRender` yet and the sidebar drawer (N4) isn't done. | M |
| M2 | ✅ Done | M | First load is 874 KB of CSS plus 1.2 MB of vendor JS. On 3G that is 8 to 12 seconds. | build report | Carbon trim and font self-hosting were already done before this pass (selective `@use` per component in `index.scss`, self-hosted IBM Plex in `_fonts.scss`). The skeleton shell is now done too (F1). | M |
| M3 | ⬜ Not started | M | No offline or stale indication. A parent on a train sees a spinner forever. | React Query defaults | "Showing saved data from 10:42" banner when a fetch fails but cache exists; retry button. Later: PWA cache of the timetable. | M |
| M4 | ⬜ Not started | L | Teachers print class registers. No print stylesheet. | none | `@media print` for rosters, timetables and attendance sessions: hide chrome, black text, page breaks per class. | S |

---

## 4. Prioritised backlog

| Now (next two sprints) | Next | Later |
|------------------------|------|-------|
| ✅ N1 five hubs, ✅ T1 phone attendance, 🚧 T3 list defaults (done)/sorting (⬜ blocked on backend `sort` param)/bulk actions (⬜), ✅ F1 skeleton shell, 🚧 F2 toast system (core done, not every mutation wired), 🚧 I1 labels (done)/`FilterBar` enforcement (⬜), ✅ C1 spelling, ✅ C2 dates, ✅ A1 focus ring, 🚧 T5 idle warning (was already ✅)/sessionStorage draft restore (⬜) | ⬜ N4 responsive shell, ⬜ D1 parent and student summary, ⬜ T2 marks keyboard flow, ⬜ D2 setup checklist, ⬜ F3 error copy map, ⬜ I4 unsaved guard, 🚧 M1 mobile tables (proven on attendance, not on `DataGrid` generically), ⬜ A5 axe in CI, ⬜ M4 print | ⬜ L1 Sinhala and Tamil, ⬜ M3 offline banner and PWA timetable, ⬜ D3 leadership overview |

Every item ships with its own before-and-after screenshot in the PR.

Not from the original backlog, added ad hoc: ✅ a global "Create sessions for all classes"
button on the admin Attendance page (`Attendance.tsx`) — bulk-creates a session for every
current class that doesn't already have one for the selected date, skipping classes that
do (the backend's `UNIQUE (class_id, date)` constraint is the safety net either way).

---

## 5. Rules going forward

Copy: sentence case, active voice, British spelling, no jargon a clerk would not
use, no em-dashes, errors as "What happened. What to do." in one sentence, buttons
as verb plus object ("Add student"), no exclamation marks.

Layout: one page header pattern (title, one-line subtitle, primary action on the
right), one card pattern (`SectionCard`), one list pattern (`FilterBar`,
`ActiveFilterTags`, `DataGrid`), one confirm pattern, one toast pattern.

Type and spacing: 12 px minimum body text, 14 px default, the 0.25 rem spacing
scale from `_utilities.scss`, colours only from `_tokens.scss`.

Touch: 44 px minimum target on phones, primary action reachable with the thumb
(sticky bottom bar on long forms and rosters).

States: every list has loading, empty, error and "filtered to nothing" states, and
every empty state leads to the action that fills it.

---

## 6. How to measure

| Metric | Target |
|--------|--------|
| Time to mark attendance for a class of 40 on a phone | under 90 seconds |
| Time to enrol a new student including guardian | under 3 minutes |
| Lighthouse accessibility on the six main pages | 95 or higher |
| Lighthouse performance on a mid-range Android profile | 90 or higher |
| Support questions about "where is X" in the first month after N1 | half of today |
| Parents who open the portal twice in a term after L1 | double |

Test with two clerks, three teachers and three parents before and after each
"Now" item. Ten minutes each, watching them do the top tasks, is enough.
