package audit

import (
	"encoding/json"

	"github.com/google/uuid"
)

// AuditLogResponse is the API representation of an audit entry.
type AuditLogResponse struct {
	ID         uuid.UUID       `json:"id"`
	EntityType string          `json:"entity_type"`
	EntityID   uuid.UUID       `json:"entity_id"`
	Action     string          `json:"action"`
	ActorID    *uuid.UUID      `json:"actor_id"`
	ActorName  *string         `json:"actor_name"`
	Before     json.RawMessage `json:"before"`
	After      json.RawMessage `json:"after"`
	Reason     *string         `json:"reason"`
	CreatedAt  string          `json:"created_at"`
}
