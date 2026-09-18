// Package dashboard owns school-wide analytics aggregation and HTTP delivery.
package dashboard

import (
	"context"
	"sort"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type countByGrade struct {
	GradeName    string
	StudentCount int64
}
type countByClass struct {
	GradeName, ClassName string
	StudentCount         int64
}
type genderCount struct {
	Gender       string
	StudentCount int64
}
type houseCount struct {
	HouseName, HouseColor string
	StudentCount          int64
}
type attendanceTrend struct {
	Date                     pgtype.Date
	PresentCount, TotalCount int64
}
type staffCounts struct{ AcademicStaffCount, NonAcademicStaffCount int64 }
type attendanceTotals struct{ PresentCount, LateCount, AbsentCount, LeaveCount int64 }
type subjectPerformance struct {
	SubjectName       string
	AveragePercentage pgtype.Numeric
	Entries           int64
}
type examinationSummary struct {
	AveragePercentage          pgtype.Numeric
	Entries, StudentsWithMarks int64
}
type gradePerformance struct {
	GradeName         string
	AveragePercentage pgtype.Numeric
}
type classPerformance struct {
	GradeName, ClassName string
	AveragePercentage    pgtype.Numeric
}
type attendancePercentage struct{ PresentCount, TotalCount int64 }
type studentGrowth struct {
	AcademicYearLabel string
	StudentCount      int64
}
type staffGrowth struct {
	Year         int32
	TeacherCount int64
}
type timetableCompletion struct{ TotalClasses, PublishedClasses int64 }
type recentStudent struct {
	ID                                          uuid.UUID
	FullName, IndexNumber, GradeName, ClassName string
	CreatedAt                                   pgtype.Timestamptz
}
type recentTeacher struct {
	ID                       uuid.UUID
	FullName, EmployeeNumber string
	CreatedAt                pgtype.Timestamptz
}

type store interface {
	studentCountByGrade(context.Context) ([]countByGrade, error)
	studentCountByClass(context.Context) ([]countByClass, error)
	studentGenderDistribution(context.Context) ([]genderCount, error)
	studentHouseDistribution(context.Context) ([]houseCount, error)
	studentAttendanceTrend(context.Context) ([]attendanceTrend, error)
	staffCounts(context.Context) (staffCounts, error)
	staffAttendanceThisMonth(context.Context) (attendanceTotals, error)
	subjectPerformance(context.Context) ([]subjectPerformance, error)
	examinationSummary(context.Context) (examinationSummary, error)
	gradeWisePerformance(context.Context) ([]gradePerformance, error)
	classWisePerformance(context.Context) ([]classPerformance, error)
	attendancePercentage(context.Context) (attendancePercentage, error)
	studentGrowth(context.Context) ([]studentGrowth, error)
	staffGrowth(context.Context) ([]staffGrowth, error)
	notificationsSentCount(context.Context) (int64, error)
	timetableCompletion(context.Context) (timetableCompletion, error)
	recentStudents(context.Context, int32) ([]recentStudent, error)
	recentTeachers(context.Context, int32) ([]recentTeacher, error)
}

// recentActivityLimit bounds how many of each of the newest students/teachers
// feed the dashboard's combined "recent activity" list.
const recentActivityLimit = 6

type Service struct{ store store }

func NewService(store store) *Service { return &Service{store: store} }

func (s *Service) Analytics(ctx context.Context) (DashboardAnalyticsResponse, error) {
	response := emptyResponse()
	byGrade, err := s.store.studentCountByGrade(ctx)
	if err != nil {
		return response, err
	}
	for _, value := range byGrade {
		response.Student.ByGrade = append(response.Student.ByGrade, CountRow{Label: value.GradeName, Count: value.StudentCount})
		response.Student.Total += value.StudentCount
	}
	byClass, err := s.store.studentCountByClass(ctx)
	if err != nil {
		return response, err
	}
	for _, value := range byClass {
		response.Student.ByClass = append(response.Student.ByClass, CountRow{Label: value.GradeName + " " + value.ClassName, Count: value.StudentCount})
	}
	gender, err := s.store.studentGenderDistribution(ctx)
	if err != nil {
		return response, err
	}
	for _, value := range gender {
		response.Student.GenderDistribution = append(response.Student.GenderDistribution, CountRow{Label: value.Gender, Count: value.StudentCount})
	}
	houses, err := s.store.studentHouseDistribution(ctx)
	if err != nil {
		return response, err
	}
	for _, value := range houses {
		response.Student.HouseDistribution = append(response.Student.HouseDistribution, HouseCountRow{Name: value.HouseName, Color: value.HouseColor, Count: value.StudentCount})
	}
	trend, err := s.store.studentAttendanceTrend(ctx)
	if err != nil {
		return response, err
	}
	for _, value := range trend {
		response.Student.AttendanceTrend = append(response.Student.AttendanceTrend, AttendanceTrendPoint{Date: value.Date.Time.Format("2006-01-02"), PresentCount: value.PresentCount, TotalCount: value.TotalCount})
	}
	staff, err := s.store.staffCounts(ctx)
	if err != nil {
		return response, err
	}
	response.Staff.AcademicStaffCount, response.Staff.NonAcademicStaffCount = staff.AcademicStaffCount, staff.NonAcademicStaffCount
	staffAttendance, err := s.store.staffAttendanceThisMonth(ctx)
	if err != nil {
		return response, err
	}
	response.Staff.AttendanceThisMonth = StaffAttendanceTotals(staffAttendance)
	subjects, err := s.store.subjectPerformance(ctx)
	if err != nil {
		return response, err
	}
	for _, value := range subjects {
		response.Academic.SubjectPerformance = append(response.Academic.SubjectPerformance, PerformanceRow{Label: value.SubjectName, AveragePercentage: numericToFloat64(value.AveragePercentage), Entries: value.Entries})
	}
	exams, err := s.store.examinationSummary(ctx)
	if err != nil {
		return response, err
	}
	response.Academic.ExaminationAverage, response.Academic.ExaminationEntries, response.Academic.StudentsWithMarks = numericToFloat64(exams.AveragePercentage), exams.Entries, exams.StudentsWithMarks
	grades, err := s.store.gradeWisePerformance(ctx)
	if err != nil {
		return response, err
	}
	for _, value := range grades {
		response.Academic.GradeWisePerformance = append(response.Academic.GradeWisePerformance, PerformanceRow{Label: value.GradeName, AveragePercentage: numericToFloat64(value.AveragePercentage)})
	}
	classes, err := s.store.classWisePerformance(ctx)
	if err != nil {
		return response, err
	}
	for _, value := range classes {
		response.Academic.ClassWisePerformance = append(response.Academic.ClassWisePerformance, PerformanceRow{Label: value.GradeName + " " + value.ClassName, AveragePercentage: numericToFloat64(value.AveragePercentage)})
	}
	attendance, err := s.store.attendancePercentage(ctx)
	if err != nil {
		return response, err
	}
	if attendance.TotalCount > 0 {
		response.Academic.AttendancePercentage = float64(attendance.PresentCount) / float64(attendance.TotalCount) * 100
	}
	students, err := s.store.studentGrowth(ctx)
	if err != nil {
		return response, err
	}
	for _, value := range students {
		response.School.StudentGrowth = append(response.School.StudentGrowth, GrowthPoint{Label: value.AcademicYearLabel, Count: value.StudentCount})
	}
	staffHistory, err := s.store.staffGrowth(ctx)
	if err != nil {
		return response, err
	}
	for _, value := range staffHistory {
		response.School.StaffGrowth = append(response.School.StaffGrowth, GrowthPoint{Label: strconv.Itoa(int(value.Year)), Count: value.TeacherCount})
	}
	response.School.NotificationsSentCount, err = s.store.notificationsSentCount(ctx)
	if err != nil {
		return response, err
	}
	timetable, err := s.store.timetableCompletion(ctx)
	if err != nil {
		return response, err
	}
	response.School.TotalClasses, response.School.PublishedClasses = timetable.TotalClasses, timetable.PublishedClasses
	if timetable.TotalClasses > 0 {
		response.School.TimetableCompletionPct = float64(timetable.PublishedClasses) / float64(timetable.TotalClasses) * 100
	}

	recentActivity, err := s.recentActivity(ctx)
	if err != nil {
		return response, err
	}
	response.School.RecentActivity = recentActivity

	return response, nil
}

// recentActivity merges the newest students and teachers into one
// server-sorted feed — a dedicated query per entity (recentStudents/
// recentTeachers), not a slice of the capped, alphabetically-sorted list
// pages, which could otherwise miss anyone enrolled/added after the first
// page once more than one page of records exists.
func (s *Service) recentActivity(ctx context.Context) ([]RecentActivityItem, error) {
	students, err := s.store.recentStudents(ctx, recentActivityLimit)
	if err != nil {
		return nil, err
	}
	teachers, err := s.store.recentTeachers(ctx, recentActivityLimit)
	if err != nil {
		return nil, err
	}

	items := make([]RecentActivityItem, 0, len(students)+len(teachers))
	for _, st := range students {
		sub := st.GradeName
		if st.ClassName != "" {
			sub += " " + st.ClassName
		}
		if st.IndexNumber != "" {
			if sub != "" {
				sub += " · "
			}
			sub += st.IndexNumber
		}
		items = append(items, RecentActivityItem{
			Key: "student-" + st.ID.String(), Text: st.FullName + " enrolled", Sub: sub,
			Time: formatTimestamptz(st.CreatedAt), Path: "/students/" + st.ID.String(), Kind: "student",
		})
	}
	for _, t := range teachers {
		items = append(items, RecentActivityItem{
			Key: "teacher-" + t.ID.String(), Text: t.FullName + " added as a teacher", Sub: t.EmployeeNumber,
			Time: formatTimestamptz(t.CreatedAt), Path: "/teachers/" + t.ID.String(), Kind: "teacher",
		})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Time > items[j].Time })
	if len(items) > recentActivityLimit {
		items = items[:recentActivityLimit]
	}
	return items, nil
}

func formatTimestamptz(value pgtype.Timestamptz) string {
	if !value.Valid {
		return ""
	}
	return value.Time.Format(time.RFC3339)
}

func emptyResponse() DashboardAnalyticsResponse {
	return DashboardAnalyticsResponse{
		Student:  DashboardStudentAnalytics{ByGrade: []CountRow{}, ByClass: []CountRow{}, GenderDistribution: []CountRow{}, HouseDistribution: []HouseCountRow{}, AttendanceTrend: []AttendanceTrendPoint{}},
		Academic: DashboardAcademicAnalytics{SubjectPerformance: []PerformanceRow{}, GradeWisePerformance: []PerformanceRow{}, ClassWisePerformance: []PerformanceRow{}},
		School:   DashboardSchoolAnalytics{StudentGrowth: []GrowthPoint{}, StaffGrowth: []GrowthPoint{}, RecentActivity: []RecentActivityItem{}},
	}
}

func numericToFloat64(value pgtype.Numeric) float64 {
	result, err := value.Float64Value()
	if err != nil || !result.Valid {
		return 0
	}
	return result.Float64
}
