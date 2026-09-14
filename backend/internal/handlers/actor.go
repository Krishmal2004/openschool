package handlers

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/openschool-org/openschool/internal/identity"
	"github.com/openschool-org/openschool/internal/middleware"
	"github.com/openschool-org/openschool/internal/services"
)

// actorFromContext supports the remaining legacy Society handler. Modules
// resolve their own caller types at their HTTP boundary.
func actorFromContext(c *gin.Context) (services.Actor, error) {
	id, err := middleware.UserIDFromContext(c)
	if err != nil {
		return services.Actor{}, fmt.Errorf("invalid caller identity")
	}
	var roles []string
	if value, ok := c.Get("roles"); ok {
		roles, _ = value.([]string)
	}
	return services.Actor{
		ID:       id,
		Email:    c.GetString("email"),
		FullName: strings.TrimSpace(c.GetString("given_name") + " " + c.GetString("family_name")),
		Role:     identity.ResolveAppRole(roles),
	}, nil
}
