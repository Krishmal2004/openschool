// Package ports contains small application-facing interfaces shared across
// modules. Implementations stay in infrastructure packages.
package ports

import (
	"context"

	"github.com/google/uuid"
)

// CurrentAcademicYearReader exposes only the school-context data other
// modules need. It intentionally does not leak sqlc-generated types.
type CurrentAcademicYearReader interface {
	CurrentAcademicYearID(ctx context.Context) (uuid.UUID, error)
}
