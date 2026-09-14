package routes

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	peoplemodule "github.com/openschool-org/openschool/internal/modules/people"
	"github.com/openschool-org/openschool/internal/ports"
	"github.com/openschool-org/openschool/internal/repositories"
	"github.com/openschool-org/openschool/internal/services"
)

func RegisterStudentRoutes(admin *gin.RouterGroup, teacherOrAdmin *gin.RouterGroup, houseAssignments ports.HouseAssignments, pool *pgxpool.Pool) {
	auditSvc := services.NewAuditService(repositories.NewAuditRepository(pool))
	studentStore := peoplemodule.NewStudentStore(pool)
	schoolRepo := repositories.NewSchoolRepository(pool)
	service := peoplemodule.NewStudentService(studentStore, newIdentityProvider(), houseAssignments, auditSvc, func(ctx context.Context) (string, error) {
		school, err := schoolRepo.Get(ctx)
		return school.SchoolType, err
	})
	peoplemodule.RegisterStudentRoutes(admin, teacherOrAdmin, service, studentStore, service)
}
