package people

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/middleware"
	"github.com/openschool-org/openschool/internal/models"
)

type studentPortfolioService interface {
	TeacherProfileIDForUser(context.Context, uuid.UUID) *uuid.UUID
	CreateProgressReport(context.Context, uuid.UUID, models.CreateProgressReportRequest, *uuid.UUID) (any, error)
	ListProgressReports(context.Context, uuid.UUID) (any, error)
	UpdateProgressReport(context.Context, uuid.UUID, uuid.UUID, models.UpdateProgressReportRequest) (any, error)
	DeleteProgressReport(context.Context, uuid.UUID, uuid.UUID) error
	CreateActivity(context.Context, uuid.UUID, models.CreateActivityRequest) (any, error)
	ListActivities(context.Context, uuid.UUID) (any, error)
	UpdateActivity(context.Context, uuid.UUID, uuid.UUID, models.UpdateActivityRequest) (any, error)
	DeleteActivity(context.Context, uuid.UUID, uuid.UUID) error
	CreateLeadershipRole(context.Context, uuid.UUID, models.CreateLeadershipRoleRequest) (any, error)
	ListLeadershipRoles(context.Context, uuid.UUID) (any, error)
	DeleteLeadershipRole(context.Context, uuid.UUID, uuid.UUID) error
	CreateAward(context.Context, uuid.UUID, models.CreateAwardRequest) (any, error)
	ListAwards(context.Context, uuid.UUID) (any, error)
	DeleteAward(context.Context, uuid.UUID, uuid.UUID) error
	CreateDisciplinaryRecord(context.Context, uuid.UUID, models.CreateDisciplinaryRecordRequest, *uuid.UUID) (any, error)
	ListDisciplinaryRecords(context.Context, uuid.UUID) (any, error)
	DeleteDisciplinaryRecord(context.Context, uuid.UUID, uuid.UUID) error
}

// StudentPortfolioHandler exposes HTTP endpoints for a student's activities, awards, and disciplinary records.
type StudentPortfolioHandler struct {
	service studentPortfolioService
}

// NewStudentPortfolioHandler constructs a StudentPortfolioHandler with its service dependency.
func NewStudentPortfolioHandler(service studentPortfolioService) *StudentPortfolioHandler {
	return &StudentPortfolioHandler{service: service}
}

// studentIDParam parses and validates the :id URL parameter as a student ID.
func studentIDParam(c *gin.Context) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student id"})
		return uuid.UUID{}, false
	}
	return id, true
}

// recordIDParam parses and validates the :record_id URL parameter.
func recordIDParam(c *gin.Context) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param("record_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid record id"})
		return uuid.UUID{}, false
	}
	return id, true
}

// optionalActorID returns the caller's user ID, or nil if the claim is missing or invalid.
func optionalActorID(c *gin.Context) *uuid.UUID {
	id, err := middleware.UserIDFromContext(c)
	if err != nil {
		return nil
	}
	return &id
}

// --- Progress reports -------------------------------------------------------

// CreateProgressReport adds a progress report.
func (h *StudentPortfolioHandler) CreateProgressReport(c *gin.Context) {
	studentID, ok := studentIDParam(c)
	if !ok {
		return
	}
	var req models.CreateProgressReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var writtenBy *uuid.UUID
	if callerID, err := middleware.UserIDFromContext(c); err == nil {
		writtenBy = h.service.TeacherProfileIDForUser(c.Request.Context(), callerID)
	}

	report, err := h.service.CreateProgressReport(c.Request.Context(), studentID, req, writtenBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, report)
}

// ListProgressReports lists a student's progress reports.
func (h *StudentPortfolioHandler) ListProgressReports(c *gin.Context) {
	studentID, ok := studentIDParam(c)
	if !ok {
		return
	}
	reports, err := h.service.ListProgressReports(c.Request.Context(), studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, reports)
}

// UpdateProgressReport updates a progress report.
func (h *StudentPortfolioHandler) UpdateProgressReport(c *gin.Context) {
	studentID, ok := studentIDParam(c)
	if !ok {
		return
	}
	id, ok := recordIDParam(c)
	if !ok {
		return
	}
	var req models.UpdateProgressReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	report, err := h.service.UpdateProgressReport(c.Request.Context(), id, studentID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, report)
}

// DeleteProgressReport deletes a progress report.
func (h *StudentPortfolioHandler) DeleteProgressReport(c *gin.Context) {
	studentID, ok := studentIDParam(c)
	if !ok {
		return
	}
	id, ok := recordIDParam(c)
	if !ok {
		return
	}
	if err := h.service.DeleteProgressReport(c.Request.Context(), id, studentID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "progress report deleted"})
}

// --- Activities --------------------------------------------------------------

// CreateActivity adds a student activity (club/sport/society/competition).
func (h *StudentPortfolioHandler) CreateActivity(c *gin.Context) {
	studentID, ok := studentIDParam(c)
	if !ok {
		return
	}
	var req models.CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	activity, err := h.service.CreateActivity(c.Request.Context(), studentID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, activity)
}

// ListActivities lists a student's activities.
func (h *StudentPortfolioHandler) ListActivities(c *gin.Context) {
	studentID, ok := studentIDParam(c)
	if !ok {
		return
	}
	activities, err := h.service.ListActivities(c.Request.Context(), studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, activities)
}

// UpdateActivity updates a student activity.
func (h *StudentPortfolioHandler) UpdateActivity(c *gin.Context) {
	studentID, ok := studentIDParam(c)
	if !ok {
		return
	}
	id, ok := recordIDParam(c)
	if !ok {
		return
	}
	var req models.UpdateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	activity, err := h.service.UpdateActivity(c.Request.Context(), id, studentID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, activity)
}

// DeleteActivity deletes a student activity.
func (h *StudentPortfolioHandler) DeleteActivity(c *gin.Context) {
	studentID, ok := studentIDParam(c)
	if !ok {
		return
	}
	id, ok := recordIDParam(c)
	if !ok {
		return
	}
	if err := h.service.DeleteActivity(c.Request.Context(), id, studentID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "activity deleted"})
}

// --- Leadership roles ----------------------------------------------------------

// CreateLeadershipRole adds a student leadership role.
func (h *StudentPortfolioHandler) CreateLeadershipRole(c *gin.Context) {
	studentID, ok := studentIDParam(c)
	if !ok {
		return
	}
	var req models.CreateLeadershipRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	role, err := h.service.CreateLeadershipRole(c.Request.Context(), studentID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, role)
}

// ListLeadershipRoles lists a student's leadership roles.
func (h *StudentPortfolioHandler) ListLeadershipRoles(c *gin.Context) {
	studentID, ok := studentIDParam(c)
	if !ok {
		return
	}
	roles, err := h.service.ListLeadershipRoles(c.Request.Context(), studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, roles)
}

// DeleteLeadershipRole deletes a student leadership role.
func (h *StudentPortfolioHandler) DeleteLeadershipRole(c *gin.Context) {
	studentID, ok := studentIDParam(c)
	if !ok {
		return
	}
	id, ok := recordIDParam(c)
	if !ok {
		return
	}
	if err := h.service.DeleteLeadershipRole(c.Request.Context(), id, studentID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "leadership role deleted"})
}

// --- Awards ------------------------------------------------------------------

// CreateAward adds a student award.
func (h *StudentPortfolioHandler) CreateAward(c *gin.Context) {
	studentID, ok := studentIDParam(c)
	if !ok {
		return
	}
	var req models.CreateAwardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	award, err := h.service.CreateAward(c.Request.Context(), studentID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, award)
}

// ListAwards lists a student's awards.
func (h *StudentPortfolioHandler) ListAwards(c *gin.Context) {
	studentID, ok := studentIDParam(c)
	if !ok {
		return
	}
	awards, err := h.service.ListAwards(c.Request.Context(), studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, awards)
}

// DeleteAward deletes a student award.
func (h *StudentPortfolioHandler) DeleteAward(c *gin.Context) {
	studentID, ok := studentIDParam(c)
	if !ok {
		return
	}
	id, ok := recordIDParam(c)
	if !ok {
		return
	}
	if err := h.service.DeleteAward(c.Request.Context(), id, studentID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "award deleted"})
}

// --- Disciplinary records -------------------------------------------------------

// CreateDisciplinaryRecord adds a disciplinary record.
func (h *StudentPortfolioHandler) CreateDisciplinaryRecord(c *gin.Context) {
	studentID, ok := studentIDParam(c)
	if !ok {
		return
	}
	var req models.CreateDisciplinaryRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	record, err := h.service.CreateDisciplinaryRecord(c.Request.Context(), studentID, req, optionalActorID(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, record)
}

// ListDisciplinaryRecords lists a student's disciplinary records.
func (h *StudentPortfolioHandler) ListDisciplinaryRecords(c *gin.Context) {
	studentID, ok := studentIDParam(c)
	if !ok {
		return
	}
	records, err := h.service.ListDisciplinaryRecords(c.Request.Context(), studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, records)
}

// DeleteDisciplinaryRecord deletes a disciplinary record.
func (h *StudentPortfolioHandler) DeleteDisciplinaryRecord(c *gin.Context) {
	studentID, ok := studentIDParam(c)
	if !ok {
		return
	}
	id, ok := recordIDParam(c)
	if !ok {
		return
	}
	if err := h.service.DeleteDisciplinaryRecord(c.Request.Context(), id, studentID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "disciplinary record deleted"})
}

// RegisterStudentPortfolioRoutes owns all student portfolio HTTP endpoints.
func RegisterStudentPortfolioRoutes(teacherOrAdmin, studentAccess *gin.RouterGroup, service *StudentPortfolioService) {
	handler := NewStudentPortfolioHandler(service)
	teacherOrAdmin.POST("/students/:id/progress-reports", handler.CreateProgressReport)
	studentAccess.GET("/students/:id/progress-reports", handler.ListProgressReports)
	teacherOrAdmin.PUT("/students/:id/progress-reports/:record_id", handler.UpdateProgressReport)
	teacherOrAdmin.DELETE("/students/:id/progress-reports/:record_id", handler.DeleteProgressReport)
	teacherOrAdmin.POST("/students/:id/activities", handler.CreateActivity)
	studentAccess.GET("/students/:id/activities", handler.ListActivities)
	teacherOrAdmin.PUT("/students/:id/activities/:record_id", handler.UpdateActivity)
	teacherOrAdmin.DELETE("/students/:id/activities/:record_id", handler.DeleteActivity)
	teacherOrAdmin.POST("/students/:id/leadership-roles", handler.CreateLeadershipRole)
	studentAccess.GET("/students/:id/leadership-roles", handler.ListLeadershipRoles)
	teacherOrAdmin.DELETE("/students/:id/leadership-roles/:record_id", handler.DeleteLeadershipRole)
	teacherOrAdmin.POST("/students/:id/awards", handler.CreateAward)
	studentAccess.GET("/students/:id/awards", handler.ListAwards)
	teacherOrAdmin.DELETE("/students/:id/awards/:record_id", handler.DeleteAward)
	teacherOrAdmin.POST("/students/:id/disciplinary-records", handler.CreateDisciplinaryRecord)
	studentAccess.GET("/students/:id/disciplinary-records", handler.ListDisciplinaryRecords)
	teacherOrAdmin.DELETE("/students/:id/disciplinary-records/:record_id", handler.DeleteDisciplinaryRecord)
}
