package attendance

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/openschool-org/openschool/db/sqlc"
)

type MarkInput struct {
	StudentID uuid.UUID
	Status    string
	Note      string
}

type Repository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool, queries: db.New(pool)}
}

func (r *Repository) createSession(ctx context.Context, classID, takenBy uuid.UUID, date time.Time) (Session, error) {
	row, err := r.queries.CreateAttendanceSession(ctx, db.CreateAttendanceSessionParams{ClassID: classID, TakenBy: takenBy, Date: dateValue(date)})
	return mapSession(row), err
}
func (r *Repository) getSession(ctx context.Context, id uuid.UUID) (Session, error) {
	row, err := r.queries.GetAttendanceSessionByID(ctx, id)
	return mapSession(row), err
}
func (r *Repository) findSession(ctx context.Context, classID uuid.UUID, date time.Time) (Session, error) {
	row, err := r.queries.GetAttendanceSessionByClassAndDate(ctx, db.GetAttendanceSessionByClassAndDateParams{ClassID: classID, Date: dateValue(date)})
	return mapSession(row), err
}
func (r *Repository) listSessionsByClass(ctx context.Context, classID uuid.UUID) ([]Session, error) {
	rows, err := r.queries.ListAttendanceSessionsByClass(ctx, classID)
	if err != nil {
		return nil, err
	}
	out := make([]Session, len(rows))
	for i, row := range rows {
		out[i] = mapSession(row)
	}
	return out, nil
}
func (r *Repository) listSessionsByDate(ctx context.Context, date time.Time, gradeIDs []uuid.UUID) ([]DailySession, error) {
	rows, err := r.queries.ListAttendanceSessionsByDate(ctx, db.ListAttendanceSessionsByDateParams{Date: dateValue(date), GradeIds: gradeIDs})
	if err != nil {
		return nil, err
	}
	out := make([]DailySession, len(rows))
	for i, row := range rows {
		out[i] = DailySession{
			Session:       Session{ID: row.ID, ClassID: row.ClassID, TakenBy: row.TakenBy, Date: row.Date, CreatedAt: row.CreatedAt},
			ClassName:     row.ClassName,
			GradeName:     row.GradeName,
			TeacherName:   row.TeacherName,
			EnrolledCount: row.EnrolledCount,
			MarkedCount:   row.MarkedCount,
		}
	}
	return out, nil
}
func (r *Repository) deleteSession(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteAttendanceSession(ctx, id)
}
func (r *Repository) listRecords(ctx context.Context, sessionID uuid.UUID) ([]Record, error) {
	rows, err := r.queries.ListAttendanceRecordsBySession(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	out := make([]Record, len(rows))
	for i, row := range rows {
		out[i] = mapRecord(row)
	}
	return out, nil
}
func (r *Repository) markBatch(ctx context.Context, sessionID uuid.UUID, records []MarkInput) ([]Record, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	qtx := r.queries.WithTx(tx)
	out := make([]Record, len(records))
	for i, record := range records {
		row, writeErr := qtx.MarkAttendance(ctx, db.MarkAttendanceParams{SessionID: sessionID, StudentID: record.StudentID, Status: record.Status, Note: textValue(record.Note)})
		err = writeErr
		if err != nil {
			return nil, err
		}
		out[i] = mapRecord(row)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return out, nil
}
func (r *Repository) listBySession(ctx context.Context, id uuid.UUID) ([]SessionRecord, error) {
	rows, err := r.queries.ListAttendanceBySession(ctx, id)
	if err != nil {
		return nil, err
	}
	out := make([]SessionRecord, len(rows))
	for i, row := range rows {
		out[i] = SessionRecord{Record: Record{ID: row.ID, SessionID: row.SessionID, StudentID: row.StudentID, Status: row.Status, Note: row.Note}, StudentName: row.StudentName, StudentIndex: row.StudentIndex}
	}
	return out, nil
}
func (r *Repository) listByStudent(ctx context.Context, id uuid.UUID) ([]StudentRecord, error) {
	rows, err := r.queries.ListAttendanceByStudent(ctx, id)
	if err != nil {
		return nil, err
	}
	out := make([]StudentRecord, len(rows))
	for i, row := range rows {
		out[i] = StudentRecord{Record: Record{ID: row.ID, SessionID: row.SessionID, StudentID: row.StudentID, Status: row.Status, Note: row.Note}, SessionDate: row.SessionDate, ClassName: row.ClassName}
	}
	return out, nil
}
func (r *Repository) summary(ctx context.Context, studentID, classID uuid.UUID) (Summary, error) {
	row, err := r.queries.GetAttendanceSummaryByStudent(ctx, db.GetAttendanceSummaryByStudentParams{StudentID: studentID, ClassID: classID})
	return Summary{TotalDays: row.TotalDays, Present: row.Present, Absent: row.Absent, Late: row.Late, Excused: row.Excused}, err
}
func (r *Repository) userExists(ctx context.Context, id uuid.UUID) bool {
	_, err := r.queries.GetUserByID(ctx, id)
	return err == nil
}
func (r *Repository) userByEmail(ctx context.Context, email string) (uuid.UUID, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)
	return user.ID, err
}
func (r *Repository) createUser(ctx context.Context, actor Actor) (uuid.UUID, error) {
	user, err := r.queries.CreateUser(ctx, db.CreateUserParams{ID: actor.ID, Email: actor.Email, FullName: actor.FullName, Role: actor.Role})
	return user.ID, err
}
func (r *Repository) teacherByUser(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	teacher, err := r.queries.GetTeacherByUserID(ctx, userID)
	return teacher.ID, err
}
func (r *Repository) teacherAssigned(ctx context.Context, classID, teacherID uuid.UUID) (bool, error) {
	return r.queries.IsTeacherAssignedToClass(ctx, db.IsTeacherAssignedToClassParams{ID: classID, FormTeacherID: pgtype.UUID{Bytes: teacherID, Valid: true}, TeacherID: teacherID})
}
func (r *Repository) studentClass(ctx context.Context, studentID uuid.UUID) (uuid.UUID, error) {
	class, err := r.queries.GetStudentCurrentClass(ctx, studentID)
	return class.ID, err
}
func (r *Repository) guardianUsers(ctx context.Context, studentID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.queries.ListGuardianUserIDsByStudentIDs(ctx, []uuid.UUID{studentID})
	if err != nil {
		return nil, err
	}
	ids := make([]uuid.UUID, 0, len(rows))
	for _, row := range rows {
		if row.Valid {
			ids = append(ids, uuid.UUID(row.Bytes))
		}
	}
	return ids, nil
}
func (r *Repository) studentName(ctx context.Context, studentID uuid.UUID) (string, error) {
	student, err := r.queries.GetStudentByID(ctx, studentID)
	return student.FullName, err
}
func (r *Repository) className(ctx context.Context, classID uuid.UUID) (string, error) {
	class, err := r.queries.GetClassByID(ctx, classID)
	return class.Name, err
}
func (r *Repository) currentAcademicYear(ctx context.Context) (uuid.UUID, error) {
	year, err := r.queries.GetCurrentAcademicYear(ctx)
	return year.ID, err
}
func (r *Repository) listReportRows(ctx context.Context, classID uuid.UUID, from, to time.Time) ([]ReportRow, error) {
	rows, err := r.queries.ListAttendanceRecordsForClassInRange(ctx, db.ListAttendanceRecordsForClassInRangeParams{ClassID: classID, Date: dateValue(from), Date_2: dateValue(to)})
	if err != nil {
		return nil, err
	}
	out := make([]ReportRow, len(rows))
	for i, row := range rows {
		out[i] = ReportRow{ID: row.ID, StudentName: row.StudentName, StudentIndex: row.StudentIndex, SessionDate: row.SessionDate, Status: row.Status, Note: row.Note}
	}
	return out, nil
}

func dateValue(value time.Time) pgtype.Date { return pgtype.Date{Time: value, Valid: true} }
func textValue(value string) pgtype.Text    { return pgtype.Text{String: value, Valid: value != ""} }
func mapSession(row db.AttendanceSession) Session {
	return Session{ID: row.ID, ClassID: row.ClassID, TakenBy: row.TakenBy, Date: row.Date, CreatedAt: row.CreatedAt}
}
func mapRecord(row db.AttendanceRecord) Record {
	return Record{ID: row.ID, SessionID: row.SessionID, StudentID: row.StudentID, Status: row.Status, Note: row.Note}
}
