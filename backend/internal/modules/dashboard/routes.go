package dashboard

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(admin *gin.RouterGroup, service *Service) {
	admin.GET("/dashboard/analytics", func(c *gin.Context) {
		analytics, err := service.Analytics(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, analytics)
	})
}
