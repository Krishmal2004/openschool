package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/models"
	"github.com/openschool-org/openschool/internal/services"
)

// SchoolHandler exposes HTTP endpoints for the single-row school profile.
type SchoolHandler struct {
	service *services.SchoolService
}

// NewSchoolHandler constructs a SchoolHandler with its service dependency.
func NewSchoolHandler(service *services.SchoolService) *SchoolHandler {
	return &SchoolHandler{service: service}
}

// Create creates the school record for this instance.
func (h *SchoolHandler) Create(c *gin.Context) {
	var req models.CreateSchoolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	school, err := h.service.CreateSchool(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, school)
}

// Get retrieves the school record for this instance.
func (h *SchoolHandler) Get(c *gin.Context) {
	school, err := h.service.GetSchool(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "school not found"})
		return
	}

	c.JSON(http.StatusOK, school)
}

// Update updates the school record by ID.
func (h *SchoolHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req models.UpdateSchoolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	school, err := h.service.UpdateSchool(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, school)
}

// CreateAcademicYear creates a new academic year.
func (h *SchoolHandler) CreateAcademicYear(c *gin.Context) {
	var req models.CreateAcademicYearRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	year, err := h.service.CreateAcademicYear(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, year)
}

// ListAcademicYears retrieves all academic years ordered by start date.
func (h *SchoolHandler) ListAcademicYears(c *gin.Context) {
	years, err := h.service.ListAcademicYears(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, years)
}

// GetCurrentAcademicYear retrieves the academic year marked as current.
func (h *SchoolHandler) GetCurrentAcademicYear(c *gin.Context) {
	year, err := h.service.GetCurrentAcademicYear(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no current academic year"})
		return
	}

	c.JSON(http.StatusOK, year)
}

// SetCurrentAcademicYear marks an academic year as the current one; clears the previous current.
func (h *SchoolHandler) SetCurrentAcademicYear(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.SetCurrentAcademicYear(c.Request.Context(), id); err != nil {
		if errors.Is(err, services.ErrAcademicYearNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "current academic year updated"})
}

// DeleteAcademicYear deletes an academic year by ID.
func (h *SchoolHandler) DeleteAcademicYear(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DeleteAcademicYear(c.Request.Context(), id); err != nil {
		if errors.Is(err, services.ErrAcademicYearNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "academic year deleted"})
}
