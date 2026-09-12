package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openschool-org/openschool/internal/services"
)

// IdentityReconciliationHandler exposes the admin endpoint for finding identity-provider accounts with no matching local user.
type IdentityReconciliationHandler struct {
	service *services.IdentityReconciliationService
}

// NewIdentityReconciliationHandler constructs an IdentityReconciliationHandler with its service dependency.
func NewIdentityReconciliationHandler(service *services.IdentityReconciliationService) *IdentityReconciliationHandler {
	return &IdentityReconciliationHandler{service: service}
}

// ListOrphaned returns identity-provider accounts with no matching local user row, left behind by a failed signup rollback.
func (h *IdentityReconciliationHandler) ListOrphaned(c *gin.Context) {
	orphaned, err := h.service.FindOrphaned(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orphaned)
}

// DeleteOrphaned re-verifies an account is still orphaned, then deletes it.
func (h *IdentityReconciliationHandler) DeleteOrphaned(c *gin.Context) {
	id := c.Param("id")

	actor, err := actorFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.DeleteOrphaned(c.Request.Context(), id, actor.ID); err != nil {
		if errors.Is(err, services.ErrOrphanNoLongerOrphaned) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "orphaned account deleted"})
}
