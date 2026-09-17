package people

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/idp"
)

type eraseStoreStub struct {
	record          studentRecord
	recordErr       error
	anonymizeCalled bool
	anonymizeID     uuid.UUID
	anonymizeErr    error
	eraseUserCalled bool
	eraseUserID     uuid.UUID
}

func (s *eraseStoreStub) Create(context.Context, studentCreate) (any, error) { return nil, nil }
func (s *eraseStoreStub) GetStudentRecord(context.Context, uuid.UUID) (studentRecord, error) {
	return s.record, s.recordErr
}
func (s *eraseStoreStub) GetUser(context.Context, uuid.UUID) (studentUser, error) {
	return studentUser{}, nil
}
func (s *eraseStoreStub) FindByIndex(context.Context, string) error { return errors.New("not found") }
func (s *eraseStoreStub) Update(context.Context, uuid.UUID, studentUpdate) (any, error) {
	return nil, nil
}
func (s *eraseStoreStub) UpdateStatus(context.Context, uuid.UUID, string) (any, error) {
	return nil, nil
}
func (s *eraseStoreStub) Delete(context.Context, uuid.UUID) error                    { return nil }
func (s *eraseStoreStub) CreateStudentUser(context.Context, studentUserCreate) error { return nil }
func (s *eraseStoreStub) DeleteUser(context.Context, uuid.UUID) error                { return nil }
func (s *eraseStoreStub) AnonymizeProfile(_ context.Context, id uuid.UUID) error {
	s.anonymizeCalled, s.anonymizeID = true, id
	return s.anonymizeErr
}
func (s *eraseStoreStub) EraseUser(_ context.Context, id uuid.UUID) error {
	s.eraseUserCalled, s.eraseUserID = true, id
	return nil
}

type eraseIdentityStub struct {
	deletedUserID string
}

func (i *eraseIdentityStub) CreateUser(context.Context, string, map[string]any) (*idp.User, error) {
	return nil, nil
}
func (*eraseIdentityStub) UpdateUser(context.Context, string, string, map[string]any) error {
	return nil
}
func (i *eraseIdentityStub) DeleteUser(_ context.Context, userID string) error {
	i.deletedUserID = userID
	return nil
}
func (*eraseIdentityStub) AssignRole(context.Context, string, string) error { return nil }
func (*eraseIdentityStub) ListUsers(context.Context) ([]idp.User, error)    { return nil, nil }

type eraseAuditStub struct {
	calls      int
	entityType string
	action     string
	reason     string
}

func (a *eraseAuditStub) Record(_ context.Context, entityType string, _ uuid.UUID, action string, _ uuid.UUID, _, _ any, reason string) error {
	a.calls++
	a.entityType, a.action, a.reason = entityType, action, reason
	return nil
}

func TestEraseAnonymisesProfileAndRemovesIdentity(t *testing.T) {
	studentID := uuid.New()
	userID := uuid.New()
	actorID := uuid.New()
	store := &eraseStoreStub{record: studentRecord{ID: studentID, UserID: userID}}
	identity := &eraseIdentityStub{}
	audit := &eraseAuditStub{}
	service := NewStudentService(store, identity, nil, audit, nil)

	if err := service.Erase(context.Background(), studentID, actorID, "guardian requested erasure"); err != nil {
		t.Fatal(err)
	}
	if !store.anonymizeCalled || store.anonymizeID != studentID {
		t.Fatal("Erase did not anonymise the profile")
	}
	if !store.eraseUserCalled || store.eraseUserID != userID {
		t.Fatal("Erase did not scrub the local user record")
	}
	if identity.deletedUserID != userID.String() {
		t.Fatal("Erase did not delete the identity provider user")
	}
	if audit.calls != 1 || audit.entityType != "student_profile" || audit.action != "erased" || audit.reason != "guardian requested erasure" {
		t.Fatalf("unexpected audit record: %+v", audit)
	}
}

func TestEraseFailsClosedWhenStudentNotFound(t *testing.T) {
	store := &eraseStoreStub{recordErr: errors.New("no rows")}
	service := NewStudentService(store, &eraseIdentityStub{}, nil, &eraseAuditStub{}, nil)

	if err := service.Erase(context.Background(), uuid.New(), uuid.New(), "reason"); err == nil {
		t.Fatal("Erase should fail when the student cannot be found")
	}
	if store.anonymizeCalled {
		t.Fatal("Erase anonymised a profile it never confirmed exists")
	}
}
