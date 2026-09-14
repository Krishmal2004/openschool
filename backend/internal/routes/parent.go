package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/openschool-org/openschool/internal/handlers"
	attendancemodule "github.com/openschool-org/openschool/internal/modules/attendance"
	peoplemodule "github.com/openschool-org/openschool/internal/modules/people"
	timetablemodule "github.com/openschool-org/openschool/internal/modules/timetable"
)

func RegisterParentRoutes(parent *gin.RouterGroup, timetables *timetablemodule.Reader, pool *pgxpool.Pool) {
	attendanceReader := attendancemodule.NewReader(attendancemodule.NewRepository(pool))
	handler := handlers.NewParentHandler(peoplemodule.NewGuardianAccess(pool), attendanceReader, newTermMarkRunner(pool), timetables)

	parent.GET("/me/children", handler.ListChildren)
	parent.GET("/me/children/:id/attendance", handler.ChildAttendance)
	parent.GET("/me/children/:id/marks", handler.ChildMarks)
	parent.GET("/me/children/:id/timetable", handler.ChildTimetable)
}
