package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openschool-org/openschool/internal/services"
)

// DashboardHandler exposes the admin dashboard's summary-statistics endpoint.
type DashboardHandler struct {
	service *services.DashboardService
}

// NewDashboardHandler constructs a DashboardHandler with its service dependency.
func NewDashboardHandler(service *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{service: service}
}

// Analytics returns composed student/staff/academic/school analytics for the current academic year and term.
func (h *DashboardHandler) Analytics(c *gin.Context) {
	analytics, err := h.service.Analytics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analytics)
}
