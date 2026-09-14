// Package audit owns append-only audit recording and administrative audit-log reads.
package audit

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/openschool-org/openschool/internal/models"
)

type createCommand struct {
	EntityType string
	EntityID   uuid.UUID
	Action     string
	ActorID    pgtype.UUID
	Before     []byte
	After      []byte
	Reason     pgtype.Text
}

type row struct {
	ID         uuid.UUID
	EntityType string
	EntityID   uuid.UUID
	Action     string
	ActorID    pgtype.UUID
	Before     []byte
	After      []byte
	Reason     pgtype.Text
	CreatedAt  pgtype.Timestamptz
	ActorName  pgtype.Text
}

type store interface {
	create(context.Context, createCommand) error
	list(context.Context, string, *uuid.UUID) ([]row, error)
}

type Service struct{ store store }

func NewService(store store) *Service { return &Service{store: store} }

func toJSON(value interface{}) []byte {
	if value == nil {
		return nil
	}
	data, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	return data
}

// Record appends one best-effort audit event through the shared AuditRecorder port.
func (s *Service) Record(ctx context.Context, entityType string, entityID uuid.UUID, action string, actorID uuid.UUID, before, after interface{}, reason string) error {
	return s.store.create(ctx, createCommand{
		EntityType: entityType,
		EntityID:   entityID,
		Action:     action,
		ActorID:    pgtype.UUID{Bytes: actorID, Valid: actorID != uuid.Nil},
		Before:     toJSON(before),
		After:      toJSON(after),
		Reason:     pgtype.Text{String: reason, Valid: reason != ""},
	})
}

func (s *Service) List(ctx context.Context, entityType string, entityID *uuid.UUID) ([]models.AuditLogResponse, error) {
	rows, err := s.store.list(ctx, entityType, entityID)
	if err != nil {
		return nil, err
	}
	result := make([]models.AuditLogResponse, len(rows))
	for i, value := range rows {
		result[i] = models.AuditLogResponse{ID: value.ID, EntityType: value.EntityType, EntityID: value.EntityID, Action: value.Action, Before: value.Before, After: value.After, CreatedAt: value.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00")}
		if value.ActorID.Valid {
			id := uuid.UUID(value.ActorID.Bytes)
			result[i].ActorID = &id
		}
		if value.ActorName.Valid {
			name := value.ActorName.String
			result[i].ActorName = &name
		}
		if value.Reason.Valid {
			reason := value.Reason.String
			result[i].Reason = &reason
		}
	}
	return result, nil
}
