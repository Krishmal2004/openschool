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

type classRepository struct{ queries *db.Queries }

func newClassRepository(pool *pgxpool.Pool) *classRepository {
	return &classRepository{queries: db.New(pool)}
}

func classUUID(v pgtype.UUID) *uuid.UUID {
	if !v.Valid {
		return nil
	}
	id := uuid.UUID(v.Bytes)
	return &id
}
func classText(v pgtype.Text) *string {
	if !v.Valid {
		return nil
	}
	value := v.String
	return &value
}
func mapClass(v db.Class) Class {
	return Class{ID: v.ID, GradeID: v.GradeID, AcademicYearID: v.AcademicYearID, FormTeacherID: classUUID(v.FormTeacherID), StreamID: classUUID(v.StreamID), StreamGroupID: classUUID(v.StreamGroupID), Name: v.Name, CreatedAt: v.CreatedAt.Time.String(), GirlMonitorID: classUUID(v.GirlMonitorID), BoyMonitorID: classUUID(v.BoyMonitorID), MediumID: classUUID(v.MediumID), HomeClassroomID: classUUID(v.HomeClassroomID)}
}
func mapDetails(id, grade, year uuid.UUID, form, stream, streamGroup, girl, boy, medium, room pgtype.UUID, name string, created pgtype.Timestamptz, gradeName, yearLabel string, mediumName, roomName pgtype.Text) ClassDetails {
	return ClassDetails{Class: Class{ID: id, GradeID: grade, AcademicYearID: year, FormTeacherID: classUUID(form), StreamID: classUUID(stream), StreamGroupID: classUUID(streamGroup), GirlMonitorID: classUUID(girl), BoyMonitorID: classUUID(boy), MediumID: classUUID(medium), HomeClassroomID: classUUID(room), Name: name, CreatedAt: created.Time.String()}, GradeName: gradeName, AcademicYearLabel: yearLabel, MediumName: classText(mediumName), HomeClassroomName: classText(roomName)}
}
func classUUIDParam(v *uuid.UUID) pgtype.UUID {
	if v == nil {
		return pgtype.UUID{}
	}
	return pgtype.UUID{Bytes: *v, Valid: true}
}

func (r *classRepository) create(ctx context.Context, v createClassRequest) (Class, error) {
	row, e := r.queries.CreateClass(ctx, db.CreateClassParams{GradeID: v.GradeID, AcademicYearID: v.AcademicYearID, Name: v.Name, FormTeacherID: classUUIDParam(v.FormTeacherID), StreamID: classUUIDParam(v.StreamID), StreamGroupID: classUUIDParam(v.StreamGroupID), MediumID: classUUIDParam(v.MediumID), HomeClassroomID: classUUIDParam(v.HomeClassroomID)})
	return mapClass(row), e
}
func (r *classRepository) get(ctx context.Context, id uuid.UUID) (Class, error) {
	v, e := r.queries.GetClassByID(ctx, id)
	return mapClass(v), e
}
func (r *classRepository) listCurrent(ctx context.Context) ([]ClassDetails, error) {
	rows, e := r.queries.ListCurrentClasses(ctx)
	if e != nil {
		return nil, e
	}
	out := make([]ClassDetails, len(rows))
	for i, v := range rows {
		out[i] = mapDetails(v.ID, v.GradeID, v.AcademicYearID, v.FormTeacherID, v.StreamID, v.StreamGroupID, v.GirlMonitorID, v.BoyMonitorID, v.MediumID, v.HomeClassroomID, v.Name, v.CreatedAt, v.GradeName, v.AcademicYearLabel, v.MediumName, v.HomeClassroomName)
	}
	return out, nil
}
func (r *classRepository) listByYear(ctx context.Context, id uuid.UUID) ([]ClassDetails, error) {
	rows, e := r.queries.ListClassesByAcademicYear(ctx, id)
	if e != nil {
		return nil, e
	}
	out := make([]ClassDetails, len(rows))
	for i, v := range rows {
		out[i] = mapDetails(v.ID, v.GradeID, v.AcademicYearID, v.FormTeacherID, v.StreamID, v.StreamGroupID, v.GirlMonitorID, v.BoyMonitorID, v.MediumID, v.HomeClassroomID, v.Name, v.CreatedAt, v.GradeName, v.AcademicYearLabel, v.MediumName, v.HomeClassroomName)
	}
	return out, nil
}
func (r *classRepository) update(ctx context.Context, id uuid.UUID, v updateClassRequest) (Class, error) {
	row, e := r.queries.UpdateClass(ctx, db.UpdateClassParams{ID: id, Name: v.Name, FormTeacherID: classUUIDParam(v.FormTeacherID), MediumID: classUUIDParam(v.MediumID), HomeClassroomID: classUUIDParam(v.HomeClassroomID)})
	return mapClass(row), e
}
func (r *classRepository) delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteClass(ctx, id)
}
func (r *classRepository) studentCount(ctx context.Context, id uuid.UUID) (int64, error) {
	return r.queries.GetClassStudentCount(ctx, id)
}
func (r *classRepository) assignFormTeacher(ctx context.Context, id, teacher uuid.UUID) (Class, error) {
	row, e := r.queries.AssignFormTeacher(ctx, db.AssignFormTeacherParams{ID: id, FormTeacherID: classUUIDParam(&teacher)})
	return mapClass(row), e
}
func (r *classRepository) assignMonitors(ctx context.Context, id uuid.UUID, girl, boy *uuid.UUID) (Class, error) {
	row, e := r.queries.AssignClassMonitors(ctx, db.AssignClassMonitorsParams{ID: id, GirlMonitorID: classUUIDParam(girl), BoyMonitorID: classUUIDParam(boy)})
	return mapClass(row), e
}
func (r *classRepository) qualified(ctx context.Context, teacher, subject uuid.UUID) (bool, error) {
	rows, e := r.queries.ListSubjectsByTeacher(ctx, teacher)
	if e != nil {
		return false, e
	}
	for _, v := range rows {
		if v.ID == subject {
			return true, nil
		}
	}
	return false, nil
}
func (r *classRepository) assignSubjectTeacher(ctx context.Context, id, subject, teacher uuid.UUID) error {
	return r.queries.AssignSubjectTeacherToClass(ctx, db.AssignSubjectTeacherToClassParams{ClassID: id, SubjectID: subject, TeacherID: teacher})
}
func (r *classRepository) listSubjectTeachers(ctx context.Context, id uuid.UUID) ([]SubjectTeacher, error) {
	rows, e := r.queries.ListSubjectTeachersByClass(ctx, id)
	if e != nil {
		return nil, e
	}
	out := make([]SubjectTeacher, len(rows))
	for i, v := range rows {
		out[i] = SubjectTeacher{SubjectID: v.SubjectID, SubjectName: v.SubjectName, SubjectCode: v.SubjectCode, TeacherID: v.TeacherID, TeacherName: v.TeacherName}
	}
	return out, nil
}
func (r *classRepository) enroll(ctx context.Context, id, student uuid.UUID) error {
	return r.queries.EnrollStudentInClass(ctx, db.EnrollStudentInClassParams{ClassID: id, StudentID: student})
}
func (r *classRepository) unenroll(ctx context.Context, id, student uuid.UUID) error {
	return r.queries.UnenrollStudentFromClass(ctx, db.UnenrollStudentFromClassParams{ClassID: id, StudentID: student})
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
