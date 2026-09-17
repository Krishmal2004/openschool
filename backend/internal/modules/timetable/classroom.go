package timetable

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/apierror"
)

var (
	errClassroomNotFound  = errors.New("classroom not found, or still referenced by a timetable")
	errLabRequiresSubject = errors.New("a lab classroom must be tagged to a subject")
)

type Classroom struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Code      *string    `json:"code"`
	Capacity  *int32     `json:"capacity"`
	CreatedAt time.Time  `json:"created_at"`
	RoomType  string     `json:"room_type"`
	SubjectID *uuid.UUID `json:"subject_id"`
}

type ClassroomListItem struct {
	Classroom
	SubjectName *string `json:"subject_name"`
}

type classroomCommand struct {
	Name      string     `json:"name" binding:"required"`
	Code      string     `json:"code"`
	Capacity  *int32     `json:"capacity"`
	RoomType  string     `json:"room_type" binding:"required,oneof=regular lab eca"`
	SubjectID *uuid.UUID `json:"subject_id"`
}

type classroomStore interface {
	create(context.Context, classroomCommand) (Classroom, error)
	list(context.Context) ([]ClassroomListItem, error)
	update(context.Context, uuid.UUID, classroomCommand) (Classroom, error)
	delete(context.Context, uuid.UUID) (int64, error)
}

type classroomService struct{ classrooms classroomStore }

func validateClassroom(command classroomCommand) error {
	if command.RoomType == "lab" && command.SubjectID == nil {
		return errLabRequiresSubject
	}
	return nil
}
func (s *classroomService) create(ctx context.Context, command classroomCommand) (Classroom, error) {
	if err := validateClassroom(command); err != nil {
		return Classroom{}, err
	}
	return s.classrooms.create(ctx, command)
}
func (s *classroomService) update(ctx context.Context, id uuid.UUID, command classroomCommand) (Classroom, error) {
	if err := validateClassroom(command); err != nil {
		return Classroom{}, err
	}
	return s.classrooms.update(ctx, id, command)
}
func (s *classroomService) delete(ctx context.Context, id uuid.UUID) error {
	rows, err := s.classrooms.delete(ctx, id)
	if err != nil {
		return err
	}
	if rows == 0 {
		return errClassroomNotFound
	}
	return nil
}

type classroomHandler struct{ service *classroomService }

func newClassroomHandler(store classroomStore) *classroomHandler {
	return &classroomHandler{service: &classroomService{classrooms: store}}
}
func (h *classroomHandler) create(c *gin.Context) {
	var command classroomCommand
	if err := c.ShouldBindJSON(&command); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	classroom, err := h.service.create(c.Request.Context(), command)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, classroom)
}
func (h *classroomHandler) list(c *gin.Context) {
	classrooms, err := h.service.classrooms.list(c.Request.Context())
	if err != nil {
		apierror.RespondInternal(c, err)
		return
	}
	c.JSON(http.StatusOK, classrooms)
}
func (h *classroomHandler) update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var command classroomCommand
	if err := c.ShouldBindJSON(&command); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	classroom, err := h.service.update(c.Request.Context(), id, command)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, classroom)
}
func (h *classroomHandler) delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.service.delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, errClassroomNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "classroom deleted"})
}
