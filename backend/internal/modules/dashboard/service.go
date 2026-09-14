// Package dashboard owns school-wide analytics aggregation and HTTP delivery.
package dashboard

import (
	"context"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/openschool-org/openschool/internal/models"
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
}

type Service struct{ store store }

func NewService(store store) *Service { return &Service{store: store} }

func (s *Service) Analytics(ctx context.Context) (models.DashboardAnalyticsResponse, error) {
	response := emptyResponse()
	byGrade, err := s.store.studentCountByGrade(ctx)
	if err != nil {
		return response, err
	}
	for _, value := range byGrade {
		response.Student.ByGrade = append(response.Student.ByGrade, models.CountRow{Label: value.GradeName, Count: value.StudentCount})
		response.Student.Total += value.StudentCount
	}
	byClass, err := s.store.studentCountByClass(ctx)
	if err != nil {
		return response, err
	}
	for _, value := range byClass {
		response.Student.ByClass = append(response.Student.ByClass, models.CountRow{Label: value.GradeName + " " + value.ClassName, Count: value.StudentCount})
	}
	gender, err := s.store.studentGenderDistribution(ctx)
	if err != nil {
		return response, err
	}
	for _, value := range gender {
		response.Student.GenderDistribution = append(response.Student.GenderDistribution, models.CountRow{Label: value.Gender, Count: value.StudentCount})
	}
	houses, err := s.store.studentHouseDistribution(ctx)
	if err != nil {
		return response, err
	}
	for _, value := range houses {
		response.Student.HouseDistribution = append(response.Student.HouseDistribution, models.HouseCountRow{Name: value.HouseName, Color: value.HouseColor, Count: value.StudentCount})
	}
	trend, err := s.store.studentAttendanceTrend(ctx)
	if err != nil {
		return response, err
	}
	for _, value := range trend {
		response.Student.AttendanceTrend = append(response.Student.AttendanceTrend, models.AttendanceTrendPoint{Date: value.Date.Time.Format("2006-01-02"), PresentCount: value.PresentCount, TotalCount: value.TotalCount})
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
	response.Staff.AttendanceThisMonth = models.StaffAttendanceTotals{PresentCount: staffAttendance.PresentCount, LateCount: staffAttendance.LateCount, AbsentCount: staffAttendance.AbsentCount, LeaveCount: staffAttendance.LeaveCount}
	subjects, err := s.store.subjectPerformance(ctx)
	if err != nil {
		return response, err
	}
	for _, value := range subjects {
		response.Academic.SubjectPerformance = append(response.Academic.SubjectPerformance, models.PerformanceRow{Label: value.SubjectName, AveragePercentage: numericToFloat64(value.AveragePercentage), Entries: value.Entries})
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
		response.Academic.GradeWisePerformance = append(response.Academic.GradeWisePerformance, models.PerformanceRow{Label: value.GradeName, AveragePercentage: numericToFloat64(value.AveragePercentage)})
	}
	classes, err := s.store.classWisePerformance(ctx)
	if err != nil {
		return response, err
	}
	for _, value := range classes {
		response.Academic.ClassWisePerformance = append(response.Academic.ClassWisePerformance, models.PerformanceRow{Label: value.GradeName + " " + value.ClassName, AveragePercentage: numericToFloat64(value.AveragePercentage)})
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
		response.School.StudentGrowth = append(response.School.StudentGrowth, models.GrowthPoint{Label: value.AcademicYearLabel, Count: value.StudentCount})
	}
	staffHistory, err := s.store.staffGrowth(ctx)
	if err != nil {
		return response, err
	}
	for _, value := range staffHistory {
		response.School.StaffGrowth = append(response.School.StaffGrowth, models.GrowthPoint{Label: strconv.Itoa(int(value.Year)), Count: value.TeacherCount})
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
	return response, nil
}

func emptyResponse() models.DashboardAnalyticsResponse {
	return models.DashboardAnalyticsResponse{
		Student:  models.DashboardStudentAnalytics{ByGrade: []models.CountRow{}, ByClass: []models.CountRow{}, GenderDistribution: []models.CountRow{}, HouseDistribution: []models.HouseCountRow{}, AttendanceTrend: []models.AttendanceTrendPoint{}},
		Academic: models.DashboardAcademicAnalytics{SubjectPerformance: []models.PerformanceRow{}, GradeWisePerformance: []models.PerformanceRow{}, ClassWisePerformance: []models.PerformanceRow{}},
		School:   models.DashboardSchoolAnalytics{StudentGrowth: []models.GrowthPoint{}, StaffGrowth: []models.GrowthPoint{}},
	}
}

func numericToFloat64(value pgtype.Numeric) float64 {
	result, err := value.Float64Value()
	if err != nil || !result.Valid {
		return 0
	}
	return result.Float64
}
