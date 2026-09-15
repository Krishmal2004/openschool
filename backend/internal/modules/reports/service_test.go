package reports

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	attendancemodule "github.com/openschool-org/openschool/internal/modules/attendance"
)

type reportStore struct {
	class      string
	term       string
	subject    string
	marks      []markRow
	classID    uuid.UUID
	termID     uuid.UUID
	subjectID  uuid.UUID
	marksCalls int
}

func (s *reportStore) className(_ context.Context, id uuid.UUID) (string, error) {
	s.classID = id
	return s.class, nil
}

func (s *reportStore) termName(_ context.Context, id uuid.UUID) (string, error) {
	s.termID = id
	return s.term, nil
}

func (s *reportStore) subjectName(_ context.Context, id uuid.UUID) (string, error) {
	s.subjectID = id
	return s.subject, nil
}

func (s *reportStore) classMarks(_ context.Context, classID, termID, subjectID uuid.UUID) ([]markRow, error) {
	s.classID, s.termID, s.subjectID = classID, termID, subjectID
	s.marksCalls++
	return s.marks, nil
}

type attendanceReader struct {
	rows    []attendancemodule.ReportRow
	classID uuid.UUID
	from    time.Time
	to      time.Time
}

func (r *attendanceReader) ListForClassInRange(_ context.Context, classID uuid.UUID, from, to time.Time) ([]attendancemodule.ReportRow, error) {
	r.classID, r.from, r.to = classID, from, to
	return r.rows, nil
}

func numeric(t *testing.T, value string) pgtype.Numeric {
	t.Helper()
	var result pgtype.Numeric
	if err := result.Scan(value); err != nil {
		t.Fatal(err)
	}
	return result
}

func TestResolveColumnsPreservesRequestedOrderAndFallsBack(t *testing.T) {
	columns := resolveColumns([]string{"status", "unknown", "student"}, attendanceReportColumns)
	if len(columns) != 2 || columns[0] != "status" || columns[1] != "student" {
		t.Fatalf("unexpected columns: %v", columns)
	}
	fallback := resolveColumns([]string{"unknown"}, attendanceReportColumns)
	if len(fallback) != len(attendanceReportColumns) || fallback[0] != "date" {
		t.Fatalf("unexpected fallback: %v", fallback)
	}
}

func TestFormatNumericPreservesReportDisplay(t *testing.T) {
	if got := formatNumeric(pgtype.Numeric{}); got != "-" {
		t.Fatalf("null numeric=%q, want -", got)
	}
	if got := formatNumeric(numeric(t, "85.00")); got != "85" {
		t.Fatalf("numeric=%q, want 85", got)
	}
}

func TestExportAttendanceReadsRangeAndReturnsPDF(t *testing.T) {
	classID := uuid.New()
	from := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 9, 30, 0, 0, 0, 0, time.UTC)
	store := &reportStore{class: "10-A"}
	attendance := &attendanceReader{rows: []attendancemodule.ReportRow{{
		StudentName: "Student One", StudentIndex: "S001", SessionDate: pgtype.Date{Time: from, Valid: true}, Status: "present",
	}}}
	data, err := NewService(store, attendance).ExportAttendance(context.Background(), AttendanceReportRequest{ClassID: classID.String(), From: from, To: to, Columns: []string{"student", "status"}})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(data, []byte("%PDF")) || attendance.classID != classID || attendance.from != from || attendance.to != to {
		t.Fatalf("invalid export: prefix=%q class=%s range=%s..%s", data[:4], attendance.classID, attendance.from, attendance.to)
	}
}

func TestExportMarksReadsEntitiesAndReturnsPDF(t *testing.T) {
	classID, termID, subjectID := uuid.New(), uuid.New(), uuid.New()
	store := &reportStore{class: "10-A", term: "Term 1", subject: "Mathematics", marks: []markRow{
		{StudentName: "Marked", IndexNumber: "S001", TermMarkID: pgtype.UUID{Bytes: uuid.New(), Valid: true}, Marks: numeric(t, "75"), MaxMarks: numeric(t, "100")},
		{StudentName: "Absent", IndexNumber: "S002", TermMarkID: pgtype.UUID{Bytes: uuid.New(), Valid: true}, MaxMarks: numeric(t, "100"), IsAbsent: pgtype.Bool{Bool: true, Valid: true}},
		{StudentName: "Unmarked", IndexNumber: "S003"},
	}}
	data, err := NewService(store, &attendanceReader{}).ExportMarks(context.Background(), MarksReportRequest{ClassID: classID.String(), TermID: termID.String(), SubjectID: subjectID.String()})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(data, []byte("%PDF")) || store.marksCalls != 1 || store.classID != classID || store.termID != termID || store.subjectID != subjectID {
		t.Fatalf("invalid marks export: calls=%d class=%s term=%s subject=%s", store.marksCalls, store.classID, store.termID, store.subjectID)
	}
}
