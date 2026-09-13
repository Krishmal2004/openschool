package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openschool-org/openschool/internal/models"
	"github.com/openschool-org/openschool/internal/services"
)

// JobsHandler exposes the admin Automation panel's endpoints for the background agent scheduler.
type JobsHandler struct {
	service *services.JobsService
}

// NewJobsHandler constructs a JobsHandler with its scheduler and settings dependencies.
func NewJobsHandler(service *services.JobsService) *JobsHandler {
	return &JobsHandler{service: service}
}

// List lists every registered background job and its current state.
func (h *JobsHandler) List(c *gin.Context) {
	ctx := c.Request.Context()
	status, err := h.service.List(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, status)
}

// SetEnabled enables or disables a background agent.
func (h *JobsHandler) SetEnabled(c *gin.Context) {
	name := c.Param("name")
	var req models.SetJobEnabledRequest
	if err := bindStrict(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.SetEnabled(c.Request.Context(), name, req.Enabled); err != nil {
		if errors.Is(err, services.ErrUnknownJob) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrSystemHealthCannotStop) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// RunNow runs a background agent immediately, outside its schedule, honoring its enabled/disabled setting.
func (h *JobsHandler) RunNow(c *gin.Context) {
	name := c.Param("name")
	result, err := h.service.RunNow(c.Request.Context(), name)
	if err != nil {
		if errors.Is(err, services.ErrUnknownJob) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "summary": result.Summary, "findings": result.Findings})
		return
	}
	c.JSON(http.StatusOK, gin.H{"summary": result.Summary, "findings": result.Findings})
}
