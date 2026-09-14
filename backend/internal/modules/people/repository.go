package people

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/openschool-org/openschool/db/sqlc"
)

type studentRepository struct{ queries *db.Queries }

func NewStudentStore(pool *pgxpool.Pool) *studentRepository {
	return &studentRepository{queries: db.New(pool)}
}

func NewStudentReader(pool *pgxpool.Pool) StudentReader {
	return &studentRepository{queries: db.New(pool)}
}
func NewStudentWriter(pool *pgxpool.Pool) StudentStatusWriter {
	return &studentRepository{queries: db.New(pool)}
}
func (r *studentRepository) Get(c context.Context, id uuid.UUID) (any, error) {
	return r.queries.GetStudentByID(c, id)
}
func (r *studentRepository) GetWithClass(c context.Context, id uuid.UUID) (any, error) {
	return r.queries.GetStudentWithClass(c, id)
}
func (r *studentRepository) List(c context.Context) (any, error) { return r.queries.ListStudents(c) }
func (r *studentRepository) ListByClass(c context.Context, id uuid.UUID) (any, error) {
	return r.queries.ListStudentsByClass(c, id)
}
func (r *studentRepository) UpdateStatus(c context.Context, id uuid.UUID, status string) (any, error) {
	return r.queries.UpdateStudentEnrollmentStatus(c, db.UpdateStudentEnrollmentStatusParams{ID: id, EnrollmentStatus: status})
}

func (r *studentRepository) FindByIndex(c context.Context, index string) error {
	_, err := r.queries.GetStudentByIndexNumber(c, index)
	return err
}
func (r *studentRepository) GetUser(c context.Context, id uuid.UUID) (studentUser, error) {
	u, err := r.queries.GetUserByID(c, id)
	return studentUser{Email: u.Email}, err
}
func (r *studentRepository) GetStudentRecord(c context.Context, id uuid.UUID) (studentRecord, error) {
	p, err := r.queries.GetStudentByID(c, id)
	return studentRecord{ID: p.ID, UserID: uuid.UUID(p.UserID.Bytes), IndexNumber: p.IndexNumber, FullName: p.FullName, Phone: p.Phone, Whatsapp: p.Whatsapp, Address: p.Address, SpecialRemark: p.SpecialRemarks, Gender: p.Gender}, err
}
func (r *studentRepository) CreateStudentUser(c context.Context, p studentUserCreate) error {
	_, err := r.queries.CreateUser(c, db.CreateUserParams{ID: p.ID, Email: p.Email, FullName: p.FullName, Role: "student", MustChangePassword: p.MustChangePassword})
	return err
}
func (r *studentRepository) CreateUser(c context.Context, p teacherUserCreate) error {
	_, err := r.queries.CreateUser(c, db.CreateUserParams{ID: p.ID, Email: p.Email, FullName: p.FullName, Role: "teacher", MustChangePassword: p.MustChangePassword})
	return err
}
func (r *studentRepository) Create(c context.Context, p studentCreate) (any, error) {
	return r.queries.CreateStudentProfile(c, db.CreateStudentProfileParams{
		UserID: pgtype.UUID{Bytes: p.UserID, Valid: true}, FullName: p.FullName, IndexNumber: p.Index,
		Address: pgtype.Text{String: p.Address, Valid: p.Address != ""}, Phone: pgtype.Text{String: p.Phone, Valid: p.Phone != ""},
		Whatsapp: pgtype.Text{String: p.WhatsApp, Valid: p.WhatsApp != ""}, SpecialRemarks: pgtype.Text{String: p.Remarks, Valid: p.Remarks != ""},
		Gender: pgtype.Text{String: p.Gender, Valid: p.Gender != ""}, HouseID: pgtype.UUID{Bytes: p.HouseID, Valid: p.HouseID != uuid.Nil},
	})
}
func (r *studentRepository) Update(c context.Context, id uuid.UUID, p studentUpdate) (any, error) {
	return r.queries.UpdateStudentProfile(c, db.UpdateStudentProfileParams{ID: id, FullName: p.FullName,
		Address: pgtype.Text{String: p.Address, Valid: p.Address != ""}, Phone: pgtype.Text{String: p.Phone, Valid: p.Phone != ""},
		Whatsapp: pgtype.Text{String: p.WhatsApp, Valid: p.WhatsApp != ""}, SpecialRemarks: pgtype.Text{String: p.Remarks, Valid: p.Remarks != ""}, Gender: pgtype.Text{String: p.Gender, Valid: p.Gender != ""}})
}
func (r *studentRepository) Delete(c context.Context, id uuid.UUID) error {
	return r.queries.DeleteStudentProfile(c, id)
}
func (r *studentRepository) DeleteUser(c context.Context, id uuid.UUID) error {
	return r.queries.DeleteUser(c, id)
}

func (r *studentRepository) NextEmployee(c context.Context) (string, error) {
	return r.queries.NextEmployeeNumber(c)
}
func (r *studentRepository) GetTeacher(c context.Context, id uuid.UUID) (teacherRecord, error) {
	p, err := r.queries.GetTeacherByID(c, id)
	return teacherRecord{ID: p.ID, UserID: p.UserID, EmployeeNumber: p.EmployeeNumber}, err
}
func (r *studentRepository) GetUserEmail(c context.Context, id uuid.UUID) (string, error) {
	u, err := r.queries.GetUserByID(c, id)
	return u.Email, err
}
func (r *studentRepository) CreateTeacher(c context.Context, p teacherCreate) (any, error) {
	return r.queries.CreateTeacherProfile(c, db.CreateTeacherProfileParams{UserID: p.UserID, FullName: p.FullName, EmployeeNumber: p.EmployeeNumber, NicNumber: p.NIC, JoinedDate: pgtype.Date{Time: p.JoinedDate, Valid: true}, Phone: pgtype.Text{String: p.Phone, Valid: p.Phone != ""}, Title: pgtype.Text{String: p.Title, Valid: p.Title != ""}, Gender: pgtype.Text{String: p.Gender, Valid: p.Gender != ""}, HouseID: pgtype.UUID{Bytes: p.HouseID, Valid: p.HouseID != uuid.Nil}})
}
func (r *studentRepository) UpdateTeacher(c context.Context, id uuid.UUID, p teacherUpdate) (any, error) {
	return r.queries.UpdateTeacherProfile(c, db.UpdateTeacherProfileParams{ID: id, FullName: p.FullName, Phone: pgtype.Text{String: p.Phone, Valid: p.Phone != ""}, Title: pgtype.Text{String: p.Title, Valid: p.Title != ""}, Gender: pgtype.Text{String: p.Gender, Valid: p.Gender != ""}, NicNumber: p.NIC})
}
func (r *studentRepository) UpdateTeacherStatus(c context.Context, id uuid.UUID, status string) (any, error) {
	return r.queries.UpdateTeacherEmploymentStatus(c, db.UpdateTeacherEmploymentStatusParams{ID: id, EmploymentStatus: status})
}
func (r *studentRepository) DeleteTeacher(c context.Context, id uuid.UUID) (int64, error) {
	return r.queries.DeleteTeacher(c, id)
}
func (r *studentRepository) AssignSubject(c context.Context, id, subject uuid.UUID) error {
	return r.queries.AssignSubjectToTeacher(c, db.AssignSubjectToTeacherParams{TeacherID: id, SubjectID: subject})
}
func (r *studentRepository) RemoveSubject(c context.Context, id, subject uuid.UUID) error {
	return r.queries.RemoveSubjectFromTeacher(c, db.RemoveSubjectFromTeacherParams{TeacherID: id, SubjectID: subject})
}
func (r *studentRepository) CountSubjects(c context.Context, id uuid.UUID) (int64, error) {
	return r.queries.CountSubjectsByTeacher(c, id)
}
func (r *studentRepository) HasWorkload(c context.Context, id uuid.UUID) (bool, error) {
	rows, err := r.queries.ListTeacherWorkload(c, id)
	return len(rows) > 0, err
}
func (r *studentRepository) SetActive(c context.Context, id uuid.UUID, active bool) error {
	return r.queries.SetTeacherActiveStatus(c, db.SetTeacherActiveStatusParams{ID: id, IsActive: active})
}

type teacherReader struct{ queries *db.Queries }

func NewTeacherReader(pool *pgxpool.Pool) TeacherReader { return &teacherReader{queries: db.New(pool)} }
func (r *teacherReader) Get(c context.Context, id uuid.UUID) (any, error) {
	return r.queries.GetTeacherByID(c, id)
}
func (r *teacherReader) List(c context.Context) (any, error) { return r.queries.ListTeachers(c) }
func (r *teacherReader) Subjects(c context.Context, id uuid.UUID) (any, error) {
	return r.queries.ListSubjectsByTeacher(c, id)
}
func (r *teacherReader) Workload(c context.Context, id uuid.UUID) (any, error) {
	return r.queries.ListTeacherWorkload(c, id)
}
func (r *teacherReader) BySubject(c context.Context, id uuid.UUID) (any, error) {
	return r.queries.ListTeachersBySubject(c, id)
}
