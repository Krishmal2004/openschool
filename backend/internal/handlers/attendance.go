package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/middleware"
	"github.com/openschool-org/openschool/internal/models"
	"github.com/openschool-org/openschool/internal/services"
)

// statusForAttendanceError maps a not-assigned-to-class failure to 403 instead of a genuine 404/400.
func statusForAttendanceError(err error, notFoundStatus int) int {
	if errors.Is(err, services.ErrNotAssignedToClass) || errors.Is(err, services.ErrInsufficientRank) {
		return http.StatusForbidden
	}
	if errors.Is(err, services.ErrSessionLocked) {
		return http.StatusLocked
	}
	return notFoundStatus
}

// AttendanceHandler exposes HTTP endpoints for marking and viewing class attendance.
type AttendanceHandler struct {
	service *services.AttendanceService
}

// NewAttendanceHandler constructs an AttendanceHandler with its service dependency.
func NewAttendanceHandler(service *services.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{service: service}
}

// actorFromContext builds a services.Actor from the JWT claims AuthMiddleware already set on the context.
func actorFromContext(c *gin.Context) (services.Actor, error) {
	id, err := middleware.UserIDFromContext(c)
	if err != nil {
		return services.Actor{}, fmt.Errorf("invalid caller identity")
	}

	var roleList []string
	if roles, ok := c.Get("roles"); ok {
		roleList, _ = roles.([]string)
	}

	return services.Actor{
		ID:       id,
		Email:    c.GetString("email"),
		FullName: strings.TrimSpace(c.GetString("given_name") + " " + c.GetString("family_name")),
		Role:     services.ResolveAppRole(roleList),
	}, nil
}

// CreateSession creates a new attendance session for a class on a specific date.
func (h *AttendanceHandler) CreateSession(c *gin.Context) {
	var req models.CreateAttendanceSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	actor, err := actorFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	session, err := h.service.CreateSession(c.Request.Context(), actor, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, session)
}

// GetSession gets an attendance session by ID.
func (h *AttendanceHandler) GetSession(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	actor, err := actorFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	session, err := h.service.GetSession(c.Request.Context(), actor, id)
	if err != nil {
		c.JSON(statusForAttendanceError(err, http.StatusNotFound), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, session)
}

// DeleteSession deletes a session and every attendance record already written for it.
func (h *AttendanceHandler) DeleteSession(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	actor, err := actorFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.DeleteSession(c.Request.Context(), actor, id); err != nil {
		c.JSON(statusForAttendanceError(err, http.StatusNotFound), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "attendance session deleted"})
}

// ListSessionsByClass gets all attendance sessions for a specific class.
func (h *AttendanceHandler) ListSessionsByClass(c *gin.Context) {
	classID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid class id"})
		return
	}

	actor, err := actorFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	sessions, err := h.service.ListSessionsByClass(c.Request.Context(), actor, classID)
	if err != nil {
		c.JSON(statusForAttendanceError(err, http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sessions)
}

// ListSessionsByDate returns every attendance session on a given date across all classes, with counts.
func (h *AttendanceHandler) ListSessionsByDate(c *gin.Context) {
	date := c.Query("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date query parameter is required"})
		return
	}

	actor, err := actorFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	sessions, err := h.service.ListSessionsByDate(c.Request.Context(), actor, date)
	if err != nil {
		c.JSON(statusForAttendanceError(err, http.StatusBadRequest), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sessions)
}

// MarkAttendance marks attendance for all students in a session.
func (h *AttendanceHandler) MarkAttendance(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}

	var req models.MarkAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	actor, err := actorFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.MarkAttendance(c.Request.Context(), actor, id, req); err != nil {
		c.JSON(statusForAttendanceError(err, http.StatusBadRequest), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "attendance marked"})
}

// ListBySession gets all attendance records for a session.
func (h *AttendanceHandler) ListBySession(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}

	actor, err := actorFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	records, err := h.service.ListBySession(c.Request.Context(), actor, id)
	if err != nil {
		c.JSON(statusForAttendanceError(err, http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, records)
}

// ListByStudent gets all attendance records for a student.
func (h *AttendanceHandler) ListByStudent(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student id"})
		return
	}

	actor, err := actorFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	records, err := h.service.ListByStudentForTeacher(c.Request.Context(), actor, id)
	if err != nil {
		c.JSON(statusForAttendanceError(err, http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, records)
}

// GetSummary gets attendance summary for a student in a specific class.
func (h *AttendanceHandler) GetSummary(c *gin.Context) {
	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student id"})
		return
	}

	classID, err := uuid.Parse(c.Param("class_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid class id"})
		return
	}

	actor, err := actorFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	summary, err := h.service.GetSummaryForTeacher(c.Request.Context(), actor, studentID, classID)
	if err != nil {
		c.JSON(statusForAttendanceError(err, http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}
