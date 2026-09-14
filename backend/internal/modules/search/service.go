// Package search owns cross-entity administrative search.
package search

import (
	"context"

	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/models"
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

func (s *Service) Global(ctx context.Context, term string) (models.GlobalSearchResponse, error) {
	response := models.GlobalSearchResponse{
		Students: []models.SearchResultItem{}, Teachers: []models.SearchResultItem{},
		Guardians: []models.SearchResultItem{}, NonAcademicStaff: []models.SearchResultItem{},
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

func mapCandidates(values []candidate) []models.SearchResultItem {
	result := make([]models.SearchResultItem, len(values))
	for i, value := range values {
		result[i] = models.SearchResultItem{ID: value.ID.String(), Name: value.Name, Subtitle: value.Subtitle}
	}
	return result
}
