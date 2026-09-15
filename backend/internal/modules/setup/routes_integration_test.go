//go:build integration

package setup

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/idp"
	"github.com/openschool-org/openschool/internal/testutil/testdb"
)

type integrationIdentity struct {
	userID      uuid.UUID
	createCalls int
	assigned    bool
}

func (i *integrationIdentity) CreateUser(context.Context, string, map[string]any) (*idp.User, error) {
	i.createCalls++
	return &idp.User{ID: i.userID.String(), Email: "admin@example.test"}, nil
}
func (*integrationIdentity) UpdateUser(context.Context, string, string, map[string]any) error {
	return nil
}
func (*integrationIdentity) DeleteUser(context.Context, string) error { return nil }
func (i *integrationIdentity) AssignRole(context.Context, string, string) error {
	i.assigned = true
	return nil
}
func (*integrationIdentity) ListUsers(context.Context) ([]idp.User, error) { return nil, nil }

func TestSetupAPIWithPostgres(t *testing.T) {
	gin.SetMode(gin.TestMode)
	pool := testdb.Open(t)
	provider := &integrationIdentity{userID: uuid.New()}
	router := gin.New()
	RegisterRoutes(router.Group("/api/v1"), NewService(NewRepository(pool), provider))

	status := performSetupRequest(t, router, http.MethodGet, "/api/v1/setup/status", nil)
	if status.Code != http.StatusOK || !bytes.Contains(status.Body.Bytes(), []byte(`"needs_setup":true`)) {
		t.Fatalf("initial setup status: code=%d body=%s", status.Code, status.Body.String())
	}

	payload := RegisterAdminRequest{
		Email: "admin@example.test", Username: "admin", GivenName: "Open",
		FamilyName: "School", PhoneNumber: "0700000000", Password: "password123",
	}
	created := performSetupRequest(t, router, http.MethodPost, "/api/v1/setup/admin", payload)
	if created.Code != http.StatusCreated {
		t.Fatalf("create first admin: code=%d body=%s", created.Code, created.Body.String())
	}
	var response struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}
	if err := json.Unmarshal(created.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.ID != provider.userID.String() || response.Email != payload.Email || !provider.assigned {
		t.Fatalf("unexpected setup response/provider state: response=%+v provider=%+v", response, provider)
	}

	var role, fullName string
	if err := pool.QueryRow(context.Background(), "SELECT role, full_name FROM users WHERE id = $1", provider.userID).Scan(&role, &fullName); err != nil {
		t.Fatalf("read persisted admin: %v", err)
	}
	if role != "admin" || fullName != "Open School" {
		t.Fatalf("persisted admin role=%q full_name=%q", role, fullName)
	}

	status = performSetupRequest(t, router, http.MethodGet, "/api/v1/setup/status", nil)
	if status.Code != http.StatusOK || !bytes.Contains(status.Body.Bytes(), []byte(`"needs_setup":false`)) {
		t.Fatalf("completed setup status: code=%d body=%s", status.Code, status.Body.String())
	}
	duplicate := performSetupRequest(t, router, http.MethodPost, "/api/v1/setup/admin", payload)
	if duplicate.Code != http.StatusForbidden || provider.createCalls != 1 {
		t.Fatalf("repeat setup: code=%d create_calls=%d body=%s", duplicate.Code, provider.createCalls, duplicate.Body.String())
	}
}

func performSetupRequest(t *testing.T, router http.Handler, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var payload []byte
	if body != nil {
		var err error
		payload, err = json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
	}
	request := httptest.NewRequest(method, path, bytes.NewReader(payload))
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}
