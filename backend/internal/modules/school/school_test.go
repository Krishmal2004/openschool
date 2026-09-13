package school

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
)

type schoolStoreStub struct {
	school     School
	schoolErr  error
	year       AcademicYear
	yearErr    error
	deleteRows int64
	created    schoolValues
}

func (s *schoolStoreStub) createSchool(_ context.Context, values schoolValues) (School, error) {
	s.created = values
	return School{}, nil
}
func (s *schoolStoreStub) getSchool(context.Context) (School, error) { return s.school, s.schoolErr }
func (s *schoolStoreStub) updateSchool(context.Context, schoolValues) (School, error) {
	return School{}, nil
}
func (s *schoolStoreStub) createYear(context.Context, yearValues) (AcademicYear, error) {
	return AcademicYear{}, nil
}
func (s *schoolStoreStub) getYear(context.Context, uuid.UUID) (AcademicYear, error) {
	return s.year, s.yearErr
}
func (s *schoolStoreStub) currentYear(context.Context) (AcademicYear, error) {
	return AcademicYear{}, nil
}
func (s *schoolStoreStub) listYears(context.Context) ([]AcademicYear, error) { return nil, nil }
func (s *schoolStoreStub) setCurrentYear(context.Context, uuid.UUID) error   { return nil }
func (s *schoolStoreStub) deleteYear(context.Context, uuid.UUID) (int64, error) {
	return s.deleteRows, nil
}

func TestCreateSchoolInvariants(t *testing.T) {
	service := &schoolService{store: &schoolStoreStub{school: School{ID: uuid.NewString()}}}
	if _, err := service.createSchool(context.Background(), schoolCommand{Name: "School"}); err == nil || err.Error() != "school already exists" {
		t.Fatalf("existing school error = %v", err)
	}

	store := &schoolStoreStub{schoolErr: errors.New("missing")}
	service = &schoolService{store: store}
	if _, err := service.createSchool(context.Background(), schoolCommand{Name: "School", LogoURL: "https://example.test/logo.png"}); !errors.Is(err, errInvalidLogoURL) {
		t.Fatalf("logo error = %v", err)
	}
	if _, err := service.createSchool(context.Background(), schoolCommand{Name: "School"}); err != nil {
		t.Fatal(err)
	}
	if store.created.SchoolType != "mixed" {
		t.Fatalf("default school type = %q", store.created.SchoolType)
	}
}

func TestAcademicYearDeletionClassification(t *testing.T) {
	id := uuid.New()
	for _, test := range []struct {
		name  string
		store *schoolStoreStub
		want  error
	}{
		{name: "deleted", store: &schoolStoreStub{deleteRows: 1}},
		{name: "missing", store: &schoolStoreStub{yearErr: errors.New("missing")}, want: errAcademicYearNotFound},
		{name: "in use", store: &schoolStoreStub{year: AcademicYear{ID: id.String()}}, want: errAcademicYearInUse},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := (&schoolService{store: test.store}).deleteYear(context.Background(), id)
			if !errors.Is(err, test.want) {
				t.Fatalf("delete error = %v, want %v", err, test.want)
			}
		})
	}
}
