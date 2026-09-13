package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/models"
	"github.com/openschool-org/openschool/internal/services"
)

// CurriculumHandler exposes HTTP endpoints for mediums, levels, and subject-selection groups.
type CurriculumHandler struct {
	service *services.CurriculumService
}

// NewCurriculumHandler constructs a CurriculumHandler with its service dependency.
func NewCurriculumHandler(service *services.CurriculumService) *CurriculumHandler {
	return &CurriculumHandler{service: service}
}

// toMediumResponse converts a generated medium row into its JSON response shape.
// toSelectionGroupResponse converts a generated selection-group row into its JSON response shape.
// ── mediums ─────────────────────────────────────────────────────────────────

// CreateMedium creates a school-defined medium of instruction.
func (h *CurriculumHandler) CreateMedium(c *gin.Context) {
	var req models.CreateMediumRequest
	if err := bindStrict(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	medium, err := h.service.CreateMedium(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, medium)
}

// ListMediums lists mediums.
func (h *CurriculumHandler) ListMediums(c *gin.Context) {
	mediums, err := h.service.ListMediums(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mediums)
}

// UpdateMedium updates a medium of instruction.
func (h *CurriculumHandler) UpdateMedium(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req models.UpdateMediumRequest
	if err := bindStrict(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	medium, err := h.service.UpdateMedium(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, medium)
}

// DeleteMedium is blocked while the medium is referenced by a group subject or enrollment.
func (h *CurriculumHandler) DeleteMedium(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DeleteMedium(c.Request.Context(), id); err != nil {
		if errors.Is(err, services.ErrMediumNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "medium deleted"})
}

// ── levels ──────────────────────────────────────────────────────────────────

// CreateLevel creates an admin-defined curriculum level.
func (h *CurriculumHandler) CreateLevel(c *gin.Context) {
	var req models.CreateLevelRequest
	if err := bindStrict(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	level, err := h.service.CreateLevel(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, level)
}

// GetLevel returns a single curriculum level by ID.
func (h *CurriculumHandler) GetLevel(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	level, err := h.service.GetLevel(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "level not found"})
		return
	}

	c.JSON(http.StatusOK, level)
}

// ListLevels returns curriculum levels, optionally filtered by grade via ?grade_id=.
func (h *CurriculumHandler) ListLevels(c *gin.Context) {
	var (
		levels []models.LevelResponse
		err    error
	)

	if gradeParam := c.Query("grade_id"); gradeParam != "" {
		gradeID, parseErr := uuid.Parse(gradeParam)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid grade_id"})
			return
		}
		levels, err = h.service.ListLevelsByGrade(c.Request.Context(), gradeID)
	} else {
		levels, err = h.service.ListLevels(c.Request.Context())
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, levels)
}

// UpdateLevel updates a curriculum level.
func (h *CurriculumHandler) UpdateLevel(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req models.UpdateLevelRequest
	if err := bindStrict(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	level, err := h.service.UpdateLevel(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, level)
}

// DuplicateLevel copies a level with all its selection groups and their subjects, under a new label.
func (h *CurriculumHandler) DuplicateLevel(c *gin.Context) {
	sourceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req models.DuplicateLevelRequest
	if err := bindStrict(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	level, err := h.service.DuplicateLevel(c.Request.Context(), sourceID, req)
	if err != nil {
		if errors.Is(err, services.ErrLevelNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, level)
}

// DeleteLevel is blocked while any of the level's groups still carry enrollments.
func (h *CurriculumHandler) DeleteLevel(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DeleteLevel(c.Request.Context(), id); err != nil {
		if errors.Is(err, services.ErrLevelNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "level deleted"})
}

// GetCurriculumTree returns a level with its selection groups and each group's subjects nested.
func (h *CurriculumHandler) GetCurriculumTree(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	tree, err := h.service.GetCurriculumTree(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrLevelNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tree)
}

// ── selection groups ────────────────────────────────────────────────────────

// CreateSelectionGroup creates a selection group — a pool the student picks between min_select and max_select subjects from.
func (h *CurriculumHandler) CreateSelectionGroup(c *gin.Context) {
	levelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid level id"})
		return
	}

	var req models.CreateSelectionGroupRequest
	if err := bindStrict(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	group, err := h.service.CreateSelectionGroup(c.Request.Context(), levelID, req)
	if err != nil {
		if errors.Is(err, services.ErrLevelNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, group)
}

// ListSelectionGroups lists selection groups of a level.
func (h *CurriculumHandler) ListSelectionGroups(c *gin.Context) {
	levelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid level id"})
		return
	}

	groups, err := h.service.ListSelectionGroupsByLevel(c.Request.Context(), levelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, groups)
}

// UpdateSelectionGroup updates a subject-selection group.
func (h *CurriculumHandler) UpdateSelectionGroup(c *gin.Context) {
	id, err := uuid.Parse(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	var req models.UpdateSelectionGroupRequest
	if err := bindStrict(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	group, err := h.service.UpdateSelectionGroup(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, group)
}

// DeleteSelectionGroup is blocked while enrollments reference the group.
func (h *CurriculumHandler) DeleteSelectionGroup(c *gin.Context) {
	id, err := uuid.Parse(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	if err := h.service.DeleteSelectionGroup(c.Request.Context(), id); err != nil {
		if errors.Is(err, services.ErrSelectionGroupNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "selection group deleted"})
}

// ── group subjects ──────────────────────────────────────────────────────────

// AddGroupSubject upserts the subject's medium restriction and prerequisite note.
func (h *CurriculumHandler) AddGroupSubject(c *gin.Context) {
	groupID, err := uuid.Parse(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	var req models.AddGroupSubjectRequest
	if err := bindStrict(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AddGroupSubject(c.Request.Context(), groupID, req); err != nil {
		if errors.Is(err, services.ErrSelectionGroupNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "subject added to group"})
}

// ListGroupSubjects lists a group's subjects.
func (h *CurriculumHandler) ListGroupSubjects(c *gin.Context) {
	groupID, err := uuid.Parse(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	subjects, err := h.service.ListGroupSubjects(c.Request.Context(), groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subjects)
}

// RemoveGroupSubject removes a subject from a selection group.
func (h *CurriculumHandler) RemoveGroupSubject(c *gin.Context) {
	groupID, err := uuid.Parse(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	subjectID, err := uuid.Parse(c.Param("subject_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subject id"})
		return
	}

	if err := h.service.RemoveGroupSubject(c.Request.Context(), groupID, subjectID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "subject removed from group"})
}
