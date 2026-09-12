package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/models"
	"github.com/openschool-org/openschool/internal/services"
)

// GradeHandler exposes HTTP endpoints for managing grades.
type GradeHandler struct {
	service *services.GradeService
}

// NewGradeHandler constructs a GradeHandler with its service dependency.
func NewGradeHandler(service *services.GradeService) *GradeHandler {
	return &GradeHandler{service: service}
}

// Create creates a new grade level.
func (h *GradeHandler) Create(c *gin.Context) {
	var req models.CreateGradeRequest
	if err := bindStrict(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grade, err := h.service.CreateGrade(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, grade)
}

// GetByID retrieves a grade by ID.
func (h *GradeHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	grade, err := h.service.GetGrade(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "grade not found"})
		return
	}

	c.JSON(http.StatusOK, grade)
}

// List retrieves all grades ordered by sort_order then name.
func (h *GradeHandler) List(c *gin.Context) {
	grades, err := h.service.ListGrades(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, grades)
}

// Update updates a grade's name or sort order by ID.
func (h *GradeHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req models.UpdateGradeRequest
	if err := bindStrict(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grade, err := h.service.UpdateGrade(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, grade)
}

// Delete deletes a grade by ID; blocked if the grade is assigned to any class.
func (h *GradeHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DeleteGrade(c.Request.Context(), id); err != nil {
		if errors.Is(err, services.ErrGradeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "grade deleted"})
}
