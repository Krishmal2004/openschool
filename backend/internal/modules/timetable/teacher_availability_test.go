package timetable

import (
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"
)

type teacherAvailabilityStoreStub struct {
	teacherID uuid.UUID
	command   teacherAvailabilityCommand
}

func (s *teacherAvailabilityStoreStub) createAvailability(_ context.Context, teacherID uuid.UUID, command teacherAvailabilityCommand) (TeacherAvailability, error) {
	s.teacherID, s.command = teacherID, command
	return TeacherAvailability{
		TeacherID: teacherID, AcademicYearID: command.AcademicYearID,
		DayOfWeek: command.DayOfWeek, PeriodNumber: command.PeriodNumber,
	}, nil
}
func (s *teacherAvailabilityStoreStub) listAvailability(context.Context, uuid.UUID, uuid.UUID) ([]TeacherAvailability, error) {
	return nil, nil
}
func (s *teacherAvailabilityStoreStub) deleteAvailability(context.Context, uuid.UUID) error {
	return nil
}

func TestTeacherAvailabilityRejectsWeekendAndInvalidDays(t *testing.T) {
	service := &teacherAvailabilityService{availability: &teacherAvailabilityStoreStub{}}
	for _, day := range []int16{-1, 0, 6, 7} {
		_, err := service.create(context.Background(), uuid.New(), teacherAvailabilityCommand{DayOfWeek: day})
		if err == nil || !strings.Contains(err.Error(), "between 1 (Monday) and 5 (Friday)") {
			t.Fatalf("create(day=%d) error = %v", day, err)
		}
	}
}

func TestTeacherAvailabilityAcceptsWeekdays(t *testing.T) {
	for _, day := range []int16{1, 3, 5} {
		store := &teacherAvailabilityStoreStub{}
		service := &teacherAvailabilityService{availability: store}
		teacherID, yearID := uuid.New(), uuid.New()
		command := teacherAvailabilityCommand{AcademicYearID: yearID, DayOfWeek: day, PeriodNumber: 4}
		result, err := service.create(context.Background(), teacherID, command)
		if err != nil {
			t.Fatalf("create(day=%d): %v", day, err)
		}
		if store.teacherID != teacherID || store.command != command {
			t.Fatalf("repository received teacher=%s command=%+v", store.teacherID, store.command)
		}
		if result.TeacherID != teacherID || result.DayOfWeek != day || result.PeriodNumber != 4 {
			t.Fatalf("result = %+v", result)
		}
	}
}
