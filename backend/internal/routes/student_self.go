package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/openschool-org/openschool/internal/handlers"
	academicsmodule "github.com/openschool-org/openschool/internal/modules/academics"
	"github.com/openschool-org/openschool/internal/repositories"
	notificationsrepositories "github.com/openschool-org/openschool/internal/repositories/notifications"
	timetablerepositories "github.com/openschool-org/openschool/internal/repositories/timetable"
	"github.com/openschool-org/openschool/internal/services"
	notificationsservices "github.com/openschool-org/openschool/internal/services/notifications"
)

func RegisterStudentSelfRoutes(student *gin.RouterGroup, pool *pgxpool.Pool) {
	studentsRepo := repositories.NewStudentRepository(pool)
	classRepo := repositories.NewClassRepository(pool)
	teacherRepo := repositories.NewTeacherRepository(pool)
	guardianRepo := repositories.NewGuardianRepository(pool)
	notifications := notificationsservices.NewNotificationService(
		notificationsrepositories.NewNotificationRepository(pool),
		classRepo,
		timetablerepositories.NewGradeSectionRepository(pool),
		repositories.NewSectionHeadRepository(pool),
		teacherRepo,
		studentsRepo,
		guardianRepo,
		repositories.NewSchoolRepository(pool),
		repositories.NewPositionRepository(pool),
	)
	auditSvc := services.NewAuditService(repositories.NewAuditRepository(pool))
	positionSvc := services.NewPositionService(repositories.NewPositionRepository(pool), repositories.NewSectionHeadRepository(pool), nil)
	attendanceService := services.NewAttendanceService(repositories.NewAttendanceRepository(pool), repositories.NewUserRepository(pool), teacherRepo, classRepo, studentsRepo, guardianRepo, notifications, auditSvc, positionSvc, repositories.NewSchoolRepository(pool))
	enrollmentService := academicsmodule.NewStudentEnrollment(academicsmodule.NewEnrollmentRepository(pool))

	handler := handlers.NewStudentSelfHandler(services.NewStudentSelfService(studentsRepo), attendanceService, newTermMarkRunner(pool), enrollmentService)

	student.GET("/me/student", handler.Profile)
	student.GET("/me/student/attendance", handler.Attendance)
	student.GET("/me/student/marks", handler.Marks)
	student.GET("/me/student/enrollments", handler.ListEnrollments)
	student.POST("/me/student/enrollments", handler.SubmitEnrollment)
	student.POST("/me/student/enrollments/confirm", handler.ConfirmEnrollment)
}
