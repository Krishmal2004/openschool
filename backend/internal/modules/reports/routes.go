package reports

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/apierror"
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
		// PDPA access log: this PDF carries per-student attendance data
		// (S11). Recorded before the PDF is sent — nothing has been
		// committed yet, so a logging failure can safely refuse the export
		// instead of shipping PDPA-sensitive data with no access record.
		if err := auditReportExport(c, audit, "attendance_report", req.ClassID, fmt.Sprintf("attendance %s to %s", req.From.Format("2006-01-02"), req.To.Format("2006-01-02"))); err != nil {
			apierror.RespondInternal(c, fmt.Errorf("failed to record report access audit: %w", err))
			return
		}
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
		// See the attendance handler above: logged before sending, so a
		// failure here refuses the export rather than losing the access record.
		if err := auditReportExport(c, audit, "marks_report", req.ClassID, fmt.Sprintf("term %s subject %s", req.TermID, req.SubjectID)); err != nil {
			apierror.RespondInternal(c, fmt.Errorf("failed to record report access audit: %w", err))
			return
		}
		c.Header("Content-Disposition", `attachment; filename="marks-report.pdf"`)
		c.Data(http.StatusOK, "application/pdf", data)
	})
}

// auditReportExport logs a class-scoped report export. A missing actor or an
// unparsable classID is treated as "nothing to log" rather than an error —
// both are pre-existing route/auth conditions, not audit failures.
func auditReportExport(c *gin.Context, audit ports.AuditRecorder, action, classID, reason string) error {
	if audit == nil {
		return fmt.Errorf("audit recorder not configured")
	}
	actor, err := middleware.UserIDFromContext(c)
	if err != nil {
		return nil
	}
	entityID, err := uuid.Parse(classID)
	if err != nil {
		return nil
	}
	return audit.Record(c.Request.Context(), "class_report_export", entityID, action, actor, nil, nil, reason)
}
