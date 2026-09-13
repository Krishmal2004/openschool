package school

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/openschool-org/openschool/db/sqlc"
	"github.com/openschool-org/openschool/internal/ports"
)

type houseRepository struct{ queries *db.Queries }

func newHouseRepository(pool *pgxpool.Pool) *houseRepository {
	return &houseRepository{queries: db.New(pool)}
}
func (r *houseRepository) create(ctx context.Context, command houseCommand) (House, error) {
	row, err := r.queries.CreateHouse(ctx, db.CreateHouseParams{Name: command.Name, Code: optionalText(command.Code), Color: command.Color})
	return mapHouse(row), err
}
func (r *houseRepository) get(ctx context.Context, id uuid.UUID) (House, error) {
	row, err := r.queries.GetHouseByID(ctx, id)
	return mapHouse(row), err
}
func (r *houseRepository) list(ctx context.Context) ([]House, error) {
	rows, err := r.queries.ListHouses(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]House, len(rows))
	for i, row := range rows {
		result[i] = mapHouse(row)
	}
	return result, nil
}
func (r *houseRepository) update(ctx context.Context, id uuid.UUID, command houseCommand) (House, error) {
	row, err := r.queries.UpdateHouse(ctx, db.UpdateHouseParams{ID: id, Name: command.Name, Code: optionalText(command.Code), Color: command.Color})
	return mapHouse(row), err
}
func (r *houseRepository) delete(ctx context.Context, id uuid.UUID) (int64, error) {
	return r.queries.DeleteHouse(ctx, id)
}
func (r *houseRepository) pickForStudent(ctx context.Context) (uuid.UUID, bool, error) {
	id, err := r.queries.PickBalancedHouseForStudent(ctx)
	if errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, false, nil
	}
	return id, err == nil, err
}
func (r *houseRepository) pickForTeacher(ctx context.Context) (uuid.UUID, bool, error) {
	id, err := r.queries.PickBalancedHouseForTeacher(ctx)
	if errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, false, nil
	}
	return id, err == nil, err
}
func (r *houseRepository) listStudentsMissingHouse(ctx context.Context) ([]uuid.UUID, error) {
	rows, err := r.queries.ListStudentsMissingHouse(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]uuid.UUID, len(rows))
	for i, row := range rows {
		result[i] = row.ID
	}
	return result, nil
}
func (r *houseRepository) listTeachersMissingHouse(ctx context.Context) ([]uuid.UUID, error) {
	rows, err := r.queries.ListTeachersMissingHouse(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]uuid.UUID, len(rows))
	for i, row := range rows {
		result[i] = row.ID
	}
	return result, nil
}
func (r *houseRepository) getStudent(ctx context.Context, id uuid.UUID) (ports.StudentHouseProfile, error) {
	row, err := r.queries.GetStudentByID(ctx, id)
	return mapStudentHouseProfile(row), err
}
func (r *houseRepository) getTeacher(ctx context.Context, id uuid.UUID) (ports.TeacherHouseProfile, error) {
	row, err := r.queries.GetTeacherByID(ctx, id)
	return mapTeacherByIDHouseProfile(row), err
}
func nullableUUID(id *uuid.UUID) pgtype.UUID {
	if id == nil {
		return pgtype.UUID{}
	}
	return pgtype.UUID{Bytes: *id, Valid: true}
}
func (r *houseRepository) updateStudentHouse(ctx context.Context, id uuid.UUID, houseID *uuid.UUID) (ports.StudentHouseProfile, error) {
	row, err := r.queries.UpdateStudentHouse(ctx, db.UpdateStudentHouseParams{ID: id, HouseID: nullableUUID(houseID)})
	return mapStudentHouseProfile(row), err
}
func (r *houseRepository) updateTeacherHouse(ctx context.Context, id uuid.UUID, houseID *uuid.UUID) (ports.TeacherHouseProfile, error) {
	row, err := r.queries.UpdateTeacherHouse(ctx, db.UpdateTeacherHouseParams{ID: id, HouseID: nullableUUID(houseID)})
	return mapTeacherHouseProfile(row), err
}
func mapHouse(row db.House) House {
	return House{ID: row.ID, Name: row.Name, Code: row.Code, CreatedAt: row.CreatedAt, Color: row.Color}
}
func mapStudentHouseProfile(row db.StudentProfile) ports.StudentHouseProfile {
	return ports.StudentHouseProfile{
		ID: row.ID, UserID: row.UserID, FullName: row.FullName, IndexNumber: row.IndexNumber,
		Address: row.Address, Phone: row.Phone, Whatsapp: row.Whatsapp, SpecialRemarks: row.SpecialRemarks,
		CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt, Gender: row.Gender, HouseID: row.HouseID,
		EnrollmentStatus: row.EnrollmentStatus,
	}
}
func mapTeacherHouseProfile(row db.TeacherProfile) ports.TeacherHouseProfile {
	return ports.TeacherHouseProfile{
		ID: row.ID, UserID: row.UserID, FullName: row.FullName, EmployeeNumber: row.EmployeeNumber,
		JoinedDate: row.JoinedDate, Phone: row.Phone, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
		Title: row.Title, Gender: row.Gender, IsActive: row.IsActive, HouseID: row.HouseID,
		EmploymentStatus: row.EmploymentStatus, NicNumber: row.NicNumber,
	}
}

func mapTeacherByIDHouseProfile(row db.GetTeacherByIDRow) ports.TeacherHouseProfile {
	return ports.TeacherHouseProfile{
		ID: row.ID, UserID: row.UserID, FullName: row.FullName, EmployeeNumber: row.EmployeeNumber,
		JoinedDate: row.JoinedDate, Phone: row.Phone, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
		Title: row.Title, Gender: row.Gender, IsActive: row.IsActive, HouseID: row.HouseID,
		EmploymentStatus: row.EmploymentStatus, NicNumber: row.NicNumber,
	}
}

type gradeRepository struct{ queries *db.Queries }

func newGradeRepository(pool *pgxpool.Pool) *gradeRepository {
	return &gradeRepository{queries: db.New(pool)}
}

func (r *gradeRepository) create(ctx context.Context, command gradeCommand) (Grade, error) {
	row, err := r.queries.CreateGrade(ctx, db.CreateGradeParams{Name: command.Name, SortOrder: command.SortOrder})
	return mapGrade(row), err
}

func (r *gradeRepository) get(ctx context.Context, id uuid.UUID) (Grade, error) {
	row, err := r.queries.GetGradeByID(ctx, id)
	return mapGrade(row), err
}

func (r *gradeRepository) list(ctx context.Context) ([]Grade, error) {
	rows, err := r.queries.ListGrades(ctx)
	if err != nil {
		return nil, err
	}
	grades := make([]Grade, len(rows))
	for i, row := range rows {
		grades[i] = mapGrade(row)
	}
	return grades, nil
}

func (r *gradeRepository) update(ctx context.Context, id uuid.UUID, command gradeCommand) (Grade, error) {
	row, err := r.queries.UpdateGrade(ctx, db.UpdateGradeParams{ID: id, Name: command.Name, SortOrder: command.SortOrder})
	return mapGrade(row), err
}

func (r *gradeRepository) delete(ctx context.Context, id uuid.UUID) (int64, error) {
	return r.queries.DeleteGrade(ctx, id)
}

func mapGrade(row db.Grade) Grade {
	return Grade{ID: row.ID, Name: row.Name, SortOrder: row.SortOrder, CreatedAt: row.CreatedAt.Time}
}

type termRepository struct{ queries *db.Queries }

func newTermRepository(pool *pgxpool.Pool) *termRepository {
	return &termRepository{queries: db.New(pool)}
}
func (r *termRepository) create(ctx context.Context, values termValues) (Term, error) {
	row, err := r.queries.CreateTerm(ctx, db.CreateTermParams{AcademicYearID: values.AcademicYearID, Name: values.Name, StartDate: pgtype.Date{Time: values.StartDate, Valid: true}, EndDate: pgtype.Date{Time: values.EndDate, Valid: true}, SortOrder: values.SortOrder})
	return mapTerm(row), err
}
func (r *termRepository) get(ctx context.Context, id uuid.UUID) (Term, error) {
	row, err := r.queries.GetTermByID(ctx, id)
	return mapTerm(row), err
}
func (r *termRepository) list(ctx context.Context, yearID uuid.UUID) ([]Term, error) {
	rows, err := r.queries.ListTermsByAcademicYear(ctx, yearID)
	if err != nil {
		return nil, err
	}
	result := make([]Term, len(rows))
	for i, row := range rows {
		result[i] = mapTerm(row)
	}
	return result, nil
}
func (r *termRepository) current(ctx context.Context) (Term, error) {
	row, err := r.queries.GetCurrentTerm(ctx)
	return mapTerm(row), err
}
func (r *termRepository) setCurrent(ctx context.Context, id uuid.UUID) error {
	return r.queries.SetCurrentTerm(ctx, id)
}
func (r *termRepository) update(ctx context.Context, id uuid.UUID, values termValues) (Term, error) {
	row, err := r.queries.UpdateTerm(ctx, db.UpdateTermParams{ID: id, Name: values.Name, StartDate: pgtype.Date{Time: values.StartDate, Valid: true}, EndDate: pgtype.Date{Time: values.EndDate, Valid: true}, SortOrder: values.SortOrder})
	return mapTerm(row), err
}
func (r *termRepository) delete(ctx context.Context, id uuid.UUID) (int64, error) {
	return r.queries.DeleteTerm(ctx, id)
}
func mapTerm(row db.Term) Term {
	return Term{ID: row.ID.String(), AcademicYearID: row.AcademicYearID.String(), Name: row.Name, StartDate: row.StartDate.Time.Format("2006-01-02"), EndDate: row.EndDate.Time.Format("2006-01-02"), IsCurrent: row.IsCurrent, SortOrder: row.SortOrder, CreatedAt: row.CreatedAt.Time.Format(time.RFC3339Nano)}
}

type schoolRepository struct{ queries *db.Queries }

func newSchoolRepository(pool *pgxpool.Pool) *schoolRepository {
	return &schoolRepository{queries: db.New(pool)}
}
func optionalText(value string) pgtype.Text { return pgtype.Text{String: value, Valid: value != ""} }
func optionalInt(value *int32) pgtype.Int4 {
	if value == nil {
		return pgtype.Int4{}
	}
	return pgtype.Int4{Int32: *value, Valid: true}
}
func (r *schoolRepository) createSchool(ctx context.Context, values schoolValues) (School, error) {
	row, err := r.queries.CreateSchool(ctx, db.CreateSchoolParams{Name: values.Name, Address: optionalText(values.Address), Phone: optionalText(values.Phone), Email: optionalText(values.Email), LogoUrl: optionalText(values.LogoURL), GradeFrom: optionalInt(values.GradeFrom), GradeTo: optionalInt(values.GradeTo), SchoolType: values.SchoolType})
	return mapSchool(row), err
}
func (r *schoolRepository) getSchool(ctx context.Context) (School, error) {
	row, err := r.queries.GetSchool(ctx)
	return mapSchool(row), err
}
func (r *schoolRepository) updateSchool(ctx context.Context, values schoolValues) (School, error) {
	row, err := r.queries.UpdateSchool(ctx, db.UpdateSchoolParams{ID: values.ID, Name: values.Name, Address: optionalText(values.Address), Phone: optionalText(values.Phone), Email: optionalText(values.Email), LogoUrl: optionalText(values.LogoURL), GradeFrom: optionalInt(values.GradeFrom), GradeTo: optionalInt(values.GradeTo), SchoolType: optionalText(values.SchoolType)})
	return mapSchool(row), err
}
func (r *schoolRepository) createYear(ctx context.Context, values yearValues) (AcademicYear, error) {
	row, err := r.queries.CreateAcademicYear(ctx, db.CreateAcademicYearParams{Label: values.Label, StartDate: pgtype.Date{Time: values.StartDate, Valid: true}, EndDate: pgtype.Date{Time: values.EndDate, Valid: true}, IsCurrent: values.IsCurrent})
	return mapAcademicYear(row), err
}
func (r *schoolRepository) getYear(ctx context.Context, id uuid.UUID) (AcademicYear, error) {
	row, err := r.queries.GetAcademicYearByID(ctx, id)
	return mapAcademicYear(row), err
}
func (r *schoolRepository) currentYear(ctx context.Context) (AcademicYear, error) {
	row, err := r.queries.GetCurrentAcademicYear(ctx)
	return mapAcademicYear(row), err
}
func (r *schoolRepository) listYears(ctx context.Context) ([]AcademicYear, error) {
	rows, err := r.queries.ListAcademicYears(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]AcademicYear, len(rows))
	for i, row := range rows {
		result[i] = mapAcademicYear(row)
	}
	return result, nil
}
func (r *schoolRepository) setCurrentYear(ctx context.Context, id uuid.UUID) error {
	return r.queries.SetCurrentAcademicYear(ctx, id)
}
func (r *schoolRepository) deleteYear(ctx context.Context, id uuid.UUID) (int64, error) {
	return r.queries.DeleteAcademicYear(ctx, id)
}
func nullableString(value pgtype.Text) *string {
	if !value.Valid {
		return nil
	}
	result := value.String
	return &result
}
func nullableInt(value pgtype.Int4) *int32 {
	if !value.Valid {
		return nil
	}
	result := value.Int32
	return &result
}
func mapSchool(row db.School) School {
	return School{ID: row.ID.String(), Name: row.Name, Address: nullableString(row.Address), Phone: nullableString(row.Phone), Email: nullableString(row.Email), LogoURL: nullableString(row.LogoUrl), GradeFrom: nullableInt(row.GradeFrom), GradeTo: nullableInt(row.GradeTo), SchoolType: row.SchoolType, CreatedAt: row.CreatedAt.Time.Format(time.RFC3339Nano)}
}
func mapAcademicYear(row db.AcademicYear) AcademicYear {
	return AcademicYear{ID: row.ID.String(), Label: row.Label, StartDate: row.StartDate.Time.Format("2006-01-02"), EndDate: row.EndDate.Time.Format("2006-01-02"), IsCurrent: row.IsCurrent, CreatedAt: row.CreatedAt.Time.Format(time.RFC3339Nano)}
}
