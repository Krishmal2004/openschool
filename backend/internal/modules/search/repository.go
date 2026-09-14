package search

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/openschool-org/openschool/db/sqlc"
)

type Repository struct{ queries *db.Queries }

func NewRepository(pool *pgxpool.Pool) *Repository { return &Repository{queries: db.New(pool)} }

func searchTerm(term string) pgtype.Text { return pgtype.Text{String: term, Valid: true} }

func (r *Repository) students(ctx context.Context, term string) ([]candidate, error) {
	rows, err := r.queries.SearchStudents(ctx, searchTerm(term))
	if err != nil {
		return nil, err
	}
	result := make([]candidate, len(rows))
	for i, value := range rows {
		result[i] = candidate{ID: value.ID, Name: value.FullName, Subtitle: value.IndexNumber}
	}
	return result, nil
}

func (r *Repository) teachers(ctx context.Context, term string) ([]candidate, error) {
	rows, err := r.queries.SearchTeachers(ctx, searchTerm(term))
	if err != nil {
		return nil, err
	}
	result := make([]candidate, len(rows))
	for i, value := range rows {
		result[i] = candidate{ID: value.ID, Name: value.FullName, Subtitle: value.EmployeeNumber}
	}
	return result, nil
}

func (r *Repository) guardians(ctx context.Context, term string) ([]candidate, error) {
	rows, err := r.queries.SearchGuardians(ctx, searchTerm(term))
	if err != nil {
		return nil, err
	}
	result := make([]candidate, len(rows))
	for i, value := range rows {
		result[i] = candidate{ID: value.ID, Name: value.FullName, Subtitle: value.Phone}
	}
	return result, nil
}

func (r *Repository) nonAcademicStaff(ctx context.Context, term string) ([]candidate, error) {
	rows, err := r.queries.SearchNonAcademicStaff(ctx, searchTerm(term))
	if err != nil {
		return nil, err
	}
	result := make([]candidate, len(rows))
	for i, value := range rows {
		result[i] = candidate{ID: value.ID, Name: value.FullName, Subtitle: value.EmployeeNumber}
	}
	return result, nil
}
