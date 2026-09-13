package school

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
)

type termStoreStub struct {
	created    termValues
	getTerm    Term
	getErr     error
	deleteRows int64
	setCalled  bool
}

func (s *termStoreStub) create(_ context.Context, values termValues) (Term, error) {
	s.created = values
	return Term{}, nil
}
func (s *termStoreStub) get(context.Context, uuid.UUID) (Term, error)    { return s.getTerm, s.getErr }
func (s *termStoreStub) list(context.Context, uuid.UUID) ([]Term, error) { return nil, nil }
func (s *termStoreStub) current(context.Context) (Term, error)           { return Term{}, nil }
func (s *termStoreStub) setCurrent(context.Context, uuid.UUID) error     { s.setCalled = true; return nil }
func (s *termStoreStub) update(context.Context, uuid.UUID, termValues) (Term, error) {
	return Term{}, nil
}
func (s *termStoreStub) delete(context.Context, uuid.UUID) (int64, error) { return s.deleteRows, nil }

func TestCreateTermRejectsInvalidAcademicYear(t *testing.T) {
	_, err := (&termService{terms: &termStoreStub{}}).create(context.Background(), createTermCommand{AcademicYearID: "bad-id"})
	if err == nil || err.Error() != "invalid academic_year_id" {
		t.Fatalf("create() error = %v", err)
	}
}

func TestSetCurrentRequiresExistingTerm(t *testing.T) {
	store := &termStoreStub{getErr: errors.New("missing")}
	err := (&termService{terms: store}).setCurrent(context.Background(), uuid.New())
	if !errors.Is(err, errTermNotFound) || store.setCalled {
		t.Fatalf("setCurrent() error = %v, called = %v", err, store.setCalled)
	}
}

func TestDeleteTermClassifiesMissingAndInUse(t *testing.T) {
	id := uuid.New()
	for _, test := range []struct {
		name  string
		store *termStoreStub
		want  error
	}{
		{name: "deleted", store: &termStoreStub{deleteRows: 1}},
		{name: "missing", store: &termStoreStub{getErr: errors.New("missing")}, want: errTermNotFound},
		{name: "in use", store: &termStoreStub{getTerm: Term{ID: id.String()}}, want: errTermInUse},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := (&termService{terms: test.store}).delete(context.Background(), id)
			if !errors.Is(err, test.want) {
				t.Fatalf("delete() error = %v, want %v", err, test.want)
			}
		})
	}
}
