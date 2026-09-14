package app

import (
	"context"
	"errors"

	"github.com/google/uuid"
	attendancemodule "github.com/openschool-org/openschool/internal/modules/attendance"
	"github.com/openschool-org/openschool/internal/services"
)

// attendanceLeadership adapts the legacy position capability to Attendance's
// narrow authorization port until Positions is migrated into its own module.
type attendanceLeadership struct{ positions *services.PositionService }

func (a attendanceLeadership) LeadershipScope(ctx context.Context, teacherID, academicYearID uuid.UUID) (bool, []uuid.UUID, error) {
	wholeSchool, gradeIDs, err := a.positions.LeadershipScope(ctx, teacherID, academicYearID)
	if errors.Is(err, services.ErrInsufficientRank) {
		return false, nil, attendancemodule.ErrInsufficientRank
	}
	return wholeSchool, gradeIDs, err
}
