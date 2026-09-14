// This file defines the RegisterGuardianRoutes function, mapping guardian management and student-guardian linking endpoints to their handlers and middlewares.

package routes

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/openschool-org/openschool/internal/models"
	peoplemodule "github.com/openschool-org/openschool/internal/modules/people"
	"github.com/openschool-org/openschool/internal/repositories"
	notificationsrepositories "github.com/openschool-org/openschool/internal/repositories/notifications"
	timetablerepositories "github.com/openschool-org/openschool/internal/repositories/timetable"
	"github.com/openschool-org/openschool/internal/services"
	notificationsservices "github.com/openschool-org/openschool/internal/services/notifications"
)

type guardianCompatibilityWriter struct{ service *services.GuardianService }

func (w guardianCompatibilityWriter) Create(c context.Context, r models.CreateGuardianRequest) (any, any, error) {
	return w.service.CreateGuardian(c, r)
}
func (w guardianCompatibilityWriter) Update(c context.Context, id uuid.UUID, r models.UpdateGuardianRequest) (any, error) {
	return w.service.UpdateGuardian(c, id, r)
}
func (w guardianCompatibilityWriter) Delete(c context.Context, id, actor uuid.UUID) error {
	return w.service.DeleteGuardian(c, id, actor)
}
func (w guardianCompatibilityWriter) Link(c context.Context, id uuid.UUID, r models.LinkGuardianRequest) error {
	return w.service.LinkToStudent(c, id, r)
}
func (w guardianCompatibilityWriter) Unlink(c context.Context, student, guardian uuid.UUID) error {
	return w.service.UnlinkFromStudent(c, student, guardian)
}
func (w guardianCompatibilityWriter) SetPrimary(c context.Context, student, guardian uuid.UUID) error {
	return w.service.SetPrimaryContact(c, student, guardian)
}
func (w guardianCompatibilityWriter) Provision(c context.Context, id uuid.UUID, r models.ProvisionGuardianLoginRequest, actor uuid.UUID) (any, error) {
	return w.service.ProvisionLogin(c, id, r, actor)
}

func RegisterGuardianRoutes(admin *gin.RouterGroup, teacherOrAdmin *gin.RouterGroup, studentAccess *gin.RouterGroup, pool *pgxpool.Pool) {
	repo := repositories.NewGuardianRepository(pool)
	notifications := notificationsservices.NewNotificationService(
		notificationsrepositories.NewNotificationRepository(pool),
		repositories.NewClassRepository(pool),
		timetablerepositories.NewGradeSectionRepository(pool),
		repositories.NewSectionHeadRepository(pool),
		repositories.NewTeacherRepository(pool),
		repositories.NewStudentRepository(pool),
		repo,
		repositories.NewSchoolRepository(pool),
		repositories.NewPositionRepository(pool),
	)
	service := services.NewGuardianService(repo, repositories.NewUserRepository(pool), newIdentityProvider(), notifications, services.NewAuditService(repositories.NewAuditRepository(pool)))
	peoplemodule.RegisterGuardianReadRoutes(teacherOrAdmin, studentAccess, peoplemodule.NewGuardianReader(pool))
	peoplemodule.RegisterGuardianWriteRoutes(admin, guardianCompatibilityWriter{service: service})
	peoplemodule.RegisterGuardianNotificationRoute(admin, peoplemodule.NewGuardianNotificationReader(pool))

}
