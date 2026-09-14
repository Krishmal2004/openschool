package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	peoplemodule "github.com/openschool-org/openschool/internal/modules/people"
)

func RegisterNonAcademicStaffRoutes(admin *gin.RouterGroup, teacherOrAdmin *gin.RouterGroup, pool *pgxpool.Pool) {
	auditSvc := newAuditRecorder(pool)
	peoplemodule.RegisterNonAcademicStaffRoutes(admin, teacherOrAdmin, peoplemodule.NewNonAcademicStaffService(peoplemodule.NewNonAcademicStaffStore(pool), auditSvc))
}
