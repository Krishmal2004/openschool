// Package reports owns server-rendered administrative report exports.
package reports

import (
	"bytes"
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jung-kurt/gofpdf"
	attendancemodule "github.com/openschool-org/openschool/internal/modules/attendance"
)

type markRow struct {
	StudentName string
	IndexNumber string
	TermMarkID  pgtype.UUID
	Marks       pgtype.Numeric
	MaxMarks    pgtype.Numeric
	IsAbsent    pgtype.Bool
}

type store interface {
	className(context.Context, uuid.UUID) (string, error)
	termName(context.Context, uuid.UUID) (string, error)
	subjectName(context.Context, uuid.UUID) (string, error)
	classMarks(context.Context, uuid.UUID, uuid.UUID, uuid.UUID) ([]markRow, error)
}

type Service struct {
	store      store
	attendance attendancemodule.ReportReader
}

func NewService(store store, attendance attendancemodule.ReportReader) *Service {
	return &Service{store: store, attendance: attendance}
}

var attendanceReportColumns = []string{"date", "student", "index_number", "status", "note"}
var attendanceReportHeaders = map[string]string{
	"date": "Date", "student": "Student", "index_number": "Index No.", "status": "Status", "note": "Note",
}

var marksReportColumns = []string{"student", "index_number", "marks", "max_marks", "percentage"}
var marksReportHeaders = map[string]string{
	"student": "Student", "index_number": "Index No.", "marks": "Marks", "max_marks": "Max Marks", "percentage": "%",
}

func resolveColumns(requested, all []string) []string {
	if len(requested) == 0 {
		return all
	}
	valid := make(map[string]bool, len(all))
	for _, column := range all {
		valid[column] = true
	}
	columns := make([]string, 0, len(requested))
	for _, column := range requested {
		if valid[column] {
			columns = append(columns, column)
		}
	}
	if len(columns) == 0 {
		return all
	}
	return columns
}

func (s *Service) ExportAttendance(ctx context.Context, req AttendanceReportRequest) ([]byte, error) {
	classID, err := uuid.Parse(req.ClassID)
	if err != nil {
		return nil, fmt.Errorf("invalid class id")
	}
	className, err := s.store.className(ctx, classID)
	if err != nil {
		return nil, fmt.Errorf("class not found")
	}
	records, err := s.attendance.ListForClassInRange(ctx, classID, req.From, req.To)
	if err != nil {
		return nil, err
	}

	columns := resolveColumns(req.Columns, attendanceReportColumns)
	rows := make([][]string, 0, len(records))
	for _, record := range records {
		values := map[string]string{
			"date": record.SessionDate.Time.Format("2006-01-02"), "student": record.StudentName,
			"index_number": record.StudentIndex, "status": record.Status, "note": record.Note.String,
		}
		row := make([]string, len(columns))
		for i, column := range columns {
			row[i] = values[column]
		}
		rows = append(rows, row)
	}

	headers := reportHeaders(columns, attendanceReportHeaders)
	title := fmt.Sprintf("Attendance Report — %s (%s to %s)", className, req.From.Format("2006-01-02"), req.To.Format("2006-01-02"))
	return renderTablePDF(title, headers, rows)
}

func (s *Service) ExportMarks(ctx context.Context, req MarksReportRequest) ([]byte, error) {
	classID, err := uuid.Parse(req.ClassID)
	if err != nil {
		return nil, fmt.Errorf("invalid class id")
	}
	termID, err := uuid.Parse(req.TermID)
	if err != nil {
		return nil, fmt.Errorf("invalid term id")
	}
	subjectID, err := uuid.Parse(req.SubjectID)
	if err != nil {
		return nil, fmt.Errorf("invalid subject id")
	}
	className, err := s.store.className(ctx, classID)
	if err != nil {
		return nil, fmt.Errorf("class not found")
	}
	termName, err := s.store.termName(ctx, termID)
	if err != nil {
		return nil, fmt.Errorf("term not found")
	}
	subjectName, err := s.store.subjectName(ctx, subjectID)
	if err != nil {
		return nil, fmt.Errorf("subject not found")
	}
	marks, err := s.store.classMarks(ctx, classID, termID, subjectID)
	if err != nil {
		return nil, err
	}

	columns := resolveColumns(req.Columns, marksReportColumns)
	rows := make([][]string, 0, len(marks))
	for _, mark := range marks {
		if !mark.TermMarkID.Valid {
			continue
		}
		marksDisplay := formatNumeric(mark.Marks)
		percentageDisplay := "0.0"
		if mark.IsAbsent.Valid && mark.IsAbsent.Bool {
			marksDisplay, percentageDisplay = "AB", "AB"
		} else if mark.MaxMarks.Valid && mark.Marks.Valid {
			maxMarks, marks := numericToFloat64(mark.MaxMarks), numericToFloat64(mark.Marks)
			percentage := 0.0
			if maxMarks > 0 {
				percentage = marks / maxMarks * 100
			}
			percentageDisplay = fmt.Sprintf("%.1f", percentage)
		}
		values := map[string]string{
			"student": mark.StudentName, "index_number": mark.IndexNumber, "marks": marksDisplay,
			"max_marks": formatNumeric(mark.MaxMarks), "percentage": percentageDisplay,
		}
		row := make([]string, len(columns))
		for i, column := range columns {
			row[i] = values[column]
		}
		rows = append(rows, row)
	}

	headers := reportHeaders(columns, marksReportHeaders)
	return renderTablePDF(fmt.Sprintf("Marks Report — %s | %s | %s", className, termName, subjectName), headers, rows)
}

func reportHeaders(columns []string, labels map[string]string) []string {
	headers := make([]string, len(columns))
	for i, column := range columns {
		headers[i] = labels[column]
	}
	return headers
}

func renderTablePDF(title string, headers []string, rows [][]string) ([]byte, error) {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Helvetica", "B", 14)
	pdf.CellFormat(0, 10, title, "", 1, "L", false, 0, "")
	pdf.Ln(2)
	pageWidth, _ := pdf.GetPageSize()
	marginLeft, _, marginRight, _ := pdf.GetMargins()
	usableWidth := pageWidth - marginLeft - marginRight
	columnWidth := usableWidth / float64(len(headers))
	pdf.SetFont("Helvetica", "B", 10)
	pdf.SetFillColor(230, 230, 230)
	for _, header := range headers {
		pdf.CellFormat(columnWidth, 8, header, "1", 0, "L", true, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetFont("Helvetica", "", 9)
	for _, row := range rows {
		for _, cell := range row {
			pdf.CellFormat(columnWidth, 7, cell, "1", 0, "L", false, 0, "")
		}
		pdf.Ln(-1)
	}
	if len(rows) == 0 {
		pdf.SetFont("Helvetica", "I", 9)
		pdf.CellFormat(usableWidth, 8, "No data for the selected range.", "1", 1, "C", false, 0, "")
	}
	var buffer bytes.Buffer
	if err := pdf.Output(&buffer); err != nil {
		return nil, fmt.Errorf("failed to render PDF: %w", err)
	}
	return buffer.Bytes(), nil
}

func formatNumeric(value pgtype.Numeric) string {
	if !value.Valid {
		return "-"
	}
	return strconv.FormatFloat(numericToFloat64(value), 'f', -1, 64)
}

func numericToFloat64(value pgtype.Numeric) float64 {
	converted, err := value.Float64Value()
	if err != nil || !converted.Valid {
		return 0
	}
	return converted.Float64
}
