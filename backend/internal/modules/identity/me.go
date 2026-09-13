package identity

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	identitycore "github.com/openschool-org/openschool/internal/identity"
)

type userProvisioner interface {
	ensureExists(ctx context.Context, command ensureUserCommand) (provisionedUser, error)
}

type ensureUserCommand struct {
	ID       uuid.UUID
	Email    string
	FullName string
	Role     string
}

type provisionedUser struct {
	MustChangePassword bool
}

type meService struct{ users userProvisioner }

func newMeService(users userProvisioner) *meService { return &meService{users: users} }

func (s *meService) ensureProvisioned(ctx context.Context, command ensureUserCommand) (provisionedUser, error) {
	if command.Role == "" {
		return provisionedUser{}, nil
	}
	return s.users.ensureExists(ctx, command)
}

type meHandler struct{ service *meService }

func newMeHandler(service *meService) *meHandler { return &meHandler{service: service} }

func (h *meHandler) get(c *gin.Context) {
	userID := c.GetString("userID")
	email := c.GetString("email")
	givenName := c.GetString("given_name")
	familyName := c.GetString("family_name")
	tokenRoles, _ := c.Get("roles")
	roleList, _ := tokenRoles.([]string)

	mustChangePassword := false
	if parsedID, err := uuid.Parse(userID); err == nil {
		user, provisionErr := h.service.ensureProvisioned(c.Request.Context(), ensureUserCommand{
			ID: parsedID, Email: email, FullName: givenName + " " + familyName, Role: identitycore.ResolveAppRole(roleList),
		})
		if provisionErr != nil {
			log.Printf("/me: failed to provision local user %s: %v", parsedID, provisionErr)
		} else {
			mustChangePassword = user.MustChangePassword
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID, "email": email, "username": c.GetString("username"),
		"given_name": givenName, "family_name": familyName,
		"phone_number": c.GetString("phone_number"), "roles": tokenRoles,
		"must_change_password": mustChangePassword,
	})
}
