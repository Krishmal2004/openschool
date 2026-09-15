// Package selfservice owns authenticated student, parent, and teacher portal endpoints.
package selfservice

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/middleware"
	academicsmodule "github.com/openschool-org/openschool/internal/modules/academics"
	attendancemodule "github.com/openschool-org/openschool/internal/modules/attendance"
)

// StudentSelfHandler serves a signed-in student's own profile/attendance/marks endpoints, resolved from their token.
type StudentSelfHandler struct {
	studentSelf *StudentProfiles
	attendance  attendancemodule.Reader
	marks       academicsmodule.TermMarkReader
	enrollments academicsmodule.StudentEnrollment
}

// NewStudentSelfHandler constructs a StudentSelfHandler with its service dependencies.
func NewStudentSelfHandler(studentSelf *StudentProfiles, attendance attendancemodule.Reader, marks academicsmodule.TermMarkReader, enrollments academicsmodule.StudentEnrollment) *StudentSelfHandler {
	return &StudentSelfHandler{studentSelf: studentSelf, attendance: attendance, marks: marks, enrollments: enrollments}
}

// resolveStudentID resolves the signed-in student's own profile ID from the request's JWT.
func (h *StudentSelfHandler) resolveStudentID(c *gin.Context) (uuid.UUID, bool) {
	callerID, err := middleware.UserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid caller identity"})
		return uuid.UUID{}, false
	}
	studentID, err := h.studentSelf.Resolve(c.Request.Context(), callerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no student profile linked to this account"})
		return uuid.UUID{}, false
	}
	return studentID, true
}

// Profile returns the signed-in student's own profile.
func (h *StudentSelfHandler) Profile(c *gin.Context) {
	studentID, ok := h.resolveStudentID(c)
	if !ok {
		return
	}

	profile, err := h.studentSelf.Profile(c.Request.Context(), studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// Attendance returns the signed-in student's own attendance history.
func (h *StudentSelfHandler) Attendance(c *gin.Context) {
	studentID, ok := h.resolveStudentID(c)
	if !ok {
		return
	}

	records, err := h.attendance.ListByStudent(c.Request.Context(), studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, records)
}

// Marks returns the signed-in student's own term marks.
func (h *StudentSelfHandler) Marks(c *gin.Context) {
	studentID, ok := h.resolveStudentID(c)
	if !ok {
		return
	}
	termID, err := uuid.Parse(c.Query("term_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or missing term_id"})
		return
	}

	marks, err := h.marks.ListStudentMarks(c.Request.Context(), studentID, termID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, marks)
}

// ListEnrollments returns the signed-in student's own subject picks for a level.
func (h *StudentSelfHandler) ListEnrollments(c *gin.Context) {
	studentID, ok := h.resolveStudentID(c)
	if !ok {
		return
	}

	levelID, err := uuid.Parse(c.Query("level_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or missing level_id"})
		return
	}
	academicYearID, err := uuid.Parse(c.Query("academic_year_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or missing academic_year_id"})
		return
	}

	picks, err := h.enrollments.ListByStudentAndLevel(c.Request.Context(), studentID, levelID, academicYearID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	locked, err := h.enrollments.IsLocked(c.Request.Context(), studentID, levelID, academicYearID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"picks": picks, "locked": locked})
}

// SubmitEnrollment validates then replaces the student's picks for a level+year; rejected if already confirmed/locked.
func (h *StudentSelfHandler) SubmitEnrollment(c *gin.Context) {
	studentID, ok := h.resolveStudentID(c)
	if !ok {
		return
	}

	var req academicsmodule.SubmitEnrollmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validationErrs, err := h.enrollments.Submit(c.Request.Context(), studentID, req)
	if err != nil {
		switch {
		case errors.Is(err, academicsmodule.ErrEnrollmentInvalid):
			c.JSON(http.StatusUnprocessableEntity, academicsmodule.EnrollmentValidationResponse{Valid: false, Errors: validationErrs})
		case errors.Is(err, academicsmodule.ErrEnrollmentLocked):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, academicsmodule.ErrLevelHasNoGroups):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, academicsmodule.EnrollmentValidationResponse{Valid: true, Errors: nil})
}

// ConfirmEnrollment locks the student's picks for this level and year until an admin unlocks them.
func (h *StudentSelfHandler) ConfirmEnrollment(c *gin.Context) {
	studentID, ok := h.resolveStudentID(c)
	if !ok {
		return
	}

	var req academicsmodule.ConfirmEnrollmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	levelID, err := uuid.Parse(req.LevelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid level_id"})
		return
	}
	academicYearID, err := uuid.Parse(req.AcademicYearID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid academic_year_id"})
		return
	}

	if err := h.enrollments.Confirm(c.Request.Context(), studentID, levelID, academicYearID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "subject selection confirmed and locked"})
}
