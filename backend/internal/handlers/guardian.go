package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/models"
	"github.com/openschool-org/openschool/internal/services"
)

// GuardianHandler exposes HTTP endpoints for managing student guardians.
type GuardianHandler struct {
	service *services.GuardianService
}

// NewGuardianHandler constructs a GuardianHandler with its service dependency.
func NewGuardianHandler(service *services.GuardianService) *GuardianHandler {
	return &GuardianHandler{service: service}
}

// Create creates a new guardian record.
func (h *GuardianHandler) Create(c *gin.Context) {
	var req models.CreateGuardianRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	guardian, duplicates, err := h.service.CreateGuardian(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"guardian": guardian, "possible_duplicates": duplicates})
}

// GetByID gets a guardian record by ID.
func (h *GuardianHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	guardian, err := h.service.GetGuardian(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "guardian not found"})
		return
	}

	c.JSON(http.StatusOK, guardian)
}

// List returns every guardian on file, optionally filtered by name, phone, or email.
func (h *GuardianHandler) List(c *gin.Context) {
	orphansOnly := c.Query("orphans") == "true"
	guardians, err := h.service.ListGuardians(c.Request.Context(), c.Query("search"), orphansOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, guardians)
}

// ListStudents returns every student linked to this guardian, regardless of portal-login status — for the guardian directory.
func (h *GuardianHandler) ListStudents(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	students, err := h.service.ListStudentsFor(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, students)
}

// ListNotifications returns notifications sent to this guardian's portal account, empty if they have none.
func (h *GuardianHandler) ListNotifications(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	notifications, err := h.service.ListNotifications(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

// Update updates a guardian record.
func (h *GuardianHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req models.UpdateGuardianRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	guardian, err := h.service.UpdateGuardian(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, guardian)
}

// Delete deletes a guardian record outright — blocked while linked to any student.
func (h *GuardianHandler) Delete(c *gin.Context) {
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

	if err := h.service.DeleteGuardian(c.Request.Context(), id, actor.ID); err != nil {
		switch {
		case errors.Is(err, services.ErrGuardianNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrGuardianInUse):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "guardian deleted"})
}

// LinkToStudent links an existing guardian to a student.
func (h *GuardianHandler) LinkToStudent(c *gin.Context) {
	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student id"})
		return
	}

	var req models.LinkGuardianRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.LinkToStudent(c.Request.Context(), studentID, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "guardian linked to student"})
}

// UnlinkFromStudent removes a guardian link from a student.
func (h *GuardianHandler) UnlinkFromStudent(c *gin.Context) {
	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student id"})
		return
	}

	guardianID, err := uuid.Parse(c.Param("guardian_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid guardian id"})
		return
	}

	if err := h.service.UnlinkFromStudent(c.Request.Context(), studentID, guardianID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "guardian unlinked from student"})
}

// ListByStudent gets all guardians linked to a student.
func (h *GuardianHandler) ListByStudent(c *gin.Context) {
	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student id"})
		return
	}

	guardians, err := h.service.ListByStudent(c.Request.Context(), studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, guardians)
}

// SetPrimaryContact sets a guardian as the primary contact for a student.
func (h *GuardianHandler) SetPrimaryContact(c *gin.Context) {
	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student id"})
		return
	}

	guardianID, err := uuid.Parse(c.Param("guardian_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid guardian id"})
		return
	}

	if err := h.service.SetPrimaryContact(c.Request.Context(), studentID, guardianID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "primary contact updated"})
}

// ProvisionLogin creates an identity-provider account for an existing guardian and links it, giving them parent-portal access.
func (h *GuardianHandler) ProvisionLogin(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req models.ProvisionGuardianLoginRequest
	if err := bindStrict(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	actor, err := actorFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	guardian, err := h.service.ProvisionLogin(c.Request.Context(), id, req, actor.ID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrGuardianNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrGuardianAlreadyProvisioned), errors.Is(err, services.ErrGuardianMissingEmail), errors.Is(err, services.ErrGuardianMissingNIC):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, guardian)
}
