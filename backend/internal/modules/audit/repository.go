package audit

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/openschool-org/openschool/db/sqlc"
)

type Repository struct{ queries *db.Queries }

func NewRepository(pool *pgxpool.Pool) *Repository { return &Repository{queries: db.New(pool)} }

func (r *Repository) create(ctx context.Context, command createCommand) error {
	_, err := r.queries.CreateAuditLog(ctx, db.CreateAuditLogParams{EntityType: command.EntityType, EntityID: command.EntityID, Action: command.Action, ActorID: command.ActorID, Before: command.Before, After: command.After, Reason: command.Reason})
	return err
}

func (r *Repository) list(ctx context.Context, entityType string, entityID *uuid.UUID, limit, offset int32) ([]row, int64, error) {
	params := db.ListAuditLogsParams{EntityType: pgtype.Text{String: entityType, Valid: entityType != ""}, PageLimit: limit, PageOffset: offset}
	if entityID != nil {
		params.EntityID = pgtype.UUID{Bytes: *entityID, Valid: true}
	}
	rows, err := r.queries.ListAuditLogs(ctx, params)
	if err != nil {
		return nil, 0, err
	}
	result := make([]row, len(rows))
	var total int64
	for i, value := range rows {
		result[i] = row{ID: value.ID, EntityType: value.EntityType, EntityID: value.EntityID, Action: value.Action, ActorID: value.ActorID, Before: value.Before, After: value.After, Reason: value.Reason, CreatedAt: value.CreatedAt, ActorName: value.ActorName}
		total = value.Total
	}
	return result, total, nil
}
