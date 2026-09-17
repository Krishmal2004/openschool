package search

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openschool-org/openschool/internal/apierror"
)

func RegisterRoutes(admin *gin.RouterGroup, service *Service) {
	admin.GET("/admin/search", func(c *gin.Context) {
		result, err := service.Global(c.Request.Context(), c.Query("q"))
		if err != nil {
			apierror.RespondInternal(c, err)
			return
		}
		c.JSON(http.StatusOK, result)
	})
}
