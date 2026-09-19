package dashboard

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openschool-org/openschool/internal/apierror"
)

func RegisterRoutes(admin *gin.RouterGroup, service *Service) {
	admin.GET("/dashboard/analytics", func(c *gin.Context) {
		analytics, err := service.Analytics(c.Request.Context())
		if err != nil {
			apierror.RespondInternal(c, err)
			return
		}
		c.JSON(http.StatusOK, analytics)
	})
}
