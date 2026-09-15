package search

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(admin *gin.RouterGroup, service *Service) {
	admin.GET("/admin/search", func(c *gin.Context) {
		result, err := service.Global(c.Request.Context(), c.Query("q"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	})
}
