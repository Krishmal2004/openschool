package school

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/validation"
)

var (
	errAcademicYearNotFound = errors.New("academic year not found")
	errAcademicYearInUse    = errors.New("academic year is used by a class and cannot be deleted")
	errInvalidLogoURL       = errors.New("logo must be a data:image/... URL under 500KB")
)

const maxLogoDataURLChars = 700_000

type School struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Address    *string `json:"address"`
	Phone      *string `json:"phone"`
	Email      *string `json:"email"`
	LogoURL    *string `json:"logo_url"`
	GradeFrom  *int32  `json:"grade_from"`
	GradeTo    *int32  `json:"grade_to"`
	SchoolType string  `json:"school_type"`
	CreatedAt  string  `json:"created_at"`
}
type AcademicYear struct {
	ID        string `json:"id"`
	Label     string `json:"label"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	IsCurrent bool   `json:"is_current"`
	CreatedAt string `json:"created_at"`
}
type schoolCommand struct {
	Name       string `json:"name" binding:"required"`
	Address    string `json:"address"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	LogoURL    string `json:"logo_url"`
	GradeFrom  *int32 `json:"grade_from"`
	GradeTo    *int32 `json:"grade_to"`
	SchoolType string `json:"school_type" binding:"omitempty,oneof=boys girls mixed"`
}
type yearCommand struct {
	Label     string    `json:"label" binding:"required"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
	IsCurrent bool      `json:"is_current"`
}
type schoolValues struct {
	ID                                               uuid.UUID
	Name, Address, Phone, Email, LogoURL, SchoolType string
	GradeFrom, GradeTo                               *int32
}
type yearValues struct {
	Label              string
	StartDate, EndDate time.Time
	IsCurrent          bool
}
type schoolStore interface {
	createSchool(context.Context, schoolValues) (School, error)
	getSchool(context.Context) (School, error)
	updateSchool(context.Context, schoolValues) (School, error)
	createYear(context.Context, yearValues) (AcademicYear, error)
	getYear(context.Context, uuid.UUID) (AcademicYear, error)
	currentYear(context.Context) (AcademicYear, error)
	listYears(context.Context) ([]AcademicYear, error)
	setCurrentYear(context.Context, uuid.UUID) error
	deleteYear(context.Context, uuid.UUID) (int64, error)
}
type schoolService struct{ store schoolStore }

func validateLogoURL(value string) error {
	if value != "" && (!strings.HasPrefix(value, "data:image/") || len(value) > maxLogoDataURLChars) {
		return errInvalidLogoURL
	}
	return nil
}
func validateSchool(command schoolCommand) error {
	if err := validateLogoURL(command.LogoURL); err != nil {
		return err
	}
	if !validation.IsValidSriLankanPhone(command.Phone) {
		return validation.ErrInvalidPhone
	}
	return nil
}
func schoolTypeOrDefault(value string) string {
	if value == "" {
		return "mixed"
	}
	return value
}
func schoolValuesFrom(command schoolCommand) schoolValues {
	return schoolValues{Name: command.Name, Address: command.Address, Phone: command.Phone, Email: command.Email, LogoURL: command.LogoURL, GradeFrom: command.GradeFrom, GradeTo: command.GradeTo, SchoolType: command.SchoolType}
}
func (s *schoolService) createSchool(ctx context.Context, command schoolCommand) (School, error) {
	if existing, err := s.store.getSchool(ctx); err == nil && existing.ID != "" {
		return School{}, fmt.Errorf("school already exists")
	}
	if err := validateSchool(command); err != nil {
		return School{}, err
	}
	values := schoolValuesFrom(command)
	values.SchoolType = schoolTypeOrDefault(values.SchoolType)
	return s.store.createSchool(ctx, values)
}
func (s *schoolService) updateSchool(ctx context.Context, id uuid.UUID, command schoolCommand) (School, error) {
	if err := validateSchool(command); err != nil {
		return School{}, err
	}
	values := schoolValuesFrom(command)
	values.ID = id
	return s.store.updateSchool(ctx, values)
}
func (s *schoolService) setCurrentYear(ctx context.Context, id uuid.UUID) error {
	if _, err := s.store.getYear(ctx, id); err != nil {
		return errAcademicYearNotFound
	}
	return s.store.setCurrentYear(ctx, id)
}
func (s *schoolService) deleteYear(ctx context.Context, id uuid.UUID) error {
	rows, err := s.store.deleteYear(ctx, id)
	if err != nil {
		return err
	}
	if rows != 0 {
		return nil
	}
	if _, err := s.store.getYear(ctx, id); err != nil {
		return errAcademicYearNotFound
	}
	return errAcademicYearInUse
}

type schoolHandler struct{ service *schoolService }

func newSchoolHandler(store schoolStore) *schoolHandler {
	return &schoolHandler{service: &schoolService{store: store}}
}
func (h *schoolHandler) createSchool(c *gin.Context) {
	var command schoolCommand
	if err := c.ShouldBindJSON(&command); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	value, err := h.service.createSchool(c.Request.Context(), command)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, value)
}
func (h *schoolHandler) getSchool(c *gin.Context) {
	value, err := h.service.store.getSchool(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "school not found"})
		return
	}
	c.JSON(http.StatusOK, value)
}
func (h *schoolHandler) updateSchool(c *gin.Context) {
	id, ok := parseSchoolID(c)
	if !ok {
		return
	}
	var command schoolCommand
	if err := c.ShouldBindJSON(&command); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	value, err := h.service.updateSchool(c.Request.Context(), id, command)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, value)
}
func (h *schoolHandler) createYear(c *gin.Context) {
	var command yearCommand
	if err := c.ShouldBindJSON(&command); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	value, err := h.service.store.createYear(c.Request.Context(), yearValues{Label: command.Label, StartDate: command.StartDate, EndDate: command.EndDate, IsCurrent: command.IsCurrent})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, value)
}
func (h *schoolHandler) listYears(c *gin.Context) {
	values, err := h.service.store.listYears(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, values)
}
func (h *schoolHandler) currentYear(c *gin.Context) {
	value, err := h.service.store.currentYear(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no current academic year"})
		return
	}
	c.JSON(http.StatusOK, value)
}
func (h *schoolHandler) setCurrentYear(c *gin.Context) {
	id, ok := parseSchoolID(c)
	if !ok {
		return
	}
	if err := h.service.setCurrentYear(c.Request.Context(), id); err != nil {
		if errors.Is(err, errAcademicYearNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "current academic year updated"})
}
func (h *schoolHandler) deleteYear(c *gin.Context) {
	id, ok := parseSchoolID(c)
	if !ok {
		return
	}
	if err := h.service.deleteYear(c.Request.Context(), id); err != nil {
		if errors.Is(err, errAcademicYearNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "academic year deleted"})
}
func parseSchoolID(c *gin.Context) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return uuid.Nil, false
	}
	return id, true
}
