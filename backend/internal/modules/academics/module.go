// Package academics owns curriculum, class structure, enrollment, marks, and promotion use cases.
package academics

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/openschool-org/openschool/internal/middleware"
)

// referenceDataMaxAgeSeconds is how long a client may cache rarely-changing
// reference data (subjects) before revalidating (section 5).
const referenceDataMaxAgeSeconds = 300

// RegisterSubjectRoutes mounts the subject vertical slice.
func RegisterSubjectRoutes(admin, teacherOrAdmin *gin.RouterGroup, pool *pgxpool.Pool) {
	handler := newSubjectHandler(newSubjectRepository(pool))
	admin.POST("/subjects", handler.create)
	teacherOrAdmin.GET("/subjects", middleware.CacheReference(referenceDataMaxAgeSeconds), handler.list)
	teacherOrAdmin.GET("/subjects/:id", handler.get)
	admin.PUT("/subjects/:id", handler.update)
	admin.DELETE("/subjects/:id", handler.delete)
}

// RegisterStreamRoutes mounts A/L streams and their subgroup endpoints.
func RegisterStreamRoutes(admin, teacherOrAdmin *gin.RouterGroup, pool *pgxpool.Pool) {
	handler := newStreamHandler(newStreamRepository(pool))
	admin.POST("/streams", handler.create)
	teacherOrAdmin.GET("/streams", handler.list)
	teacherOrAdmin.GET("/streams/:id", handler.get)
	admin.PUT("/streams/:id", handler.update)
	admin.DELETE("/streams/:id", handler.delete)
	admin.POST("/streams/:id/groups", handler.createGroup)
	teacherOrAdmin.GET("/streams/:id/groups", handler.listGroups)
	teacherOrAdmin.GET("/streams/:id/groups/:groupId", handler.getGroup)
	admin.PUT("/streams/:id/groups/:groupId", handler.updateGroup)
	admin.DELETE("/streams/:id/groups/:groupId", handler.deleteGroup)
}
