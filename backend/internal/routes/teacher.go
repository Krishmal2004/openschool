package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/openschool-org/openschool/internal/handlers"
	"github.com/openschool-org/openschool/internal/ports"
	"github.com/openschool-org/openschool/internal/repositories"
	"github.com/openschool-org/openschool/internal/services"
)

func RegisterTeacherRoutes(admin *gin.RouterGroup, teacherOrAdmin *gin.RouterGroup, houseAssignments ports.HouseAssignments, pool *pgxpool.Pool) {
	repo := repositories.NewTeacherRepository(pool)
	auditSvc := services.NewAuditService(repositories.NewAuditRepository(pool))
	service := services.NewTeacherService(repo, newIdentityProvider(), houseAssignments, auditSvc)
	handler := handlers.NewTeacherHandler(service)

	admin.POST("/teachers", handler.Create)
	teacherOrAdmin.GET("/teachers", handler.List)
	teacherOrAdmin.GET("/teachers/:id", handler.GetByID)
	admin.PUT("/teachers/:id", handler.Update)
	admin.PUT("/teachers/:id/house", handler.UpdateHouse)
	admin.PUT("/teachers/:id/employment-status", handler.UpdateEmploymentStatus)
	admin.DELETE("/teachers/:id", handler.Delete)
	admin.POST("/teachers/:id/subjects", handler.AssignSubject)
	admin.DELETE("/teachers/:id/subjects/:subject_id", handler.RemoveSubject)
	teacherOrAdmin.GET("/teachers/:id/subjects", handler.ListSubjects)
	teacherOrAdmin.GET("/teachers/:id/workload", handler.Workload)
	admin.GET("/subjects/:id/teachers", handler.ListBySubject)
}
