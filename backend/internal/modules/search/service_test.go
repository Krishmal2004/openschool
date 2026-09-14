package search

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
)

type searchStore struct {
	studentsResult  []candidate
	teachersResult  []candidate
	guardiansResult []candidate
	staffResult     []candidate
	teachersErr     error
	terms           []string
}

func (s *searchStore) students(_ context.Context, term string) ([]candidate, error) {
	s.terms = append(s.terms, term)
	return s.studentsResult, nil
}

func (s *searchStore) teachers(_ context.Context, term string) ([]candidate, error) {
	s.terms = append(s.terms, term)
	return s.teachersResult, s.teachersErr
}

func (s *searchStore) guardians(_ context.Context, term string) ([]candidate, error) {
	s.terms = append(s.terms, term)
	return s.guardiansResult, nil
}

func (s *searchStore) nonAcademicStaff(_ context.Context, term string) ([]candidate, error) {
	s.terms = append(s.terms, term)
	return s.staffResult, nil
}

func TestGlobalEmptyQueryReturnsEmptyArraysWithoutSearching(t *testing.T) {
	store := &searchStore{}
	result, err := NewService(store).Global(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}
	if len(store.terms) != 0 {
		t.Fatalf("empty query reached repository: %v", store.terms)
	}
	data, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), ":null") {
		t.Fatalf("search response contains null collection: %s", data)
	}
}

func TestGlobalMapsCategorizedResults(t *testing.T) {
	studentID, teacherID, guardianID, staffID := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	store := &searchStore{
		studentsResult:  []candidate{{ID: studentID, Name: "Student", Subtitle: "S001"}},
		teachersResult:  []candidate{{ID: teacherID, Name: "Teacher", Subtitle: "T001"}},
		guardiansResult: []candidate{{ID: guardianID, Name: "Guardian", Subtitle: "0700000000"}},
		staffResult:     []candidate{{ID: staffID, Name: "Staff", Subtitle: "N001"}},
	}
	result, err := NewService(store).Global(context.Background(), "test")
	if err != nil {
		t.Fatal(err)
	}
	if len(store.terms) != 4 || result.Students[0].ID != studentID.String() || result.Teachers[0].Subtitle != "T001" || result.Guardians[0].ID != guardianID.String() || result.NonAcademicStaff[0].ID != staffID.String() {
		t.Fatalf("unexpected search response: %#v; terms=%v", result, store.terms)
	}
}

func TestGlobalStopsAfterRepositoryError(t *testing.T) {
	want := errors.New("search unavailable")
	store := &searchStore{teachersErr: want}
	_, err := NewService(store).Global(context.Background(), "name")
	if !errors.Is(err, want) {
		t.Fatalf("error=%v, want %v", err, want)
	}
	if len(store.terms) != 2 {
		t.Fatalf("search continued after error: %v", store.terms)
	}
}
