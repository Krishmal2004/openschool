package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/models"
	"github.com/openschool-org/openschool/internal/services"
)

// PositionHandler exposes HTTP endpoints for the in-app teacher position hierarchy.
type PositionHandler struct {
	service *services.PositionService
}

// NewPositionHandler constructs a PositionHandler with its service dependency.
func NewPositionHandler(service *services.PositionService) *PositionHandler {
	return &PositionHandler{service: service}
}

// AssignPrincipal makes a permanent appointment, replacing any existing Principal.
func (h *PositionHandler) AssignPrincipal(c *gin.Context) {
	var req models.AssignPrincipalRequest
	if err := bindStrict(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	actor, err := actorFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	position, err := h.service.AssignPrincipal(c.Request.Context(), req, actor.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, position)
}

// AssignVicePrincipal makes a permanent Vice Principal appointment with a whole-school or grade-scoped notification grant.
func (h *PositionHandler) AssignVicePrincipal(c *gin.Context) {
	var req models.AssignVicePrincipalRequest
	if err := bindStrict(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	actor, err := actorFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	position, err := h.service.AssignVicePrincipal(c.Request.Context(), req, actor.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, position)
}

// List lists leadership positions (Principal, Vice Principals).
func (h *PositionHandler) List(c *gin.Context) {
	list, err := h.service.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}

// Delete removes a position assignment.
func (h *PositionHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	actor, err := actorFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id, actor.ID); err != nil {
		if errors.Is(err, services.ErrPositionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "position removed"})
}
