package app

import (
	"context"
	"errors"

	"github.com/google/uuid"
	attendancemodule "github.com/openschool-org/openschool/internal/modules/attendance"
	leadershipmodule "github.com/openschool-org/openschool/internal/modules/leadership"
)

// attendanceLeadership adapts Leadership's rank error to Attendance's
// domain-specific authorization error.
type attendanceLeadership struct{ positions *leadershipmodule.Service }

func (a attendanceLeadership) LeadershipScope(ctx context.Context, teacherID, academicYearID uuid.UUID) (bool, []uuid.UUID, error) {
	wholeSchool, gradeIDs, err := a.positions.LeadershipScope(ctx, teacherID, academicYearID)
	if errors.Is(err, leadershipmodule.ErrInsufficientRank) {
		return false, nil, attendancemodule.ErrInsufficientRank
	}
	return wholeSchool, gradeIDs, err
}
