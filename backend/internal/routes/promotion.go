package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	academicsmodule "github.com/openschool-org/openschool/internal/modules/academics"
)

func RegisterPromotionRoutes(admin *gin.RouterGroup, pool *pgxpool.Pool) {
	academicsmodule.RegisterPromotionRoutes(admin, academicsmodule.NewPromotionService(academicsmodule.NewPromotionRepository(pool)))
}
