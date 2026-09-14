// This file defines the RegisterStudentPortfolioRoutes function, mapping student portfolio endpoints (progress reports, activities, leadership, awards, disciplinary records) to their corresponding handlers and middlewares.

package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	peoplemodule "github.com/openschool-org/openschool/internal/modules/people"
)

func RegisterStudentPortfolioRoutes(teacherOrAdmin *gin.RouterGroup, studentAccess *gin.RouterGroup, pool *pgxpool.Pool) {
	service := peoplemodule.NewStudentPortfolioService(peoplemodule.NewStudentPortfolioStore(pool))
	peoplemodule.RegisterStudentPortfolioRoutes(teacherOrAdmin, studentAccess, service)
}
