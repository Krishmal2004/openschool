package academics

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/apierror"
	"github.com/openschool-org/openschool/internal/platform/httpx"
)

var (
	errSubjectNotFound = errors.New("subject not found")
	errSubjectInUse    = errors.New("subject is in use and cannot be deleted")
)

type Subject struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Code      string  `json:"code"`
	Type      *string `json:"type"`
	MaxMarks  float64 `json:"max_marks"`
	CreatedAt string  `json:"created_at"`
}

type subjectCommand struct {
	Name     string   `json:"name" binding:"required"`
	Code     string   `json:"code" binding:"required"`
	Type     string   `json:"type"`
	MaxMarks *float64 `json:"max_marks"`
}

type normalizedSubjectCommand struct {
	Name, Code, Type string
	MaxMarks         float64
}

type subjectStore interface {
	create(context.Context, normalizedSubjectCommand) (Subject, error)
	get(context.Context, uuid.UUID) (Subject, error)
	list(context.Context) ([]Subject, error)
	update(context.Context, uuid.UUID, normalizedSubjectCommand) (Subject, error)
	delete(context.Context, uuid.UUID) (int64, error)
}

type subjectService struct{ subjects subjectStore }

func normalizeSubject(command subjectCommand) normalizedSubjectCommand {
	maxMarks := 100.0
	if command.MaxMarks != nil {
		maxMarks = *command.MaxMarks
	}
	return normalizedSubjectCommand{Name: command.Name, Code: command.Code, Type: command.Type, MaxMarks: maxMarks}
}

func (s *subjectService) create(ctx context.Context, command subjectCommand) (Subject, error) {
	return s.subjects.create(ctx, normalizeSubject(command))
}
func (s *subjectService) get(ctx context.Context, id uuid.UUID) (Subject, error) {
	return s.subjects.get(ctx, id)
}
func (s *subjectService) list(ctx context.Context) ([]Subject, error) { return s.subjects.list(ctx) }
func (s *subjectService) update(ctx context.Context, id uuid.UUID, command subjectCommand) (Subject, error) {
	return s.subjects.update(ctx, id, normalizeSubject(command))
}
func (s *subjectService) delete(ctx context.Context, id uuid.UUID) error {
	rows, err := s.subjects.delete(ctx, id)
	if err != nil {
		return err
	}
	if rows != 0 {
		return nil
	}
	if _, err := s.subjects.get(ctx, id); err != nil {
		return errSubjectNotFound
	}
	return errSubjectInUse
}

type subjectHandler struct{ service *subjectService }

func newSubjectHandler(store subjectStore) *subjectHandler {
	return &subjectHandler{service: &subjectService{subjects: store}}
}

func (h *subjectHandler) create(c *gin.Context) {
	var command subjectCommand
	if err := httpx.BindStrict(c, &command); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	subject, err := h.service.create(c.Request.Context(), command)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, subject)
}

func (h *subjectHandler) get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	subject, err := h.service.get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "subject not found"})
		return
	}
	c.JSON(http.StatusOK, subject)
}

func (h *subjectHandler) list(c *gin.Context) {
	subjects, err := h.service.list(c.Request.Context())
	if err != nil {
		apierror.RespondInternal(c, err)
		return
	}
	c.JSON(http.StatusOK, subjects)
}

func (h *subjectHandler) update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var command subjectCommand
	if err := httpx.BindStrict(c, &command); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	subject, err := h.service.update(c.Request.Context(), id, command)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subject)
}

func (h *subjectHandler) delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.service.delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, errSubjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "subject deleted"})
}
