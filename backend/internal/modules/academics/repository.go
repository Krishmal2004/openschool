package academics

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/openschool-org/openschool/db/sqlc"
)

type subjectRepository struct{ queries *db.Queries }

func newSubjectRepository(pool *pgxpool.Pool) *subjectRepository {
	return &subjectRepository{queries: db.New(pool)}
}

func (r *subjectRepository) create(ctx context.Context, command normalizedSubjectCommand) (Subject, error) {
	params, err := subjectParams(command)
	if err != nil {
		return Subject{}, err
	}
	row, err := r.queries.CreateSubject(ctx, db.CreateSubjectParams{Name: params.Name, Code: params.Code, Type: params.Type, MaxMarks: params.MaxMarks})
	return mapSubject(row), err
}
func (r *subjectRepository) get(ctx context.Context, id uuid.UUID) (Subject, error) {
	row, err := r.queries.GetSubjectByID(ctx, id)
	return mapSubject(row), err
}
func (r *subjectRepository) list(ctx context.Context) ([]Subject, error) {
	rows, err := r.queries.ListSubjects(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]Subject, len(rows))
	for i, row := range rows {
		result[i] = mapSubject(row)
	}
	return result, nil
}
func (r *subjectRepository) update(ctx context.Context, id uuid.UUID, command normalizedSubjectCommand) (Subject, error) {
	params, err := subjectParams(command)
	if err != nil {
		return Subject{}, err
	}
	row, err := r.queries.UpdateSubject(ctx, db.UpdateSubjectParams{ID: id, Name: params.Name, Code: params.Code, Type: params.Type, MaxMarks: params.MaxMarks})
	return mapSubject(row), err
}
func (r *subjectRepository) delete(ctx context.Context, id uuid.UUID) (int64, error) {
	return r.queries.DeleteSubject(ctx, id)
}

func subjectParams(command normalizedSubjectCommand) (db.CreateSubjectParams, error) {
	var numeric pgtype.Numeric
	if err := numeric.Scan(strconv.FormatFloat(command.MaxMarks, 'f', 2, 64)); err != nil {
		return db.CreateSubjectParams{}, fmt.Errorf("failed to convert %v to numeric: %w", command.MaxMarks, err)
	}
	return db.CreateSubjectParams{Name: command.Name, Code: command.Code, Type: pgtype.Text{String: command.Type, Valid: command.Type != ""}, MaxMarks: numeric}, nil
}

func mapSubject(row db.Subject) Subject {
	var subjectType *string
	if row.Type.Valid {
		value := row.Type.String
		subjectType = &value
	}
	maxMarks := float64(0)
	if value, err := row.MaxMarks.Float64Value(); err == nil && value.Valid {
		maxMarks = value.Float64
	}
	return Subject{ID: row.ID.String(), Name: row.Name, Code: row.Code, Type: subjectType, MaxMarks: maxMarks, CreatedAt: row.CreatedAt.Time.String()}
}

type streamRepository struct{ queries *db.Queries }

func newStreamRepository(pool *pgxpool.Pool) *streamRepository {
	return &streamRepository{queries: db.New(pool)}
}
func (r *streamRepository) createStream(ctx context.Context, name string) (Stream, error) {
	row, err := r.queries.CreateStream(ctx, name)
	return mapStream(row), err
}
func (r *streamRepository) getStream(ctx context.Context, id uuid.UUID) (Stream, error) {
	row, err := r.queries.GetStreamByID(ctx, id)
	return mapStream(row), err
}
func (r *streamRepository) listStreams(ctx context.Context) ([]Stream, error) {
	rows, err := r.queries.ListStreams(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]Stream, len(rows))
	for i, row := range rows {
		result[i] = mapStream(row)
	}
	return result, nil
}
func (r *streamRepository) updateStream(ctx context.Context, id uuid.UUID, name string) (Stream, error) {
	row, err := r.queries.UpdateStream(ctx, db.UpdateStreamParams{ID: id, Name: name})
	return mapStream(row), err
}
func (r *streamRepository) deleteStream(ctx context.Context, id uuid.UUID) (int64, error) {
	return r.queries.DeleteStream(ctx, id)
}
func (r *streamRepository) createGroup(ctx context.Context, streamID uuid.UUID, name string) (StreamGroup, error) {
	row, err := r.queries.CreateStreamGroup(ctx, db.CreateStreamGroupParams{StreamID: streamID, Name: name})
	return mapStreamGroup(row), err
}
func (r *streamRepository) getGroup(ctx context.Context, id uuid.UUID) (StreamGroup, error) {
	row, err := r.queries.GetStreamGroupByID(ctx, id)
	return mapStreamGroup(row), err
}
func (r *streamRepository) listGroups(ctx context.Context, streamID uuid.UUID) ([]StreamGroup, error) {
	rows, err := r.queries.ListStreamGroupsByStream(ctx, streamID)
	if err != nil {
		return nil, err
	}
	result := make([]StreamGroup, len(rows))
	for i, row := range rows {
		result[i] = mapStreamGroup(row)
	}
	return result, nil
}
func (r *streamRepository) updateGroup(ctx context.Context, id uuid.UUID, name string) (StreamGroup, error) {
	row, err := r.queries.UpdateStreamGroup(ctx, db.UpdateStreamGroupParams{ID: id, Name: name})
	return mapStreamGroup(row), err
}
func (r *streamRepository) deleteGroup(ctx context.Context, id uuid.UUID) (int64, error) {
	return r.queries.DeleteStreamGroup(ctx, id)
}
func mapStream(row db.Stream) Stream {
	return Stream{ID: row.ID, Name: row.Name, CreatedAt: row.CreatedAt.Time}
}
func mapStreamGroup(row db.StreamGroup) StreamGroup {
	return StreamGroup{ID: row.ID, StreamID: row.StreamID, Name: row.Name, CreatedAt: row.CreatedAt.Time}
}
