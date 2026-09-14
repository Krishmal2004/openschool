package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/openschool-org/openschool/internal/handlers"
	timetablemodule "github.com/openschool-org/openschool/internal/modules/timetable"
	"github.com/openschool-org/openschool/internal/repositories"
	"github.com/openschool-org/openschool/internal/services"
)

// RegisterTeacherSelfRoutes wires the signed-in teacher's self-service
// endpoints. It receives the timetable module's read contract so this
// aggregate does not depend on the legacy timetable service.
func RegisterTeacherSelfRoutes(teacher *gin.RouterGroup, pool *pgxpool.Pool, timetableService *timetablemodule.Reader) {
	teacherRepo := repositories.NewTeacherRepository(pool)
	schoolRepo := repositories.NewSchoolRepository(pool)
	positionService := services.NewPositionService(repositories.NewPositionRepository(pool), repositories.NewSectionHeadRepository(pool), nil)
	societyService := services.NewSocietyService(repositories.NewSocietyRepository(pool), teacherRepo)
	dashboardService := services.NewDashboardService(repositories.NewDashboardRepository(pool))
	staffAttendanceService := services.NewStaffAttendanceService(repositories.NewStaffAttendanceRepository(pool))
	handler := handlers.NewTeacherSelfHandler(services.NewTeacherSelfService(teacherRepo), schoolRepo, positionService, societyService, dashboardService, timetableService, staffAttendanceService)

	teacher.GET("/me/teacher", handler.Profile)
	teacher.GET("/me/teacher/position", handler.Position)
	teacher.GET("/me/teacher/leadership-overview", handler.LeadershipOverview)
	teacher.GET("/me/teacher/society", handler.Society)
	teacher.GET("/me/teacher/analytics", handler.Analytics)
	teacher.GET("/me/teacher/timetables", handler.Timetables)
	teacher.GET("/me/teacher/attendance", handler.Attendance)
}
