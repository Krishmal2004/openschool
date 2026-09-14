package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/middleware"
)

type requestActor struct{ ID uuid.UUID }

// actorFromContext supports the remaining legacy authentication and identity
// reconciliation handlers. Modules resolve their callers at their HTTP boundary.
func actorFromContext(c *gin.Context) (requestActor, error) {
	id, err := middleware.UserIDFromContext(c)
	if err != nil {
		return requestActor{}, fmt.Errorf("invalid caller identity")
	}
	return requestActor{ID: id}, nil
}
