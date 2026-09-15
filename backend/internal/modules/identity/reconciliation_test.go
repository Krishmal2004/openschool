package identity

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	idp "github.com/openschool-org/openschool/internal/idp"
)

var errReconciliationTest = errors.New("reconciliation test error")

type reconciliationProviderStub struct {
	users        []idp.User
	listErr      error
	deleteErr    error
	deletedID    string
	deleteCalled bool
}

func (s *reconciliationProviderStub) ListUsers(context.Context) ([]idp.User, error) {
	return s.users, s.listErr
}

func (s *reconciliationProviderStub) DeleteUser(_ context.Context, id string) error {
	s.deleteCalled = true
	s.deletedID = id
	return s.deleteErr
}

type reconciliationUsersStub struct {
	ids          []uuid.UUID
	listErr      error
	existsResult bool
	existsErr    error
}

func (s *reconciliationUsersStub) listIDs(context.Context) ([]uuid.UUID, error) {
	return s.ids, s.listErr
}

func (s *reconciliationUsersStub) exists(context.Context, uuid.UUID) (bool, error) {
	return s.existsResult, s.existsErr
}

type reconciliationAuditStub struct {
	called     bool
	entityType string
	action     string
	actorID    uuid.UUID
}

func (s *reconciliationAuditStub) Record(_ context.Context, entityType string, _ uuid.UUID, action string, actorID uuid.UUID, _, _ interface{}, _ string) error {
	s.called = true
	s.entityType = entityType
	s.action = action
	s.actorID = actorID
	return nil
}

func TestFindOrphanedReturnsOnlyProviderAccountsWithoutLocalUsers(t *testing.T) {
	localID := uuid.New()
	provider := &reconciliationProviderStub{users: []idp.User{
		{ID: localID.String(), Username: "local"},
		{ID: "provider-only", Username: "orphan", Email: "orphan@example.test"},
	}}
	service := newReconciliationService(provider, &reconciliationUsersStub{ids: []uuid.UUID{localID}}, nil)

	result, err := service.findOrphaned(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 || result[0].ID != "provider-only" || result[0].Username != "orphan" || result[0].Email != "orphan@example.test" {
		t.Fatalf("unexpected orphaned accounts: %#v", result)
	}
}

func TestFindOrphanedPropagatesProviderAndRepositoryFailures(t *testing.T) {
	tests := []struct {
		name     string
		provider *reconciliationProviderStub
		users    *reconciliationUsersStub
	}{
		{name: "provider", provider: &reconciliationProviderStub{listErr: errReconciliationTest}, users: &reconciliationUsersStub{}},
		{name: "repository", provider: &reconciliationProviderStub{}, users: &reconciliationUsersStub{listErr: errReconciliationTest}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := newReconciliationService(test.provider, test.users, nil).findOrphaned(context.Background())
			if !errors.Is(err, errReconciliationTest) {
				t.Fatalf("error = %v, want reconciliation failure", err)
			}
		})
	}
}

func TestDeleteOrphanedRechecksLocalUserBeforeProviderDeletion(t *testing.T) {
	provider := &reconciliationProviderStub{}
	service := newReconciliationService(provider, &reconciliationUsersStub{existsResult: true}, nil)

	err := service.deleteOrphaned(context.Background(), uuid.NewString(), uuid.New())
	if !errors.Is(err, ErrOrphanNoLongerOrphaned) {
		t.Fatalf("error = %v, want conflict", err)
	}
	if provider.deleteCalled {
		t.Fatal("provider account was deleted despite a matching local user")
	}
}

func TestDeleteOrphanedFailsClosedWhenLocalLookupFails(t *testing.T) {
	provider := &reconciliationProviderStub{}
	service := newReconciliationService(provider, &reconciliationUsersStub{existsErr: errReconciliationTest}, nil)

	err := service.deleteOrphaned(context.Background(), uuid.NewString(), uuid.New())
	if !errors.Is(err, errReconciliationTest) {
		t.Fatalf("error = %v, want repository failure", err)
	}
	if provider.deleteCalled {
		t.Fatal("provider account was deleted while local identity state was unavailable")
	}
}

func TestDeleteOrphanedDeletesAndAuditsConfirmedOrphan(t *testing.T) {
	provider := &reconciliationProviderStub{}
	audit := &reconciliationAuditStub{}
	actorID := uuid.New()
	providerID := uuid.NewString()
	service := newReconciliationService(provider, &reconciliationUsersStub{}, audit)

	if err := service.deleteOrphaned(context.Background(), providerID, actorID); err != nil {
		t.Fatal(err)
	}
	if !provider.deleteCalled || provider.deletedID != providerID {
		t.Fatalf("deleted provider ID = %q, want %q", provider.deletedID, providerID)
	}
	if !audit.called || audit.entityType != "orphaned_identity" || audit.action != "deleted" || audit.actorID != actorID {
		t.Fatalf("unexpected audit call: %#v", audit)
	}
}

func TestDeleteOrphanedDoesNotAuditFailedProviderDeletion(t *testing.T) {
	provider := &reconciliationProviderStub{deleteErr: errReconciliationTest}
	audit := &reconciliationAuditStub{}
	service := newReconciliationService(provider, &reconciliationUsersStub{}, audit)

	err := service.deleteOrphaned(context.Background(), uuid.NewString(), uuid.New())
	if !errors.Is(err, errReconciliationTest) {
		t.Fatalf("error = %v, want provider failure", err)
	}
	if audit.called {
		t.Fatal("failed provider deletion was audited as successful")
	}
}
