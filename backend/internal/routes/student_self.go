package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/openschool-org/openschool/internal/handlers"
	academicsmodule "github.com/openschool-org/openschool/internal/modules/academics"
	attendancemodule "github.com/openschool-org/openschool/internal/modules/attendance"
	"github.com/openschool-org/openschool/internal/repositories"
	"github.com/openschool-org/openschool/internal/services"
)

func RegisterStudentSelfRoutes(student *gin.RouterGroup, pool *pgxpool.Pool) {
	studentsRepo := repositories.NewStudentRepository(pool)
	attendanceReader := attendancemodule.NewReader(attendancemodule.NewRepository(pool))
	enrollmentService := academicsmodule.NewStudentEnrollment(academicsmodule.NewEnrollmentRepository(pool))

	handler := handlers.NewStudentSelfHandler(services.NewStudentSelfService(studentsRepo), attendanceReader, newTermMarkRunner(pool), enrollmentService)

	student.GET("/me/student", handler.Profile)
	student.GET("/me/student/attendance", handler.Attendance)
	student.GET("/me/student/marks", handler.Marks)
	student.GET("/me/student/enrollments", handler.ListEnrollments)
	student.POST("/me/student/enrollments", handler.SubmitEnrollment)
	student.POST("/me/student/enrollments/confirm", handler.ConfirmEnrollment)
}
