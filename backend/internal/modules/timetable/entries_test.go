package timetable

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
)

type timetableEntryStoreStub struct {
	statusValue string
	statusErr   error
	upserts     []timetableEntryCommand
	deleted     bool
	entries     []TimetableEntry
}

func (s *timetableEntryStoreStub) status(context.Context, uuid.UUID) (string, error) {
	return s.statusValue, s.statusErr
}
func (s *timetableEntryStoreStub) listEntries(context.Context, uuid.UUID) ([]TimetableEntry, error) {
	return s.entries, nil
}
func (s *timetableEntryStoreStub) upsertEntry(_ context.Context, _ uuid.UUID, command timetableEntryCommand) error {
	s.upserts = append(s.upserts, command)
	return nil
}
func (s *timetableEntryStoreStub) deleteEntry(context.Context, uuid.UUID, int16, int16) error {
	s.deleted = true
	return nil
}

func TestTimetableEntriesOnlyAllowDrafts(t *testing.T) {
	for _, status := range []string{"under_review", "approved", "published", "archived"} {
		t.Run(status, func(t *testing.T) {
			store := &timetableEntryStoreStub{statusValue: status}
			service := &timetableEntryService{entries: store}
			err := service.save(context.Background(), uuid.New(), saveTimetableEntriesCommand{})
			if !errors.Is(err, errTimetableEntriesDraftOnly) {
				t.Fatalf("save() error = %v, want %v", err, errTimetableEntriesDraftOnly)
			}
			if err := service.delete(context.Background(), uuid.New(), 1, 1); !errors.Is(err, errTimetableEntriesDraftOnly) {
				t.Fatalf("delete() error = %v, want %v", err, errTimetableEntriesDraftOnly)
			}
		})
	}
}

func TestTimetableEntriesSavePassesNullableAssignments(t *testing.T) {
	store := &timetableEntryStoreStub{statusValue: timetableStatusDraft}
	subjectID, teacherID := uuid.New(), uuid.New()
	command := timetableEntryCommand{DayOfWeek: 1, PeriodNumber: 2, SubjectID: &subjectID, TeacherID: &teacherID}
	if err := (&timetableEntryService{entries: store}).save(context.Background(), uuid.New(), saveTimetableEntriesCommand{Entries: []timetableEntryCommand{command}}); err != nil {
		t.Fatal(err)
	}
	if len(store.upserts) != 1 || store.upserts[0].DayOfWeek != 1 || store.upserts[0].PeriodNumber != 2 {
		t.Fatalf("upserts = %+v", store.upserts)
	}
	if store.upserts[0].SubjectID == nil || *store.upserts[0].SubjectID != subjectID || store.upserts[0].TeacherID == nil || *store.upserts[0].TeacherID != teacherID {
		t.Fatalf("nullable assignments were not preserved: %+v", store.upserts[0])
	}
}

func TestTimetableEntriesDeleteDraft(t *testing.T) {
	store := &timetableEntryStoreStub{statusValue: timetableStatusDraft}
	if err := (&timetableEntryService{entries: store}).delete(context.Background(), uuid.New(), 2, 5); err != nil {
		t.Fatal(err)
	}
	if !store.deleted {
		t.Fatal("delete operation was not forwarded")
	}
}
