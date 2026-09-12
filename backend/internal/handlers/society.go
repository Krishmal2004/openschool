package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/models"
	"github.com/openschool-org/openschool/internal/services"
)

// SocietyHandler exposes HTTP endpoints for managing societies and their members.
type SocietyHandler struct {
	service *services.SocietyService
}

// NewSocietyHandler constructs a SocietyHandler with its service dependency.
func NewSocietyHandler(service *services.SocietyService) *SocietyHandler {
	return &SocietyHandler{service: service}
}

// societyServiceError maps a society-service error to the appropriate HTTP status code.
func societyServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrSocietyNotFound), errors.Is(err, services.ErrSocietyMemberNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, services.ErrNotTeacherInCharge):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

// Create creates a society.
func (h *SocietyHandler) Create(c *gin.Context) {
	var req models.CreateSocietyRequest
	if err := bindStrict(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	society, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		societyServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, society)
}

// Update renames a society or reassigns its Teacher-in-Charge.
func (h *SocietyHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid society id"})
		return
	}

	var req models.UpdateSocietyRequest
	if err := bindStrict(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	society, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		societyServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, society)
}

// Delete deletes a society.
func (h *SocietyHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid society id"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		societyServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "society deleted"})
}

// List lists societies for an academic year.
func (h *SocietyHandler) List(c *gin.Context) {
	yearID, err := uuid.Parse(c.Query("academic_year_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "a valid academic_year_id is required"})
		return
	}

	list, err := h.service.ListByYear(c.Request.Context(), yearID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}

// ListYears lists academic years that have at least one society.
func (h *SocietyHandler) ListYears(c *gin.Context) {
	years, err := h.service.ListYears(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, years)
}

// ListMembers lists a society's roster.
func (h *SocietyHandler) ListMembers(c *gin.Context) {
	societyID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid society id"})
		return
	}

	members, err := h.service.ListMembers(c.Request.Context(), societyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, members)
}

// AssignMember adds a student to the society roster (Teacher-in-Charge or admin only).
func (h *SocietyHandler) AssignMember(c *gin.Context) {
	societyID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid society id"})
		return
	}

	var req models.AssignSocietyMemberRequest
	if err := bindStrict(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	actor, err := actorFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	member, err := h.service.AssignMember(c.Request.Context(), actor, societyID, req)
	if err != nil {
		societyServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, member)
}

// RemoveMember removes a student from the society roster (Teacher-in-Charge or admin only).
func (h *SocietyHandler) RemoveMember(c *gin.Context) {
	societyID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid society id"})
		return
	}

	memberID, err := uuid.Parse(c.Param("memberId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid member id"})
		return
	}

	actor, err := actorFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.RemoveMember(c.Request.Context(), actor, societyID, memberID); err != nil {
		societyServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "society member removed"})
}

// ListByStudent lists a student's society memberships across all years.
func (h *SocietyHandler) ListByStudent(c *gin.Context) {
	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student id"})
		return
	}

	list, err := h.service.ListMembershipsByStudent(c.Request.Context(), studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}
