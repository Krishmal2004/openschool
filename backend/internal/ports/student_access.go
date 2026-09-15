package ports

import (
	"context"

	"github.com/google/uuid"
)

// StudentAccessAuthorizer exposes the identity-to-student relationships needed
// by student-resource authorization middleware without exposing persistence.
type StudentAccessAuthorizer interface {
	StudentIDForUser(context.Context, uuid.UUID) (uuid.UUID, error)
	IsGuardianOfStudent(context.Context, uuid.UUID, uuid.UUID) (bool, error)
}
