package leadership

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/apierror"
	"github.com/openschool-org/openschool/internal/middleware"
	"github.com/openschool-org/openschool/internal/platform/httpx"
)

func RegisterRoutes(admin, teacherOrAdmin *gin.RouterGroup, service *Service) {
	admin.PUT("/positions/principal", func(c *gin.Context) {
		var req AssignPrincipalRequest
		if err := httpx.BindStrict(c, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		actorID, err := middleware.UserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid caller identity"})
			return
		}
		position, err := service.AssignPrincipal(c.Request.Context(), req, actorID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, position)
	})

	admin.PUT("/positions/vice-principal", func(c *gin.Context) {
		var req AssignVicePrincipalRequest
		if err := httpx.BindStrict(c, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		actorID, err := middleware.UserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid caller identity"})
			return
		}
		position, err := service.AssignVicePrincipal(c.Request.Context(), req, actorID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, position)
	})

	teacherOrAdmin.GET("/positions", func(c *gin.Context) {
		positions, err := service.ListPositions(c.Request.Context())
		if err != nil {
			apierror.RespondInternal(c, err)
			return
		}
		c.JSON(http.StatusOK, positions)
	})

	admin.DELETE("/positions/:id", func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		actorID, err := middleware.UserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid caller identity"})
			return
		}
		if err := service.DeletePosition(c.Request.Context(), id, actorID); err != nil {
			if errors.Is(err, ErrPositionNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			} else {
				apierror.RespondInternal(c, err)
			}
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "position removed"})
	})

	admin.PUT("/section-heads", func(c *gin.Context) {
		var req AssignSectionHeadRequest
		if err := httpx.BindStrict(c, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		sectionHead, err := service.AssignSectionHead(c.Request.Context(), req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, sectionHead)
	})

	teacherOrAdmin.GET("/section-heads", func(c *gin.Context) {
		yearID, err := uuid.Parse(c.Query("academic_year_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "a valid academic_year_id is required"})
			return
		}
		sectionHeads, err := service.ListSectionHeads(c.Request.Context(), yearID)
		if err != nil {
			apierror.RespondInternal(c, err)
			return
		}
		c.JSON(http.StatusOK, sectionHeads)
	})

	admin.DELETE("/section-heads/:id", func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		if err := service.DeleteSectionHead(c.Request.Context(), id); err != nil {
			if errors.Is(err, ErrSectionHeadNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			} else {
				apierror.RespondInternal(c, err)
			}
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "section head removed"})
	})
}
