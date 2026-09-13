package timetable

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
)

type gradeSectionStoreStub struct {
	deleteRows    int64
	deleteErr     error
	deletedPeriod bool
	createdSlots  []periodValues
	periods       []TimetablePeriod
}

func (s *gradeSectionStoreStub) createSection(context.Context, gradeSectionCommand, int64, int64) (gradeSectionRecord, error) {
	return gradeSectionRecord{}, nil
}
func (s *gradeSectionStoreStub) getSection(context.Context, uuid.UUID) (gradeSectionRecord, error) {
	return gradeSectionRecord{}, nil
}
func (s *gradeSectionStoreStub) listSections(context.Context, uuid.UUID) ([]gradeSectionRecord, error) {
	return nil, nil
}
func (s *gradeSectionStoreStub) updateSection(context.Context, uuid.UUID, updateGradeSectionCommand, int64, int64) error {
	return nil
}
func (s *gradeSectionStoreStub) deleteSection(context.Context, uuid.UUID) (int64, error) {
	return s.deleteRows, s.deleteErr
}
func (s *gradeSectionStoreStub) assignGrade(context.Context, uuid.UUID, uuid.UUID, uuid.UUID) error {
	return nil
}
func (s *gradeSectionStoreStub) removeGrade(context.Context, uuid.UUID, uuid.UUID) error {
	return nil
}
func (s *gradeSectionStoreStub) listGradeIDs(context.Context, uuid.UUID) ([]uuid.UUID, error) {
	return nil, nil
}
func (s *gradeSectionStoreStub) createPeriod(_ context.Context, values periodValues) (TimetablePeriod, error) {
	s.createdSlots = append(s.createdSlots, values)
	return TimetablePeriod{}, nil
}
func (s *gradeSectionStoreStub) listPeriods(context.Context, uuid.UUID) ([]TimetablePeriod, error) {
	return s.periods, nil
}
func (s *gradeSectionStoreStub) deletePeriods(context.Context, uuid.UUID) error {
	s.deletedPeriod = true
	return nil
}
func (s *gradeSectionStoreStub) getSettingsValues(context.Context, uuid.UUID) (settingsValues, error) {
	return settingsValues{}, errors.New("not configured")
}

func clockMicros(t *testing.T, value string) int64 {
	t.Helper()
	result, err := parseClock(value)
	if err != nil {
		t.Fatal(err)
	}
	return result
}

func TestValidateGradeSectionInterval(t *testing.T) {
	start, end, err := validateInterval("10:30", "11:00")
	if err != nil || start != clockMicros(t, "10:30") || end != clockMicros(t, "11:00") {
		t.Fatalf("validateInterval() = %d, %d, %v", start, end, err)
	}
	for _, values := range [][2]string{{"11:00", "10:30"}, {"10:30", "10:30"}, {"bad", "11:00"}} {
		if _, _, err := validateInterval(values[0], values[1]); err == nil {
			t.Fatalf("validateInterval(%q, %q) succeeded", values[0], values[1])
		}
	}
}

func TestSplitPeriodSlots(t *testing.T) {
	sectionID := uuid.New()
	section := gradeSectionRecord{ID: sectionID, IntervalStartMicroseconds: clockMicros(t, "10:00"), IntervalEndMicroseconds: clockMicros(t, "10:20")}
	settings := settingsValues{StartMicroseconds: clockMicros(t, "08:00"), EndMicroseconds: clockMicros(t, "12:20"), NumberOfPeriods: 6, PeriodDurationMinutes: 40}
	slots := splitPeriodSlots(section, settings)
	if len(slots) != 7 {
		t.Fatalf("slot count = %d, want 7", len(slots))
	}
	periodNumber := int32(1)
	for i, slot := range slots {
		if slot.SortOrder != int32(i) {
			t.Fatalf("slot %d sort order = %d", i, slot.SortOrder)
		}
		if slot.SlotType == "interval" {
			if slot.PeriodNumber != nil || slot.StartMicroseconds != section.IntervalStartMicroseconds || slot.EndMicroseconds != section.IntervalEndMicroseconds {
				t.Fatalf("invalid interval slot: %+v", slot)
			}
			continue
		}
		if slot.PeriodNumber == nil || *slot.PeriodNumber != periodNumber {
			t.Fatalf("slot %d period number = %v, want %d", i, slot.PeriodNumber, periodNumber)
		}
		periodNumber++
	}
	if periodNumber != 7 {
		t.Fatalf("generated %d periods, want 6", periodNumber-1)
	}
}

func TestSequentialPeriodSlotsFallback(t *testing.T) {
	section := gradeSectionRecord{ID: uuid.New(), IntervalStartMicroseconds: clockMicros(t, "09:20"), IntervalEndMicroseconds: clockMicros(t, "09:40")}
	settings := settingsValues{StartMicroseconds: clockMicros(t, "08:00"), NumberOfPeriods: 3, PeriodDurationMinutes: 40}
	slots := sequentialPeriodSlots(section, settings)
	if len(slots) != 4 {
		t.Fatalf("slot count = %d, want 4", len(slots))
	}
	if slots[2].SlotType != "interval" {
		t.Fatalf("interval position = %+v, want index 2", slots)
	}
}

func TestSavePeriodsValidation(t *testing.T) {
	periodNumber := int32(1)
	tests := []struct {
		name  string
		entry periodCommand
		want  string
	}{
		{name: "invalid slot", entry: periodCommand{SlotType: "break", StartTime: "08:00", EndTime: "08:40"}, want: "slot_type must be 'period' or 'interval'"},
		{name: "missing period number", entry: periodCommand{SlotType: "period", StartTime: "08:00", EndTime: "08:40"}, want: "period_number is required for slot_type 'period'"},
		{name: "valid period", entry: periodCommand{SlotType: "period", StartTime: "08:00", EndTime: "08:40", PeriodNumber: &periodNumber}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := &gradeSectionStoreStub{}
			_, err := (&gradeSectionService{sections: store}).savePeriods(context.Background(), uuid.New(), savePeriodsCommand{Periods: []periodCommand{test.entry}})
			if test.want == "" && err != nil {
				t.Fatal(err)
			}
			if test.want != "" && (err == nil || err.Error() != test.want) {
				t.Fatalf("savePeriods() error = %v, want %q", err, test.want)
			}
		})
	}
}

func TestDeleteGradeSection(t *testing.T) {
	service := &gradeSectionService{sections: &gradeSectionStoreStub{}}
	if err := service.delete(context.Background(), uuid.New()); !errors.Is(err, errGradeSectionNotFound) {
		t.Fatalf("delete() error = %v, want %v", err, errGradeSectionNotFound)
	}
	service.sections = &gradeSectionStoreStub{deleteRows: 1}
	if err := service.delete(context.Background(), uuid.New()); err != nil {
		t.Fatal(err)
	}
}
