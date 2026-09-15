package timetable

import (
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"
)

type settingsStoreStub struct{ values settingsValues }

func (s *settingsStoreStub) upsert(_ context.Context, values settingsValues) (Settings, error) {
	s.values = values
	return Settings{}, nil
}
func (s *settingsStoreStub) getByYear(context.Context, uuid.UUID) (Settings, error) {
	return Settings{}, nil
}

func TestParseClockAcceptsSupportedFormats(t *testing.T) {
	for _, value := range []string{"08:30", "08:30:15"} {
		if _, err := parseClock(value); err != nil {
			t.Fatalf("parseClock(%q): %v", value, err)
		}
	}
	if _, err := parseClock("8am"); err == nil {
		t.Fatal("parseClock accepted invalid time")
	}
}

func TestUpsertSettingsValidation(t *testing.T) {
	tests := []struct {
		name      string
		request   settingsRequest
		errorPart string
	}{
		{name: "end before start", request: settingsRequest{SchoolStartTime: "09:00", SchoolEndTime: "08:00", NumberOfPeriods: 1}, errorPart: "after start"},
		{name: "no periods", request: settingsRequest{SchoolStartTime: "08:00", SchoolEndTime: "09:00"}, errorPart: "must be positive"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := (&settingsService{store: &settingsStoreStub{}}).upsert(context.Background(), test.request)
			if err == nil || !strings.Contains(err.Error(), test.errorPart) {
				t.Fatalf("upsert() error = %v", err)
			}
		})
	}
}

func TestUpsertSettingsPassesPlainValuesToRepository(t *testing.T) {
	store := &settingsStoreStub{}
	request := settingsRequest{AcademicYearID: uuid.New(), SchoolStartTime: "08:00", SchoolEndTime: "14:00", NumberOfPeriods: 8, PeriodDurationMinutes: 40, IntervalDurationMinutes: 20}
	if _, err := (&settingsService{store: store}).upsert(context.Background(), request); err != nil {
		t.Fatal(err)
	}
	if store.values.StartMicroseconds != 8*3600*1_000_000 || store.values.EndMicroseconds != 14*3600*1_000_000 {
		t.Fatalf("unexpected clock values: %+v", store.values)
	}
}
