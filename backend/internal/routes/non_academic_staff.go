package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	peoplemodule "github.com/openschool-org/openschool/internal/modules/people"
	"github.com/openschool-org/openschool/internal/repositories"
	"github.com/openschool-org/openschool/internal/services"
)

func RegisterNonAcademicStaffRoutes(admin *gin.RouterGroup, teacherOrAdmin *gin.RouterGroup, pool *pgxpool.Pool) {
	auditSvc := services.NewAuditService(repositories.NewAuditRepository(pool))
	peoplemodule.RegisterNonAcademicStaffRoutes(admin, teacherOrAdmin, peoplemodule.NewNonAcademicStaffService(peoplemodule.NewNonAcademicStaffStore(pool), auditSvc))
}
