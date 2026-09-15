package reports

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/openschool-org/openschool/db/sqlc"
)

type Repository struct{ queries *db.Queries }

func NewRepository(pool *pgxpool.Pool) *Repository { return &Repository{queries: db.New(pool)} }

func (r *Repository) className(ctx context.Context, id uuid.UUID) (string, error) {
	value, err := r.queries.GetClassByID(ctx, id)
	return value.Name, err
}

func (r *Repository) termName(ctx context.Context, id uuid.UUID) (string, error) {
	value, err := r.queries.GetTermByID(ctx, id)
	return value.Name, err
}

func (r *Repository) subjectName(ctx context.Context, id uuid.UUID) (string, error) {
	value, err := r.queries.GetSubjectByID(ctx, id)
	return value.Name, err
}

func (r *Repository) classMarks(ctx context.Context, classID, termID, subjectID uuid.UUID) ([]markRow, error) {
	rows, err := r.queries.ListClassMarksForTermSubject(ctx, db.ListClassMarksForTermSubjectParams{ClassID: classID, TermID: termID, SubjectID: subjectID})
	if err != nil {
		return nil, err
	}
	result := make([]markRow, len(rows))
	for i, value := range rows {
		result[i] = markRow{StudentName: value.StudentName, IndexNumber: value.IndexNumber, TermMarkID: value.TermMarkID, Marks: value.Marks, MaxMarks: value.MaxMarks, IsAbsent: value.IsAbsent}
	}
	return result, nil
}
