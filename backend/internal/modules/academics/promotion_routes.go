package academics

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/models"
	"github.com/openschool-org/openschool/internal/platform/httpx"
)

type PromotionRunner interface {
	Preview(context.Context, uuid.UUID, uuid.UUID, *uuid.UUID) ([]models.PromotionPreviewRow, error)
	CommitAssignments(context.Context, models.CommitAssignmentsRequest) (int, error)
}

// RegisterPromotionRoutes owns the promotion HTTP boundary. The runner is
// injected so its persistence implementation can move behind the module
// repository boundary without changing the API contract.
func RegisterPromotionRoutes(admin *gin.RouterGroup, runner PromotionRunner) {
	admin.GET("/promotion/preview", func(c *gin.Context) {
		source, err := uuid.Parse(c.Query("source_year_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "a valid source_year_id is required"})
			return
		}
		target, err := uuid.Parse(c.Query("target_year_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "a valid target_year_id is required"})
			return
		}
		var term *uuid.UUID
		if raw := c.Query("rank_by_term_id"); raw != "" {
			id, e := uuid.Parse(raw)
			if e != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rank_by_term_id"})
				return
			}
			term = &id
		}
		rows, err := runner.Preview(c.Request.Context(), source, target, term)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, rows)
	})
	admin.POST("/promotion/commit", func(c *gin.Context) {
		var request models.CommitAssignmentsRequest
		if err := httpx.BindStrict(c, &request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		count, err := runner.CommitAssignments(c.Request.Context(), request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"assigned": count})
	})
}
