//go:build integration

package automation

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/openschool-org/openschool/internal/testutil/testdb"
)

func TestAutomationSettingsWithPostgres(t *testing.T) {
	gin.SetMode(gin.TestMode)
	pool := testdb.Open(t)
	router := gin.New()
	RegisterRoutes(router.Group(""), pool)

	list := automationRequest(router, http.MethodGet, "/jobs")
	if list.Code != http.StatusOK || !strings.Contains(list.Body.String(), PeopleComplianceAgentName) || !strings.Contains(list.Body.String(), SystemHealthAgentName) {
		t.Fatalf("list jobs: code=%d body=%s", list.Code, list.Body.String())
	}
	disabled := automationRequest(router, http.MethodPut, "/jobs/"+PeopleComplianceAgentName+"/enabled", []byte(`{"enabled":false}`))
	if disabled.Code != http.StatusOK {
		t.Fatalf("disable people-compliance job: code=%d body=%s", disabled.Code, disabled.Body.String())
	}
	var enabled bool
	if err := pool.QueryRow(context.Background(), "SELECT enabled FROM job_settings WHERE job_name = $1", PeopleComplianceAgentName).Scan(&enabled); err != nil || enabled {
		t.Fatalf("persisted automation setting: enabled=%t err=%v", enabled, err)
	}
	cannotDisableHealth := automationRequest(router, http.MethodPut, "/jobs/"+SystemHealthAgentName+"/enabled", []byte(`{"enabled":false}`))
	if cannotDisableHealth.Code != http.StatusBadRequest {
		t.Fatalf("disable system health job: code=%d body=%s", cannotDisableHealth.Code, cannotDisableHealth.Body.String())
	}
}

func automationRequest(router http.Handler, method, path string, body ...[]byte) *httptest.ResponseRecorder {
	var payload []byte
	if len(body) > 0 {
		payload = body[0]
	}
	request := httptest.NewRequest(method, path, bytes.NewReader(payload))
	if len(payload) > 0 {
		request.Header.Set("Content-Type", "application/json")
	}
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}
