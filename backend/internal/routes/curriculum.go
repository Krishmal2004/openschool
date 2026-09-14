package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	curriculummodule "github.com/openschool-org/openschool/internal/modules/curriculum"
)

func RegisterCurriculumRoutes(admin *gin.RouterGroup, protected *gin.RouterGroup, pool *pgxpool.Pool) {
	curriculummodule.RegisterPresetRoutes(admin, pool)
}
