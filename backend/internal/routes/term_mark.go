package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	academicsmodule "github.com/openschool-org/openschool/internal/modules/academics"
)

func RegisterTermMarkRoutes(teacherOrAdmin *gin.RouterGroup, pool *pgxpool.Pool) {
	academicsmodule.RegisterTermMarkRoutes(teacherOrAdmin, academicsmodule.NewTermMarkService(academicsmodule.NewTermMarkRepository(pool)))
}
