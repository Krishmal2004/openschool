// Package identity owns authenticated account identity and password lifecycle use cases.
package identity

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Register mounts identity-owned endpoints and constructs their private dependencies.
func Register(protected *gin.RouterGroup, pool *pgxpool.Pool) {
	repository := newUserRepository(pool)
	handler := newMeHandler(newMeService(repository))
	protected.GET("/me", handler.get)
}
