package identity

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	idp "github.com/openschool-org/openschool/internal/idp"
)

type reconciliationRouteProvider struct{ reconciliationProviderStub }

func (*reconciliationRouteProvider) CreateUser(context.Context, string, map[string]any) (*idp.User, error) {
	return nil, nil
}

func (*reconciliationRouteProvider) UpdateUser(context.Context, string, string, map[string]any) error {
	return nil
}

func (*reconciliationRouteProvider) AssignRole(context.Context, string, string) error { return nil }

func TestRegisterReconciliationPreservesAdminEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api := router.Group("/api/v1")
	RegisterReconciliation(api, nil, &reconciliationRouteProvider{}, nil)

	want := map[string]bool{
		"GET /api/v1/admin/orphaned-accounts":        true,
		"DELETE /api/v1/admin/orphaned-accounts/:id": true,
	}
	for _, route := range router.Routes() {
		delete(want, route.Method+" "+route.Path)
	}
	if len(want) != 0 {
		t.Fatalf("missing identity reconciliation routes: %v", want)
	}
}
