package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/middleware"
	"github.com/openschool-org/openschool/internal/models"
	"github.com/openschool-org/openschool/internal/services"
)

// parseYearMonth reads and validates the year/month query parameters shared by this handler's endpoints.
func parseYearMonth(c *gin.Context) (int, int, error) {
	year, err := strconv.Atoi(c.Query("year"))
	if err != nil {
		return 0, 0, fmt.Errorf("invalid or missing year")
	}
	month, err := strconv.Atoi(c.Query("month"))
	if err != nil || month < 1 || month > 12 {
		return 0, 0, fmt.Errorf("invalid or missing month (expected 1-12)")
	}
	return year, month, nil
}

// StaffAttendanceHandler exposes HTTP endpoints for teacher/staff self-attendance.
type StaffAttendanceHandler struct {
	service *services.StaffAttendanceService
}

// NewStaffAttendanceHandler constructs a StaffAttendanceHandler with its service dependency.
func NewStaffAttendanceHandler(service *services.StaffAttendanceService) *StaffAttendanceHandler {
	return &StaffAttendanceHandler{service: service}
}

// Mark upserts a single teacher's or non-academic staff member's attendance for a date.
func (h *StaffAttendanceHandler) Mark(c *gin.Context) {
	var req models.MarkStaffAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	markedBy, err := middleware.UserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid caller identity"})
		return
	}

	record, err := h.service.MarkAttendance(c.Request.Context(), req, markedBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, record)
}

// ListByDate returns every staff member's attendance status for a given date.
func (h *StaffAttendanceHandler) ListByDate(c *gin.Context) {
	date, err := time.Parse("2006-01-02", c.Query("date"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or missing date (expected YYYY-MM-DD)"})
		return
	}

	teachers, staff, err := h.service.ListByDate(c.Request.Context(), date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"teachers": teachers, "non_academic_staff": staff})
}

// MonthlySummary returns per-status attendance counts for every active staff member in a given month.
func (h *StaffAttendanceHandler) MonthlySummary(c *gin.Context) {
	year, month, err := parseYearMonth(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	from := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	to := from.AddDate(0, 1, -1)

	teachers, staff, err := h.service.MonthlySummary(c.Request.Context(), from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"teachers": teachers, "non_academic_staff": staff})
}

// TeacherHistory returns a teacher's staff-attendance history for a month.
func (h *StaffAttendanceHandler) TeacherHistory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	year, month, err := parseYearMonth(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	from := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	to := from.AddDate(0, 1, -1)

	records, err := h.service.TeacherHistory(c.Request.Context(), id, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, records)
}

// NonAcademicStaffHistory returns a non-academic staff member's attendance history for a month.
func (h *StaffAttendanceHandler) NonAcademicStaffHistory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	year, month, err := parseYearMonth(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	from := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	to := from.AddDate(0, 1, -1)

	records, err := h.service.NonAcademicStaffHistory(c.Request.Context(), id, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, records)
}
