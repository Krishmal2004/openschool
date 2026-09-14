package jobs

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/openschool-org/openschool/internal/modules/notifications"
	"github.com/openschool-org/openschool/internal/repositories"
)

// BuildAll constructs and wires up all five registered agents — the single place that lists every agent that exists.
func BuildAll(pool *pgxpool.Pool) []Job {
	checks := repositories.NewJobChecksRepository(pool)

	notifSvc := notifications.NewNotificationService(notifications.NewNotificationRepository(pool))

	return []Job{
		NewSystemHealthAgent(pool, checks, notifSvc),
		NewStructuralIntegrityAgent(checks, notifSvc),
		NewPeopleComplianceAgent(checks, notifSvc),
		NewAcademicDeliveryAgent(checks, notifSvc),
		NewSecurityAuditAgent(checks, notifSvc),
	}
}
