package app

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	attendancemodule "github.com/openschool-org/openschool/internal/modules/attendance"
	auditmodule "github.com/openschool-org/openschool/internal/modules/audit"
	leadershipmodule "github.com/openschool-org/openschool/internal/modules/leadership"
	notificationmodule "github.com/openschool-org/openschool/internal/modules/notifications"
	reportsmodule "github.com/openschool-org/openschool/internal/modules/reports"
	selfservicemodule "github.com/openschool-org/openschool/internal/modules/selfservice"
)

// attendanceLeadership adapts Leadership's rank error to Attendance's
// domain-specific authorization error.
type attendanceLeadership struct{ positions *leadershipmodule.Service }

func (a attendanceLeadership) LeadershipScope(ctx context.Context, teacherID, academicYearID uuid.UUID) (bool, []uuid.UUID, error) {
	wholeSchool, gradeIDs, err := a.positions.LeadershipScope(ctx, teacherID, academicYearID)
	if errors.Is(err, leadershipmodule.ErrInsufficientRank) {
		return false, nil, attendancemodule.ErrInsufficientRank
	}
	return wholeSchool, gradeIDs, err
}

func registerAttendanceAndReports(
	groups HTTPGroups,
	pool *pgxpool.Pool,
	notifications *notificationmodule.NotificationService,
	audit *auditmodule.Service,
	leadership *leadershipmodule.Service,
	teacherProfiles *selfservicemodule.TeacherProfiles,
) {
	attendanceService := attendancemodule.NewService(
		attendancemodule.NewRepository(pool), notifications, audit,
		attendanceLeadership{positions: leadership},
	)
	attendancemodule.RegisterRoutes(groups.TeacherOrAdmin, attendanceService)
	reportsmodule.RegisterRoutes(groups.Admin, reportsmodule.NewService(reportsmodule.NewRepository(pool), attendanceService))
	staffService := attendancemodule.NewStaffService(attendancemodule.NewRepository(pool))
	attendancemodule.RegisterStaffRoutes(groups.Admin, groups.Teacher, staffService, teacherProfiles)
}
