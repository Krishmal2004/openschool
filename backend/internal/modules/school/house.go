package school

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/openschool-org/openschool/internal/apierror"
	"github.com/openschool-org/openschool/internal/platform/httpx"
	"github.com/openschool-org/openschool/internal/ports"
)

var (
	errHouseNotFound = errors.New("house not found")
	errHouseInUse    = errors.New("house is assigned to a student or teacher and cannot be deleted")
)

type House struct {
	ID        uuid.UUID          `json:"id"`
	Name      string             `json:"name"`
	Code      pgtype.Text        `json:"code"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	Color     string             `json:"color"`
}

type houseCommand struct {
	Name  string `json:"name" binding:"required"`
	Code  string `json:"code"`
	Color string `json:"color" binding:"required"`
}

type houseStore interface {
	create(context.Context, houseCommand) (House, error)
	get(context.Context, uuid.UUID) (House, error)
	list(context.Context) ([]House, error)
	update(context.Context, uuid.UUID, houseCommand) (House, error)
	delete(context.Context, uuid.UUID) (int64, error)
	pickForStudent(context.Context) (uuid.UUID, bool, error)
	pickForTeacher(context.Context) (uuid.UUID, bool, error)
	listStudentsMissingHouse(context.Context) ([]uuid.UUID, error)
	listTeachersMissingHouse(context.Context) ([]uuid.UUID, error)
	getStudent(context.Context, uuid.UUID) (ports.StudentHouseProfile, error)
	getTeacher(context.Context, uuid.UUID) (ports.TeacherHouseProfile, error)
	updateStudentHouse(context.Context, uuid.UUID, *uuid.UUID) (ports.StudentHouseProfile, error)
	updateTeacherHouse(context.Context, uuid.UUID, *uuid.UUID) (ports.TeacherHouseProfile, error)
}

// HouseService owns house CRUD, balancing, and audited member reassignment.
// It is exported only as the implementation of ports.HouseAssignments.
type HouseService struct {
	houses houseStore
	audit  ports.AuditRecorder
}

func NewHouseService(pool *pgxpool.Pool, audit ports.AuditRecorder) *HouseService {
	return &HouseService{houses: newHouseRepository(pool), audit: audit}
}

func (s *HouseService) create(ctx context.Context, command houseCommand) (House, error) {
	return s.houses.create(ctx, command)
}
func (s *HouseService) get(ctx context.Context, id uuid.UUID) (House, error) {
	return s.houses.get(ctx, id)
}
func (s *HouseService) list(ctx context.Context) ([]House, error) { return s.houses.list(ctx) }
func (s *HouseService) update(ctx context.Context, id uuid.UUID, command houseCommand) (House, error) {
	return s.houses.update(ctx, id, command)
}
func (s *HouseService) delete(ctx context.Context, id uuid.UUID) error {
	rows, err := s.houses.delete(ctx, id)
	if err != nil {
		return err
	}
	if rows != 0 {
		return nil
	}
	if _, err := s.houses.get(ctx, id); err != nil {
		return errHouseNotFound
	}
	return errHouseInUse
}

func (s *HouseService) PickForStudent(ctx context.Context) (uuid.UUID, bool) {
	id, ok, err := s.houses.pickForStudent(ctx)
	return id, ok && err == nil
}
func (s *HouseService) PickForTeacher(ctx context.Context) (uuid.UUID, bool) {
	id, ok, err := s.houses.pickForTeacher(ctx)
	return id, ok && err == nil
}

func (s *HouseService) reassignMissing(ctx context.Context) (int, error) {
	students, err := s.houses.listStudentsMissingHouse(ctx)
	if err != nil {
		return 0, err
	}
	assigned := 0
	for _, studentID := range students {
		houseID, ok := s.PickForStudent(ctx)
		if !ok {
			break
		}
		if _, err := s.houses.updateStudentHouse(ctx, studentID, &houseID); err != nil {
			return assigned, err
		}
		assigned++
	}
	return assigned, nil
}

func (s *HouseService) reassignMissingStaff(ctx context.Context) (int, error) {
	teachers, err := s.houses.listTeachersMissingHouse(ctx)
	if err != nil {
		return 0, err
	}
	assigned := 0
	for _, teacherID := range teachers {
		houseID, ok := s.PickForTeacher(ctx)
		if !ok {
			break
		}
		if _, err := s.houses.updateTeacherHouse(ctx, teacherID, &houseID); err != nil {
			return assigned, err
		}
		assigned++
	}
	return assigned, nil
}

type houseState struct {
	HouseID *uuid.UUID `json:"house_id"`
}

func parseOptionalHouseID(raw string) (*uuid.UUID, error) {
	if raw == "" {
		return nil, nil
	}
	id, err := uuid.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("invalid house id")
	}
	return &id, nil
}

func houseIDPtr(id pgtype.UUID) *uuid.UUID {
	if !id.Valid {
		return nil
	}
	value := uuid.UUID(id.Bytes)
	return &value
}

func (s *HouseService) ChangeStudentHouse(ctx context.Context, studentID uuid.UUID, rawHouseID string, actorID uuid.UUID) (ports.StudentHouseProfile, error) {
	houseID, err := parseOptionalHouseID(rawHouseID)
	if err != nil {
		return ports.StudentHouseProfile{}, err
	}
	before, err := s.houses.getStudent(ctx, studentID)
	if err != nil {
		return ports.StudentHouseProfile{}, err
	}
	updated, err := s.houses.updateStudentHouse(ctx, studentID, houseID)
	if err != nil {
		return ports.StudentHouseProfile{}, err
	}
	if s.audit != nil {
		_ = s.audit.Record(ctx, "student_house", studentID, "house_changed", actorID,
			houseState{HouseID: houseIDPtr(before.HouseID)}, houseState{HouseID: houseID}, "")
	}
	return updated, nil
}

func (s *HouseService) ChangeTeacherHouse(ctx context.Context, teacherID uuid.UUID, rawHouseID string, actorID uuid.UUID) (ports.TeacherHouseProfile, error) {
	houseID, err := parseOptionalHouseID(rawHouseID)
	if err != nil {
		return ports.TeacherHouseProfile{}, err
	}
	before, err := s.houses.getTeacher(ctx, teacherID)
	if err != nil {
		return ports.TeacherHouseProfile{}, err
	}
	updated, err := s.houses.updateTeacherHouse(ctx, teacherID, houseID)
	if err != nil {
		return ports.TeacherHouseProfile{}, err
	}
	if s.audit != nil {
		_ = s.audit.Record(ctx, "teacher_house", teacherID, "house_changed", actorID,
			houseState{HouseID: houseIDPtr(before.HouseID)}, houseState{HouseID: houseID}, "")
	}
	return updated, nil
}

type houseHandler struct{ service *HouseService }

func (h *houseHandler) create(c *gin.Context) {
	var command houseCommand
	if err := httpx.BindStrict(c, &command); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	house, err := h.service.create(c.Request.Context(), command)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, house)
}

func (h *houseHandler) get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	house, err := h.service.get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "house not found"})
		return
	}
	c.JSON(http.StatusOK, house)
}

func (h *houseHandler) list(c *gin.Context) {
	houses, err := h.service.list(c.Request.Context())
	if err != nil {
		apierror.RespondInternal(c, err)
		return
	}
	c.JSON(http.StatusOK, houses)
}

func (h *houseHandler) update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var command houseCommand
	if err := httpx.BindStrict(c, &command); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	house, err := h.service.update(c.Request.Context(), id, command)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, house)
}

func (h *houseHandler) delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.service.delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, errHouseNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "house deleted"})
}

func (h *houseHandler) reassignMissing(c *gin.Context) {
	assigned, err := h.service.reassignMissing(c.Request.Context())
	if err != nil {
		apierror.RespondInternal(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"assigned": assigned})
}

func (h *houseHandler) reassignMissingStaff(c *gin.Context) {
	assigned, err := h.service.reassignMissingStaff(c.Request.Context())
	if err != nil {
		apierror.RespondInternal(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"assigned": assigned})
}
