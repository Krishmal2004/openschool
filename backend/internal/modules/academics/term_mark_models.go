package academics

type MarkEntry struct {
	StudentID string  `json:"student_id" binding:"required"`
	Marks     float64 `json:"marks"`
	MaxMarks  float64 `json:"max_marks"`
	// IsAbsent records "absent for the test" instead of a numeric mark.
	// Marks is ignored (stored as 0) when true.
	IsAbsent bool `json:"is_absent"`
}

type BulkUpsertMarksRequest struct {
	TermID    string      `json:"term_id" binding:"required"`
	SubjectID string      `json:"subject_id" binding:"required"`
	Entries   []MarkEntry `json:"entries" binding:"required,min=1,dive"`
}
