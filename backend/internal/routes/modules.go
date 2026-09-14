package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/openschool-org/openschool/internal/jobs"
	timetablemodule "github.com/openschool-org/openschool/internal/modules/timetable"
	"github.com/openschool-org/openschool/internal/ports"
	"github.com/openschool-org/openschool/internal/repositories"
)

// The functions in this file are the composition-root modules. Each module
// owns route registration for one business area; endpoint implementations
// remain in the existing route files.
func RegisterCoreModule(admin, protected *gin.RouterGroup, pool *pgxpool.Pool) {
	RegisterCurriculumRoutes(admin, protected, pool)
}

func RegisterAcademicModule(admin, teacherOrAdmin, studentAccess, protected *gin.RouterGroup, pool *pgxpool.Pool) {
	RegisterPromotionRoutes(admin, pool)
	RegisterTermMarkRoutes(teacherOrAdmin, pool)
}

func RegisterPeopleModule(admin, teacherOrAdmin, studentAccess *gin.RouterGroup, houseAssignments ports.HouseAssignments, pool *pgxpool.Pool) {
	RegisterStudentRoutes(admin, teacherOrAdmin, houseAssignments, pool)
	RegisterTeacherRoutes(admin, teacherOrAdmin, houseAssignments, pool)
	RegisterGuardianRoutes(admin, teacherOrAdmin, studentAccess, pool)
	RegisterNonAcademicStaffRoutes(admin, teacherOrAdmin, pool)
	RegisterStudentPortfolioRoutes(teacherOrAdmin, studentAccess, pool)
	RegisterSearchRoutes(admin, pool)
}

func RegisterSelfServiceModule(student *gin.RouterGroup, pool *pgxpool.Pool) {
	RegisterStudentSelfRoutes(student, pool)
}

func RegisterParentAndTeacherSelfModule(parent, teacher *gin.RouterGroup, timetableService *timetablemodule.Reader, pool *pgxpool.Pool) {
	RegisterParentRoutes(parent, timetableService, pool)
	RegisterTeacherSelfRoutes(teacher, pool, timetableService)
}

func RegisterAutomationModule(admin *gin.RouterGroup, pool *pgxpool.Pool) *jobs.Scheduler {
	scheduler := jobs.NewScheduler(jobs.BuildAll(pool), repositories.NewJobSchedulerRepository(pool))
	RegisterJobRoutes(admin, pool, scheduler)
	return scheduler
}

func RegisterAdminOperationsModule(admin, teacherOrAdmin, studentAccess *gin.RouterGroup, pool *pgxpool.Pool) {
	RegisterPrefectRoutes(admin, teacherOrAdmin, studentAccess, pool)
	RegisterSocietyRoutes(admin, teacherOrAdmin, studentAccess, pool)
	RegisterPositionRoutes(admin, teacherOrAdmin, pool)
	RegisterIdentityReconciliationRoutes(admin, pool)
	RegisterDashboardRoutes(admin, pool)
	RegisterReportExportRoutes(admin, pool)
}
