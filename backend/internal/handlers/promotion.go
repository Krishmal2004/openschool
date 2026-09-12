package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/models"
	"github.com/openschool-org/openschool/internal/services"
)

// PromotionHandler exposes HTTP endpoints for the year-end student promotion workflow.
type PromotionHandler struct {
	service *services.PromotionService
}

// NewPromotionHandler constructs a PromotionHandler with its service dependency.
func NewPromotionHandler(service *services.PromotionService) *PromotionHandler {
	return &PromotionHandler{service: service}
}

// Preview computes next-grade + suggested target class per student — no writes.
func (h *PromotionHandler) Preview(c *gin.Context) {
	sourceYearID, err := uuid.Parse(c.Query("source_year_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "a valid source_year_id is required"})
		return
	}
	targetYearID, err := uuid.Parse(c.Query("target_year_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "a valid target_year_id is required"})
		return
	}

	var rankByTermID *uuid.UUID
	if raw := c.Query("rank_by_term_id"); raw != "" {
		id, err := uuid.Parse(raw)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rank_by_term_id"})
			return
		}
		rankByTermID = &id
	}

	rows, err := h.service.Preview(c.Request.Context(), sourceYearID, targetYearID, rankByTermID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rows)
}

// Commit bulk-writes class_students rows — shared by promotion-commit and general reassignment/shuffle.
func (h *PromotionHandler) Commit(c *gin.Context) {
	var req models.CommitAssignmentsRequest
	if err := bindStrict(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	count, err := h.service.CommitAssignments(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"assigned": count})
}
