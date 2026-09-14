// This file defines the RegisterGuardianRoutes function, mapping guardian management and student-guardian linking endpoints to their handlers and middlewares.

package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	peoplemodule "github.com/openschool-org/openschool/internal/modules/people"
)

func RegisterGuardianRoutes(admin *gin.RouterGroup, teacherOrAdmin *gin.RouterGroup, studentAccess *gin.RouterGroup, pool *pgxpool.Pool) {
	store := peoplemodule.NewGuardianStore(pool)
	service := peoplemodule.NewGuardianService(store, newIdentityProvider(), newAuditRecorder(pool))
	peoplemodule.RegisterGuardianReadRoutes(teacherOrAdmin, studentAccess, peoplemodule.NewGuardianReader(pool))
	peoplemodule.RegisterGuardianWriteRoutes(admin, service)
	peoplemodule.RegisterGuardianNotificationRoute(admin, peoplemodule.NewGuardianNotificationReader(pool))

}
