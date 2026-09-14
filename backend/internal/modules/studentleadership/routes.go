package studentleadership

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/identity"
	"github.com/openschool-org/openschool/internal/middleware"
	"github.com/openschool-org/openschool/internal/platform/httpx"
)

func societyError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrSocietyNotFound), errors.Is(err, ErrSocietyMemberNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, ErrNotTeacherInCharge):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func routeActor(c *gin.Context) (Actor, error) {
	id, err := middleware.UserIDFromContext(c)
	if err != nil {
		return Actor{}, err
	}
	roles, _ := c.Get("roles")
	roleList, _ := roles.([]string)
	return Actor{ID: id, Role: identity.ResolveAppRole(roleList)}, nil
}

func RegisterRoutes(admin, teacherOrAdmin, studentAccess *gin.RouterGroup, service *Service) {
	admin.PUT("/prefects", func(c *gin.Context) {
		var req AssignPrefectRequest
		if err := httpx.BindStrict(c, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		prefect, err := service.AssignPrefect(c.Request.Context(), req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, prefect)
	})

	teacherOrAdmin.GET("/prefects", func(c *gin.Context) {
		yearID, err := uuid.Parse(c.Query("academic_year_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "a valid academic_year_id is required"})
			return
		}
		prefects, err := service.ListPrefectsByYear(c.Request.Context(), yearID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, prefects)
	})

	teacherOrAdmin.GET("/prefects/years", func(c *gin.Context) {
		years, err := service.ListPrefectYears(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, years)
	})

	studentAccess.GET("/students/:id/prefect-appointments", func(c *gin.Context) {
		studentID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student id"})
			return
		}
		appointments, err := service.ListPrefectsByStudent(c.Request.Context(), studentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, appointments)
	})

	admin.DELETE("/prefects/:id", func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		if err := service.DeletePrefect(c.Request.Context(), id); err != nil {
			if errors.Is(err, ErrPrefectNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "prefect appointment removed"})
	})

	admin.POST("/societies", func(c *gin.Context) {
		var req CreateSocietyRequest
		if err := httpx.BindStrict(c, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		society, err := service.CreateSociety(c.Request.Context(), req)
		if err != nil {
			societyError(c, err)
			return
		}
		c.JSON(http.StatusCreated, society)
	})

	admin.PUT("/societies/:id", func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid society id"})
			return
		}
		var req UpdateSocietyRequest
		if err := httpx.BindStrict(c, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		society, err := service.UpdateSociety(c.Request.Context(), id, req)
		if err != nil {
			societyError(c, err)
			return
		}
		c.JSON(http.StatusOK, society)
	})

	admin.DELETE("/societies/:id", func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid society id"})
			return
		}
		if err := service.DeleteSociety(c.Request.Context(), id); err != nil {
			societyError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "society deleted"})
	})

	teacherOrAdmin.GET("/societies", func(c *gin.Context) {
		yearID, err := uuid.Parse(c.Query("academic_year_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "a valid academic_year_id is required"})
			return
		}
		societies, err := service.ListSocietiesByYear(c.Request.Context(), yearID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, societies)
	})

	teacherOrAdmin.GET("/societies/years", func(c *gin.Context) {
		years, err := service.ListSocietyYears(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, years)
	})

	teacherOrAdmin.GET("/societies/:id/members", func(c *gin.Context) {
		societyID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid society id"})
			return
		}
		members, err := service.ListSocietyMembers(c.Request.Context(), societyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, members)
	})

	teacherOrAdmin.PUT("/societies/:id/members", func(c *gin.Context) {
		societyID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid society id"})
			return
		}
		var req AssignSocietyMemberRequest
		if err := httpx.BindStrict(c, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		actor, err := routeActor(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid caller identity"})
			return
		}
		member, err := service.AssignSocietyMember(c.Request.Context(), actor, societyID, req)
		if err != nil {
			societyError(c, err)
			return
		}
		c.JSON(http.StatusOK, member)
	})

	teacherOrAdmin.DELETE("/societies/:id/members/:memberId", func(c *gin.Context) {
		societyID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid society id"})
			return
		}
		memberID, err := uuid.Parse(c.Param("memberId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid member id"})
			return
		}
		actor, err := routeActor(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid caller identity"})
			return
		}
		if err := service.RemoveSocietyMember(c.Request.Context(), actor, societyID, memberID); err != nil {
			societyError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "society member removed"})
	})

	studentAccess.GET("/students/:id/society-memberships", func(c *gin.Context) {
		studentID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student id"})
			return
		}
		memberships, err := service.ListSocietyMembershipsByStudent(c.Request.Context(), studentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, memberships)
	})
}
