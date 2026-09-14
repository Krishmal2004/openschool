package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	academicsmodule "github.com/openschool-org/openschool/internal/modules/academics"
)

func newTermMarkRunner(pool *pgxpool.Pool) academicsmodule.TermMarkReader {
	return academicsmodule.NewTermMarkService(academicsmodule.NewTermMarkRepository(pool))
}

func RegisterTermMarkRoutes(teacherOrAdmin *gin.RouterGroup, pool *pgxpool.Pool) {
	academicsmodule.RegisterTermMarkRoutes(teacherOrAdmin, academicsmodule.NewTermMarkService(academicsmodule.NewTermMarkRepository(pool)))
}
