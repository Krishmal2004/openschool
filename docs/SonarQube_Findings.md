# SonarQube Findings Summary

**Analysis Date:** 20 September 2026
**Scope:** Frontend & Backend

## 1. High-Severity Findings

|  # | File                                                                       | Issue                                                                     | Line |  Effort |
| -: | -------------------------------------------------------------------------- | ------------------------------------------------------------------------- | ---: | ------: |
|  1 | `frontend/src/shared/styles/_layout.scss`                                  | Overridden `background-color` by shorthand `background`                   |  235 |    5min |
|  2 | `frontend/src/shared/styles/_layout.scss`                                  | Overridden `background-color` by shorthand `background`                   |  244 |    5min |
|  3 | `backend/cmd/api/main.go`                                                  | Reduce Cognitive Complexity **18 → 15**                                   |   79 |    8min |
|  4 | `backend/internal/architecture/dependencies_test.go`                       | Reduce Cognitive Complexity **28 → 15**                                   |   16 |   18min |
|  5 | `backend/internal/middleware/ratelimit.go`                                 | Reduce Cognitive Complexity **18 → 15**                                   |   57 |    8min |
|  6 | `backend/internal/middleware/student_access.go`                            | Reduce Cognitive Complexity **19 → 15**                                   |   16 |    9min |
|  7 | `backend/internal/modules/academics/classes.go`                            | Extract duplicated `"/classes/:id"` constant                              |  138 |    6min |
|  8 | `backend/internal/modules/academics/enrollment.go`                         | Reduce Cognitive Complexity **29 → 15**                                   |   87 |   19min |
|  9 | `backend/internal/modules/academics/enrollment.go`                         | Extract duplicated `"invalid student id"` constant                        |  247 |    8min |
| 10 | `backend/internal/modules/academics/module.go`                             | Extract duplicated `"/subjects/:id"` constant                             |   19 |    6min |
| 11 | `backend/internal/modules/academics/module.go`                             | Extract duplicated `"/streams/:id"` constant                              |   29 |    6min |
| 12 | `backend/internal/modules/academics/module.go`                             | Extract duplicated `"/streams/:id/groups/:groupId"` constant              |   34 |    6min |
| 13 | `backend/internal/modules/academics/promotion.go`                          | Reduce Cognitive Complexity **33 → 15**                                   |   47 |   23min |
| 14 | `backend/internal/modules/academics/stream.go`                             | Extract duplicated `"invalid id"` constant                                |  110 |    6min |
| 15 | `backend/internal/modules/academics/stream.go`                             | Extract duplicated `"invalid group id"` constant                          |  174 |    6min |
| 16 | `backend/internal/modules/academics/subject.go`                            | Extract duplicated `"invalid id"` constant                                |  105 |    6min |
| 17 | `backend/internal/modules/academics/term_mark_routes.go`                   | Reduce Cognitive Complexity **24 → 15**                                   |   31 |   14min |
| 18 | `backend/internal/modules/academics/term_mark_routes.go`                   | Extract duplicated `"invalid id"` constant                                |   33 |    8min |
| 19 | `backend/internal/modules/academics/workflows_integration_test.go`         | Reduce Cognitive Complexity **16 → 15**                                   |  123 |    6min |
| 20 | `backend/internal/modules/attendance/service.go`                           | Extract duplicated `"2006-01-02"` constant                                |  171 |    6min |
| 21 | `backend/internal/modules/attendance/service.go`                           | Reduce Cognitive Complexity **16 → 15**                                   |  279 |    6min |
| 22 | `backend/internal/modules/attendance/workflows_integration_test.go`        | Reduce Cognitive Complexity **17 → 15**                                   |   50 |    7min |
| 23 | `backend/internal/modules/auth/service.go`                                 | Extract duplicated `"secret mismatch"` constant                           |   98 |    6min |
| 24 | `backend/internal/modules/automation/data_retention.go`                    | Reduce Cognitive Complexity **17 → 15**                                   |   47 |    7min |
| 25 | `backend/internal/modules/automation/identity_erasure_retry.go`            | Reduce Cognitive Complexity **19 → 15**                                   |   48 |    9min |
| 26 | `backend/internal/modules/automation/security_audit.go`                    | Reduce Cognitive Complexity **16 → 15**                                   |   77 |    6min |
| 27 | `backend/internal/modules/curriculum/levels.go`                            | Extract duplicated `"/levels/:id"` constant                               |  216 |    6min |
| 28 | `backend/internal/modules/curriculum/preset.go`                            | Reduce Cognitive Complexity **78 → 15**                                   |   45 | 1h 8min |
| 29 | `backend/internal/modules/curriculum/preset_data.go`                       | Extract duplicated `"Health & Physical Education"` constant               |   52 |    6min |
| 30 | `backend/internal/modules/curriculum/workflows_integration_test.go`        | Reduce Cognitive Complexity **21 → 15**                                   |   19 |   11min |
| 31 | `backend/internal/modules/leadership/routes.go`                            | Reduce Cognitive Complexity **30 → 15**                                   |   14 |   20min |
| 32 | `backend/internal/modules/leadership/routes.go`                            | Extract duplicated `"invalid caller identity"` constant                   |   23 |    6min |
| 33 | `backend/internal/modules/leadership/workflows_integration_test.go`        | Reduce Cognitive Complexity **21 → 15**                                   |   28 |   11min |
| 34 | `backend/internal/modules/notifications/routes.go`                         | Extract duplicated `"invalid id"` constant                                |  102 |   12min |
| 35 | `backend/internal/modules/notifications/service.go`                        | Reduce Cognitive Complexity **46 → 15**                                   |  166 |   36min |
| 36 | `backend/internal/modules/notifications/service.go`                        | Reduce Cognitive Complexity **47 → 15**                                   |  265 |   37min |
| 37 | `backend/internal/modules/notifications/workflows_integration_test.go`     | Reduce Cognitive Complexity **24 → 15**                                   |   25 |   14min |
| 38 | `backend/internal/modules/people/guardians_routes.go`                      | Reduce Cognitive Complexity **28 → 15**                                   |   54 |   18min |
| 39 | `backend/internal/modules/people/guardians_routes.go`                      | Extract duplicated `"/guardians/:id"` constant                            |   68 |    6min |
| 40 | `backend/internal/modules/people/non_academic_staff.go`                    | Reduce Cognitive Complexity **24 → 15**                                   |  141 |   14min |
| 41 | `backend/internal/modules/people/non_academic_staff.go`                    | Extract duplicated `"/non-academic-staff/:id"` constant                   |  151 |    6min |
| 42 | `backend/internal/modules/people/student_portfolio.go`                     | Extract duplicated `"academic year"` constant                             |   73 |    8min |
| 43 | `backend/internal/modules/people/students_routes.go`                       | Reduce Cognitive Complexity **36 → 15**                                   |   50 |   26min |
| 44 | `backend/internal/modules/people/students_routes.go`                       | Extract duplicated `"/students/:id"` constant                             |   80 |    6min |
| 45 | `backend/internal/modules/people/teachers_guardians_integration_test.go`   | Reduce Cognitive Complexity **20 → 15**                                   |   79 |   10min |
| 46 | `backend/internal/modules/people/teachers_routes.go`                       | Reduce Cognitive Complexity **30 → 15**                                   |   39 |   20min |
| 47 | `backend/internal/modules/people/teachers_routes.go`                       | Extract duplicated `"/teachers/:id"` constant                             |   57 |    6min |
| 48 | `backend/internal/modules/reports/service.go`                              | Extract duplicated `"2006-01-02"` constant                                |   89 |    6min |
| 49 | `backend/internal/modules/reports/service.go`                              | Reduce Cognitive Complexity **20 → 15**                                   |  104 |   10min |
| 50 | `backend/internal/modules/reports/workflows_integration_test.go`           | Reduce Cognitive Complexity **18 → 15**                                   |   16 |    8min |
| 51 | `backend/internal/modules/school/grade.go`                                 | Extract duplicated `"invalid id"` constant                                |   94 |    6min |
| 52 | `backend/internal/modules/school/house.go`                                 | Extract duplicated `"invalid id"` constant                                |  217 |    6min |
| 53 | `backend/internal/modules/school/module.go`                                | Extract duplicated `"/houses/:id"` constant                               |   18 |    6min |
| 54 | `backend/internal/modules/school/module.go`                                | Extract duplicated `"/grades/:id"` constant                               |   30 |    6min |
| 55 | `backend/internal/modules/school/repository.go`                            | Extract duplicated `"2006-01-02"` constant                                |  214 |    8min |
| 56 | `backend/internal/modules/school/workflows_integration_test.go`            | Reduce Cognitive Complexity **25 → 15**                                   |   20 |   15min |
| 57 | `backend/internal/modules/selfservice/parent.go`                           | Extract duplicated `"invalid id"` constant                                |   76 |    6min |
| 58 | `backend/internal/modules/selfservice/teacher.go`                          | Extract duplicated `"invalid caller identity"` constant                   |   45 |   10min |
| 59 | `backend/internal/modules/selfservice/teacher.go`                          | Extract duplicated `"no teacher profile linked to this account"` constant |   51 |   10min |
| 60 | `backend/internal/modules/selfservice/teacher.go`                          | Extract duplicated `"no current academic year configured"` constant       |   57 |   10min |
| 61 | `backend/internal/modules/studentleadership/routes.go`                     | Reduce Cognitive Complexity **51 → 15**                                   |   36 |   41min |
| 62 | `backend/internal/modules/studentleadership/routes.go`                     | Extract duplicated `"invalid society id"` constant                        |  122 |   10min |
| 63 | `backend/internal/modules/studentleadership/workflows_integration_test.go` | Reduce Cognitive Complexity **23 → 15**                                   |   27 |   13min |
| 64 | `backend/internal/modules/timetable/generation.go`                         | Reduce Cognitive Complexity **58 → 15**                                   |   92 |   48min |
| 65 | `backend/internal/modules/timetable/generation.go`                         | Reduce Cognitive Complexity **57 → 15**                                   |  240 |   47min |
| 66 | `backend/internal/modules/timetable/grade_section.go`                      | Extract duplicated `"invalid id"` constant                                |  384 |   16min |
| 67 | `backend/internal/modules/timetable/module.go`                             | Extract duplicated `"/grade-sections/:id"` constant                       |   41 |    6min |
| 68 | `backend/internal/modules/timetable/validation.go`                         | Reduce Cognitive Complexity **54 → 15**                                   |   56 |   44min |
| 69 | `backend/internal/modules/timetable/workflows_integration_test.go`         | Reduce Cognitive Complexity **18 → 15**                                   |   38 |    8min |
| 70 | `frontend/.../features/academics/pages/admin/Promotion.tsx`                | Reduce Cognitive Complexity **18 → 15**                                   |   15 |    8min |
| 71 | `frontend/.../features/guardians/components/GuardianDetail.tsx`            | Reduce Cognitive Complexity **21 → 15**                                   |   14 |   11min |
| 72 | `frontend/src/features/school/hooks/useSchoolSetupSubmit.ts`               | Reduce Cognitive Complexity **78 → 15**                                   |   70 | 1h 8min |
| 73 | `frontend/src/features/students/hooks/useEnrollmentPicker.ts`              | Reduce Cognitive Complexity **25 → 15**                                   |   21 |   15min |
| 74 | `frontend/src/shared/testing/setup.ts`                                     | Unexpected empty method `observe`                                         |    2 |    5min |
| 75 | `frontend/src/shared/testing/setup.ts`                                     | Unexpected empty method `unobserve`                                       |    3 |    5min |
| 76 | `frontend/src/shared/testing/setup.ts`                                     | Unexpected empty method `disconnect`                                      |    — |    5min |

---

## 2. Medium-Severity Findings

The main medium findings are:

| Area                        | Examples                                                                              |
| --------------------------- | ------------------------------------------------------------------------------------- |
| Security                    | Unsafe pseudorandom number generator in `PromotionGroup.tsx:20`                       |
| Accessibility / Reliability | Non-native interactive elements in `NotificationsBell.tsx:56` and `Timetables.tsx:82` |
| Maintainability             | Excessive function parameters                                                         |
| Maintainability             | Nested ternary expressions                                                            |
| Maintainability             | Nested template literals                                                              |
| Reliability                 | Complex regular expression                                                            |
| Maintainability             | `NaN` instead of `Number.NaN`                                                         |

### Backend Parameter Issues

| File                          | Issue                      | Line |
| ----------------------------- | -------------------------- | ---: |
| `academics/repository.go`     | Function has 16 parameters |  106 |
| `academics/repository.go`     | Function has 8 parameters  |  294 |
| `notifications/repository.go` | Function has 11 parameters |  224 |
| `people/repository.go`        | Function has 8 parameters  |  259 |
| `selfservice/routes.go`       | Function has 11 parameters |   11 |
| `timetable/generation.go`     | Function has 10 parameters |  240 |

### Frontend Issues

Most remaining frontend findings are **nested ternaries**, primarily in:

* `Promotion.tsx`
* `GuardianDetail.tsx`
* `Attendance.tsx`
* `StaffAttendance.tsx`
* `TeacherMyAttendance.tsx`
* `AddSubjectModal.tsx`
* `ClassMarks.tsx`
* `StudentMarks.tsx`
* `TeacherMarks.tsx`
* `NotificationCenter.tsx`
* `StudentRecordsRollup.tsx`
* `Positions.tsx`
* `AttendanceByClassSection.tsx`

---

## 3. Code Duplication

| File                                        | Duplicated Lines |         % |
| ------------------------------------------- | ---------------: | --------: |
| `selfservice/routes_test.go`                |               26 | **59.1%** |
| `dashboard/workflows_integration_test.go`   |               16 | **30.8%** |
| `db/sqlc/job_scheduler.sql.go`              |               50 | **23.4%** |
| `selfservice/teacher.go`                    |               51 | **23.1%** |
| `db/sqlc/notification_recipients.sql.go`    |               44 | **20.8%** |
| `reports/workflows_integration_test.go`     |               15 | **20.5%** |
| `selfservice/workflows_integration_test.go` |               10 | **18.9%** |
| `db/sqlc/houses.sql.go`                     |               63 | **18.7%** |
| `db/sqlc/teachers.sql.go`                   |              120 | **16.7%** |
| `db/sqlc/users.sql.go`                      |               54 | **16.7%** |
| `db/sqlc/timetables.sql.go`                 |               87 | **14.0%** |
| `db/sqlc/staff_attendance.sql.go`           |               54 | **13.7%** |
| `db/sqlc/subjects.sql.go`                   |               24 | **13.7%** |
| `automation/people_compliance.go`           |               19 | **13.6%** |
| `db/sqlc/notifications.sql.go`              |               46 | **13.5%** |

---

## 4. Reliability & Coverage

### Reliability

| File                         | Rating | Issues |
| ---------------------------- | :----: | -----: |
| `shared/styles/_layout.scss` |  **D** |      2 |
| `NotificationsBell.tsx`      |  **C** |      2 |
| `Timetables.tsx`             |  **C** |      2 |
| `usePersistedPageSize.ts`    |  **C** |      1 |
| `validation.ts`              |  **C** |      1 |
| `AuditLog.tsx`               |  **B** |      2 |
| `Automation.tsx`             |  **B** |      1 |
| `ClassesStep.tsx`            |  **B** |      1 |
| `jwt.ts`                     |  **B** |      2 |
| `useSchoolSetupSubmit.ts`    |  **B** |      1 |

### New-Code Coverage

**277 lines** currently require coverage.

Largest areas:

| File                        | Lines |
| --------------------------- | ----: |
| `useAttendanceMarking.ts`   |    29 |
| `useUnsavedChangesGuard.ts` |    22 |
| `shared/lib/date.ts`        |    21 |
| `eslint.config.js`          |    18 |
| `StudentDashboard.tsx`      |    17 |
| `MarksEntryTable.tsx`       |    16 |
| `Attendance.tsx`            |    15 |
| `usePersistedPageSize.ts`   |    13 |
| `ToastContext.tsx`          |    12 |

---

## 5. Summary

### Main Areas Requiring Work

1. **76 high-severity issues**
2. **74 medium-severity issues**
3. Cognitive complexity hotspots, including values of **78, 58, 57, 54, and 51**
4. Code duplication reaching **59.1%**
5. New-code reliability rating of **D**
6. **277 lines** requiring test coverage
7. Multiple functions with excessive parameter counts
8. Numerous frontend nested ternary expressions

### Suggested Work Sequence

**Reliability fixes → High-complexity functions → Duplicated literals → Parameter reduction → Frontend maintainability → Test coverage → Generated SQLC review**

> **Note:** `backend/db/sqlc/*.sql.go` files appear to be generated code. These should generally be reviewed at the SQL/schema/code-generation level rather than manually edited.
