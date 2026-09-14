package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/openschool-org/openschool/internal/handlers"
	attendancemodule "github.com/openschool-org/openschool/internal/modules/attendance"
	"github.com/openschool-org/openschool/internal/repositories"
	"github.com/openschool-org/openschool/internal/services"
)

func RegisterReportExportRoutes(admin *gin.RouterGroup, pool *pgxpool.Pool) {
	service := services.NewReportExportService(
		attendancemodule.NewService(attendancemodule.NewRepository(pool), nil, nil, nil),
		repositories.NewTermMarkRepository(pool),
		repositories.NewClassRepository(pool),
		repositories.NewTermRepository(pool),
		repositories.NewSubjectRepository(pool),
	)
	handler := handlers.NewReportExportHandler(service)

	admin.GET("/reports/attendance", handler.ExportAttendance)
	admin.GET("/reports/marks", handler.ExportMarks)
}
