package timetable

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Settings struct {
	AcademicYearID          uuid.UUID `json:"academic_year_id"`
	SchoolStartTime         string    `json:"school_start_time"`
	SchoolEndTime           string    `json:"school_end_time"`
	NumberOfPeriods         int32     `json:"number_of_periods"`
	PeriodDurationMinutes   int32     `json:"period_duration_minutes"`
	IntervalDurationMinutes int32     `json:"interval_duration_minutes"`
}
type settingsRequest struct {
	AcademicYearID          uuid.UUID `json:"academic_year_id" binding:"required"`
	SchoolStartTime         string    `json:"school_start_time" binding:"required"`
	SchoolEndTime           string    `json:"school_end_time" binding:"required"`
	NumberOfPeriods         int32     `json:"number_of_periods" binding:"required"`
	PeriodDurationMinutes   int32     `json:"period_duration_minutes" binding:"required"`
	IntervalDurationMinutes int32     `json:"interval_duration_minutes" binding:"required"`
}
type settingsValues struct {
	AcademicYearID                                                  uuid.UUID
	StartMicroseconds, EndMicroseconds                              int64
	NumberOfPeriods, PeriodDurationMinutes, IntervalDurationMinutes int32
}
type settingsStore interface {
	upsert(context.Context, settingsValues) (Settings, error)
	getByYear(context.Context, uuid.UUID) (Settings, error)
}
type settingsService struct{ store settingsStore }

func parseClock(value string) (int64, error) {
	parsed, err := time.Parse("15:04", value)
	if err != nil {
		parsed, err = time.Parse("15:04:05", value)
		if err != nil {
			return 0, fmt.Errorf("invalid time %q, expected HH:MM", value)
		}
	}
	return int64(parsed.Hour()*3600+parsed.Minute()*60+parsed.Second()) * 1_000_000, nil
}
func (s *settingsService) upsert(ctx context.Context, request settingsRequest) (Settings, error) {
	start, err := parseClock(request.SchoolStartTime)
	if err != nil {
		return Settings{}, err
	}
	end, err := parseClock(request.SchoolEndTime)
	if err != nil {
		return Settings{}, err
	}
	if end <= start {
		return Settings{}, fmt.Errorf("school end time must be after start time")
	}
	if request.NumberOfPeriods <= 0 {
		return Settings{}, fmt.Errorf("number_of_periods must be positive")
	}
	return s.store.upsert(ctx, settingsValues{AcademicYearID: request.AcademicYearID, StartMicroseconds: start, EndMicroseconds: end, NumberOfPeriods: request.NumberOfPeriods, PeriodDurationMinutes: request.PeriodDurationMinutes, IntervalDurationMinutes: request.IntervalDurationMinutes})
}

type settingsHandler struct{ service *settingsService }

func newSettingsHandler(store settingsStore) *settingsHandler {
	return &settingsHandler{service: &settingsService{store: store}}
}
func (h *settingsHandler) upsert(c *gin.Context) {
	var request settingsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	settings, err := h.service.upsert(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, settings)
}
func (h *settingsHandler) getByYear(c *gin.Context) {
	yearID, err := uuid.Parse(c.Param("academic_year_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid academic year id"})
		return
	}
	settings, err := h.service.store.getByYear(c.Request.Context(), yearID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no timetable settings configured for this academic year"})
		return
	}
	c.JSON(http.StatusOK, settings)
}
