//go:build integration

package timetable

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/openschool-org/openschool/internal/testutil/testdb"
)

type timetableFixture struct {
	adminUserID, reviewerUserID uuid.UUID
	reviewerTeacherID           uuid.UUID
	yearID, gradeID, classID    uuid.UUID
	subjectID                   uuid.UUID
}

type timetableNotification struct {
	title      string
	recipients []uuid.UUID
}

type recordingTimetableNotifier struct{ calls []timetableNotification }

func (n *recordingTimetableNotifier) SendDirect(_ context.Context, title, _, _, _ string, _ uuid.UUID, recipients []uuid.UUID) error {
	n.calls = append(n.calls, timetableNotification{title: title, recipients: append([]uuid.UUID(nil), recipients...)})
	return nil
}

func TestTimetableDraftReviewAndPublicationAPIWithPostgres(t *testing.T) {
	gin.SetMode(gin.TestMode)
	pool := testdb.Open(t)
	fixture := seedTimetableFixture(t, pool)
	notifier := &recordingTimetableNotifier{}
	adminRouter := timetableAdminRouter(pool, fixture.adminUserID, notifier)
	reviewerRouter := timetableReviewerRouter(pool, fixture.reviewerUserID, notifier)

	created := performTimetableRequest(t, adminRouter, http.MethodPost, "/timetables", timetableCreateRequest{AcademicYearID: fixture.yearID, ClassID: fixture.classID})
	if created.Code != http.StatusCreated {
		t.Fatalf("create timetable: code=%d body=%s", created.Code, created.Body.String())
	}
	var timetable Timetable
	if err := json.Unmarshal(created.Body.Bytes(), &timetable); err != nil {
		t.Fatal(err)
	}
	if timetable.ID == uuid.Nil || timetable.Status != statusDraft || timetable.Version != 1 || timetable.CreatedBy != fixture.adminUserID {
		t.Fatalf("unexpected timetable: %+v", timetable)
	}

	entriesPath := "/timetables/" + timetable.ID.String() + "/entries"
	entryCommand := saveTimetableEntriesCommand{Entries: []timetableEntryCommand{{DayOfWeek: 1, PeriodNumber: 1, SubjectID: &fixture.subjectID, TeacherID: &fixture.reviewerTeacherID}}}
	saved := performTimetableRequest(t, adminRouter, http.MethodPut, entriesPath, entryCommand)
	if saved.Code != http.StatusOK || !bytes.Contains(saved.Body.Bytes(), []byte(`"subject_name":"Timetable Mathematics"`)) {
		t.Fatalf("save timetable entry: code=%d body=%s", saved.Code, saved.Body.String())
	}
	assertTimetableEntry(t, pool, timetable.ID, fixture.subjectID, fixture.reviewerTeacherID, 1)

	validated := performTimetableRequest(t, adminRouter, http.MethodGet, "/timetables/"+timetable.ID.String()+"/validate", nil)
	if validated.Code != http.StatusOK || !bytes.Contains(validated.Body.Bytes(), []byte(`"valid":true`)) {
		t.Fatalf("validate timetable: code=%d body=%s", validated.Code, validated.Body.String())
	}
	submitted := performTimetableRequest(t, adminRouter, http.MethodPost, "/timetables/"+timetable.ID.String()+"/submit", nil)
	if submitted.Code != http.StatusOK || !bytes.Contains(submitted.Body.Bytes(), []byte(`"status":"under_review"`)) {
		t.Fatalf("submit timetable: code=%d body=%s", submitted.Code, submitted.Body.String())
	}
	lockedEdit := performTimetableRequest(t, adminRouter, http.MethodPut, entriesPath, entryCommand)
	if lockedEdit.Code != http.StatusBadRequest {
		t.Fatalf("edit submitted timetable: code=%d body=%s", lockedEdit.Code, lockedEdit.Body.String())
	}

	approved := performTimetableRequest(t, reviewerRouter, http.MethodPost, "/timetables/"+timetable.ID.String()+"/approve", map[string]string{"comment": "ready"})
	if approved.Code != http.StatusOK || !bytes.Contains(approved.Body.Bytes(), []byte(`"status":"approved"`)) {
		t.Fatalf("approve timetable: code=%d body=%s", approved.Code, approved.Body.String())
	}
	published := performTimetableRequest(t, adminRouter, http.MethodPost, "/timetables/"+timetable.ID.String()+"/publish", nil)
	if published.Code != http.StatusOK || !bytes.Contains(published.Body.Bytes(), []byte(`"status":"published"`)) {
		t.Fatalf("publish timetable: code=%d body=%s", published.Code, published.Body.String())
	}
	assertTimetableState(t, pool, timetable.ID, statusPublished, 3)
	assertTimetableNotifications(t, notifier, fixture.reviewerUserID, fixture.adminUserID)

	revised := performTimetableRequest(t, adminRouter, http.MethodPost, "/timetables/"+timetable.ID.String()+"/revise", nil)
	if revised.Code != http.StatusCreated {
		t.Fatalf("revise timetable: code=%d body=%s", revised.Code, revised.Body.String())
	}
	var revision Timetable
	if err := json.Unmarshal(revised.Body.Bytes(), &revision); err != nil {
		t.Fatal(err)
	}
	if revision.Status != statusDraft || revision.Version != 2 || revision.ParentTimetableID == nil || *revision.ParentTimetableID != timetable.ID {
		t.Fatalf("unexpected revision: %+v", revision)
	}
	assertTimetableEntry(t, pool, revision.ID, fixture.subjectID, fixture.reviewerTeacherID, 1)
	deletedEntry := performTimetableRequest(t, adminRouter, http.MethodDelete, "/timetables/"+revision.ID.String()+"/entries/1/1", nil)
	if deletedEntry.Code != http.StatusOK {
		t.Fatalf("delete revision entry: code=%d body=%s", deletedEntry.Code, deletedEntry.Body.String())
	}
	assertTimetableEntryCount(t, pool, revision.ID, 0)
}

func timetableAdminRouter(pool *pgxpool.Pool, actorID uuid.UUID, notifier timetableNotifier) *gin.Engine {
	router := gin.New()
	admin := router.Group("")
	admin.Use(timetableActor(actorID))
	RegisterTimetableCRUDRoutes(admin, admin, pool)
	RegisterTimetableEntryRoutes(admin, admin, pool)
	RegisterTimetableValidationRoute(admin, pool)
	RegisterTimetableWorkflowRoutes(admin, admin, router.Group("/teacher"), NewWorkflowRepository(pool), NewWorkflowValidator(pool), notifier)
	return router
}

func timetableReviewerRouter(pool *pgxpool.Pool, actorID uuid.UUID, notifier timetableNotifier) *gin.Engine {
	router := gin.New()
	teacher := router.Group("")
	teacher.Use(timetableActor(actorID))
	RegisterTimetableWorkflowRoutes(router.Group("/admin"), teacher, teacher, NewWorkflowRepository(pool), NewWorkflowValidator(pool), notifier)
	return router
}

func timetableActor(userID uuid.UUID) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userID", userID.String())
		c.Next()
	}
}

func performTimetableRequest(t *testing.T, router http.Handler, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var payload []byte
	if body != nil {
		var err error
		payload, err = json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
	}
	request := httptest.NewRequest(method, path, bytes.NewReader(payload))
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

func seedTimetableFixture(t *testing.T, pool *pgxpool.Pool) timetableFixture {
	t.Helper()
	ctx := context.Background()
	f := timetableFixture{adminUserID: uuid.New(), reviewerUserID: uuid.New()}
	if _, err := pool.Exec(ctx, "INSERT INTO users (id, email, full_name, role) VALUES ($1, 'timetable-admin@example.test', 'Timetable Admin', 'admin'), ($2, 'timetable-reviewer@example.test', 'Timetable Reviewer', 'teacher')", f.adminUserID, f.reviewerUserID); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, "INSERT INTO teacher_profiles (user_id, full_name, employee_number, joined_date, nic_number) VALUES ($1, 'Timetable Reviewer', 'TIME-T-001', '2020-01-01', 'TIME-NIC-T-001') RETURNING id", f.reviewerUserID).Scan(&f.reviewerTeacherID); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, "INSERT INTO academic_years (label, start_date, end_date, is_current) VALUES ('2026 Timetable', '2026-01-01', '2026-12-31', TRUE) RETURNING id").Scan(&f.yearID); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, "INSERT INTO grades (name, sort_order) VALUES ('Grade 11 Timetable', 11) RETURNING id").Scan(&f.gradeID); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, "INSERT INTO classes (grade_id, academic_year_id, form_teacher_id, name) VALUES ($1, $2, $3, 'A') RETURNING id", f.gradeID, f.yearID, f.reviewerTeacherID).Scan(&f.classID); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, "INSERT INTO subjects (name, code) VALUES ('Timetable Mathematics', 'TIME-MATH') RETURNING id").Scan(&f.subjectID); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, "INSERT INTO teacher_subjects (teacher_id, subject_id) VALUES ($1, $2)", f.reviewerTeacherID, f.subjectID); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, "INSERT INTO class_subject_teachers (class_id, subject_id, teacher_id) VALUES ($1, $2, $3)", f.classID, f.subjectID, f.reviewerTeacherID); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, "INSERT INTO section_heads (academic_year_id, grade_id, teacher_id) VALUES ($1, $2, $3)", f.yearID, f.gradeID, f.reviewerTeacherID); err != nil {
		t.Fatal(err)
	}
	return f
}

func assertTimetableEntry(t *testing.T, pool *pgxpool.Pool, timetableID, subjectID, teacherID uuid.UUID, wantCount int) {
	t.Helper()
	var count int
	var gotSubjectID, gotTeacherID uuid.UUID
	if err := pool.QueryRow(context.Background(), "SELECT COUNT(*) OVER (), subject_id, teacher_id FROM timetable_entries WHERE timetable_id = $1", timetableID).Scan(&count, &gotSubjectID, &gotTeacherID); err != nil {
		t.Fatal(err)
	}
	if count != wantCount || gotSubjectID != subjectID || gotTeacherID != teacherID {
		t.Fatalf("entry count=%d subject=%s teacher=%s", count, gotSubjectID, gotTeacherID)
	}
}

func assertTimetableEntryCount(t *testing.T, pool *pgxpool.Pool, timetableID uuid.UUID, want int) {
	t.Helper()
	var count int
	if err := pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM timetable_entries WHERE timetable_id = $1", timetableID).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != want {
		t.Fatalf("entry count=%d, want %d", count, want)
	}
}

func assertTimetableState(t *testing.T, pool *pgxpool.Pool, timetableID uuid.UUID, wantStatus string, wantHistory int) {
	t.Helper()
	var status string
	if err := pool.QueryRow(context.Background(), "SELECT status FROM timetables WHERE id = $1", timetableID).Scan(&status); err != nil {
		t.Fatal(err)
	}
	if status != wantStatus {
		t.Fatalf("timetable status=%q, want %q", status, wantStatus)
	}
	var history int
	if err := pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM timetable_status_history WHERE timetable_id = $1", timetableID).Scan(&history); err != nil {
		t.Fatal(err)
	}
	if history != wantHistory {
		t.Fatalf("status history count=%d, want %d", history, wantHistory)
	}
}

func assertTimetableNotifications(t *testing.T, notifier *recordingTimetableNotifier, reviewerID, creatorID uuid.UUID) {
	t.Helper()
	if len(notifier.calls) != 3 {
		t.Fatalf("notification calls=%d, want 3", len(notifier.calls))
	}
	if notifier.calls[0].title != "Timetable Submitted for Review" || len(notifier.calls[0].recipients) != 1 || notifier.calls[0].recipients[0] != reviewerID {
		t.Fatalf("unexpected submission notification: %+v", notifier.calls[0])
	}
	if notifier.calls[1].title != "Timetable Approved" || len(notifier.calls[1].recipients) != 1 || notifier.calls[1].recipients[0] != creatorID {
		t.Fatalf("unexpected approval notification: %+v", notifier.calls[1])
	}
	if notifier.calls[2].title != "Timetable Published" || len(notifier.calls[2].recipients) != 1 || notifier.calls[2].recipients[0] != reviewerID {
		t.Fatalf("unexpected publication notification: %+v", notifier.calls[2])
	}
}
