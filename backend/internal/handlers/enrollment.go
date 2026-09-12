package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/models"
	"github.com/openschool-org/openschool/internal/services"
)

// EnrollmentHandler exposes HTTP endpoints for class enrollment and subject selection.
type EnrollmentHandler struct {
	service *services.EnrollmentService
}

// NewEnrollmentHandler constructs an EnrollmentHandler with its service dependency.
func NewEnrollmentHandler(service *services.EnrollmentService) *EnrollmentHandler {
	return &EnrollmentHandler{service: service}
}

// requireAcademicYear reads the mandatory ?academic_year_id= query parameter.
func requireAcademicYear(c *gin.Context) (uuid.UUID, bool) {
	raw := c.Query("academic_year_id")
	if raw == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "academic_year_id is required"})
		return uuid.UUID{}, false
	}

	id, err := uuid.Parse(raw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid academic_year_id"})
		return uuid.UUID{}, false
	}

	return id, true
}

// Validate dry-runs a student's picks against every group of the level without saving.
func (h *EnrollmentHandler) Validate(c *gin.Context) {
	var req models.SubmitEnrollmentRequest
	if err := bindStrict(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	levelID, err := uuid.Parse(req.LevelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid level_id"})
		return
	}

	validationErrs, err := h.service.Validate(c.Request.Context(), levelID, req.Picks)
	if err != nil {
		if errors.Is(err, services.ErrLevelHasNoGroups) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.EnrollmentValidationResponse{
		Valid:  len(validationErrs) == 0,
		Errors: validationErrs,
	})
}

// Submit validates then replaces the student's picks for that level and academic year.
func (h *EnrollmentHandler) Submit(c *gin.Context) {
	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student id"})
		return
	}

	var req models.SubmitEnrollmentRequest
	if err := bindStrict(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validationErrs, err := h.service.Submit(c.Request.Context(), studentID, req)
	if err != nil {
		if errors.Is(err, services.ErrEnrollmentInvalid) {
			c.JSON(http.StatusUnprocessableEntity, models.EnrollmentValidationResponse{
				Valid:  false,
				Errors: validationErrs,
			})
			return
		}
		if errors.Is(err, services.ErrLevelHasNoGroups) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrEnrollmentLocked) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.EnrollmentValidationResponse{Valid: true, Errors: nil})
}

// ListByStudent lists a student's subjects.
func (h *EnrollmentHandler) ListByStudent(c *gin.Context) {
	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student id"})
		return
	}

	academicYearID, ok := requireAcademicYear(c)
	if !ok {
		return
	}

	enrollments, err := h.service.ListByStudent(c.Request.Context(), studentID, academicYearID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, enrollments)
}

// Delete removes one enrollment pick.
func (h *EnrollmentHandler) Delete(c *gin.Context) {
	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student id"})
		return
	}

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

	academicYearID, ok := requireAcademicYear(c)
	if !ok {
		return
	}

	if err := h.service.Delete(c.Request.Context(), studentID, academicYearID, groupID, subjectID); err != nil {
		if errors.Is(err, services.ErrEnrollmentLocked) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "enrollment removed"})
}

// Unlock removes the enrollment lock so the student's picks for this level and year can be changed again (admin-only).
func (h *EnrollmentHandler) Unlock(c *gin.Context) {
	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student id"})
		return
	}

	levelID, err := uuid.Parse(c.Param("level_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid level id"})
		return
	}

	academicYearID, ok := requireAcademicYear(c)
	if !ok {
		return
	}

	if err := h.service.Unlock(c.Request.Context(), studentID, levelID, academicYearID); err != nil {
		if errors.Is(err, services.ErrEnrollmentNotLocked) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "subject selection unlocked"})
}

// ListStudentsBySubject lists students taking a subject.
func (h *EnrollmentHandler) ListStudentsBySubject(c *gin.Context) {
	subjectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subject id"})
		return
	}

	academicYearID, ok := requireAcademicYear(c)
	if !ok {
		return
	}

	students, err := h.service.ListStudentsBySubject(c.Request.Context(), subjectID, academicYearID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, students)
}

// ListStudentsByGroup lists students enrolled through a group.
func (h *EnrollmentHandler) ListStudentsByGroup(c *gin.Context) {
	groupID, err := uuid.Parse(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	academicYearID, ok := requireAcademicYear(c)
	if !ok {
		return
	}

	students, err := h.service.ListStudentsByGroup(c.Request.Context(), groupID, academicYearID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, students)
}
