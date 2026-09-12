package jobs

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/openschool-org/openschool/internal/repositories"
	notificationrepositories "github.com/openschool-org/openschool/internal/repositories/notifications"
	timetablerepositories "github.com/openschool-org/openschool/internal/repositories/timetable"
	"github.com/openschool-org/openschool/internal/services/notifications"
)

// BuildAll constructs and wires up all five registered agents — the single place that lists every agent that exists.
func BuildAll(pool *pgxpool.Pool) []Job {
	checks := repositories.NewJobChecksRepository(pool)

	// Same construction notifications.RegisterNotificationRoutes uses, so agents reuse the normal notification pipeline.
	notifSvc := notifications.NewNotificationService(
		notificationrepositories.NewNotificationRepository(pool),
		repositories.NewClassRepository(pool),
		timetablerepositories.NewGradeSectionRepository(pool),
		repositories.NewSectionHeadRepository(pool),
		repositories.NewTeacherRepository(pool),
		repositories.NewStudentRepository(pool),
		repositories.NewGuardianRepository(pool),
		repositories.NewSchoolRepository(pool),
		repositories.NewPositionRepository(pool),
	)

	return []Job{
		NewSystemHealthAgent(pool, checks, notifSvc),
		NewStructuralIntegrityAgent(checks, notifSvc),
		NewPeopleComplianceAgent(checks, notifSvc),
		NewAcademicDeliveryAgent(checks, notifSvc),
		NewSecurityAuditAgent(checks, notifSvc),
	}
}
