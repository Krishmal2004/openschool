package dashboard

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/openschool-org/openschool/db/sqlc"
)

type Repository struct{ queries *db.Queries }

func NewRepository(pool *pgxpool.Pool) *Repository { return &Repository{queries: db.New(pool)} }

func (r *Repository) studentCountByGrade(ctx context.Context) ([]countByGrade, error) {
	rows, err := r.queries.DashboardStudentCountByGrade(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]countByGrade, len(rows))
	for i, value := range rows {
		result[i] = countByGrade{GradeName: value.GradeName, StudentCount: value.StudentCount}
	}
	return result, nil
}

func (r *Repository) studentCountByClass(ctx context.Context) ([]countByClass, error) {
	rows, err := r.queries.DashboardStudentCountByClass(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]countByClass, len(rows))
	for i, value := range rows {
		result[i] = countByClass{GradeName: value.GradeName, ClassName: value.ClassName, StudentCount: value.StudentCount}
	}
	return result, nil
}

func (r *Repository) studentGenderDistribution(ctx context.Context) ([]genderCount, error) {
	rows, err := r.queries.DashboardStudentGenderDistribution(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]genderCount, len(rows))
	for i, value := range rows {
		result[i] = genderCount{Gender: value.Gender, StudentCount: value.StudentCount}
	}
	return result, nil
}

func (r *Repository) studentHouseDistribution(ctx context.Context) ([]houseCount, error) {
	rows, err := r.queries.DashboardStudentHouseDistribution(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]houseCount, len(rows))
	for i, value := range rows {
		result[i] = houseCount{HouseName: value.HouseName, HouseColor: value.HouseColor, StudentCount: value.StudentCount}
	}
	return result, nil
}

func (r *Repository) studentAttendanceTrend(ctx context.Context) ([]attendanceTrend, error) {
	rows, err := r.queries.DashboardStudentAttendanceTrend(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]attendanceTrend, len(rows))
	for i, value := range rows {
		result[i] = attendanceTrend{Date: value.Date, PresentCount: value.PresentCount, TotalCount: value.TotalCount}
	}
	return result, nil
}

func (r *Repository) staffCounts(ctx context.Context) (staffCounts, error) {
	value, err := r.queries.DashboardStaffCounts(ctx)
	return staffCounts{AcademicStaffCount: value.AcademicStaffCount, NonAcademicStaffCount: value.NonAcademicStaffCount}, err
}

func (r *Repository) staffAttendanceThisMonth(ctx context.Context) (attendanceTotals, error) {
	value, err := r.queries.DashboardStaffAttendanceThisMonth(ctx)
	return attendanceTotals{PresentCount: value.PresentCount, LateCount: value.LateCount, AbsentCount: value.AbsentCount, LeaveCount: value.LeaveCount}, err
}

func (r *Repository) subjectPerformance(ctx context.Context) ([]subjectPerformance, error) {
	rows, err := r.queries.DashboardSubjectPerformance(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]subjectPerformance, len(rows))
	for i, value := range rows {
		result[i] = subjectPerformance{SubjectName: value.SubjectName, AveragePercentage: value.AveragePercentage, Entries: value.Entries}
	}
	return result, nil
}

func (r *Repository) examinationSummary(ctx context.Context) (examinationSummary, error) {
	value, err := r.queries.DashboardExaminationSummary(ctx)
	return examinationSummary{AveragePercentage: value.AveragePercentage, Entries: value.Entries, StudentsWithMarks: value.StudentsWithMarks}, err
}

func (r *Repository) gradeWisePerformance(ctx context.Context) ([]gradePerformance, error) {
	rows, err := r.queries.DashboardGradeWisePerformance(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]gradePerformance, len(rows))
	for i, value := range rows {
		result[i] = gradePerformance{GradeName: value.GradeName, AveragePercentage: value.AveragePercentage}
	}
	return result, nil
}

func (r *Repository) classWisePerformance(ctx context.Context) ([]classPerformance, error) {
	rows, err := r.queries.DashboardClassWisePerformance(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]classPerformance, len(rows))
	for i, value := range rows {
		result[i] = classPerformance{GradeName: value.GradeName, ClassName: value.ClassName, AveragePercentage: value.AveragePercentage}
	}
	return result, nil
}

func (r *Repository) attendancePercentage(ctx context.Context) (attendancePercentage, error) {
	value, err := r.queries.DashboardAttendancePercentage(ctx)
	return attendancePercentage{PresentCount: value.PresentCount, TotalCount: value.TotalCount}, err
}

func (r *Repository) studentGrowth(ctx context.Context) ([]studentGrowth, error) {
	rows, err := r.queries.DashboardStudentGrowth(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]studentGrowth, len(rows))
	for i, value := range rows {
		result[i] = studentGrowth{AcademicYearLabel: value.AcademicYearLabel, StudentCount: value.StudentCount}
	}
	return result, nil
}

func (r *Repository) staffGrowth(ctx context.Context) ([]staffGrowth, error) {
	rows, err := r.queries.DashboardStaffGrowth(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]staffGrowth, len(rows))
	for i, value := range rows {
		result[i] = staffGrowth{Year: value.Year, TeacherCount: value.TeacherCount}
	}
	return result, nil
}

func (r *Repository) notificationsSentCount(ctx context.Context) (int64, error) {
	return r.queries.DashboardNotificationsSentCount(ctx)
}

func (r *Repository) timetableCompletion(ctx context.Context) (timetableCompletion, error) {
	value, err := r.queries.DashboardTimetableCompletion(ctx)
	return timetableCompletion{TotalClasses: value.TotalClasses, PublishedClasses: value.PublishedClasses}, err
}
