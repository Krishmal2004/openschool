//go:build integration

package reports

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	attendancemodule "github.com/openschool-org/openschool/internal/modules/attendance"
	"github.com/openschool-org/openschool/internal/testutil/testdb"
)

func TestReportExportsWithPostgres(t *testing.T) {
	pool := testdb.Open(t)
	ctx := context.Background()
	adminID, teacherUserID := uuid.New(), uuid.New()
	for _, user := range []struct {
		id                uuid.UUID
		email, name, role string
	}{{adminID, "report-admin@example.test", "Report Admin", "admin"}, {teacherUserID, "report-teacher@example.test", "Report Teacher", "teacher"}} {
		if _, err := pool.Exec(ctx, "INSERT INTO users (id, email, full_name, role) VALUES ($1, $2, $3, $4)", user.id, user.email, user.name, user.role); err != nil {
			t.Fatal(err)
		}
	}
	var yearID, gradeID, classID, studentID, termID, subjectID, sessionID uuid.UUID
	if _, err := pool.Exec(ctx, "INSERT INTO teacher_profiles (user_id, full_name, employee_number, joined_date, nic_number) VALUES ($1, 'Report Teacher', 'REPORT-T-1', '2020-01-01', '199000000031')", teacherUserID); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, "INSERT INTO academic_years (label, start_date, end_date, is_current) VALUES ('2026 Reports', '2026-01-01', '2026-12-31', TRUE) RETURNING id").Scan(&yearID); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, "INSERT INTO grades (name, sort_order) VALUES ('Reports Grade 10', 10) RETURNING id").Scan(&gradeID); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, "INSERT INTO classes (grade_id, academic_year_id, name) VALUES ($1, $2, 'A') RETURNING id", gradeID, yearID).Scan(&classID); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, "INSERT INTO student_profiles (full_name, index_number) VALUES ('Report Student', 'REPORT-S-1') RETURNING id").Scan(&studentID); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, "INSERT INTO class_students (class_id, student_id) VALUES ($1, $2)", classID, studentID); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, "INSERT INTO terms (academic_year_id, name, start_date, end_date, is_current, sort_order) VALUES ($1, 'Report Term', '2026-01-01', '2026-04-30', TRUE, 1) RETURNING id", yearID).Scan(&termID); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, "INSERT INTO subjects (name, code) VALUES ('Report Mathematics', 'REPORT-MATH') RETURNING id").Scan(&subjectID); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, "INSERT INTO term_marks (student_id, subject_id, term_id, marks, max_marks, entered_by) VALUES ($1, $2, $3, 82, 100, $4)", studentID, subjectID, termID, adminID); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, "INSERT INTO attendance_sessions (class_id, taken_by, date) VALUES ($1, $2, '2026-03-10') RETURNING id", classID, teacherUserID).Scan(&sessionID); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, "INSERT INTO attendance_records (session_id, student_id, status, note) VALUES ($1, $2, 'absent', 'medical')", sessionID, studentID); err != nil {
		t.Fatal(err)
	}

	service := NewService(NewRepository(pool), attendancemodule.NewService(attendancemodule.NewRepository(pool), nil, nil, nil))
	attendancePDF, err := service.ExportAttendance(ctx, AttendanceReportRequest{ClassID: classID.String(), From: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC), To: time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC), Columns: []string{"student", "status"}})
	if err != nil || !bytes.HasPrefix(attendancePDF, []byte("%PDF")) {
		t.Fatalf("attendance export bytes=%d err=%v", len(attendancePDF), err)
	}
	marksPDF, err := service.ExportMarks(ctx, MarksReportRequest{ClassID: classID.String(), TermID: termID.String(), SubjectID: subjectID.String(), Columns: []string{"student", "percentage"}})
	if err != nil || !bytes.HasPrefix(marksPDF, []byte("%PDF")) {
		t.Fatalf("marks export bytes=%d err=%v", len(marksPDF), err)
	}
}
