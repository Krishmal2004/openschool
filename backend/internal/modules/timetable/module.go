// Package timetable owns scheduling configuration, generation, review, and publication.
package timetable

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterSettingsRoutes(admin *gin.RouterGroup, pool *pgxpool.Pool) {
	handler := newSettingsHandler(newSettingsRepository(pool))
	admin.PUT("/timetable-settings", handler.upsert)
	admin.GET("/timetable-settings/:academic_year_id", handler.getByYear)
}

func RegisterClassroomRoutes(admin, teacherOrAdmin *gin.RouterGroup, pool *pgxpool.Pool) {
	handler := newClassroomHandler(newClassroomRepository(pool))
	admin.POST("/classrooms", handler.create)
	teacherOrAdmin.GET("/classrooms", handler.list)
	admin.PUT("/classrooms/:id", handler.update)
	admin.DELETE("/classrooms/:id", handler.delete)
}

func RegisterSubjectPeriodRequirementRoutes(admin, teacherOrAdmin *gin.RouterGroup, pool *pgxpool.Pool) {
	handler := newSubjectPeriodRequirementHandler(newSubjectPeriodRequirementRepository(pool))
	admin.PUT("/subject-period-requirements", handler.upsert)
	teacherOrAdmin.GET("/subject-period-requirements", handler.listByGrade)
	admin.DELETE("/subject-period-requirements/:id", handler.delete)
}

func RegisterTeacherAvailabilityRoutes(admin, teacherOrAdmin *gin.RouterGroup, pool *pgxpool.Pool) {
	handler := newTeacherAvailabilityHandler(newTeacherAvailabilityRepository(pool))
	admin.POST("/teachers/:id/availability", handler.create)
	teacherOrAdmin.GET("/teachers/:id/availability", handler.listByTeacherYear)
	admin.DELETE("/teachers/:id/availability/:availability_id", handler.delete)
}

func RegisterGradeSectionRoutes(admin, teacherOrAdmin *gin.RouterGroup, pool *pgxpool.Pool) {
	handler := newGradeSectionHandler(newGradeSectionRepository(pool))
	admin.POST("/grade-sections", handler.create)
	teacherOrAdmin.GET("/grade-sections", handler.listByYear)
	teacherOrAdmin.GET("/grade-sections/:id", handler.get)
	admin.PUT("/grade-sections/:id", handler.update)
	admin.DELETE("/grade-sections/:id", handler.delete)
	admin.PUT("/grade-sections/:id/grades", handler.assignGrades)
	admin.DELETE("/grade-sections/:id/grades/:grade_id", handler.removeGrade)
	teacherOrAdmin.GET("/grade-sections/:id/periods", handler.getPeriods)
	admin.PUT("/grade-sections/:id/periods", handler.savePeriods)
	admin.POST("/grade-sections/:id/periods/generate", handler.regeneratePeriods)
}
