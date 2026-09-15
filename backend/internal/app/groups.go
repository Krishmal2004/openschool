package app

import (
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/openschool-org/openschool/internal/authz"
	"github.com/openschool-org/openschool/internal/middleware"
	"github.com/openschool-org/openschool/internal/ports"
)

// HTTPGroups contains the authorization-scoped route groups shared by modules.
type HTTPGroups struct {
	API            *gin.RouterGroup
	Protected      *gin.RouterGroup
	Admin          *gin.RouterGroup
	TeacherOrAdmin *gin.RouterGroup
	Parent         *gin.RouterGroup
	Student        *gin.RouterGroup
	Teacher        *gin.RouterGroup
	StudentAccess  *gin.RouterGroup
}

func newHTTPGroups(router *gin.Engine, studentAccess ports.StudentAccessAuthorizer) HTTPGroups {
	api := router.Group("/api/v1")
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())
	protected.Use(middleware.PerAccountRateLimit(
		envFloatOr("API_PER_ACCOUNT_RATE_LIMIT_RPS", 30),
		envIntOr("API_PER_ACCOUNT_RATE_LIMIT_BURST", 60),
	))

	admin := protected.Group("")
	admin.Use(middleware.RequireRole(authz.RoleAdmin))
	teacherOrAdmin := protected.Group("")
	teacherOrAdmin.Use(middleware.RequireRole(authz.RoleAdmin, authz.RoleTeacher))
	parent := protected.Group("")
	parent.Use(middleware.RequireRole(authz.RoleParent))
	student := protected.Group("")
	student.Use(middleware.RequireRole(authz.RoleStudent))
	teacher := protected.Group("")
	teacher.Use(middleware.RequireRole(authz.RoleTeacher))
	studentScoped := protected.Group("")
	studentScoped.Use(middleware.RequireStudentAccess(studentAccess))

	return HTTPGroups{
		API: api, Protected: protected, Admin: admin, TeacherOrAdmin: teacherOrAdmin,
		Parent: parent, Student: student, Teacher: teacher, StudentAccess: studentScoped,
	}
}

func envFloatOr(key string, fallback float64) float64 {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseFloat(value, 64); err == nil {
			return parsed
		}
	}
	return fallback
}

func envIntOr(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return fallback
}
