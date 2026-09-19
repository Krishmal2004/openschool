package dashboard

type CountRow struct {
	Label string `json:"label"`
	Count int64  `json:"count"`
}

type HouseCountRow struct {
	Name  string `json:"name"`
	Color string `json:"color"`
	Count int64  `json:"count"`
}

type AttendanceTrendPoint struct {
	Date         string `json:"date"`
	PresentCount int64  `json:"present_count"`
	TotalCount   int64  `json:"total_count"`
}

type PerformanceRow struct {
	Label             string  `json:"label"`
	AveragePercentage float64 `json:"average_percentage"`
	Entries           int64   `json:"entries"`
}

type GrowthPoint struct {
	Label string `json:"label"`
	Count int64  `json:"count"`
}

type StaffAttendanceTotals struct {
	PresentCount int64 `json:"present_count"`
	LateCount    int64 `json:"late_count"`
	AbsentCount  int64 `json:"absent_count"`
	LeaveCount   int64 `json:"leave_count"`
}

type DashboardStudentAnalytics struct {
	Total              int64                  `json:"total"`
	ByGrade            []CountRow             `json:"by_grade"`
	ByClass            []CountRow             `json:"by_class"`
	GenderDistribution []CountRow             `json:"gender_distribution"`
	HouseDistribution  []HouseCountRow        `json:"house_distribution"`
	AttendanceTrend    []AttendanceTrendPoint `json:"attendance_trend"`
}

type DashboardStaffAnalytics struct {
	AcademicStaffCount    int64                 `json:"academic_staff_count"`
	NonAcademicStaffCount int64                 `json:"non_academic_staff_count"`
	AttendanceThisMonth   StaffAttendanceTotals `json:"attendance_this_month"`
}

type DashboardAcademicAnalytics struct {
	SubjectPerformance   []PerformanceRow `json:"subject_performance"`
	ExaminationAverage   float64          `json:"examination_average"`
	ExaminationEntries   int64            `json:"examination_entries"`
	StudentsWithMarks    int64            `json:"students_with_marks"`
	GradeWisePerformance []PerformanceRow `json:"grade_wise_performance"`
	ClassWisePerformance []PerformanceRow `json:"class_wise_performance"`
	AttendancePercentage float64          `json:"attendance_percentage"`
}

// RecentActivityItem mirrors the frontend's RecentActivitySection item shape
// directly, so Dashboard.tsx can render this list without reshaping it.
type RecentActivityItem struct {
	Key  string `json:"key"`
	Text string `json:"text"`
	Sub  string `json:"sub"`
	Time string `json:"time"`
	Path string `json:"path"`
	Kind string `json:"kind"`
}

type DashboardSchoolAnalytics struct {
	StudentGrowth          []GrowthPoint        `json:"student_growth"`
	StaffGrowth            []GrowthPoint        `json:"staff_growth"`
	NotificationsSentCount int64                `json:"notifications_sent_count"`
	TimetableCompletionPct float64              `json:"timetable_completion_pct"`
	TotalClasses           int64                `json:"total_classes"`
	PublishedClasses       int64                `json:"published_classes"`
	RecentActivity         []RecentActivityItem `json:"recent_activity"`
}

type DashboardAnalyticsResponse struct {
	Student  DashboardStudentAnalytics  `json:"student"`
	Staff    DashboardStaffAnalytics    `json:"staff"`
	Academic DashboardAcademicAnalytics `json:"academic"`
	School   DashboardSchoolAnalytics   `json:"school"`
}
