package timetable

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/openschool-org/openschool/internal/middleware"
)

var errTimetableNotFound = errors.New("timetable not found")

type Timetable struct {
	ID                uuid.UUID  `json:"id"`
	AcademicYearID    uuid.UUID  `json:"academic_year_id"`
	ClassID           uuid.UUID  `json:"class_id"`
	Version           int32      `json:"version"`
	Status            string     `json:"status"`
	ParentTimetableID *uuid.UUID `json:"parent_timetable_id"`
	CreatedBy         uuid.UUID  `json:"created_by"`
	SubmittedAt       *time.Time `json:"submitted_at"`
	SubmittedBy       *uuid.UUID `json:"submitted_by"`
	ReviewedBy        *uuid.UUID `json:"reviewed_by"`
	ReviewedAt        *time.Time `json:"reviewed_at"`
	ReviewComments    *string    `json:"review_comments"`
	PublishedAt       *time.Time `json:"published_at"`
	PublishedBy       *uuid.UUID `json:"published_by"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type TimetableListItem struct {
	Timetable
	ClassName string `json:"class_name"`
	GradeName string `json:"grade_name"`
}

type timetableCreateRequest struct {
	AcademicYearID uuid.UUID `json:"academic_year_id" binding:"required"`
	ClassID        uuid.UUID `json:"class_id" binding:"required"`
}

type timetableCopyRequest struct {
	AcademicYearID    uuid.UUID `json:"academic_year_id" binding:"required"`
	ClassID           uuid.UUID `json:"class_id" binding:"required"`
	SourceTimetableID uuid.UUID `json:"source_timetable_id" binding:"required"`
}

type timetableCRUDStore interface {
	create(context.Context, uuid.UUID, uuid.UUID, uuid.UUID, *uuid.UUID) (Timetable, error)
	get(context.Context, uuid.UUID) (Timetable, error)
	copyEntries(context.Context, uuid.UUID, uuid.UUID) error
	listByClass(context.Context, uuid.UUID, uuid.UUID) ([]TimetableListItem, error)
	listByAcademicYear(context.Context, uuid.UUID) ([]TimetableListItem, error)
	deleteDraft(context.Context, uuid.UUID) (int64, error)
	archive(context.Context, uuid.UUID) (Timetable, error)
}

type timetableCRUDService struct{ store timetableCRUDStore }

func (s *timetableCRUDService) create(ctx context.Context, request timetableCreateRequest, actor uuid.UUID) (Timetable, error) {
	return s.store.create(ctx, request.AcademicYearID, request.ClassID, actor, nil)
}

func (s *timetableCRUDService) copy(ctx context.Context, request timetableCopyRequest, actor uuid.UUID) (Timetable, error) {
	draft, err := s.create(ctx, timetableCreateRequest{AcademicYearID: request.AcademicYearID, ClassID: request.ClassID}, actor)
	if err != nil {
		return Timetable{}, err
	}
	if err := s.store.copyEntries(ctx, request.SourceTimetableID, draft.ID); err != nil {
		return Timetable{}, err
	}
	return draft, nil
}

func (s *timetableCRUDService) revise(ctx context.Context, publishedID, actor uuid.UUID) (Timetable, error) {
	published, err := s.store.get(ctx, publishedID)
	if err != nil {
		return Timetable{}, err
	}
	if published.Status != statusPublished {
		return Timetable{}, errors.New("timetable is not in the required status for this action")
	}
	draft, err := s.store.create(ctx, published.AcademicYearID, published.ClassID, actor, &published.ID)
	if err != nil {
		return Timetable{}, err
	}
	if err := s.store.copyEntries(ctx, published.ID, draft.ID); err != nil {
		return Timetable{}, err
	}
	return draft, nil
}

type timetableCRUDHandler struct{ service *timetableCRUDService }

func newTimetableCRUDHandler(store timetableCRUDStore) *timetableCRUDHandler {
	return &timetableCRUDHandler{service: &timetableCRUDService{store: store}}
}

func RegisterTimetableCRUDRoutes(admin, teacherOrAdmin *gin.RouterGroup, pool *pgxpool.Pool) {
	handler := newTimetableCRUDHandler(newTimetableRepository(pool))
	admin.POST("/timetables", handler.create)
	admin.POST("/timetables/copy", handler.copy)
	admin.POST("/timetables/:id/revise", handler.revise)
	teacherOrAdmin.GET("/timetables/:id", handler.get)
	teacherOrAdmin.GET("/timetables", handler.listByClass)
	admin.GET("/timetables/by-year", handler.listByAcademicYear)
	admin.DELETE("/timetables/:id", handler.delete)
	admin.POST("/timetables/:id/archive", handler.archive)
}

func (h *timetableCRUDHandler) actor(c *gin.Context) (uuid.UUID, bool) {
	id, err := middleware.UserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid caller identity"})
		return uuid.Nil, false
	}
	return id, true
}

func (h *timetableCRUDHandler) writeError(c *gin.Context, err error) {
	if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, errTimetableNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": errTimetableNotFound.Error()})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}

func parseTimetableID(c *gin.Context) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return uuid.Nil, false
	}
	return id, true
}

func (h *timetableCRUDHandler) create(c *gin.Context) {
	actor, ok := h.actor(c)
	if !ok {
		return
	}
	var request timetableCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	timetable, err := h.service.create(c.Request.Context(), request, actor)
	if err != nil {
		h.writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, timetable)
}

func (h *timetableCRUDHandler) copy(c *gin.Context) {
	actor, ok := h.actor(c)
	if !ok {
		return
	}
	var request timetableCopyRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	timetable, err := h.service.copy(c.Request.Context(), request, actor)
	if err != nil {
		h.writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, timetable)
}

func (h *timetableCRUDHandler) revise(c *gin.Context) {
	actor, ok := h.actor(c)
	if !ok {
		return
	}
	id, ok := parseTimetableID(c)
	if !ok {
		return
	}
	timetable, err := h.service.revise(c.Request.Context(), id, actor)
	if err != nil {
		h.writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, timetable)
}

func (h *timetableCRUDHandler) get(c *gin.Context) {
	id, ok := parseTimetableID(c)
	if !ok {
		return
	}
	timetable, err := h.service.store.get(c.Request.Context(), id)
	if err != nil {
		h.writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, timetable)
}

func (h *timetableCRUDHandler) listByClass(c *gin.Context) {
	classID, err := uuid.Parse(c.Query("class_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "a valid class_id is required"})
		return
	}
	yearID, err := uuid.Parse(c.Query("academic_year_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "a valid academic_year_id is required"})
		return
	}
	list, err := h.service.store.listByClass(c.Request.Context(), classID, yearID)
	if err != nil {
		h.writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *timetableCRUDHandler) listByAcademicYear(c *gin.Context) {
	yearID, err := uuid.Parse(c.Query("academic_year_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "a valid academic_year_id is required"})
		return
	}
	list, err := h.service.store.listByAcademicYear(c.Request.Context(), yearID)
	if err != nil {
		h.writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *timetableCRUDHandler) delete(c *gin.Context) {
	id, ok := parseTimetableID(c)
	if !ok {
		return
	}
	count, err := h.service.store.deleteDraft(c.Request.Context(), id)
	if err != nil {
		h.writeError(c, err)
		return
	}
	if count == 0 {
		h.writeError(c, errTimetableNotFound)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "timetable deleted"})
}

func (h *timetableCRUDHandler) archive(c *gin.Context) {
	id, ok := parseTimetableID(c)
	if !ok {
		return
	}
	timetable, err := h.service.store.archive(c.Request.Context(), id)
	if err != nil {
		h.writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, timetable)
}
