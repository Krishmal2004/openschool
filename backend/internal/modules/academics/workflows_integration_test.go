//go:build integration

package academics

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

type academicFixture struct {
	studentID                   uuid.UUID
	sourceYearID, targetYearID  uuid.UUID
	sourceClassID               uuid.UUID
	targetClassID, targetClassB uuid.UUID
	levelID, groupID            uuid.UUID
	subjectAID, subjectBID      uuid.UUID
}

func TestEnrollmentWorkflowAPIWithPostgres(t *testing.T) {
	gin.SetMode(gin.TestMode)
	pool := testdb.Open(t)
	fixture := seedAcademicFixture(t, pool)
	repository := NewEnrollmentRepository(pool)
	router := gin.New()
	group := router.Group("")
	RegisterEnrollmentRoutes(group, group, group, group, repository)

	request := SubmitEnrollmentRequest{
		AcademicYearID: fixture.sourceYearID.String(),
		LevelID:        fixture.levelID.String(),
		Picks: []EnrollmentPick{{
			GroupID: fixture.groupID.String(), SubjectID: fixture.subjectAID.String(),
		}},
	}
	submitted := performAcademicRequest(t, router, http.MethodPost, "/students/"+fixture.studentID.String()+"/enrollments", request)
	if submitted.Code != http.StatusOK || !bytes.Contains(submitted.Body.Bytes(), []byte(`"valid":true`)) {
		t.Fatalf("submit enrollment: code=%d body=%s", submitted.Code, submitted.Body.String())
	}
	assertEnrollmentSubject(t, pool, fixture.studentID, fixture.sourceYearID, fixture.subjectAID)

	invalid := request
	invalid.Picks = []EnrollmentPick{{GroupID: fixture.groupID.String(), SubjectID: uuid.NewString()}}
	rejected := performAcademicRequest(t, router, http.MethodPost, "/students/"+fixture.studentID.String()+"/enrollments", invalid)
	if rejected.Code != http.StatusUnprocessableEntity {
		t.Fatalf("invalid enrollment: code=%d body=%s", rejected.Code, rejected.Body.String())
	}
	assertEnrollmentSubject(t, pool, fixture.studentID, fixture.sourceYearID, fixture.subjectAID)

	enrollments := NewStudentEnrollment(repository)
	if err := enrollments.Confirm(context.Background(), fixture.studentID, fixture.levelID, fixture.sourceYearID); err != nil {
		t.Fatalf("confirm enrollment: %v", err)
	}
	request.Picks[0].SubjectID = fixture.subjectBID.String()
	locked := performAcademicRequest(t, router, http.MethodPost, "/students/"+fixture.studentID.String()+"/enrollments", request)
	if locked.Code != http.StatusConflict {
		t.Fatalf("locked enrollment update: code=%d body=%s", locked.Code, locked.Body.String())
	}
	assertEnrollmentSubject(t, pool, fixture.studentID, fixture.sourceYearID, fixture.subjectAID)

	unlockPath := "/students/" + fixture.studentID.String() + "/enrollments/lock/" + fixture.levelID.String() + "?academic_year_id=" + fixture.sourceYearID.String()
	unlocked := performAcademicRequest(t, router, http.MethodDelete, unlockPath, nil)
	if unlocked.Code != http.StatusOK {
		t.Fatalf("unlock enrollment: code=%d body=%s", unlocked.Code, unlocked.Body.String())
	}
	replaced := performAcademicRequest(t, router, http.MethodPost, "/students/"+fixture.studentID.String()+"/enrollments", request)
	if replaced.Code != http.StatusOK {
		t.Fatalf("replace enrollment: code=%d body=%s", replaced.Code, replaced.Body.String())
	}
	assertEnrollmentSubject(t, pool, fixture.studentID, fixture.sourceYearID, fixture.subjectBID)
}

func TestPromotionWorkflowAPIWithPostgres(t *testing.T) {
	gin.SetMode(gin.TestMode)
	pool := testdb.Open(t)
	fixture := seedAcademicFixture(t, pool)
	router := gin.New()
	RegisterPromotionRoutes(router.Group(""), NewPromotionService(NewPromotionRepository(pool)))

	previewPath := "/promotion/preview?source_year_id=" + fixture.sourceYearID.String() + "&target_year_id=" + fixture.targetYearID.String()
	preview := performAcademicRequest(t, router, http.MethodGet, previewPath, nil)
	if preview.Code != http.StatusOK {
		t.Fatalf("promotion preview: code=%d body=%s", preview.Code, preview.Body.String())
	}
	var rows []PromotionPreviewRow
	if err := json.Unmarshal(preview.Body.Bytes(), &rows); err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].StudentID != fixture.studentID.String() || rows[0].SuggestedClassID == nil || *rows[0].SuggestedClassID != fixture.targetClassID.String() {
		t.Fatalf("unexpected promotion preview: %+v", rows)
	}

	commit := CommitAssignmentsRequest{AcademicYearID: fixture.targetYearID.String(), Assignments: []AssignmentEntry{{StudentID: fixture.studentID.String(), ClassID: fixture.targetClassID.String()}}}
	committed := performAcademicRequest(t, router, http.MethodPost, "/promotion/commit", commit)
	if committed.Code != http.StatusOK || !bytes.Contains(committed.Body.Bytes(), []byte(`"assigned":1`)) {
		t.Fatalf("promotion commit: code=%d body=%s", committed.Code, committed.Body.String())
	}
	assertClassAssignment(t, pool, fixture.studentID, fixture.targetYearID, fixture.targetClassID)

	commit.Assignments[0].ClassID = fixture.targetClassB.String()
	reassigned := performAcademicRequest(t, router, http.MethodPost, "/promotion/commit", commit)
	if reassigned.Code != http.StatusOK {
		t.Fatalf("promotion reassignment: code=%d body=%s", reassigned.Code, reassigned.Body.String())
	}
	assertClassAssignment(t, pool, fixture.studentID, fixture.targetYearID, fixture.targetClassB)

	commit.Assignments[0].ClassID = fixture.sourceClassID.String()
	mismatched := performAcademicRequest(t, router, http.MethodPost, "/promotion/commit", commit)
	if mismatched.Code != http.StatusBadRequest {
		t.Fatalf("cross-year class commit: code=%d body=%s", mismatched.Code, mismatched.Body.String())
	}
	assertClassAssignment(t, pool, fixture.studentID, fixture.targetYearID, fixture.targetClassB)
}

func seedAcademicFixture(t *testing.T, pool *pgxpool.Pool) academicFixture {
	t.Helper()
	ctx := context.Background()
	fixture := academicFixture{}
	for _, year := range []struct {
		label         string
		start, end    string
		destinationID *uuid.UUID
	}{
		{label: "2025", start: "2025-01-01", end: "2025-12-31", destinationID: &fixture.sourceYearID},
		{label: "2026", start: "2026-01-01", end: "2026-12-31", destinationID: &fixture.targetYearID},
	} {
		if err := pool.QueryRow(ctx, "INSERT INTO academic_years (label, start_date, end_date) VALUES ($1, $2, $3) RETURNING id", year.label, year.start, year.end).Scan(year.destinationID); err != nil {
			t.Fatalf("seed academic year: %v", err)
		}
	}
	var gradeSixID, gradeSevenID uuid.UUID
	if err := pool.QueryRow(ctx, "INSERT INTO grades (name, sort_order) VALUES ('Grade 6', 6) RETURNING id").Scan(&gradeSixID); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, "INSERT INTO grades (name, sort_order) VALUES ('Grade 7', 7) RETURNING id").Scan(&gradeSevenID); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, "INSERT INTO classes (grade_id, academic_year_id, name) VALUES ($1, $2, 'A') RETURNING id", gradeSixID, fixture.sourceYearID).Scan(&fixture.sourceClassID); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, "INSERT INTO classes (grade_id, academic_year_id, name) VALUES ($1, $2, 'A') RETURNING id", gradeSevenID, fixture.targetYearID).Scan(&fixture.targetClassID); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, "INSERT INTO classes (grade_id, academic_year_id, name) VALUES ($1, $2, 'B') RETURNING id", gradeSevenID, fixture.targetYearID).Scan(&fixture.targetClassB); err != nil {
		t.Fatal(err)
	}

	studentUserID := uuid.New()
	if _, err := pool.Exec(ctx, "INSERT INTO users (id, email, full_name, role) VALUES ($1, 'academic-student@example.test', 'Academic Student', 'student')", studentUserID); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, "INSERT INTO student_profiles (user_id, full_name, index_number) VALUES ($1, 'Academic Student', 'ACADEMIC-001') RETURNING id", studentUserID).Scan(&fixture.studentID); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, "INSERT INTO class_students (class_id, student_id) VALUES ($1, $2)", fixture.sourceClassID, fixture.studentID); err != nil {
		t.Fatal(err)
	}

	if err := pool.QueryRow(ctx, "INSERT INTO levels (label, grade_id, sort_order) VALUES ('Grade 6 Choices', $1, 6) RETURNING id", gradeSixID).Scan(&fixture.levelID); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, "INSERT INTO selection_groups (level_id, label, min_select, max_select) VALUES ($1, 'Elective', 1, 1) RETURNING id", fixture.levelID).Scan(&fixture.groupID); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, "INSERT INTO subjects (name, code) VALUES ('Art', 'ART') RETURNING id").Scan(&fixture.subjectAID); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, "INSERT INTO subjects (name, code) VALUES ('Music', 'MUS') RETURNING id").Scan(&fixture.subjectBID); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, "INSERT INTO group_subjects (group_id, subject_id) VALUES ($1, $2), ($1, $3)", fixture.groupID, fixture.subjectAID, fixture.subjectBID); err != nil {
		t.Fatal(err)
	}
	return fixture
}

func performAcademicRequest(t *testing.T, router http.Handler, method, path string, body any) *httptest.ResponseRecorder {
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

func assertEnrollmentSubject(t *testing.T, pool *pgxpool.Pool, studentID, yearID, wantSubjectID uuid.UUID) {
	t.Helper()
	var count int
	var subjectID uuid.UUID
	if err := pool.QueryRow(context.Background(), `
		SELECT COUNT(*) OVER (), subject_id
		FROM student_subject_enrollments
		WHERE student_id = $1 AND academic_year_id = $2`, studentID, yearID).Scan(&count, &subjectID); err != nil {
		t.Fatalf("read enrollment: %v", err)
	}
	if count != 1 || subjectID != wantSubjectID {
		t.Fatalf("enrollment count=%d subject=%s, want one row for %s", count, subjectID, wantSubjectID)
	}
}

func assertClassAssignment(t *testing.T, pool *pgxpool.Pool, studentID, yearID, wantClassID uuid.UUID) {
	t.Helper()
	var count int
	var classID uuid.UUID
	if err := pool.QueryRow(context.Background(), `
		SELECT COUNT(*) OVER (), class_id
		FROM class_students
		WHERE student_id = $1 AND academic_year_id = $2`, studentID, yearID).Scan(&count, &classID); err != nil {
		t.Fatalf("read class assignment: %v", err)
	}
	if count != 1 || classID != wantClassID {
		t.Fatalf("class assignment count=%d class=%s, want one row for %s", count, classID, wantClassID)
	}
}
