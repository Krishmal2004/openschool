package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/services"
)

// AuditHandler exposes the read-only audit-log endpoint for admins.
type AuditHandler struct {
	service *services.AuditService
}

// NewAuditHandler constructs an AuditHandler with its service dependency.
func NewAuditHandler(service *services.AuditService) *AuditHandler {
	return &AuditHandler{service: service}
}

// List returns the admin-only trail of manual house re-assignments and attendance edits made after the 24h lock.
func (h *AuditHandler) List(c *gin.Context) {
	entityType := c.Query("entity_type")

	var entityID *uuid.UUID
	if raw := c.Query("entity_id"); raw != "" {
		parsed, err := uuid.Parse(raw)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entity_id"})
			return
		}
		entityID = &parsed
	}

	logs, err := h.service.List(c.Request.Context(), entityType, entityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}
