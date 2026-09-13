package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// AuditRecorder is the narrow audit capability required by business modules.
type AuditRecorder interface {
	Record(context.Context, string, uuid.UUID, string, uuid.UUID, interface{}, interface{}, string) error
}

// HouseAssignments exposes house allocation to the people module without
// exposing the school module's repository implementation.
type HouseAssignments interface {
	PickForStudent(context.Context) (uuid.UUID, bool)
	PickForTeacher(context.Context) (uuid.UUID, bool)
	ChangeStudentHouse(context.Context, uuid.UUID, string, uuid.UUID) (StudentHouseProfile, error)
	ChangeTeacherHouse(context.Context, uuid.UUID, string, uuid.UUID) (TeacherHouseProfile, error)
}

// StudentHouseProfile preserves the existing student-house endpoint payload
// while keeping generated sqlc models behind the repository boundary.
type StudentHouseProfile struct {
	ID               uuid.UUID          `json:"id"`
	UserID           pgtype.UUID        `json:"user_id"`
	FullName         string             `json:"full_name"`
	IndexNumber      string             `json:"index_number"`
	Address          pgtype.Text        `json:"address"`
	Phone            pgtype.Text        `json:"phone"`
	Whatsapp         pgtype.Text        `json:"whatsapp"`
	SpecialRemarks   pgtype.Text        `json:"special_remarks"`
	CreatedAt        pgtype.Timestamptz `json:"created_at"`
	UpdatedAt        pgtype.Timestamptz `json:"updated_at"`
	Gender           pgtype.Text        `json:"gender"`
	HouseID          pgtype.UUID        `json:"house_id"`
	EnrollmentStatus string             `json:"enrollment_status"`
}

// TeacherHouseProfile preserves the existing teacher-house endpoint payload.
type TeacherHouseProfile struct {
	ID               uuid.UUID          `json:"id"`
	UserID           uuid.UUID          `json:"user_id"`
	FullName         string             `json:"full_name"`
	EmployeeNumber   string             `json:"employee_number"`
	JoinedDate       pgtype.Date        `json:"joined_date"`
	Phone            pgtype.Text        `json:"phone"`
	CreatedAt        pgtype.Timestamptz `json:"created_at"`
	UpdatedAt        pgtype.Timestamptz `json:"updated_at"`
	Title            pgtype.Text        `json:"title"`
	Gender           pgtype.Text        `json:"gender"`
	IsActive         bool               `json:"is_active"`
	HouseID          pgtype.UUID        `json:"house_id"`
	EmploymentStatus string             `json:"employment_status"`
	NicNumber        string             `json:"nic_number"`
}
