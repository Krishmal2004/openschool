package school

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/platform/httpx"
)

var (
	errTermNotFound = errors.New("term not found")
	errTermInUse    = errors.New("term has marks recorded against it and cannot be deleted")
)

type Term struct {
	ID             string `json:"id"`
	AcademicYearID string `json:"academic_year_id"`
	Name           string `json:"name"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
	IsCurrent      bool   `json:"is_current"`
	SortOrder      int32  `json:"sort_order"`
	CreatedAt      string `json:"created_at"`
}
type createTermCommand struct {
	AcademicYearID string    `json:"academic_year_id" binding:"required"`
	Name           string    `json:"name" binding:"required"`
	StartDate      time.Time `json:"start_date" binding:"required"`
	EndDate        time.Time `json:"end_date" binding:"required"`
	SortOrder      int32     `json:"sort_order"`
}
type updateTermCommand struct {
	Name      string    `json:"name" binding:"required"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
	SortOrder int32     `json:"sort_order"`
}
type termValues struct {
	AcademicYearID     uuid.UUID
	Name               string
	StartDate, EndDate time.Time
	SortOrder          int32
}
type termStore interface {
	create(context.Context, termValues) (Term, error)
	get(context.Context, uuid.UUID) (Term, error)
	list(context.Context, uuid.UUID) ([]Term, error)
	current(context.Context) (Term, error)
	setCurrent(context.Context, uuid.UUID) error
	update(context.Context, uuid.UUID, termValues) (Term, error)
	delete(context.Context, uuid.UUID) (int64, error)
}
type termService struct{ terms termStore }

func (s *termService) create(ctx context.Context, command createTermCommand) (Term, error) {
	yearID, err := uuid.Parse(command.AcademicYearID)
	if err != nil {
		return Term{}, errors.New("invalid academic_year_id")
	}
	return s.terms.create(ctx, termValues{AcademicYearID: yearID, Name: command.Name, StartDate: command.StartDate, EndDate: command.EndDate, SortOrder: command.SortOrder})
}
func (s *termService) setCurrent(ctx context.Context, id uuid.UUID) error {
	if _, err := s.terms.get(ctx, id); err != nil {
		return errTermNotFound
	}
	return s.terms.setCurrent(ctx, id)
}
func (s *termService) delete(ctx context.Context, id uuid.UUID) error {
	rows, err := s.terms.delete(ctx, id)
	if err != nil {
		return err
	}
	if rows != 0 {
		return nil
	}
	if _, err := s.terms.get(ctx, id); err != nil {
		return errTermNotFound
	}
	return errTermInUse
}

type termHandler struct{ service *termService }

func newTermHandler(store termStore) *termHandler {
	return &termHandler{service: &termService{terms: store}}
}
func (h *termHandler) create(c *gin.Context) {
	var command createTermCommand
	if err := httpx.BindStrict(c, &command); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	value, err := h.service.create(c.Request.Context(), command)
	respondTermWrite(c, http.StatusCreated, value, err)
}
func (h *termHandler) list(c *gin.Context) {
	yearID, err := uuid.Parse(c.Query("academic_year_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or missing academic_year_id"})
		return
	}
	values, err := h.service.terms.list(c.Request.Context(), yearID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, values)
}
func (h *termHandler) current(c *gin.Context) {
	value, err := h.service.terms.current(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no current term"})
		return
	}
	c.JSON(http.StatusOK, value)
}
func (h *termHandler) setCurrent(c *gin.Context) {
	id, ok := parseTermID(c)
	if !ok {
		return
	}
	if err := h.service.setCurrent(c.Request.Context(), id); err != nil {
		if errors.Is(err, errTermNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "current term updated"})
}
func (h *termHandler) update(c *gin.Context) {
	id, ok := parseTermID(c)
	if !ok {
		return
	}
	var command updateTermCommand
	if err := httpx.BindStrict(c, &command); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	value, err := h.service.terms.update(c.Request.Context(), id, termValues{Name: command.Name, StartDate: command.StartDate, EndDate: command.EndDate, SortOrder: command.SortOrder})
	respondTermWrite(c, http.StatusOK, value, err)
}
func (h *termHandler) delete(c *gin.Context) {
	id, ok := parseTermID(c)
	if !ok {
		return
	}
	err := h.service.delete(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, errTermNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "term deleted"})
}
func parseTermID(c *gin.Context) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return uuid.Nil, false
	}
	return id, true
}
func respondTermWrite(c *gin.Context, status int, value Term, err error) {
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(status, value)
}
