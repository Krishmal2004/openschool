package selfservice

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/openschool-org/openschool/internal/ports"
)

// RegisterRoutes mounts the authenticated student, parent, and teacher portal
// endpoints while preserving their existing paths and authorization groups.
func RegisterRoutes(
	student, parent, teacher *gin.RouterGroup,
	pool *pgxpool.Pool,
	students *StudentProfiles,
	teachers *TeacherProfiles,
	guardians ports.GuardianAccess,
	timetables TimetableReader,
	positions LeadershipReader,
	societies SocietyReader,
	dashboard DashboardReader,
) {
	registerStudentRoutes(student, students, pool)
	registerParentRoutes(parent, guardians, timetables, pool)
	registerTeacherRoutes(teacher, teachers, timetables, positions, societies, dashboard)
}
