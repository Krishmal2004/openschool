package audit

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/apierror"
	"github.com/openschool-org/openschool/internal/platform/httpx"
)

func RegisterRoutes(admin *gin.RouterGroup, service *Service) {
	admin.GET("/audit-logs", func(c *gin.Context) {
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
		page := httpx.ParsePage(c)
		logs, total, err := service.List(c, entityType, entityID, page.Limit, page.Offset)
		if err != nil {
			apierror.RespondInternal(c, err)
			return
		}
		c.JSON(http.StatusOK, httpx.Page[AuditLogResponse]{Items: logs, Total: total, Limit: page.Limit, Offset: page.Offset})
	})
}
