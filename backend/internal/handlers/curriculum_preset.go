package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openschool-org/openschool/internal/services"
)

// CurriculumPresetHandler exposes the HTTP endpoint for applying a curriculum preset.
type CurriculumPresetHandler struct {
	service *services.CurriculumPresetService
}

// NewCurriculumPresetHandler constructs a CurriculumPresetHandler with its service dependency.
func NewCurriculumPresetHandler(service *services.CurriculumPresetService) *CurriculumPresetHandler {
	return &CurriculumPresetHandler{service: service}
}

// Preview computes exactly what Run would create, without writing anything — for admin review before committing.
func (h *CurriculumPresetHandler) Preview(c *gin.Context) {
	summary, err := h.service.Preview(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// Run seeds the standard Sri Lankan Grade 1-13 curriculum for whichever grades exist in this school.
func (h *CurriculumPresetHandler) Run(c *gin.Context) {
	summary, err := h.service.Run(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}
