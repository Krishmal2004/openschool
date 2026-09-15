package setup

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/openschool-org/openschool/internal/idp"
	"github.com/openschool-org/openschool/internal/middleware"
)

func RegisterRoutes(public *gin.RouterGroup, service *Service) {
	public.GET("/setup/status", func(c *gin.Context) {
		needsSetup, err := service.NeedsSetup(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check setup status"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"needs_setup": needsSetup})
	})

	public.POST("/setup/admin", middleware.RateLimit(1, 3), func(c *gin.Context) {
		var req RegisterAdminRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": friendlyBindError(err)})
			return
		}
		user, err := service.RegisterFirstAdmin(c.Request.Context(), req)
		if err != nil {
			switch {
			case errors.Is(err, ErrAlreadyDone):
				c.JSON(http.StatusForbidden, gin.H{"error": "An admin account already exists. Setup can only be run once."})
			case errors.Is(err, idp.ErrDuplicateUser):
				c.JSON(http.StatusConflict, gin.H{"error": "An account with that email, username, or phone number already exists."})
			default:
				log.Printf("setup: RegisterFirstAdmin failed: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong while creating the admin account. Please try again."})
			}
			return
		}
		c.JSON(http.StatusCreated, gin.H{"id": user.ID.String(), "email": user.Email})
	})
}

func friendlyBindError(err error) string {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) || len(validationErrors) == 0 {
		return "Please check the form and try again."
	}
	switch validationErrors[0].Field() {
	case "Email":
		return "Please enter a valid email address."
	case "Password":
		return "Password must be at least 8 characters."
	case "Username", "GivenName", "FamilyName":
		return "Please fill in all required fields."
	default:
		return "Please check the form and try again."
	}
}
