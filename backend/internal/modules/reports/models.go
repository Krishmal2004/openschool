package reports

import "time"

// AttendanceReportRequest drives the attendance PDF export.
type AttendanceReportRequest struct {
	ClassID string    `form:"class_id" binding:"required"`
	From    time.Time `form:"from" binding:"required" time_format:"2006-01-02"`
	To      time.Time `form:"to" binding:"required" time_format:"2006-01-02"`
	Columns []string  `form:"columns"`
}

// MarksReportRequest drives the class, term, and subject marks PDF export.
type MarksReportRequest struct {
	ClassID   string   `form:"class_id" binding:"required"`
	TermID    string   `form:"term_id" binding:"required"`
	SubjectID string   `form:"subject_id" binding:"required"`
	Columns   []string `form:"columns"`
}
