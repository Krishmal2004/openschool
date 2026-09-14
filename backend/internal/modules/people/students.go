package people

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/openschool-org/openschool/internal/identity"
	"github.com/openschool-org/openschool/internal/models"
	"github.com/openschool-org/openschool/internal/ports"
	"github.com/openschool-org/openschool/internal/validation"
)

var ErrGenderMismatchSchoolType = errors.New("student gender does not match the school's single-sex type")

type studentRecord struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	IndexNumber   string
	Email         string
	FullName      string
	Phone         pgtype.Text
	Whatsapp      pgtype.Text
	Address       pgtype.Text
	SpecialRemark pgtype.Text
	Gender        pgtype.Text
}

type studentStore interface {
	Create(context.Context, studentCreate) (any, error)
	GetStudentRecord(context.Context, uuid.UUID) (studentRecord, error)
	GetUser(context.Context, uuid.UUID) (studentUser, error)
	FindByIndex(context.Context, string) error
	Update(context.Context, uuid.UUID, studentUpdate) (any, error)
	UpdateStatus(context.Context, uuid.UUID, string) (any, error)
	Delete(context.Context, uuid.UUID) error
	CreateStudentUser(context.Context, studentUserCreate) error
	DeleteUser(context.Context, uuid.UUID) error
}

type studentUser struct{ Email string }
type studentUserCreate struct {
	ID                 uuid.UUID
	Email              string
	FullName           string
	MustChangePassword bool
}
type studentCreate struct {
	UserID, HouseID   uuid.UUID
	FullName, Index   string
	Address, Phone    string
	WhatsApp, Remarks string
	Gender            string
}
type studentUpdate struct {
	FullName, Address, Phone, WhatsApp, Remarks, Gender string
}

type StudentService struct {
	store  studentStore
	idp    identity.Provider
	houses ports.HouseAssignments
	audit  ports.AuditRecorder
	school func(context.Context) (string, error)
}

func NewStudentService(store studentStore, idp identity.Provider, houses ports.HouseAssignments, audit ports.AuditRecorder, school func(context.Context) (string, error)) *StudentService {
	return &StudentService{store: store, idp: idp, houses: houses, audit: audit, school: school}
}

func (s *StudentService) validateGender(ctx context.Context, gender string) error {
	if s.school == nil {
		return nil
	}
	t, err := s.school(ctx)
	if err != nil {
		return nil
	}
	if (t == "boys" && gender != "male") || (t == "girls" && gender != "female") {
		return ErrGenderMismatchSchoolType
	}
	return nil
}

func (s *StudentService) Create(ctx context.Context, req models.CreateStudentRequest, actor uuid.UUID) (any, error) {
	if !validation.IsValidSriLankanPhone(req.PhoneNumber) || !validation.IsValidSriLankanPhone(req.WhatsApp) {
		return nil, validation.ErrInvalidPhone
	}
	if err := s.validateGender(ctx, req.Gender); err != nil {
		return nil, err
	}
	if err := s.store.FindByIndex(ctx, req.IndexNumber); err == nil {
		return nil, fmt.Errorf("index number already exists")
	}
	idpUser, err := s.idp.CreateUser(ctx, models.RoleStudent, map[string]any{"username": req.IndexNumber, "email": req.Email, "given_name": req.GivenName, "family_name": req.FamilyName, "phone": req.PhoneNumber, "password": req.IndexNumber})
	if err != nil {
		return nil, fmt.Errorf("failed to create identity provider user: %w", err)
	}
	uid, err := uuid.Parse(idpUser.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid identity provider user ID: %w", err)
	}
	rollback := func() { _ = s.idp.DeleteUser(ctx, idpUser.ID); _ = s.store.DeleteUser(ctx, uid) }
	fullName := req.GivenName + " " + req.FamilyName
	if err = s.store.CreateStudentUser(ctx, studentUserCreate{ID: uid, Email: req.Email, FullName: fullName, MustChangePassword: true}); err != nil {
		_ = s.idp.DeleteUser(ctx, idpUser.ID)
		return nil, fmt.Errorf("failed to create user record: %w", err)
	}
	if err = s.idp.AssignRole(ctx, identity.RoleID(models.RoleStudent), idpUser.ID); err != nil {
		rollback()
		return nil, fmt.Errorf("failed to assign student role: %w", err)
	}
	var house uuid.UUID
	if s.houses != nil {
		house, _ = s.houses.PickForStudent(ctx)
	}
	profile, err := s.store.Create(ctx, studentCreate{UserID: uid, HouseID: house, FullName: fullName, Index: req.IndexNumber, Address: req.Address, Phone: req.PhoneNumber, WhatsApp: req.WhatsApp, Remarks: req.SpecialRemarks, Gender: req.Gender})
	if err != nil {
		rollback()
		return nil, fmt.Errorf("failed to create student profile: %w", err)
	}
	if s.audit != nil {
		_ = s.audit.Record(ctx, "student_account", uid, "account_created", actor, nil, profile, "")
	}
	return profile, nil
}

func (s *StudentService) Update(ctx context.Context, id uuid.UUID, req models.UpdateStudentRequest) (any, error) {
	if !validation.IsValidSriLankanPhone(req.PhoneNumber) || !validation.IsValidSriLankanPhone(req.WhatsApp) {
		return nil, validation.ErrInvalidPhone
	}
	if err := s.validateGender(ctx, req.Gender); err != nil {
		return nil, err
	}
	student, err := s.store.GetStudentRecord(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("student not found")
	}
	user, err := s.store.GetUser(ctx, student.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	if err := s.idp.UpdateUser(ctx, student.UserID.String(), models.RoleStudent, map[string]any{"username": student.IndexNumber, "email": user.Email, "given_name": req.GivenName, "family_name": req.FamilyName, "phone": req.PhoneNumber}); err != nil {
		log.Printf("UpdateStudent: failed to update identity provider user: %v", err)
	}
	return s.store.Update(ctx, id, studentUpdate{FullName: req.GivenName + " " + req.FamilyName, Address: req.Address, Phone: req.PhoneNumber, WhatsApp: req.WhatsApp, Remarks: req.SpecialRemarks, Gender: req.Gender})
}

func (s *StudentService) UpdateHouse(ctx context.Context, id uuid.UUID, req models.UpdateStudentHouseRequest, actor uuid.UUID) (any, error) {
	return s.houses.ChangeStudentHouse(ctx, id, req.HouseID, actor)
}
func (s *StudentService) UpdateStatus(ctx context.Context, id uuid.UUID, status string) (any, error) {
	if status != "active" && status != "left" {
		return nil, fmt.Errorf("invalid status %q — must be active or left", status)
	}
	return s.store.UpdateStatus(ctx, id, status)
}
func (s *StudentService) Delete(ctx context.Context, id, actor uuid.UUID) error {
	student, err := s.store.GetStudentRecord(ctx, id)
	if err != nil {
		return fmt.Errorf("student not found")
	}
	if err = s.store.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete student profile: %w", err)
	}
	if s.audit != nil {
		_ = s.audit.Record(ctx, "student_account", id, "account_deleted", actor, student, nil, "")
	}
	if student.UserID == uuid.Nil {
		return nil
	}
	if err = s.store.DeleteUser(ctx, student.UserID); err != nil {
		return fmt.Errorf("student profile deleted but failed to delete local user record — the account (id %s) can still sign in with no profile and must be removed manually or by retrying this delete: %w", student.UserID, err)
	}
	if err = s.idp.DeleteUser(ctx, student.UserID.String()); err != nil {
		return fmt.Errorf("student profile deleted locally but failed to delete identity provider user (account is now orphaned and must be removed manually): %w", err)
	}
	return nil
}
