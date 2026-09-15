package academics

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
)

type subjectStoreStub struct {
	created    normalizedSubjectCommand
	deleteRows int64
	getSubject Subject
	getErr     error
}

func (s *subjectStoreStub) create(_ context.Context, command normalizedSubjectCommand) (Subject, error) {
	s.created = command
	return Subject{}, nil
}
func (s *subjectStoreStub) get(context.Context, uuid.UUID) (Subject, error) {
	return s.getSubject, s.getErr
}
func (s *subjectStoreStub) list(context.Context) ([]Subject, error) { return nil, nil }
func (s *subjectStoreStub) update(context.Context, uuid.UUID, normalizedSubjectCommand) (Subject, error) {
	return Subject{}, nil
}
func (s *subjectStoreStub) delete(context.Context, uuid.UUID) (int64, error) {
	return s.deleteRows, nil
}

func TestCreateSubjectDefaultsMaxMarks(t *testing.T) {
	store := &subjectStoreStub{}
	_, err := (&subjectService{subjects: store}).create(context.Background(), subjectCommand{Name: "Math", Code: "MATH"})
	if err != nil {
		t.Fatal(err)
	}
	if store.created.MaxMarks != 100 {
		t.Fatalf("default max marks = %v, want 100", store.created.MaxMarks)
	}
}

func TestDeleteSubjectConflictClassification(t *testing.T) {
	id := uuid.New()
	tests := []struct {
		name  string
		store *subjectStoreStub
		want  error
	}{
		{name: "deleted", store: &subjectStoreStub{deleteRows: 1}},
		{name: "missing", store: &subjectStoreStub{getErr: errors.New("missing")}, want: errSubjectNotFound},
		{name: "in use", store: &subjectStoreStub{getSubject: Subject{ID: id.String()}}, want: errSubjectInUse},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := (&subjectService{subjects: test.store}).delete(context.Background(), id); !errors.Is(err, test.want) {
				t.Fatalf("delete() error = %v, want %v", err, test.want)
			}
		})
	}
}
