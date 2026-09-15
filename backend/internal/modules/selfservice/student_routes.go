package selfservice

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	academicsmodule "github.com/openschool-org/openschool/internal/modules/academics"
	attendancemodule "github.com/openschool-org/openschool/internal/modules/attendance"
)

func registerStudentRoutes(student *gin.RouterGroup, profiles *StudentProfiles, pool *pgxpool.Pool) {
	attendanceReader := attendancemodule.NewReader(attendancemodule.NewRepository(pool))
	enrollmentService := academicsmodule.NewStudentEnrollment(academicsmodule.NewEnrollmentRepository(pool))
	marks := academicsmodule.NewTermMarkService(academicsmodule.NewTermMarkRepository(pool))

	handler := NewStudentSelfHandler(profiles, attendanceReader, marks, enrollmentService)

	student.GET("/me/student", handler.Profile)
	student.GET("/me/student/attendance", handler.Attendance)
	student.GET("/me/student/marks", handler.Marks)
	student.GET("/me/student/enrollments", handler.ListEnrollments)
	student.POST("/me/student/enrollments", handler.SubmitEnrollment)
	student.POST("/me/student/enrollments/confirm", handler.ConfirmEnrollment)
}
