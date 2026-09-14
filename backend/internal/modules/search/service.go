// Package search owns cross-entity administrative search.
package search

import (
	"context"

	"github.com/google/uuid"
)

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
	if term == "" {
		return response, nil
	}

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
