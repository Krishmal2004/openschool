package timetable

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/apierror"
)

var errGradeSectionNotFound = errors.New("grade section not found, or still has grades assigned")

type GradeSection struct {
	ID                   uuid.UUID   `json:"id"`
	AcademicYearID       uuid.UUID   `json:"academic_year_id"`
	Name                 string      `json:"name"`
	IntervalStartTime    string      `json:"interval_start_time"`
	IntervalEndTime      string      `json:"interval_end_time"`
	SectionHeadTeacherID *uuid.UUID  `json:"section_head_teacher_id,omitempty"`
	SectionHeadName      *string     `json:"section_head_name,omitempty"`
	SortOrder            int32       `json:"sort_order"`
	GradeIDs             []uuid.UUID `json:"grade_ids"`
}

type gradeSectionCommand struct {
	AcademicYearID       uuid.UUID   `json:"academic_year_id" binding:"required"`
	Name                 string      `json:"name" binding:"required"`
	IntervalStartTime    string      `json:"interval_start_time" binding:"required"`
	IntervalEndTime      string      `json:"interval_end_time" binding:"required"`
	SectionHeadTeacherID *uuid.UUID  `json:"section_head_teacher_id"`
	SortOrder            int32       `json:"sort_order"`
	GradeIDs             []uuid.UUID `json:"grade_ids"`
}
type updateGradeSectionCommand struct {
	Name                 string     `json:"name" binding:"required"`
	IntervalStartTime    string     `json:"interval_start_time" binding:"required"`
	IntervalEndTime      string     `json:"interval_end_time" binding:"required"`
	SectionHeadTeacherID *uuid.UUID `json:"section_head_teacher_id"`
	SortOrder            int32      `json:"sort_order"`
}
type assignGradesCommand struct {
	GradeIDs []uuid.UUID `json:"grade_ids" binding:"required"`
}
type periodCommand struct {
	PeriodNumber *int32 `json:"period_number"`
	StartTime    string `json:"start_time" binding:"required"`
	EndTime      string `json:"end_time" binding:"required"`
	SlotType     string `json:"slot_type" binding:"required"`
}
type savePeriodsCommand struct {
	Periods []periodCommand `json:"periods" binding:"required"`
}

type gradeSectionRecord struct {
	ID, AcademicYearID                                 uuid.UUID
	Name                                               string
	IntervalStartMicroseconds, IntervalEndMicroseconds int64
	SectionHeadTeacherID                               *uuid.UUID
	SectionHeadName                                    *string
	SortOrder                                          int32
}
type periodValues struct {
	GradeSectionID                     uuid.UUID
	SortOrder                          int32
	PeriodNumber                       *int32
	StartMicroseconds, EndMicroseconds int64
	SlotType                           string
}
type TimetablePeriod struct {
	ID           uuid.UUID `json:"id"`
	SortOrder    int32     `json:"sort_order"`
	PeriodNumber *int32    `json:"period_number"`
	StartTime    string    `json:"start_time"`
	EndTime      string    `json:"end_time"`
	SlotType     string    `json:"slot_type"`
}

type gradeSectionStore interface {
	createSection(context.Context, gradeSectionCommand, int64, int64) (gradeSectionRecord, error)
	getSection(context.Context, uuid.UUID) (gradeSectionRecord, error)
	listSections(context.Context, uuid.UUID) ([]gradeSectionRecord, error)
	updateSection(context.Context, uuid.UUID, updateGradeSectionCommand, int64, int64) error
	deleteSection(context.Context, uuid.UUID) (int64, error)
	assignGrade(context.Context, uuid.UUID, uuid.UUID, uuid.UUID) error
	removeGrade(context.Context, uuid.UUID, uuid.UUID) error
	listGradeIDs(context.Context, uuid.UUID) ([]uuid.UUID, error)
	createPeriod(context.Context, periodValues) (TimetablePeriod, error)
	listPeriods(context.Context, uuid.UUID) ([]TimetablePeriod, error)
	deletePeriods(context.Context, uuid.UUID) error
	getSettingsValues(context.Context, uuid.UUID) (settingsValues, error)
}

type gradeSectionService struct{ sections gradeSectionStore }

func validateInterval(startText, endText string) (int64, int64, error) {
	start, err := parseClock(startText)
	if err != nil {
		return 0, 0, err
	}
	end, err := parseClock(endText)
	if err != nil {
		return 0, 0, err
	}
	if end <= start {
		return 0, 0, fmt.Errorf("interval end time must be after start time")
	}
	return start, end, nil
}
func formatClockMicros(value int64) string {
	seconds := value / 1_000_000
	return fmt.Sprintf("%02d:%02d", seconds/3600, (seconds%3600)/60)
}
func mapGradeSection(record gradeSectionRecord, gradeIDs []uuid.UUID) GradeSection {
	if gradeIDs == nil {
		gradeIDs = []uuid.UUID{}
	}
	return GradeSection{
		ID: record.ID, AcademicYearID: record.AcademicYearID, Name: record.Name,
		IntervalStartTime:    formatClockMicros(record.IntervalStartMicroseconds),
		IntervalEndTime:      formatClockMicros(record.IntervalEndMicroseconds),
		SectionHeadTeacherID: record.SectionHeadTeacherID, SectionHeadName: record.SectionHeadName,
		SortOrder: record.SortOrder, GradeIDs: gradeIDs,
	}
}
func (s *gradeSectionService) create(ctx context.Context, command gradeSectionCommand) (GradeSection, error) {
	start, end, err := validateInterval(command.IntervalStartTime, command.IntervalEndTime)
	if err != nil {
		return GradeSection{}, err
	}
	section, err := s.sections.createSection(ctx, command, start, end)
	if err != nil {
		return GradeSection{}, err
	}
	for _, gradeID := range command.GradeIDs {
		if err := s.sections.assignGrade(ctx, section.ID, gradeID, command.AcademicYearID); err != nil {
			return GradeSection{}, err
		}
	}
	if err := s.generatePeriodsFromSettings(ctx, section); err != nil {
		return GradeSection{}, fmt.Errorf("section created, but could not generate periods: %w", err)
	}
	return s.get(ctx, section.ID)
}
func (s *gradeSectionService) get(ctx context.Context, id uuid.UUID) (GradeSection, error) {
	section, err := s.sections.getSection(ctx, id)
	if err != nil {
		return GradeSection{}, err
	}
	gradeIDs, err := s.sections.listGradeIDs(ctx, id)
	if err != nil {
		return GradeSection{}, err
	}
	return mapGradeSection(section, gradeIDs), nil
}
func (s *gradeSectionService) listByYear(ctx context.Context, yearID uuid.UUID) ([]GradeSection, error) {
	rows, err := s.sections.listSections(ctx, yearID)
	if err != nil {
		return nil, err
	}
	result := make([]GradeSection, len(rows))
	for i, row := range rows {
		gradeIDs, err := s.sections.listGradeIDs(ctx, row.ID)
		if err != nil {
			return nil, err
		}
		result[i] = mapGradeSection(row, gradeIDs)
	}
	return result, nil
}
func (s *gradeSectionService) update(ctx context.Context, id uuid.UUID, command updateGradeSectionCommand) (GradeSection, error) {
	start, end, err := validateInterval(command.IntervalStartTime, command.IntervalEndTime)
	if err != nil {
		return GradeSection{}, err
	}
	if err := s.sections.updateSection(ctx, id, command, start, end); err != nil {
		return GradeSection{}, err
	}
	return s.get(ctx, id)
}
func (s *gradeSectionService) delete(ctx context.Context, id uuid.UUID) error {
	rows, err := s.sections.deleteSection(ctx, id)
	if err != nil {
		return err
	}
	if rows == 0 {
		return errGradeSectionNotFound
	}
	return nil
}
func (s *gradeSectionService) assignGrades(ctx context.Context, sectionID, yearID uuid.UUID, gradeIDs []uuid.UUID) error {
	for _, gradeID := range gradeIDs {
		if err := s.sections.assignGrade(ctx, sectionID, gradeID, yearID); err != nil {
			return err
		}
	}
	return nil
}
func (s *gradeSectionService) getPeriods(ctx context.Context, sectionID uuid.UUID) ([]TimetablePeriod, error) {
	return s.sections.listPeriods(ctx, sectionID)
}
func (s *gradeSectionService) savePeriods(ctx context.Context, sectionID uuid.UUID, command savePeriodsCommand) ([]TimetablePeriod, error) {
	if err := s.sections.deletePeriods(ctx, sectionID); err != nil {
		return nil, err
	}
	for i, entry := range command.Periods {
		if entry.SlotType != "period" && entry.SlotType != "interval" {
			return nil, fmt.Errorf("slot_type must be 'period' or 'interval'")
		}
		start, err := parseClock(entry.StartTime)
		if err != nil {
			return nil, err
		}
		end, err := parseClock(entry.EndTime)
		if err != nil {
			return nil, err
		}
		if entry.SlotType == "period" && entry.PeriodNumber == nil {
			return nil, fmt.Errorf("period_number is required for slot_type 'period'")
		}
		if _, err := s.sections.createPeriod(ctx, periodValues{GradeSectionID: sectionID, SortOrder: int32(i), PeriodNumber: entry.PeriodNumber, StartMicroseconds: start, EndMicroseconds: end, SlotType: entry.SlotType}); err != nil {
			return nil, err
		}
	}
	return s.getPeriods(ctx, sectionID)
}

func (s *gradeSectionService) createGeneratedPeriods(ctx context.Context, section gradeSectionRecord, settings settingsValues, sequential bool) error {
	var slots []periodValues
	if sequential {
		slots = sequentialPeriodSlots(section, settings)
	} else {
		slots = splitPeriodSlots(section, settings)
	}
	for _, slot := range slots {
		if _, err := s.sections.createPeriod(ctx, slot); err != nil {
			return err
		}
	}
	return nil
}
func sequentialPeriodSlots(section gradeSectionRecord, settings settingsValues) []periodValues {
	cursor := settings.StartMicroseconds
	duration := int64(settings.PeriodDurationMinutes) * 60 * 1_000_000
	intervalInserted := false
	slots := make([]periodValues, 0, settings.NumberOfPeriods+1)
	for i := int32(0); i < settings.NumberOfPeriods; {
		if !intervalInserted && cursor >= section.IntervalStartMicroseconds {
			slots = append(slots, periodValues{GradeSectionID: section.ID, StartMicroseconds: section.IntervalStartMicroseconds, EndMicroseconds: section.IntervalEndMicroseconds, SlotType: "interval"})
			cursor, intervalInserted = section.IntervalEndMicroseconds, true
			continue
		}
		next := cursor + duration
		slots = append(slots, periodValues{GradeSectionID: section.ID, StartMicroseconds: cursor, EndMicroseconds: next, SlotType: "period"})
		cursor = next
		i++
	}
	if !intervalInserted {
		slots = append(slots, periodValues{GradeSectionID: section.ID, StartMicroseconds: section.IntervalStartMicroseconds, EndMicroseconds: section.IntervalEndMicroseconds, SlotType: "interval"})
	}
	return numberAndSortPeriodSlots(slots)
}
func splitPeriodSlots(section gradeSectionRecord, settings settingsValues) []periodValues {
	totalPeriods := int64(settings.NumberOfPeriods)
	beforeDuration := section.IntervalStartMicroseconds - settings.StartMicroseconds
	templateDuration := int64(settings.PeriodDurationMinutes) * 60 * 1_000_000
	periodsBefore := int64(math.Round(float64(beforeDuration) / float64(templateDuration)))
	if periodsBefore < 1 {
		periodsBefore = 1
	}
	if periodsBefore >= totalPeriods {
		periodsBefore = totalPeriods - 1
	}
	periodsAfter := totalPeriods - periodsBefore
	slots := make([]periodValues, 0, settings.NumberOfPeriods+1)
	cursor := settings.StartMicroseconds
	durationBefore := beforeDuration / periodsBefore
	for i := int64(0); i < periodsBefore; i++ {
		next := cursor + durationBefore
		if i == periodsBefore-1 {
			next = section.IntervalStartMicroseconds
		}
		slots = append(slots, periodValues{GradeSectionID: section.ID, StartMicroseconds: cursor, EndMicroseconds: next, SlotType: "period"})
		cursor = next
	}
	slots = append(slots, periodValues{GradeSectionID: section.ID, StartMicroseconds: section.IntervalStartMicroseconds, EndMicroseconds: section.IntervalEndMicroseconds, SlotType: "interval"})
	cursor = section.IntervalEndMicroseconds
	afterDuration := settings.EndMicroseconds - section.IntervalEndMicroseconds
	durationAfter := afterDuration / periodsAfter
	for i := int64(0); i < periodsAfter; i++ {
		next := cursor + durationAfter
		if i == periodsAfter-1 {
			next = settings.EndMicroseconds
		}
		slots = append(slots, periodValues{GradeSectionID: section.ID, StartMicroseconds: cursor, EndMicroseconds: next, SlotType: "period"})
		cursor = next
	}
	return numberAndSortPeriodSlots(slots)
}
func numberAndSortPeriodSlots(slots []periodValues) []periodValues {
	sort.SliceStable(slots, func(i, j int) bool { return slots[i].StartMicroseconds < slots[j].StartMicroseconds })
	periodNumber := int32(1)
	for i := range slots {
		slots[i].SortOrder = int32(i)
		if slots[i].SlotType == "period" {
			n := periodNumber
			slots[i].PeriodNumber = &n
			periodNumber++
		}
	}
	return slots
}
func (s *gradeSectionService) generatePeriodsFromSettings(ctx context.Context, section gradeSectionRecord) error {
	settings, err := s.sections.getSettingsValues(ctx, section.AcademicYearID)
	if err != nil {
		return nil
	}
	invalidBounds := settings.StartMicroseconds >= settings.EndMicroseconds ||
		section.IntervalStartMicroseconds >= section.IntervalEndMicroseconds ||
		section.IntervalStartMicroseconds < settings.StartMicroseconds || section.IntervalEndMicroseconds > settings.EndMicroseconds
	if invalidBounds || settings.NumberOfPeriods < 2 {
		return s.createGeneratedPeriods(ctx, section, settings, true)
	}
	return s.createGeneratedPeriods(ctx, section, settings, false)
}
func (s *gradeSectionService) regeneratePeriods(ctx context.Context, sectionID uuid.UUID) ([]TimetablePeriod, error) {
	section, err := s.sections.getSection(ctx, sectionID)
	if err != nil {
		return nil, err
	}
	if err := s.sections.deletePeriods(ctx, sectionID); err != nil {
		return nil, err
	}
	if err := s.generatePeriodsFromSettings(ctx, section); err != nil {
		return nil, err
	}
	return s.getPeriods(ctx, sectionID)
}

type gradeSectionHandler struct{ service *gradeSectionService }

func newGradeSectionHandler(store gradeSectionStore) *gradeSectionHandler {
	return &gradeSectionHandler{service: &gradeSectionService{sections: store}}
}
func parsePathID(c *gin.Context, name, message string) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param(name))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": message})
		return uuid.Nil, false
	}
	return id, true
}
func (h *gradeSectionHandler) create(c *gin.Context) {
	var command gradeSectionCommand
	if err := c.ShouldBindJSON(&command); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	section, err := h.service.create(c.Request.Context(), command)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, section)
}
func (h *gradeSectionHandler) listByYear(c *gin.Context) {
	yearID, err := uuid.Parse(c.Query("academic_year_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "a valid academic_year_id is required"})
		return
	}
	sections, err := h.service.listByYear(c.Request.Context(), yearID)
	if err != nil {
		apierror.RespondInternal(c, err)
		return
	}
	c.JSON(http.StatusOK, sections)
}
func (h *gradeSectionHandler) get(c *gin.Context) {
	id, ok := parsePathID(c, "id", "invalid id")
	if !ok {
		return
	}
	section, err := h.service.get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "grade section not found"})
		return
	}
	c.JSON(http.StatusOK, section)
}
func (h *gradeSectionHandler) update(c *gin.Context) {
	id, ok := parsePathID(c, "id", "invalid id")
	if !ok {
		return
	}
	var command updateGradeSectionCommand
	if err := c.ShouldBindJSON(&command); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	section, err := h.service.update(c.Request.Context(), id, command)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, section)
}
func (h *gradeSectionHandler) delete(c *gin.Context) {
	id, ok := parsePathID(c, "id", "invalid id")
	if !ok {
		return
	}
	if err := h.service.delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, errGradeSectionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "grade section deleted"})
}
func (h *gradeSectionHandler) assignGrades(c *gin.Context) {
	id, ok := parsePathID(c, "id", "invalid id")
	if !ok {
		return
	}
	var command assignGradesCommand
	if err := c.ShouldBindJSON(&command); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	section, err := h.service.get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "grade section not found"})
		return
	}
	if err := h.service.assignGrades(c.Request.Context(), id, section.AcademicYearID, command.GradeIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updated, err := h.service.get(c.Request.Context(), id)
	if err != nil {
		apierror.RespondInternal(c, err)
		return
	}
	c.JSON(http.StatusOK, updated)
}
func (h *gradeSectionHandler) removeGrade(c *gin.Context) {
	id, ok := parsePathID(c, "id", "invalid id")
	if !ok {
		return
	}
	gradeID, ok := parsePathID(c, "grade_id", "invalid grade id")
	if !ok {
		return
	}
	if err := h.service.sections.removeGrade(c.Request.Context(), id, gradeID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "grade removed from section"})
}
func (h *gradeSectionHandler) getPeriods(c *gin.Context) {
	id, ok := parsePathID(c, "id", "invalid id")
	if !ok {
		return
	}
	periods, err := h.service.getPeriods(c.Request.Context(), id)
	if err != nil {
		apierror.RespondInternal(c, err)
		return
	}
	c.JSON(http.StatusOK, periods)
}
func (h *gradeSectionHandler) savePeriods(c *gin.Context) {
	id, ok := parsePathID(c, "id", "invalid id")
	if !ok {
		return
	}
	var command savePeriodsCommand
	if err := c.ShouldBindJSON(&command); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	periods, err := h.service.savePeriods(c.Request.Context(), id, command)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, periods)
}
func (h *gradeSectionHandler) regeneratePeriods(c *gin.Context) {
	id, ok := parsePathID(c, "id", "invalid id")
	if !ok {
		return
	}
	periods, err := h.service.regeneratePeriods(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, periods)
}
