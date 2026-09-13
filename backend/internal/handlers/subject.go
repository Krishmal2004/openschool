package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/models"
	"github.com/openschool-org/openschool/internal/services"
)

// SubjectHandler exposes HTTP endpoints for managing subjects.
type SubjectHandler struct {
	service *services.SubjectService
}

// NewSubjectHandler constructs a SubjectHandler with its service dependency.
func NewSubjectHandler(service *services.SubjectService) *SubjectHandler {
	return &SubjectHandler{service: service}
}

// Create creates a new subject with a unique name and code.
func (h *SubjectHandler) Create(c *gin.Context) {
	var req models.CreateSubjectRequest
	if err := bindStrict(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subject, err := h.service.CreateSubject(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, subject)
}

// GetByID retrieves a subject by ID.
func (h *SubjectHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	subject, err := h.service.GetSubject(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "subject not found"})
		return
	}

	c.JSON(http.StatusOK, subject)
}

// List retrieves all subjects ordered by name.
func (h *SubjectHandler) List(c *gin.Context) {
	subjects, err := h.service.ListSubjects(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subjects)
}

// Update updates a subject's name and code by ID.
func (h *SubjectHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req models.UpdateSubjectRequest
	if err := bindStrict(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subject, err := h.service.UpdateSubject(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subject)
}

// Delete deletes a subject by ID; blocked if the subject is assigned to a grade, class, or student selection.
func (h *SubjectHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DeleteSubject(c.Request.Context(), id); err != nil {
		if errors.Is(err, services.ErrSubjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "subject deleted"})
}
