// Package search owns cross-entity administrative search.
package search

import (
	"context"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
)

// minTermLength and maxTermLength bound the search term server-side (S6):
// the frontend already enforces a minimum, but nothing stopped a direct API
// call sending "" or "%", which ILIKE '%<term>%' turns into a four-table
// scan. maxTermLength keeps an oversized query out of the log/index.
const (
	minTermLength = 2
	maxTermLength = 64
)

// escapeLikeTerm neutralises ILIKE's own wildcards in user input: '%' and
// '_' would otherwise let a caller widen their own match arbitrarily (or,
// as with a bare '%', match every row), and a literal backslash would
// escape the next character instead of being matched literally (S6).
func escapeLikeTerm(term string) string {
	replacer := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return replacer.Replace(term)
}

type candidate struct {
	ID       uuid.UUID
	Name     string
	Subtitle string
}

type store interface {
	students(context.Context, string) ([]candidate, error)
	teachers(context.Context, string) ([]candidate, error)
	guardians(context.Context, string) ([]candidate, error)
	nonAcademicStaff(context.Context, string) ([]candidate, error)
}

type Service struct{ store store }

func NewService(store store) *Service { return &Service{store: store} }

func (s *Service) Global(ctx context.Context, term string) (GlobalSearchResponse, error) {
	response := GlobalSearchResponse{
		Students: []SearchResultItem{}, Teachers: []SearchResultItem{},
		Guardians: []SearchResultItem{}, NonAcademicStaff: []SearchResultItem{},
	}
	term = strings.TrimSpace(term)
	// Count characters, not bytes: a 22-character Sinhala/Chinese term can
	// exceed 64 UTF-8 bytes while still being a reasonable search, and a
	// single multi-byte character shouldn't be able to satisfy a
	// 2-character minimum.
	termLength := utf8.RuneCountInString(term)
	if termLength < minTermLength || termLength > maxTermLength {
		return response, nil
	}
	term = escapeLikeTerm(term)

	students, err := s.store.students(ctx, term)
	if err != nil {
		return response, err
	}
	response.Students = mapCandidates(students)
	teachers, err := s.store.teachers(ctx, term)
	if err != nil {
		return response, err
	}
	response.Teachers = mapCandidates(teachers)
	guardians, err := s.store.guardians(ctx, term)
	if err != nil {
		return response, err
	}
	response.Guardians = mapCandidates(guardians)
	staff, err := s.store.nonAcademicStaff(ctx, term)
	if err != nil {
		return response, err
	}
	response.NonAcademicStaff = mapCandidates(staff)
	return response, nil
}

func mapCandidates(values []candidate) []SearchResultItem {
	result := make([]SearchResultItem, len(values))
	for i, value := range values {
		result[i] = SearchResultItem{ID: value.ID.String(), Name: value.Name, Subtitle: value.Subtitle}
	}
	return result
}
