package timetable

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
)

type classroomStoreStub struct {
	created    classroomCommand
	updated    classroomCommand
	deleteRows int64
	deleteErr  error
}

func (s *classroomStoreStub) create(_ context.Context, command classroomCommand) (Classroom, error) {
	s.created = command
	return Classroom{Name: command.Name, RoomType: command.RoomType}, nil
}
func (s *classroomStoreStub) list(context.Context) ([]ClassroomListItem, error) { return nil, nil }
func (s *classroomStoreStub) update(_ context.Context, _ uuid.UUID, command classroomCommand) (Classroom, error) {
	s.updated = command
	return Classroom{Name: command.Name, RoomType: command.RoomType}, nil
}
func (s *classroomStoreStub) delete(context.Context, uuid.UUID) (int64, error) {
	return s.deleteRows, s.deleteErr
}

func TestClassroomLabRequiresSubject(t *testing.T) {
	service := &classroomService{classrooms: &classroomStoreStub{}}
	if _, err := service.create(context.Background(), classroomCommand{Name: "Science Lab", RoomType: "lab"}); !errors.Is(err, errLabRequiresSubject) {
		t.Fatalf("create() error = %v, want %v", err, errLabRequiresSubject)
	}
	if _, err := service.update(context.Background(), uuid.New(), classroomCommand{Name: "Science Lab", RoomType: "lab"}); !errors.Is(err, errLabRequiresSubject) {
		t.Fatalf("update() error = %v, want %v", err, errLabRequiresSubject)
	}
}

func TestClassroomLabWithSubjectIsAccepted(t *testing.T) {
	subjectID := uuid.New()
	store := &classroomStoreStub{}
	command := classroomCommand{Name: "Science Lab", RoomType: "lab", SubjectID: &subjectID}
	if _, err := (&classroomService{classrooms: store}).create(context.Background(), command); err != nil {
		t.Fatal(err)
	}
	if store.created.SubjectID == nil || *store.created.SubjectID != subjectID {
		t.Fatalf("subject ID was not passed to repository: %+v", store.created)
	}
}

func TestDeleteClassroom(t *testing.T) {
	tests := []struct {
		name  string
		store *classroomStoreStub
		want  error
	}{
		{name: "deleted", store: &classroomStoreStub{deleteRows: 1}},
		{name: "missing or referenced", store: &classroomStoreStub{}, want: errClassroomNotFound},
		{name: "repository error", store: &classroomStoreStub{deleteErr: errors.New("database unavailable")}, want: errors.New("database unavailable")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := (&classroomService{classrooms: test.store}).delete(context.Background(), uuid.New())
			if test.want == nil && err != nil {
				t.Fatalf("delete() error = %v", err)
			}
			if test.want != nil && (err == nil || err.Error() != test.want.Error()) {
				t.Fatalf("delete() error = %v, want %v", err, test.want)
			}
		})
	}
}
