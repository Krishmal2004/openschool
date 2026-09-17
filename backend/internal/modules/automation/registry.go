package automation

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/openschool-org/openschool/internal/modules/notifications"
)

// BuildAll constructs and wires up all six registered agents — the single place that lists every agent that exists.
func BuildAll(pool *pgxpool.Pool) []Job {
	checks := NewRepository(pool)

	notifSvc := notifications.NewNotificationService(notifications.NewNotificationRepository(pool))

	return []Job{
		NewSystemHealthAgent(pool, checks, notifSvc),
		NewStructuralIntegrityAgent(checks, notifSvc),
		NewPeopleComplianceAgent(checks, notifSvc),
		NewAcademicDeliveryAgent(checks, notifSvc),
		NewSecurityAuditAgent(checks, notifSvc),
		NewDataRetentionAgent(checks, notifSvc),
	}
}
