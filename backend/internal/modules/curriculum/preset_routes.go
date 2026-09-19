package curriculum

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/openschool-org/openschool/internal/apierror"
)

// RegisterPresetRoutes owns the HTTP boundary for the curriculum preset.
func RegisterPresetRoutes(admin *gin.RouterGroup, pool *pgxpool.Pool) {
	service := NewPresetService(newPresetRepository(pool))
	admin.GET("/curriculum/preset/preview", func(c *gin.Context) {
		summary, err := service.Preview(c.Request.Context())
		if err != nil {
			apierror.RespondInternal(c, err)
			return
		}
		c.JSON(http.StatusOK, summary)
	})
	admin.POST("/curriculum/preset", func(c *gin.Context) {
		summary, err := service.Run(c.Request.Context())
		if err != nil {
			apierror.RespondInternal(c, err)
			return
		}
		c.JSON(http.StatusOK, summary)
	})
}
