// Package app owns the OpenSchool API composition root.
package app

import (
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/openschool-org/openschool/internal/jobs"
	"github.com/openschool-org/openschool/internal/middleware"
	"github.com/openschool-org/openschool/internal/models"
	academicsmodule "github.com/openschool-org/openschool/internal/modules/academics"
	attendancemodule "github.com/openschool-org/openschool/internal/modules/attendance"
	auditmodule "github.com/openschool-org/openschool/internal/modules/audit"
	curriculummodule "github.com/openschool-org/openschool/internal/modules/curriculum"
	identitymodule "github.com/openschool-org/openschool/internal/modules/identity"
	leadershipmodule "github.com/openschool-org/openschool/internal/modules/leadership"
	notificationmodule "github.com/openschool-org/openschool/internal/modules/notifications"
	schoolmodule "github.com/openschool-org/openschool/internal/modules/school"
	studentleadershipmodule "github.com/openschool-org/openschool/internal/modules/studentleadership"
	timetablemodule "github.com/openschool-org/openschool/internal/modules/timetable"
	"github.com/openschool-org/openschool/internal/repositories"
	"github.com/openschool-org/openschool/internal/routes"
	"github.com/openschool-org/openschool/internal/services"
)

// HTTPGroups contains the authorization-scoped route groups shared by modules.
type HTTPGroups struct {
	API            *gin.RouterGroup
	Protected      *gin.RouterGroup
	Admin          *gin.RouterGroup
	TeacherOrAdmin *gin.RouterGroup
	Parent         *gin.RouterGroup
	Student        *gin.RouterGroup
	Teacher        *gin.RouterGroup
	StudentAccess  *gin.RouterGroup
}

// Setup composes the API modules and returns the scheduler owned by the process lifecycle.
func Setup(router *gin.Engine, pool *pgxpool.Pool) *jobs.Scheduler {
	groups := newHTTPGroups(router, pool)
	notifications := notificationmodule.NewNotificationService(notificationmodule.NewNotificationRepository(pool))

	groups.API.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	routes.RegisterSetupRoutes(groups.API, pool)
	routes.RegisterAuthRoutes(groups.API, groups.Protected, pool)
	identitymodule.Register(groups.Protected, pool)
	routes.RegisterCoreModule(groups.Admin, groups.Protected, pool)
	curriculummodule.RegisterMediumRoutes(groups.Admin, groups.Protected, pool)
	curriculummodule.RegisterLevelRoutes(groups.Admin, groups.Protected, pool)
	auditService := auditmodule.NewService(auditmodule.NewRepository(pool))
	auditmodule.RegisterRoutes(groups.Admin, auditService)
	leadershipService := leadershipmodule.NewService(leadershipmodule.NewRepository(pool), auditService)
	leadershipmodule.RegisterRoutes(groups.Admin, groups.TeacherOrAdmin, leadershipService)
	studentLeadershipService := studentleadershipmodule.NewService(studentleadershipmodule.NewRepository(pool))
	studentleadershipmodule.RegisterRoutes(groups.Admin, groups.TeacherOrAdmin, groups.StudentAccess, studentLeadershipService)
	houseService := schoolmodule.NewHouseService(pool, auditService)
	schoolmodule.RegisterHouseRoutes(groups.Admin, groups.TeacherOrAdmin, houseService)
	schoolmodule.RegisterSchoolRoutes(groups.Admin, groups.TeacherOrAdmin, groups.Protected, pool)
	schoolmodule.RegisterGradeRoutes(groups.Admin, groups.TeacherOrAdmin, pool)
	schoolmodule.RegisterTermRoutes(groups.Admin, groups.Protected, pool)
	timetablemodule.RegisterSettingsRoutes(groups.Admin, pool)
	timetablemodule.RegisterClassroomRoutes(groups.Admin, groups.TeacherOrAdmin, pool)
	timetablemodule.RegisterSubjectPeriodRequirementRoutes(groups.Admin, groups.TeacherOrAdmin, pool)
	timetablemodule.RegisterTeacherAvailabilityRoutes(groups.Admin, groups.TeacherOrAdmin, pool)
	timetablemodule.RegisterGradeSectionRoutes(groups.Admin, groups.TeacherOrAdmin, pool)
	timetablemodule.RegisterTimetableEntryRoutes(groups.Admin, groups.TeacherOrAdmin, pool)
	timetablemodule.RegisterTimetableValidationRoute(groups.TeacherOrAdmin, pool)
	timetablemodule.RegisterTimetableCRUDRoutes(groups.Admin, groups.TeacherOrAdmin, pool)
	timetablemodule.RegisterTimetableStatusHistoryRoute(groups.TeacherOrAdmin, pool)
	timetableRepository := timetablemodule.NewWorkflowRepository(pool)
	timetablemodule.RegisterTimetableWorkflowRoutes(groups.Admin, groups.TeacherOrAdmin, groups.Teacher, timetableRepository, timetablemodule.NewWorkflowValidator(pool), notifications)
	timetablemodule.RegisterTimetablePortalRoutes(groups.Teacher, groups.Student, groups.TeacherOrAdmin, timetableRepository)
	timetablemodule.RegisterTimetableGenerationRoute(groups.Admin, timetableRepository)
	academicsmodule.RegisterSubjectRoutes(groups.Admin, groups.TeacherOrAdmin, pool)
	academicsmodule.RegisterStreamRoutes(groups.Admin, groups.TeacherOrAdmin, pool)
	academicsmodule.RegisterClassRoutes(groups.Admin, groups.TeacherOrAdmin, pool)
	academicsmodule.RegisterEnrollmentRoutes(groups.Admin, groups.TeacherOrAdmin, groups.StudentAccess, groups.Protected, academicsmodule.NewEnrollmentRepository(pool))
	routes.RegisterAcademicModule(groups.Admin, groups.TeacherOrAdmin, groups.StudentAccess, groups.Protected, pool)
	routes.RegisterPeopleModule(groups.Admin, groups.TeacherOrAdmin, groups.StudentAccess, houseService, pool)
	routes.RegisterSelfServiceModule(groups.Student, pool)
	attendanceService := attendancemodule.NewService(attendancemodule.NewRepository(pool), notifications, auditService, attendanceLeadership{positions: leadershipService})
	attendancemodule.RegisterRoutes(groups.TeacherOrAdmin, attendanceService)
	staffAttendanceService := attendancemodule.NewStaffService(attendancemodule.NewRepository(pool))
	attendancemodule.RegisterStaffRoutes(groups.Admin, groups.Teacher, staffAttendanceService, services.NewTeacherSelfService(repositories.NewTeacherRepository(pool)))
	routes.RegisterAdminOperationsModule(groups.Admin, groups.TeacherOrAdmin, groups.StudentAccess, pool)

	timetableReader := timetablemodule.NewReader(pool)
	routes.RegisterParentAndTeacherSelfModule(groups.Parent, groups.Teacher, timetableReader, leadershipService, studentLeadershipService, pool)
	notificationmodule.RegisterRoutes(groups.TeacherOrAdmin, groups.Protected, notifications)

	return routes.RegisterAutomationModule(groups.Admin, pool)
}

func newHTTPGroups(router *gin.Engine, pool *pgxpool.Pool) HTTPGroups {
	api := router.Group("/api/v1")
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())
	protected.Use(middleware.PerAccountRateLimit(
		envFloatOr("API_PER_ACCOUNT_RATE_LIMIT_RPS", 30),
		envIntOr("API_PER_ACCOUNT_RATE_LIMIT_BURST", 60),
	))

	admin := protected.Group("")
	admin.Use(middleware.RequireRole(models.RoleAdmin))
	teacherOrAdmin := protected.Group("")
	teacherOrAdmin.Use(middleware.RequireRole(models.RoleAdmin, models.RoleTeacher))
	parent := protected.Group("")
	parent.Use(middleware.RequireRole(models.RoleParent))
	student := protected.Group("")
	student.Use(middleware.RequireRole(models.RoleStudent))
	teacher := protected.Group("")
	teacher.Use(middleware.RequireRole(models.RoleTeacher))
	studentAccess := protected.Group("")
	studentAccess.Use(middleware.RequireStudentAccess(pool))

	return HTTPGroups{
		API: api, Protected: protected, Admin: admin, TeacherOrAdmin: teacherOrAdmin,
		Parent: parent, Student: student, Teacher: teacher, StudentAccess: studentAccess,
	}
}

func envFloatOr(key string, fallback float64) float64 {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseFloat(value, 64); err == nil {
			return parsed
		}
	}
	return fallback
}

func envIntOr(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return fallback
}
