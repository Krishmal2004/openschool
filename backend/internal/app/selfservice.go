package app

import (
	"github.com/jackc/pgx/v5/pgxpool"
	peoplemodule "github.com/openschool-org/openschool/internal/modules/people"
	selfservicemodule "github.com/openschool-org/openschool/internal/modules/selfservice"
	timetablemodule "github.com/openschool-org/openschool/internal/modules/timetable"
)

type portalProfiles struct {
	student *selfservicemodule.StudentProfiles
	teacher *selfservicemodule.TeacherProfiles
}

func newPortalProfiles(pool *pgxpool.Pool) portalProfiles {
	repository := selfservicemodule.NewRepository(pool)
	return portalProfiles{
		student: selfservicemodule.NewStudentProfiles(repository),
		teacher: selfservicemodule.NewTeacherProfiles(repository),
	}
}

func registerSelfService(groups HTTPGroups, pool *pgxpool.Pool, profiles portalProfiles, shared sharedServices) {
	selfservicemodule.RegisterRoutes(
		groups.Student, groups.Parent, groups.Teacher, pool,
		profiles.student, profiles.teacher, peoplemodule.NewGuardianAccess(pool), timetablemodule.NewReader(pool),
		shared.leadership, shared.studentLeadership, shared.dashboard,
	)
}
