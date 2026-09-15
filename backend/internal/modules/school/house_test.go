package school

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/openschool-org/openschool/internal/ports"
)

type houseStoreStub struct {
	deleteRows     int64
	deleteErr      error
	getErr         error
	studentIDs     []uuid.UUID
	pickedHouses   []uuid.UUID
	studentUpdates []uuid.UUID
	studentBefore  ports.StudentHouseProfile
	studentAfter   ports.StudentHouseProfile
}

func (s *houseStoreStub) create(context.Context, houseCommand) (House, error) { return House{}, nil }
func (s *houseStoreStub) get(context.Context, uuid.UUID) (House, error) {
	return House{}, s.getErr
}
func (s *houseStoreStub) list(context.Context) ([]House, error) { return nil, nil }
func (s *houseStoreStub) update(context.Context, uuid.UUID, houseCommand) (House, error) {
	return House{}, nil
}
func (s *houseStoreStub) delete(context.Context, uuid.UUID) (int64, error) {
	return s.deleteRows, s.deleteErr
}
func (s *houseStoreStub) pickForStudent(context.Context) (uuid.UUID, bool, error) {
	if len(s.pickedHouses) == 0 {
		return uuid.Nil, false, nil
	}
	id := s.pickedHouses[0]
	s.pickedHouses = s.pickedHouses[1:]
	return id, true, nil
}
func (s *houseStoreStub) pickForTeacher(context.Context) (uuid.UUID, bool, error) {
	return uuid.Nil, false, nil
}
func (s *houseStoreStub) listStudentsMissingHouse(context.Context) ([]uuid.UUID, error) {
	return s.studentIDs, nil
}
func (s *houseStoreStub) listTeachersMissingHouse(context.Context) ([]uuid.UUID, error) {
	return nil, nil
}
func (s *houseStoreStub) getStudent(context.Context, uuid.UUID) (ports.StudentHouseProfile, error) {
	return s.studentBefore, nil
}
func (s *houseStoreStub) getTeacher(context.Context, uuid.UUID) (ports.TeacherHouseProfile, error) {
	return ports.TeacherHouseProfile{}, nil
}
func (s *houseStoreStub) updateStudentHouse(_ context.Context, studentID uuid.UUID, _ *uuid.UUID) (ports.StudentHouseProfile, error) {
	s.studentUpdates = append(s.studentUpdates, studentID)
	return s.studentAfter, nil
}
func (s *houseStoreStub) updateTeacherHouse(context.Context, uuid.UUID, *uuid.UUID) (ports.TeacherHouseProfile, error) {
	return ports.TeacherHouseProfile{}, nil
}

type auditRecorderStub struct {
	entityType string
	action     string
	actorID    uuid.UUID
}

func (a *auditRecorderStub) Record(_ context.Context, entityType string, _ uuid.UUID, action string, actorID uuid.UUID, _, _ interface{}, _ string) error {
	a.entityType, a.action, a.actorID = entityType, action, actorID
	return nil
}

func TestDeleteHouse(t *testing.T) {
	tests := []struct {
		name  string
		store *houseStoreStub
		want  error
	}{
		{name: "deleted", store: &houseStoreStub{deleteRows: 1}},
		{name: "not found", store: &houseStoreStub{getErr: errors.New("missing")}, want: errHouseNotFound},
		{name: "in use", store: &houseStoreStub{}, want: errHouseInUse},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := (&HouseService{houses: test.store}).delete(context.Background(), uuid.New())
			if !errors.Is(err, test.want) {
				t.Fatalf("delete() error = %v, want %v", err, test.want)
			}
		})
	}
}

func TestReassignMissingStudents(t *testing.T) {
	studentIDs := []uuid.UUID{uuid.New(), uuid.New()}
	store := &houseStoreStub{studentIDs: studentIDs, pickedHouses: []uuid.UUID{uuid.New(), uuid.New()}}
	assigned, err := (&HouseService{houses: store}).reassignMissing(context.Background())
	if err != nil || assigned != 2 {
		t.Fatalf("reassignMissing() = %d, %v; want 2, nil", assigned, err)
	}
	if len(store.studentUpdates) != 2 || store.studentUpdates[0] != studentIDs[0] || store.studentUpdates[1] != studentIDs[1] {
		t.Fatalf("updated students = %v, want %v", store.studentUpdates, studentIDs)
	}
}

func TestChangeStudentHouseRecordsAudit(t *testing.T) {
	oldHouseID, newHouseID, actorID := uuid.New(), uuid.New(), uuid.New()
	store := &houseStoreStub{
		studentBefore: ports.StudentHouseProfile{HouseID: pgtype.UUID{Bytes: oldHouseID, Valid: true}},
		studentAfter:  ports.StudentHouseProfile{HouseID: pgtype.UUID{Bytes: newHouseID, Valid: true}},
	}
	audit := &auditRecorderStub{}
	updated, err := (&HouseService{houses: store, audit: audit}).ChangeStudentHouse(context.Background(), uuid.New(), newHouseID.String(), actorID)
	if err != nil {
		t.Fatal(err)
	}
	if !updated.HouseID.Valid || uuid.UUID(updated.HouseID.Bytes) != newHouseID {
		t.Fatalf("updated house = %v, want %s", updated.HouseID, newHouseID)
	}
	if audit.entityType != "student_house" || audit.action != "house_changed" || audit.actorID != actorID {
		t.Fatalf("unexpected audit record: %#v", audit)
	}
}

func TestChangeStudentHouseRejectsInvalidID(t *testing.T) {
	_, err := (&HouseService{houses: &houseStoreStub{}}).ChangeStudentHouse(context.Background(), uuid.New(), "not-a-uuid", uuid.New())
	if err == nil || err.Error() != "invalid house id" {
		t.Fatalf("ChangeStudentHouse() error = %v, want invalid house id", err)
	}
}

var _ ports.HouseAssignments = (*HouseService)(nil)
