package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/models"
	"github.com/openschool-org/openschool/internal/services"
)

// TermHandler exposes HTTP endpoints for managing academic terms.
type TermHandler struct {
	service *services.TermService
}

// NewTermHandler constructs a TermHandler with its service dependency.
func NewTermHandler(service *services.TermService) *TermHandler {
	return &TermHandler{service: service}
}

// Create creates a new term within an academic year.
func (h *TermHandler) Create(c *gin.Context) {
	var req models.CreateTermRequest
	if err := bindStrict(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	term, err := h.service.CreateTerm(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, term)
}

// ListByAcademicYear retrieves all terms for an academic year, ordered by sort_order.
func (h *TermHandler) ListByAcademicYear(c *gin.Context) {
	yearIDStr := c.Query("academic_year_id")
	yearID, err := uuid.Parse(yearIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or missing academic_year_id"})
		return
	}

	terms, err := h.service.ListTermsByAcademicYear(c.Request.Context(), yearID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, terms)
}

// GetCurrent retrieves the term marked as current.
func (h *TermHandler) GetCurrent(c *gin.Context) {
	term, err := h.service.GetCurrentTerm(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no current term"})
		return
	}

	c.JSON(http.StatusOK, term)
}

// SetCurrent marks a term as the current one; clears the previous current.
func (h *TermHandler) SetCurrent(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.SetCurrentTerm(c.Request.Context(), id); err != nil {
		if errors.Is(err, services.ErrTermNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "current term updated"})
}

// Update updates a term's name, dates, or sort order by ID.
func (h *TermHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req models.UpdateTermRequest
	if err := bindStrict(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	term, err := h.service.UpdateTerm(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, term)
}

// Delete deletes a term by ID; blocked if marks have been recorded against it.
func (h *TermHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DeleteTerm(c.Request.Context(), id); err != nil {
		if errors.Is(err, services.ErrTermNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "term deleted"})
}
