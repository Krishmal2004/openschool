package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/middleware"
	"github.com/openschool-org/openschool/internal/models"
	"github.com/openschool-org/openschool/internal/services"
)

// NonAcademicStaffHandler exposes HTTP endpoints for managing non-academic staff.
type NonAcademicStaffHandler struct {
	service *services.NonAcademicStaffService
}

// NewNonAcademicStaffHandler constructs a NonAcademicStaffHandler with its service dependency.
func NewNonAcademicStaffHandler(service *services.NonAcademicStaffService) *NonAcademicStaffHandler {
	return &NonAcademicStaffHandler{service: service}
}

// Create creates a new non-academic staff record (Lab Assistant, Librarian, Office Staff, etc.) — no login/IDP account.
func (h *NonAcademicStaffHandler) Create(c *gin.Context) {
	var req models.CreateNonAcademicStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	staff, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, staff)
}

// GetByID gets non-academic staff by ID.
func (h *NonAcademicStaffHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	staff, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, staff)
}

// List lists non-academic staff.
func (h *NonAcademicStaffHandler) List(c *gin.Context) {
	staff, err := h.service.List(c.Request.Context(), c.Query("search"), c.Query("designation"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, staff)
}

// Update updates non-academic staff.
func (h *NonAcademicStaffHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req models.UpdateNonAcademicStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	staff, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, staff)
}

// UpdateEmploymentStatus updates non-academic staff employment status.
func (h *NonAcademicStaffHandler) UpdateEmploymentStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req models.UpdateNonAcademicStaffEmploymentStatusRequest
	if err := bindStrict(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	staff, err := h.service.SetEmploymentStatus(c.Request.Context(), id, req.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, staff)
}

// UpdateHouse updates non-academic staff house (admin override, audited).
func (h *NonAcademicStaffHandler) UpdateHouse(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req models.UpdateNonAcademicStaffHouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	actorID, err := middleware.UserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid caller identity"})
		return
	}

	staff, err := h.service.UpdateHouse(c.Request.Context(), id, req.HouseID, actorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, staff)
}

// Delete permanently removes the staff record (nothing else references non-academic staff rows).
func (h *NonAcademicStaffHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, services.ErrNonAcademicStaffNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "staff member deleted"})
}
