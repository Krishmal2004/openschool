package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openschool-org/openschool/internal/services"
)

// SearchHandler exposes the admin global-search endpoint.
type SearchHandler struct {
	service *services.SearchService
}

// NewSearchHandler constructs a SearchHandler with its service dependency.
func NewSearchHandler(service *services.SearchService) *SearchHandler {
	return &SearchHandler{service: service}
}

// Global searches students, teachers, guardians, and non-academic staff for the admin header's jump-to-record search.
func (h *SearchHandler) Global(c *gin.Context) {
	result, err := h.service.Global(c.Request.Context(), c.Query("q"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
