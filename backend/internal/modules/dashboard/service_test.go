package dashboard

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type dashboardStore struct {
	grades          []countByGrade
	classes         []countByClass
	genders         []genderCount
	houses          []houseCount
	trend           []attendanceTrend
	staff           staffCounts
	staffAttendance attendanceTotals
	subjects        []subjectPerformance
	exams           examinationSummary
	gradePerf       []gradePerformance
	classPerf       []classPerformance
	attendance      attendancePercentage
	studentHistory  []studentGrowth
	staffHistory    []staffGrowth
	notifications   int64
	timetable       timetableCompletion
}

func (s *dashboardStore) studentCountByGrade(context.Context) ([]countByGrade, error) {
	return s.grades, nil
}
func (s *dashboardStore) studentCountByClass(context.Context) ([]countByClass, error) {
	return s.classes, nil
}
func (s *dashboardStore) studentGenderDistribution(context.Context) ([]genderCount, error) {
	return s.genders, nil
}
func (s *dashboardStore) studentHouseDistribution(context.Context) ([]houseCount, error) {
	return s.houses, nil
}
func (s *dashboardStore) studentAttendanceTrend(context.Context) ([]attendanceTrend, error) {
	return s.trend, nil
}
func (s *dashboardStore) staffCounts(context.Context) (staffCounts, error) { return s.staff, nil }
func (s *dashboardStore) staffAttendanceThisMonth(context.Context) (attendanceTotals, error) {
	return s.staffAttendance, nil
}
func (s *dashboardStore) subjectPerformance(context.Context) ([]subjectPerformance, error) {
	return s.subjects, nil
}
func (s *dashboardStore) examinationSummary(context.Context) (examinationSummary, error) {
	return s.exams, nil
}
func (s *dashboardStore) gradeWisePerformance(context.Context) ([]gradePerformance, error) {
	return s.gradePerf, nil
}
func (s *dashboardStore) classWisePerformance(context.Context) ([]classPerformance, error) {
	return s.classPerf, nil
}
func (s *dashboardStore) attendancePercentage(context.Context) (attendancePercentage, error) {
	return s.attendance, nil
}
func (s *dashboardStore) studentGrowth(context.Context) ([]studentGrowth, error) {
	return s.studentHistory, nil
}
func (s *dashboardStore) staffGrowth(context.Context) ([]staffGrowth, error) {
	return s.staffHistory, nil
}
func (s *dashboardStore) notificationsSentCount(context.Context) (int64, error) {
	return s.notifications, nil
}
func (s *dashboardStore) timetableCompletion(context.Context) (timetableCompletion, error) {
	return s.timetable, nil
}

func dashboardNumeric(t *testing.T, value string) pgtype.Numeric {
	t.Helper()
	var numeric pgtype.Numeric
	if err := numeric.Scan(value); err != nil {
		t.Fatal(err)
	}
	return numeric
}

func TestAnalyticsComposesAllDashboardSections(t *testing.T) {
	store := &dashboardStore{
		grades:          []countByGrade{{GradeName: "Grade 10", StudentCount: 20}, {GradeName: "Grade 11", StudentCount: 15}},
		classes:         []countByClass{{GradeName: "Grade 10", ClassName: "A", StudentCount: 20}},
		genders:         []genderCount{{Gender: "female", StudentCount: 18}},
		houses:          []houseCount{{HouseName: "Blue", HouseColor: "#00f", StudentCount: 9}},
		trend:           []attendanceTrend{{Date: pgtype.Date{Time: time.Date(2026, 9, 14, 0, 0, 0, 0, time.UTC), Valid: true}, PresentCount: 30, TotalCount: 35}},
		staff:           staffCounts{AcademicStaffCount: 12, NonAcademicStaffCount: 4},
		staffAttendance: attendanceTotals{PresentCount: 10, LateCount: 1, AbsentCount: 1, LeaveCount: 2},
		subjects:        []subjectPerformance{{SubjectName: "Math", AveragePercentage: dashboardNumeric(t, "75.5"), Entries: 30}},
		exams:           examinationSummary{AveragePercentage: dashboardNumeric(t, "72.25"), Entries: 60, StudentsWithMarks: 32},
		gradePerf:       []gradePerformance{{GradeName: "Grade 10", AveragePercentage: dashboardNumeric(t, "70")}},
		classPerf:       []classPerformance{{GradeName: "Grade 10", ClassName: "A", AveragePercentage: dashboardNumeric(t, "71")}},
		attendance:      attendancePercentage{PresentCount: 30, TotalCount: 40},
		studentHistory:  []studentGrowth{{AcademicYearLabel: "2026", StudentCount: 35}},
		staffHistory:    []staffGrowth{{Year: 2026, TeacherCount: 12}},
		notifications:   8,
		timetable:       timetableCompletion{TotalClasses: 8, PublishedClasses: 6},
	}
	response, err := NewService(store).Analytics(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if response.Student.Total != 35 || response.Student.ByClass[0].Label != "Grade 10 A" || response.Student.AttendanceTrend[0].Date != "2026-09-14" {
		t.Fatalf("unexpected student analytics: %#v", response.Student)
	}
	if response.Academic.ExaminationAverage != 72.25 || response.Academic.AttendancePercentage != 75 {
		t.Fatalf("unexpected academic analytics: %#v", response.Academic)
	}
	if response.School.TimetableCompletionPct != 75 || response.School.NotificationsSentCount != 8 || response.School.StaffGrowth[0].Label != "2026" {
		t.Fatalf("unexpected school analytics: %#v", response.School)
	}
}

func TestAnalyticsUsesEmptyArraysInsteadOfNull(t *testing.T) {
	response, err := NewService(&dashboardStore{}).Analytics(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(response)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), ":null") {
		t.Fatalf("dashboard response contains null collection: %s", data)
	}
}
