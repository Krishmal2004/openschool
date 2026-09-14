package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	peoplemodule "github.com/openschool-org/openschool/internal/modules/people"
	"github.com/openschool-org/openschool/internal/ports"
)

func RegisterTeacherRoutes(admin *gin.RouterGroup, teacherOrAdmin *gin.RouterGroup, houseAssignments ports.HouseAssignments, pool *pgxpool.Pool) {
	auditSvc := newAuditRecorder(pool)
	service := peoplemodule.NewTeacherService(peoplemodule.NewStudentStore(pool), newIdentityProvider(), houseAssignments, auditSvc)
	peoplemodule.RegisterTeacherReadRoutes(teacherOrAdmin, admin, peoplemodule.NewTeacherReader(pool))
	peoplemodule.RegisterTeacherWriteRoutes(admin, service)
}
