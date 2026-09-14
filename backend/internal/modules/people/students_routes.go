package people

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/middleware"
	"github.com/openschool-org/openschool/internal/models"
	"github.com/openschool-org/openschool/internal/platform/httpx"
)

type StudentRunner interface {
	Create(context.Context, models.CreateStudentRequest, uuid.UUID) (any, error)
	Update(context.Context, uuid.UUID, models.UpdateStudentRequest) (any, error)
	UpdateHouse(context.Context, uuid.UUID, models.UpdateStudentHouseRequest, uuid.UUID) (any, error)
	Delete(context.Context, uuid.UUID, uuid.UUID) error
}
type StudentStatusWriter interface {
	UpdateStatus(context.Context, uuid.UUID, string) (any, error)
}

type StudentReader interface {
	Get(context.Context, uuid.UUID) (any, error)
	GetWithClass(context.Context, uuid.UUID) (any, error)
	List(context.Context) (any, error)
	ListByClass(context.Context, uuid.UUID) (any, error)
}

func RegisterStudentRoutes(admin, teacherOrAdmin *gin.RouterGroup, runner StudentRunner, reader StudentReader, statusWriter StudentStatusWriter) {
	admin.POST("/students", func(c *gin.Context) {
		var r models.CreateStudentRequest
		if e := httpx.BindStrict(c, &r); e != nil {
			c.JSON(400, gin.H{"error": e.Error()})
			return
		}
		a, ok := studentActor(c)
		if !ok {
			return
		}
		v, e := runner.Create(c, r, a)
		if e != nil {
			c.JSON(400, gin.H{"error": e.Error()})
			return
		}
		c.JSON(201, v)
	})
	teacherOrAdmin.GET("/students", func(c *gin.Context) {
		v, e := reader.List(c)
		if e != nil {
			c.JSON(500, gin.H{"error": e.Error()})
			return
		}
		c.JSON(200, v)
	})
	teacherOrAdmin.GET("/students/:id", func(c *gin.Context) {
		id, ok := studentID(c)
		if !ok {
			return
		}
		v, e := reader.Get(c, id)
		if e != nil {
			c.JSON(404, gin.H{"error": "student not found"})
			return
		}
		c.JSON(200, v)
	})
	teacherOrAdmin.GET("/students/:id/class", func(c *gin.Context) {
		id, ok := studentID(c)
		if !ok {
			return
		}
		v, e := reader.GetWithClass(c, id)
		if e != nil {
			c.JSON(404, gin.H{"error": "student not found"})
			return
		}
		c.JSON(200, v)
	})
	admin.PUT("/students/:id", func(c *gin.Context) {
		id, ok := studentID(c)
		if !ok {
			return
		}
		var r models.UpdateStudentRequest
		if e := httpx.BindStrict(c, &r); e != nil {
			c.JSON(400, gin.H{"error": e.Error()})
			return
		}
		v, e := runner.Update(c, id, r)
		if e != nil {
			c.JSON(400, gin.H{"error": e.Error()})
			return
		}
		c.JSON(200, v)
	})
	admin.PUT("/students/:id/house", func(c *gin.Context) {
		id, ok := studentID(c)
		if !ok {
			return
		}
		var r models.UpdateStudentHouseRequest
		if e := httpx.BindStrict(c, &r); e != nil {
			c.JSON(400, gin.H{"error": e.Error()})
			return
		}
		a, ok := studentActor(c)
		if !ok {
			return
		}
		v, e := runner.UpdateHouse(c, id, r, a)
		if e != nil {
			c.JSON(400, gin.H{"error": e.Error()})
			return
		}
		c.JSON(200, v)
	})
	admin.PUT("/students/:id/enrollment-status", func(c *gin.Context) {
		id, ok := studentID(c)
		if !ok {
			return
		}
		var r models.UpdateStudentEnrollmentStatusRequest
		if e := httpx.BindStrict(c, &r); e != nil {
			c.JSON(400, gin.H{"error": e.Error()})
			return
		}
		v, e := statusWriter.UpdateStatus(c, id, r.Status)
		if e != nil {
			c.JSON(400, gin.H{"error": e.Error()})
			return
		}
		c.JSON(200, v)
	})
	admin.DELETE("/students/:id", func(c *gin.Context) {
		id, ok := studentID(c)
		if !ok {
			return
		}
		a, ok := studentActor(c)
		if !ok {
			return
		}
		if e := runner.Delete(c, id, a); e != nil {
			c.JSON(400, gin.H{"error": e.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "student deleted"})
	})
	teacherOrAdmin.GET("/classes/:id/students", func(c *gin.Context) {
		id, ok := studentID(c)
		if !ok {
			return
		}
		v, e := reader.ListByClass(c, id)
		if e != nil {
			c.JSON(500, gin.H{"error": e.Error()})
			return
		}
		c.JSON(200, v)
	})
}
func studentID(c *gin.Context) (uuid.UUID, bool) {
	id, e := uuid.Parse(c.Param("id"))
	if e != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return uuid.Nil, false
	}
	return id, true
}
func studentActor(c *gin.Context) (uuid.UUID, bool) {
	id, e := middleware.UserIDFromContext(c)
	if e != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("invalid caller identity")})
		return uuid.Nil, false
	}
	return id, true
}
