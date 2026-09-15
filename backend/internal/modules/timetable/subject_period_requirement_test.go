package timetable

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
)

type subjectPeriodRequirementStoreStub struct {
	command subjectPeriodRequirementCommand
}

func (s *subjectPeriodRequirementStoreStub) upsert(_ context.Context, command subjectPeriodRequirementCommand) (SubjectPeriodRequirement, error) {
	s.command = command
	return SubjectPeriodRequirement{
		AcademicYearID: command.AcademicYearID, GradeID: command.GradeID, SubjectID: command.SubjectID,
		PeriodsPerWeek: command.PeriodsPerWeek, LabPeriodsPerWeek: command.LabPeriodsPerWeek,
		DoublePeriodBlocks: command.DoublePeriodBlocks,
	}, nil
}
func (s *subjectPeriodRequirementStoreStub) listByGrade(context.Context, uuid.UUID, uuid.UUID) ([]SubjectPeriodRequirementListItem, error) {
	return nil, nil
}
func (s *subjectPeriodRequirementStoreStub) deleteRequirement(context.Context, uuid.UUID) error {
	return nil
}

func TestSubjectPeriodRequirementValidation(t *testing.T) {
	tests := []struct {
		name    string
		command subjectPeriodRequirementCommand
		want    error
	}{
		{name: "lab periods exceed total", command: subjectPeriodRequirementCommand{PeriodsPerWeek: 4, LabPeriodsPerWeek: 5}, want: errLabPeriodsExceedTotal},
		{name: "double blocks exceed total", command: subjectPeriodRequirementCommand{PeriodsPerWeek: 5, DoublePeriodBlocks: 3}, want: errDoublePeriodBlocksExceedTotal},
		{name: "valid mixed schedule", command: subjectPeriodRequirementCommand{PeriodsPerWeek: 6, LabPeriodsPerWeek: 2, DoublePeriodBlocks: 2}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := (&subjectPeriodRequirementService{requirements: &subjectPeriodRequirementStoreStub{}}).upsert(context.Background(), test.command)
			if !errors.Is(err, test.want) {
				t.Fatalf("upsert() error = %v, want %v", err, test.want)
			}
		})
	}
}

func TestSubjectPeriodRequirementPassesValidatedCommand(t *testing.T) {
	store := &subjectPeriodRequirementStoreStub{}
	command := subjectPeriodRequirementCommand{
		AcademicYearID: uuid.New(), GradeID: uuid.New(), SubjectID: uuid.New(),
		PeriodsPerWeek: 8, LabPeriodsPerWeek: 2, DoublePeriodBlocks: 3,
	}
	result, err := (&subjectPeriodRequirementService{requirements: store}).upsert(context.Background(), command)
	if err != nil {
		t.Fatal(err)
	}
	if store.command != command {
		t.Fatalf("repository command = %+v, want %+v", store.command, command)
	}
	if result.PeriodsPerWeek != command.PeriodsPerWeek || result.DoublePeriodBlocks != command.DoublePeriodBlocks {
		t.Fatalf("result = %+v", result)
	}
}
