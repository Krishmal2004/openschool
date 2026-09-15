package app

import (
	"github.com/jackc/pgx/v5/pgxpool"
	notificationmodule "github.com/openschool-org/openschool/internal/modules/notifications"
	timetablemodule "github.com/openschool-org/openschool/internal/modules/timetable"
)

func registerTimetable(groups HTTPGroups, pool *pgxpool.Pool, notifications *notificationmodule.NotificationService) {
	timetablemodule.RegisterSettingsRoutes(groups.Admin, pool)
	timetablemodule.RegisterClassroomRoutes(groups.Admin, groups.TeacherOrAdmin, pool)
	timetablemodule.RegisterSubjectPeriodRequirementRoutes(groups.Admin, groups.TeacherOrAdmin, pool)
	timetablemodule.RegisterTeacherAvailabilityRoutes(groups.Admin, groups.TeacherOrAdmin, pool)
	timetablemodule.RegisterGradeSectionRoutes(groups.Admin, groups.TeacherOrAdmin, pool)
	timetablemodule.RegisterTimetableEntryRoutes(groups.Admin, groups.TeacherOrAdmin, pool)
	timetablemodule.RegisterTimetableValidationRoute(groups.TeacherOrAdmin, pool)
	timetablemodule.RegisterTimetableCRUDRoutes(groups.Admin, groups.TeacherOrAdmin, pool)
	timetablemodule.RegisterTimetableStatusHistoryRoute(groups.TeacherOrAdmin, pool)

	repository := timetablemodule.NewWorkflowRepository(pool)
	timetablemodule.RegisterTimetableWorkflowRoutes(groups.Admin, groups.TeacherOrAdmin, groups.Teacher, repository, timetablemodule.NewWorkflowValidator(pool), notifications)
	timetablemodule.RegisterTimetablePortalRoutes(groups.Teacher, groups.Student, groups.TeacherOrAdmin, repository)
	timetablemodule.RegisterTimetableGenerationRoute(groups.Admin, repository)
}
