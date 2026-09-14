package timetable

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

type timetableCRUDStoreStub struct {
	created       []Timetable
	createdParent *uuid.UUID
	copySource    uuid.UUID
	copyTarget    uuid.UUID
	deleted       int64
	status        string
}

func (s *timetableCRUDStoreStub) create(_ context.Context, yearID, classID, actor uuid.UUID, parentID *uuid.UUID) (Timetable, error) {
	timetable := Timetable{ID: uuid.New(), AcademicYearID: yearID, ClassID: classID, CreatedBy: actor, Status: "draft"}
	s.created = append(s.created, timetable)
	s.createdParent = parentID
	return timetable, nil
}
func (s *timetableCRUDStoreStub) get(context.Context, uuid.UUID) (Timetable, error) {
	return Timetable{ID: uuid.New(), AcademicYearID: uuid.New(), ClassID: uuid.New(), Status: s.status}, nil
}
func (s *timetableCRUDStoreStub) copyEntries(_ context.Context, sourceID, targetID uuid.UUID) error {
	s.copySource, s.copyTarget = sourceID, targetID
	return nil
}
func (s *timetableCRUDStoreStub) listByClass(context.Context, uuid.UUID, uuid.UUID) ([]TimetableListItem, error) {
	return nil, nil
}
func (s *timetableCRUDStoreStub) listByAcademicYear(context.Context, uuid.UUID) ([]TimetableListItem, error) {
	return nil, nil
}
func (s *timetableCRUDStoreStub) deleteDraft(context.Context, uuid.UUID) (int64, error) {
	return s.deleted, nil
}
func (s *timetableCRUDStoreStub) archive(context.Context, uuid.UUID) (Timetable, error) {
	return Timetable{}, nil
}

func TestTimetableCRUDCopyCreatesDraftAndCopiesEntries(t *testing.T) {
	store := &timetableCRUDStoreStub{}
	request := timetableCopyRequest{AcademicYearID: uuid.New(), ClassID: uuid.New(), SourceTimetableID: uuid.New()}
	draft, err := (&timetableCRUDService{store: store}).copy(context.Background(), request, uuid.New())
	if err != nil {
		t.Fatal(err)
	}
	if len(store.created) != 1 || store.copySource != request.SourceTimetableID || store.copyTarget != draft.ID {
		t.Fatalf("copy operation was not forwarded correctly: created=%+v source=%s target=%s", store.created, store.copySource, store.copyTarget)
	}
}

func TestTimetableCRUDReviseOnlyPublishedTimetables(t *testing.T) {
	store := &timetableCRUDStoreStub{status: "draft"}
	if _, err := (&timetableCRUDService{store: store}).revise(context.Background(), uuid.New(), uuid.New()); err == nil {
		t.Fatal("revise() accepted a non-published timetable")
	}

	store.status = statusPublished
	if _, err := (&timetableCRUDService{store: store}).revise(context.Background(), uuid.New(), uuid.New()); err != nil {
		t.Fatal(err)
	}
	if store.createdParent == nil {
		t.Fatal("revision did not link the new draft to the published timetable")
	}
}

func TestTimetableCRUDDeleteReportsMissingDraft(t *testing.T) {
	store := &timetableCRUDStoreStub{}
	if count, err := (&timetableCRUDService{store: store}).store.deleteDraft(context.Background(), uuid.New()); err != nil || count != 0 {
		t.Fatalf("deleteDraft() = %d, %v", count, err)
	}
}

type statusHistoryStoreStub struct {
	history []StatusHistoryItem
}

func (s *statusHistoryStoreStub) listStatusHistory(context.Context, uuid.UUID) ([]StatusHistoryItem, error) {
	return s.history, nil
}

func TestStatusHistoryListUsesModuleOwnedValues(t *testing.T) {
	history := []StatusHistoryItem{{ID: uuid.New(), ToStatus: statusPublished, ChangedByName: "Admin"}}
	store := &statusHistoryStoreStub{history: history}
	got, err := (&statusHistoryService{store: store}).list(context.Background(), uuid.New())
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].ToStatus != statusPublished || got[0].ChangedByName != "Admin" {
		t.Fatalf("history = %+v", got)
	}
}
