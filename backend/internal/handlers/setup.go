package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/openschool-org/openschool/internal/identity"
	"github.com/openschool-org/openschool/internal/models"
	"github.com/openschool-org/openschool/internal/services"
)

// SetupHandler exposes the one-time initial-admin setup endpoint.
type SetupHandler struct {
	service *services.SetupService
}

// NewSetupHandler constructs a SetupHandler with its service dependency.
func NewSetupHandler(service *services.SetupService) *SetupHandler {
	return &SetupHandler{service: service}
}

// Status reports whether an admin account still needs to be registered.
func (h *SetupHandler) Status(c *gin.Context) {
	needsSetup, err := h.service.NeedsSetup(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check setup status"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"needs_setup": needsSetup})
}

// RegisterAdmin creates the first admin account; only succeeds while none exists yet.
func (h *SetupHandler) RegisterAdmin(c *gin.Context) {
	var req models.RegisterAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": friendlyBindError(err)})
		return
	}

	user, err := h.service.RegisterFirstAdmin(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrSetupAlreadyDone):
			c.JSON(http.StatusForbidden, gin.H{"error": "An admin account already exists. Setup can only be run once."})
		case errors.Is(err, identity.ErrDuplicateUser):
			c.JSON(http.StatusConflict, gin.H{"error": "An account with that email, username, or phone number already exists."})
		default:
			// Don't leak internal/upstream error details to an unauthenticated caller.
			log.Printf("setup: RegisterFirstAdmin failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong while creating the admin account. Please try again."})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": user.ID.String(), "email": user.Email})
}

// friendlyBindError turns a validation failure into a user-facing message using the failing field's name.
func friendlyBindError(err error) string {
	var verrs validator.ValidationErrors
	if !errors.As(err, &verrs) || len(verrs) == 0 {
		return "Please check the form and try again."
	}

	switch verrs[0].Field() {
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
