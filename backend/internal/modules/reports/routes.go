package reports

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/middleware"
	"github.com/openschool-org/openschool/internal/ports"
)

func RegisterRoutes(admin *gin.RouterGroup, service *Service, audit ports.AuditRecorder) {
	admin.GET("/reports/attendance", func(c *gin.Context) {
		var req AttendanceReportRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		data, err := service.ExportAttendance(c.Request.Context(), req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// PDPA access log: this PDF carries per-student attendance data (S11).
		auditReportExport(c, audit, "attendance_report", req.ClassID, fmt.Sprintf("attendance %s to %s", req.From.Format("2006-01-02"), req.To.Format("2006-01-02")))
		c.Header("Content-Disposition", `attachment; filename="attendance-report.pdf"`)
		c.Data(http.StatusOK, "application/pdf", data)
	})

	admin.GET("/reports/marks", func(c *gin.Context) {
		var req MarksReportRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		data, err := service.ExportMarks(c.Request.Context(), req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// PDPA access log: this PDF carries per-student marks data (S11).
		auditReportExport(c, audit, "marks_report", req.ClassID, fmt.Sprintf("term %s subject %s", req.TermID, req.SubjectID))
		c.Header("Content-Disposition", `attachment; filename="marks-report.pdf"`)
		c.Data(http.StatusOK, "application/pdf", data)
	})
}

// auditReportExport logs a class-scoped report export, best-effort: a
// logging failure must never fail an export that already succeeded.
func auditReportExport(c *gin.Context, audit ports.AuditRecorder, action, classID, reason string) {
	if audit == nil {
		return
	}
	actor, err := middleware.UserIDFromContext(c)
	if err != nil {
		return
	}
	entityID, err := uuid.Parse(classID)
	if err != nil {
		return
	}
	_ = audit.Record(c.Request.Context(), "class_report_export", entityID, action, actor, nil, nil, reason)
}
