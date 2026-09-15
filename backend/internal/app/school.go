package app

import (
	"github.com/jackc/pgx/v5/pgxpool"
	auditmodule "github.com/openschool-org/openschool/internal/modules/audit"
	curriculummodule "github.com/openschool-org/openschool/internal/modules/curriculum"
	schoolmodule "github.com/openschool-org/openschool/internal/modules/school"
	"github.com/openschool-org/openschool/internal/ports"
)

func registerSchool(groups HTTPGroups, pool *pgxpool.Pool, audit *auditmodule.Service) ports.HouseAssignments {
	curriculummodule.RegisterPresetRoutes(groups.Admin, pool)
	curriculummodule.RegisterMediumRoutes(groups.Admin, groups.Protected, pool)
	curriculummodule.RegisterLevelRoutes(groups.Admin, groups.Protected, pool)

	houseService := schoolmodule.NewHouseService(pool, audit)
	schoolmodule.RegisterHouseRoutes(groups.Admin, groups.TeacherOrAdmin, houseService)
	schoolmodule.RegisterSchoolRoutes(groups.Admin, groups.TeacherOrAdmin, groups.Protected, pool)
	schoolmodule.RegisterGradeRoutes(groups.Admin, groups.TeacherOrAdmin, pool)
	schoolmodule.RegisterTermRoutes(groups.Admin, groups.Protected, pool)
	return houseService
}
