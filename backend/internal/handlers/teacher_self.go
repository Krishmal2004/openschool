package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/openschool-org/openschool/db/sqlc"
	"github.com/openschool-org/openschool/internal/middleware"
	"github.com/openschool-org/openschool/internal/repositories"
	"github.com/openschool-org/openschool/internal/services"
	timetableservices "github.com/openschool-org/openschool/internal/services/timetable"
)

// TeacherSelfHandler resolves a signed-in teacher's own profile ID for the existing teacherOrAdmin routes to use.
type TeacherSelfHandler struct {
	teachers        *repositories.TeacherRepository
	school          *repositories.SchoolRepository
	positions       *services.PositionService
	societies       *services.SocietyService
	dashboard       *services.DashboardService
	timetables      *timetableservices.TimetableService
	staffAttendance *services.StaffAttendanceService
}

// NewTeacherSelfHandler constructs a TeacherSelfHandler with its service dependencies.
func NewTeacherSelfHandler(
	teachers *repositories.TeacherRepository,
	school *repositories.SchoolRepository,
	positions *services.PositionService,
	societies *services.SocietyService,
	dashboard *services.DashboardService,
	timetables *timetableservices.TimetableService,
	staffAttendance *services.StaffAttendanceService,
) *TeacherSelfHandler {
	return &TeacherSelfHandler{
		teachers:        teachers,
		school:          school,
		positions:       positions,
		societies:       societies,
		dashboard:       dashboard,
		timetables:      timetables,
		staffAttendance: staffAttendance,
	}
}

// leadershipTeacher resolves the caller's teacher profile and 403s unless their rank is Principal or Vice Principal.
func (h *TeacherSelfHandler) leadershipTeacher(c *gin.Context) (db.TeacherProfile, bool) {
	callerID, err := middleware.UserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid caller identity"})
		return db.TeacherProfile{}, false
	}

	teacher, err := h.teachers.GetByUserID(c.Request.Context(), callerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no teacher profile linked to this account"})
		return db.TeacherProfile{}, false
	}

	year, err := h.school.GetCurrentAcademicYear(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no current academic year configured"})
		return db.TeacherProfile{}, false
	}

	rank, _, err := h.positions.RankForTeacher(c.Request.Context(), teacher.ID, year.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return db.TeacherProfile{}, false
	}
	if !rank.IsPrincipalOrVicePrincipal() {
		c.JSON(http.StatusForbidden, gin.H{"error": "this action requires the Principal or Vice Principal position"})
		return db.TeacherProfile{}, false
	}

	return teacher, true
}

// Profile returns the signed-in teacher's own profile.
func (h *TeacherSelfHandler) Profile(c *gin.Context) {
	callerID, err := middleware.UserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid caller identity"})
		return
	}

	teacher, err := h.teachers.GetByUserID(c.Request.Context(), callerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no teacher profile linked to this account"})
		return
	}

	c.JSON(http.StatusOK, teacher)
}

// Position returns the signed-in teacher's leadership rank and notification reach.
func (h *TeacherSelfHandler) Position(c *gin.Context) {
	callerID, err := middleware.UserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid caller identity"})
		return
	}

	teacher, err := h.teachers.GetByUserID(c.Request.Context(), callerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no teacher profile linked to this account"})
		return
	}

	year, err := h.school.GetCurrentAcademicYear(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no current academic year configured"})
		return
	}

	summary, err := h.positions.SummaryForTeacher(c.Request.Context(), teacher.ID, year.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// Society returns the society the signed-in teacher is Teacher-in-Charge of, for the current academic year.
func (h *TeacherSelfHandler) Society(c *gin.Context) {
	callerID, err := middleware.UserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid caller identity"})
		return
	}

	teacher, err := h.teachers.GetByUserID(c.Request.Context(), callerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no teacher profile linked to this account"})
		return
	}

	year, err := h.school.GetCurrentAcademicYear(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no current academic year configured"})
		return
	}

	society, err := h.societies.GetForTeacher(c.Request.Context(), teacher.ID, year.ID)
	if err != nil {
		if errors.Is(err, services.ErrSocietyNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "you are not the Teacher-in-Charge of any society this year"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, society)
}

// LeadershipOverview returns the signed-in teacher's scoped leadership-panel counts (Section Head and above only).
func (h *TeacherSelfHandler) LeadershipOverview(c *gin.Context) {
	callerID, err := middleware.UserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid caller identity"})
		return
	}

	teacher, err := h.teachers.GetByUserID(c.Request.Context(), callerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no teacher profile linked to this account"})
		return
	}

	year, err := h.school.GetCurrentAcademicYear(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no current academic year configured"})
		return
	}

	overview, err := h.positions.LeadershipOverview(c.Request.Context(), teacher.ID, year.ID)
	if err != nil {
		if errors.Is(err, services.ErrInsufficientRank) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, overview)
}

// Analytics returns school-wide analytics for the signed-in Principal or Vice Principal.
func (h *TeacherSelfHandler) Analytics(c *gin.Context) {
	if _, ok := h.leadershipTeacher(c); !ok {
		return
	}

	analytics, err := h.dashboard.Analytics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// Timetables returns the same listing as the admin Timetables page, read-only, for Principal/Vice Principal.
func (h *TeacherSelfHandler) Timetables(c *gin.Context) {
	if _, ok := h.leadershipTeacher(c); !ok {
		return
	}

	year, err := h.school.GetCurrentAcademicYear(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no current academic year configured"})
		return
	}

	list, err := h.timetables.ListByAcademicYear(c.Request.Context(), year.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}

// Attendance returns the signed-in teacher's own staff-attendance history for a month.
func (h *TeacherSelfHandler) Attendance(c *gin.Context) {
	callerID, err := middleware.UserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid caller identity"})
		return
	}

	teacher, err := h.teachers.GetByUserID(c.Request.Context(), callerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no teacher profile linked to this account"})
		return
	}

	year, month, err := parseYearMonth(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	from := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	to := from.AddDate(0, 1, -1)

	records, err := h.staffAttendance.TeacherHistory(c.Request.Context(), teacher.ID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, records)
}
