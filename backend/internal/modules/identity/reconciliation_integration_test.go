//go:build integration

package identity

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	idp "github.com/openschool-org/openschool/internal/idp"
	auditmodule "github.com/openschool-org/openschool/internal/modules/audit"
	"github.com/openschool-org/openschool/internal/testutil/testdb"
)

type integrationIdentityProvider struct {
	users     []idp.User
	deletedID string
}

func (p *integrationIdentityProvider) ListUsers(context.Context) ([]idp.User, error) {
	return append([]idp.User(nil), p.users...), nil
}

func (p *integrationIdentityProvider) DeleteUser(_ context.Context, id string) error {
	p.deletedID = id
	return nil
}

func (*integrationIdentityProvider) CreateUser(context.Context, string, map[string]any) (*idp.User, error) {
	return nil, nil
}

func (*integrationIdentityProvider) UpdateUser(context.Context, string, string, map[string]any) error {
	return nil
}

func (*integrationIdentityProvider) AssignRole(context.Context, string, string) error { return nil }

func TestIdentityReconciliationAPIWithPostgres(t *testing.T) {
	gin.SetMode(gin.TestMode)
	pool := testdb.Open(t)
	actorID := uuid.New()
	localID := uuid.New()
	orphanID := uuid.New()
	if _, err := pool.Exec(context.Background(), `
		INSERT INTO users (id, email, full_name, role) VALUES
		($1, 'reconciliation-admin@example.test', 'Reconciliation Admin', 'admin'),
		($2, 'local-identity@example.test', 'Local Identity', 'student')`, actorID, localID); err != nil {
		t.Fatal(err)
	}
	provider := &integrationIdentityProvider{users: []idp.User{
		{ID: localID.String(), Username: "local", Email: "local-identity@example.test"},
		{ID: orphanID.String(), Username: "orphan", Email: "orphan@example.test"},
	}}
	router := gin.New()
	admin := router.Group("")
	admin.Use(func(c *gin.Context) {
		c.Set("userID", actorID.String())
		c.Next()
	})
	RegisterReconciliation(admin, pool, provider, auditmodule.NewService(auditmodule.NewRepository(pool)))

	listed := performReconciliationRequest(router, http.MethodGet, "/admin/orphaned-accounts")
	if listed.Code != http.StatusOK || !bytes.Contains(listed.Body.Bytes(), []byte(orphanID.String())) || bytes.Contains(listed.Body.Bytes(), []byte(localID.String())) {
		t.Fatalf("list orphaned identities: code=%d body=%s", listed.Code, listed.Body.String())
	}

	protected := performReconciliationRequest(router, http.MethodDelete, "/admin/orphaned-accounts/"+localID.String())
	if protected.Code != http.StatusConflict {
		t.Fatalf("delete linked identity: code=%d body=%s", protected.Code, protected.Body.String())
	}
	if provider.deletedID != "" {
		t.Fatalf("provider deleted linked identity %q", provider.deletedID)
	}

	deleted := performReconciliationRequest(router, http.MethodDelete, "/admin/orphaned-accounts/"+orphanID.String())
	if deleted.Code != http.StatusOK {
		t.Fatalf("delete orphaned identity: code=%d body=%s", deleted.Code, deleted.Body.String())
	}
	if provider.deletedID != orphanID.String() {
		t.Fatalf("provider deleted ID=%q, want %q", provider.deletedID, orphanID)
	}
	assertReconciliationAudit(t, pool, actorID, orphanID.String())
}

func performReconciliationRequest(router http.Handler, method, path string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, path, nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

func assertReconciliationAudit(t *testing.T, pool *pgxpool.Pool, actorID uuid.UUID, providerID string) {
	t.Helper()
	var entityType, action string
	var gotActorID uuid.UUID
	var before []byte
	if err := pool.QueryRow(context.Background(), "SELECT entity_type, action, actor_id, before FROM audit_logs WHERE entity_type = 'orphaned_identity'").Scan(&entityType, &action, &gotActorID, &before); err != nil {
		t.Fatal(err)
	}
	if entityType != "orphaned_identity" || action != "deleted" || gotActorID != actorID || !bytes.Contains(before, []byte(providerID)) {
		t.Fatalf("unexpected reconciliation audit: entity=%q action=%q actor=%s before=%s", entityType, action, gotActorID, before)
	}
}
