package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openschool-org/openschool/internal/models"
	"github.com/openschool-org/openschool/internal/services"
)

// ReportExportHandler exposes HTTP endpoints for exporting admin reports.
type ReportExportHandler struct {
	service *services.ReportExportService
}

// NewReportExportHandler constructs a ReportExportHandler with its service dependency.
func NewReportExportHandler(service *services.ReportExportService) *ReportExportHandler {
	return &ReportExportHandler{service: service}
}

// ExportAttendance returns a PDF attendance report for a class over a date range.
func (h *ReportExportHandler) ExportAttendance(c *gin.Context) {
	var req models.AttendanceReportRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pdfBytes, err := h.service.ExportAttendance(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", `attachment; filename="attendance-report.pdf"`)
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// ExportMarks returns a PDF marks report for a class/term/subject, one row per student.
func (h *ReportExportHandler) ExportMarks(c *gin.Context) {
	var req models.MarksReportRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pdfBytes, err := h.service.ExportMarks(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", `attachment; filename="marks-report.pdf"`)
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}
