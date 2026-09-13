package school

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
)

type gradeStoreStub struct {
	deleteRows int64
	deleteErr  error
	getGrade   Grade
	getErr     error
}

func (s *gradeStoreStub) create(context.Context, gradeCommand) (Grade, error) { return Grade{}, nil }
func (s *gradeStoreStub) get(context.Context, uuid.UUID) (Grade, error) {
	return s.getGrade, s.getErr
}
func (s *gradeStoreStub) list(context.Context) ([]Grade, error) { return nil, nil }
func (s *gradeStoreStub) update(context.Context, uuid.UUID, gradeCommand) (Grade, error) {
	return Grade{}, nil
}
func (s *gradeStoreStub) delete(context.Context, uuid.UUID) (int64, error) {
	return s.deleteRows, s.deleteErr
}

func TestDeleteGrade(t *testing.T) {
	tests := []struct {
		name  string
		store *gradeStoreStub
		want  error
	}{
		{name: "deleted", store: &gradeStoreStub{deleteRows: 1}},
		{name: "not found", store: &gradeStoreStub{getErr: errors.New("missing")}, want: errGradeNotFound},
		{name: "in use", store: &gradeStoreStub{getGrade: Grade{ID: uuid.New()}}, want: errGradeInUse},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := (&gradeService{grades: test.store}).delete(context.Background(), uuid.New())
			if !errors.Is(err, test.want) {
				t.Fatalf("delete() error = %v, want %v", err, test.want)
			}
		})
	}
}
