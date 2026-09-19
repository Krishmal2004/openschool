// Package school owns school structure and academic-calendar configuration.
package school

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/apierror"
	"github.com/openschool-org/openschool/internal/platform/httpx"
)

var (
	errGradeNotFound = errors.New("grade not found")
	errGradeInUse    = errors.New("grade is used by a class or curriculum level and cannot be deleted")
)

type Grade struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	SortOrder int32     `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}

type gradeCommand struct {
	Name      string `json:"name" binding:"required"`
	SortOrder int32  `json:"sort_order"`
}

type gradeStore interface {
	create(context.Context, gradeCommand) (Grade, error)
	get(context.Context, uuid.UUID) (Grade, error)
	list(context.Context) ([]Grade, error)
	update(context.Context, uuid.UUID, gradeCommand) (Grade, error)
	delete(context.Context, uuid.UUID) (int64, error)
}

type gradeService struct{ grades gradeStore }

func (s *gradeService) create(ctx context.Context, command gradeCommand) (Grade, error) {
	return s.grades.create(ctx, command)
}

func (s *gradeService) get(ctx context.Context, id uuid.UUID) (Grade, error) {
	return s.grades.get(ctx, id)
}

func (s *gradeService) list(ctx context.Context) ([]Grade, error) { return s.grades.list(ctx) }

func (s *gradeService) update(ctx context.Context, id uuid.UUID, command gradeCommand) (Grade, error) {
	return s.grades.update(ctx, id, command)
}

func (s *gradeService) delete(ctx context.Context, id uuid.UUID) error {
	rows, err := s.grades.delete(ctx, id)
	if err != nil {
		return err
	}
	if rows != 0 {
		return nil
	}
	if _, err := s.grades.get(ctx, id); err != nil {
		return errGradeNotFound
	}
	return errGradeInUse
}

type gradeHandler struct{ service *gradeService }

func newGradeHandler(store gradeStore) *gradeHandler {
	return &gradeHandler{service: &gradeService{grades: store}}
}

func (h *gradeHandler) create(c *gin.Context) {
	var command gradeCommand
	if err := httpx.BindStrict(c, &command); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	grade, err := h.service.create(c.Request.Context(), command)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, grade)
}

func (h *gradeHandler) get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	grade, err := h.service.get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "grade not found"})
		return
	}
	c.JSON(http.StatusOK, grade)
}

func (h *gradeHandler) list(c *gin.Context) {
	grades, err := h.service.list(c.Request.Context())
	if err != nil {
		apierror.RespondInternal(c, err)
		return
	}
	c.JSON(http.StatusOK, grades)
}

func (h *gradeHandler) update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var command gradeCommand
	if err := httpx.BindStrict(c, &command); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	grade, err := h.service.update(c.Request.Context(), id, command)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, grade)
}

func (h *gradeHandler) delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.service.delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, errGradeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "grade deleted"})
}
