package reports

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openschool-org/openschool/internal/models"
)

func RegisterRoutes(admin *gin.RouterGroup, service *Service) {
	admin.GET("/reports/attendance", func(c *gin.Context) {
		var req models.AttendanceReportRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		data, err := service.ExportAttendance(c.Request.Context(), req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.Header("Content-Disposition", `attachment; filename="attendance-report.pdf"`)
		c.Data(http.StatusOK, "application/pdf", data)
	})

	admin.GET("/reports/marks", func(c *gin.Context) {
		var req models.MarksReportRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		data, err := service.ExportMarks(c.Request.Context(), req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.Header("Content-Disposition", `attachment; filename="marks-report.pdf"`)
		c.Data(http.StatusOK, "application/pdf", data)
	})
}
