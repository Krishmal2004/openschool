package routes

import (
	"github.com/jackc/pgx/v5/pgxpool"
	auditmodule "github.com/openschool-org/openschool/internal/modules/audit"
)

func newAuditRecorder(pool *pgxpool.Pool) *auditmodule.Service {
	return auditmodule.NewService(auditmodule.NewRepository(pool))
}
