package notifications

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	rootmodels "github.com/openschool-org/openschool/internal/models"
	models "github.com/openschool-org/openschool/internal/models/notifications"
)

type fakeStore struct {
	store
	created    int
	recipients []uuid.UUID
	students   []pgtype.UUID
	guardians  []pgtype.UUID
	teachers   []uuid.UUID
}

func (f *fakeStore) Create(_ context.Context, command notificationCommand) (notification, error) {
	f.created++
	return notification{ID: uuid.New(), Title: command.Title, Message: command.Message, Category: command.Category, Priority: command.Priority, Status: command.Status, RecipientRules: command.RecipientRules, CreatedBy: command.CreatedBy, SentAt: command.SentAt}, nil
}
func (f *fakeStore) AddRecipient(_ context.Context, _ uuid.UUID, userID uuid.UUID) error {
	f.recipients = append(f.recipients, userID)
	return nil
}
func (f *fakeStore) ListStudentUserIDsByClass(context.Context, uuid.UUID) ([]pgtype.UUID, error) {
	return f.students, nil
}
func (f *fakeStore) ListGuardianUserIDsByClass(context.Context, uuid.UUID) ([]pgtype.UUID, error) {
	return f.guardians, nil
}
func (f *fakeStore) ListTeacherUserIDsByClass(context.Context, uuid.UUID) ([]uuid.UUID, error) {
	return f.teachers, nil
}

func pgUUID(id uuid.UUID) pgtype.UUID { return pgtype.UUID{Bytes: id, Valid: true} }

func TestNonTeacherCannotComposeNotifications(t *testing.T) {
	service := NewNotificationService(&fakeStore{})
	err := service.authorizeSender(context.Background(), rootmodels.RoleStudent, uuid.New(), []models.RecipientRule{{Type: models.RuleTeacher, TeacherID: pointer(uuid.New())}})
	if !errors.Is(err, ErrForbiddenRecipients) {
		t.Fatalf("expected forbidden recipients, got %v", err)
	}
}

func TestClassRecipientResolutionDeduplicatesUsers(t *testing.T) {
	shared, student, teacher := uuid.New(), uuid.New(), uuid.New()
	store := &fakeStore{students: []pgtype.UUID{pgUUID(student), pgUUID(shared)}, guardians: []pgtype.UUID{pgUUID(shared)}, teachers: []uuid.UUID{teacher, shared}}
	service := NewNotificationService(store)
	classID := uuid.New()
	ids, err := service.resolveRecipientUserIDs(context.Background(), []models.RecipientRule{{Type: models.RuleClass, ClassID: &classID}})
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 3 {
		t.Fatalf("expected three unique recipients, got %d: %v", len(ids), ids)
	}
}

func TestSendDirectValidatesAndMaterializesRecipients(t *testing.T) {
	store := &fakeStore{}
	service := NewNotificationService(store)
	if err := service.SendDirect(context.Background(), "Title", "Message", "bad-category", "normal", uuid.New(), []uuid.UUID{uuid.New()}); err == nil {
		t.Fatal("expected invalid category error")
	}
	users := []uuid.UUID{uuid.New(), uuid.New()}
	if err := service.SendDirect(context.Background(), "Title", "Message", "general", "normal", uuid.New(), users); err != nil {
		t.Fatal(err)
	}
	if store.created != 1 || len(store.recipients) != 2 {
		t.Fatalf("created=%d recipients=%d", store.created, len(store.recipients))
	}
}

func pointer(value uuid.UUID) *uuid.UUID { return &value }
