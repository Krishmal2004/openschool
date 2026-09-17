# UX Review Playbook

Status: **DRAFT for approval. No code changed yet.**
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

| # | Impact | Finding | Evidence | Fix | Effort |
|---|--------|---------|----------|-----|--------|
| N1 | H | The admin sidebar has 27 links in 8 groups. Timetable alone has 6 entries. New admins scan the whole list to find anything. | `layouts/nav/admin.ts` | Five hubs: Overview, People, Academics, Timetable, Operations. Timetable becomes one page with tabs (Timetables, Generate, Classrooms, Requirements, Settings). Audit Log, Orphaned Accounts and Automation move under Settings. Keep Ctrl+K search and extend it to actions ("Add student", "Mark attendance"). | M |
| N2 | M | Two routes for the same page: `/subjects` and `/curriculum` both render `SubjectsCurriculum`. `/grades` redirects. | `admin.routes.tsx` | One canonical route per page; redirect the others. | S |
| N3 | M | Breadcrumbs appear on one page only (`AddClass`). Detail pages rely on a Back button. | `AddClass.tsx`, detail pages | Breadcrumb in `os-page__header` for every nested page, generated from the route. | S |
| N4 | M | Students, teachers and parents cannot collapse the sidebar; it takes 15 rem of a phone screen. | `PortalShell` `collapsible` only for admin | Responsive shell for all portals: sidebar becomes a drawer under 66 rem, bottom tab bar for the four most-used student and parent pages. | M |
| N5 | L | The page title in the browser tab is always "OpenSchool". Tabs and history are unreadable. | `index.html` | `document.title = "<Page> · OpenSchool"` from a small `usePageTitle` hook in every page. | S |

### 3.2 Dashboards and first screens

| # | Impact | Finding | Fix | Effort |
|---|--------|---------|-----|--------|
| D1 | H | Parent and student first screens do not answer "what matters today". Key numbers sit behind tabs. | A summary card row on the first screen: attendance this month, latest marks, next class, unread notices. Tabs stay for detail. | M |
| D2 | M | After the setup wizard the admin lands on a dashboard with no next step. | A setup checklist card (academic year current, grades, subjects, teachers, first class, first students) with progress and links; hides itself when complete. | M |
| D3 | M | Principal view is split across three pages. | One leadership overview: classes not marked today, timetables awaiting review, absence spikes, with links into each. Most data already exists in `LeadershipOverviewPanel`. | M |

### 3.3 Top tasks

| # | Impact | Finding | Evidence | Fix | Effort |
|---|--------|---------|----------|-----|--------|
| T1 | H | Marking attendance on a phone: status buttons are small, the roster table scrolls horizontally, the Save button is below 40 rows. | `AttendanceMark.tsx`, `StudentAttendanceRow.tsx` | Card rows on narrow screens with four large status buttons, sticky bottom bar showing unmarked count and Save, "Mark all present then fix exceptions" as the default flow, swipe or tap to toggle. | M |
| T2 | H | Marks entry: one number input per row, no keyboard flow, absent checkbox separate. | `MarksEntryTable.tsx` | Enter moves to the next row, arrow keys move between rows, "AB" typed in the field marks absent, sticky header with subject and max marks, unsaved-change guard on navigation. | M |
| T3 | H | Lists default to 10 rows with 7,000 students. No sorting, no bulk actions, no export. | `DataGrid` default `pageSize=10` | Default 25, remember the user's choice, click-to-sort on columns, checkbox selection with "Enrol selected into class" and "Export CSV" where it matters (Students, Teachers, Attendance sessions). Server pagination from the security playbook. | M |
| T4 | M | Pickers load every student to filter in the browser. | `EntityCombobox` | Server search with 2-character minimum and recent choices at the top. | M |
| T5 | M | Idle logout after 30 minutes with no warning. A teacher mid-marking loses work. | `useIdleLogout.ts` | 60-second warning modal with "Stay signed in"; keep unsaved drafts in `sessionStorage` and restore after sign-in. | S |
| T6 | L | Filters are lost on refresh and cannot be shared. | `useListFilters` | Mirror filters into the URL query string. | S |

### 3.4 Feedback and system status

| # | Impact | Finding | Evidence | Fix | Effort |
|---|--------|---------|----------|-----|--------|
| F1 | H | Loading states between routes and during auth are blank divs. Users see a white screen for up to a second. | `App.tsx` `Loading`, `ProtectedRoute.tsx` | Skeleton shell: header, sidebar and three grey blocks, so the layout never disappears. | S |
| F2 | H | Success feedback is inconsistent: "✓ Saved" text that vanishes after 2.5 seconds in two places, inline notifications elsewhere, nothing after most deletes. No toast system exists. | `GeneralSettingsTab.tsx`, `AttendanceMark.tsx`, 0 `ToastNotification` uses | One `useToast()` with a Carbon `ToastNotification` stack in `PortalShell`. Every mutation success shows a short toast with an undo where safe (unenrol, archive). | M |
| F3 | M | 62 distinct "Failed to ..." messages, many showing the raw backend error string. | `MutationErrorNotification` fallback strings | Error copy pattern: what happened, what to do, in one sentence. Map backend error codes to copy in one file. Never show internal text. | M |
| F4 | M | Delete confirmations are good, but "Saved, redirecting" waits 1.2 seconds on a timer. | `AttendanceMark.tsx` | Redirect immediately and show the toast on the destination page. | S |
| F5 | L | Notification bell shows eight items with no "mark all read". | `NotificationsBell.tsx` | Add "Mark all read" and group by day. | S |

### 3.5 Forms and input

| # | Impact | Finding | Evidence | Fix | Effort |
|---|--------|---------|----------|-----|--------|
| I1 | H | 11 selects render with `labelText=""`. Sighted users lose context in filter bars; screen readers announce nothing. | `grep labelText=""` | Real labels, or `hideLabel` with a visible placeholder and an `aria-label`. `FilterBar` should require a label per control. | S |
| I2 | M | Long helper texts under selects (home classroom, medium) read like documentation. | `AddClass.tsx` | One short line; move the explanation to a tooltip icon. | S |
| I3 | M | Validation appears only after blur. Submitting a long wizard step with three errors scrolls nowhere. | `SchoolSetup.tsx`, `FirstRunSetup.tsx` | On submit, focus the first invalid field and show an error summary at the top of the step. | S |
| I4 | M | No unsaved-change guard when leaving marks, attendance or a profile in edit mode. | `useMarksDraft`, `useStudentProfileEditor` | `useBlocker` from React Router with a confirm dialog when `hasUnsaved` is true. | S |
| I5 | L | Date pickers accept free text and the format hint is `YYYY-MM-DD` while the UI shows `en-LK` dates elsewhere. | `TermForm.tsx` | Carbon date picker in calendar mode with locale `en-LK`; keep ISO only in the API. | S |

### 3.6 Content and copy

| # | Impact | Finding | Evidence | Fix | Effort |
|---|--------|---------|----------|-----|--------|
| C1 | H | Spelling is mixed: "Enrol" 17 times, "Enroll" 61 times. | `grep` | British spelling throughout (Enrol, Enrolment, Organisation), enforced by a word list in ESLint or a lint script. | S |
| C2 | M | Five date formats in use (`toLocaleString()`, `en-LK` short, weekday long, and more). | `grep toLocale` | `formatDate`, `formatDateTime`, `formatMonth` in `shared/lib/date.ts` with `en-LK`; ban direct `toLocale*` calls by lint. | S |
| C3 | M | Jargon and long labels: "Non-Academic Staff", "Mediums", "Leadership Positions", "Subject Requirements". | `nav/admin.ts` | "Staff", "Languages of instruction" (keep "Medium" in data), "Principal and VPs", "Periods per subject". Test labels with two clerks. | S |
| C4 | M | Em-dashes as separators in 32 UI strings (between grade and class, for example). They render inconsistently on Android fonts. | grep for the character | Use a middle dot or a hyphen with spaces; add to the copy guide. | S |
| C5 | M | Title Case and sentence case are mixed across headings and buttons. | Various | Sentence case for everything except proper nouns and nav labels. | S |
| C6 | L | Subtitles under page titles are two-line explanations on some pages. | `Classes.tsx`, `Positions.tsx` | One line, under 80 characters, or none. | S |

### 3.7 Accessibility

| # | Impact | Finding | Evidence | Fix | Effort |
|---|--------|---------|----------|-----|--------|
| A1 | H | Focus styles are defined in three places only; custom pills and select-rows rely on browser defaults. | `grep :focus` | One `:focus-visible` ring token applied to every interactive `os-*` class. | S |
| A2 | M | 11 px text (`os-text-2xs`) used nine times for labels and metadata. | `grep os-text-2xs` | Minimum 12 px; reserve 11 px for uppercase stat labels only. | S |
| A3 | M | No skip-to-content link; the sidebar is tabbed through on every page. | `PortalShell.tsx` | "Skip to main content" link as the first focusable element. | S |
| A4 | M | Colour is the only signal in a few places (house dot, progress bars). | `Houses`, `AttendanceByClassSection` | Add text or a pattern next to the colour. | S |
| A5 | M | No automated accessibility check. | CI | `axe-core` via Playwright on the six main pages; fail under 95 Lighthouse accessibility. | M |
| A6 | L | `lang="en"` only. | `index.html` | Set per locale once localisation lands. | S |

### 3.8 Localisation

| # | Impact | Finding | Fix | Effort |
|---|--------|---------|-----|--------|
| L1 | H | The UI is English only. Sinhala and Tamil exist as data (mediums) but not for users. Parents are the group most affected. | `react-i18next` with `en`, `si`, `ta`. Extract student and parent portal strings first (about 300), then teacher, then admin. Language switch in the header, remembered per user. Fonts: Noto Sans Sinhala and Noto Sans Tamil, self-hosted. | L |
| L2 | M | Numbers and dates are not locale-aware. | `Intl.NumberFormat("en-LK")` and the date helpers from C2. | S |

### 3.9 Mobile and low bandwidth

| # | Impact | Finding | Evidence | Fix | Effort |
|---|--------|---------|----------|-----|--------|
| M1 | H | Eight media queries in the whole stylesheet. Tables and the fixed sidebar do not adapt. | `grep @media` | Breakpoints at 42 rem and 66 rem; drawer sidebar; `DataGrid` switches to stacked cards under 42 rem when a `mobileRender` is provided. | M |
| M2 | M | First load is 874 KB of CSS plus 1.2 MB of vendor JS. On 3G that is 8 to 12 seconds. | build report | Carbon trim and font self-hosting from the performance playbook; show the skeleton shell immediately. | M |
| M3 | M | No offline or stale indication. A parent on a train sees a spinner forever. | React Query defaults | "Showing saved data from 10:42" banner when a fetch fails but cache exists; retry button. Later: PWA cache of the timetable. | M |
| M4 | L | Teachers print class registers. No print stylesheet. | none | `@media print` for rosters, timetables and attendance sessions: hide chrome, black text, page breaks per class. | S |

---

## 4. Prioritised backlog

| Now (next two sprints) | Next | Later |
|------------------------|------|-------|
| N1 five hubs, T1 phone attendance, T3 list defaults and sorting, F1 skeleton shell, F2 toast system, I1 labels, C1 spelling, C2 dates, A1 focus ring, T5 idle warning | N4 responsive shell, D1 parent and student summary, T2 marks keyboard flow, D2 setup checklist, F3 error copy map, I4 unsaved guard, M1 mobile tables, A5 axe in CI, M4 print | L1 Sinhala and Tamil, M3 offline banner and PWA timetable, D3 leadership overview, T4 server pickers (with pagination phase) |

Every item ships with its own before-and-after screenshot in the PR.

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
