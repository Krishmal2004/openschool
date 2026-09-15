//go:build integration

package search

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/testutil/testdb"
)

func TestGlobalSearchWithPostgres(t *testing.T) {
	gin.SetMode(gin.TestMode)
	pool := testdb.Open(t)
	ctx := context.Background()
	teacherUserID := uuid.New()
	if _, err := pool.Exec(ctx, "INSERT INTO users (id, email, full_name, role) VALUES ($1, 'search-teacher@example.test', 'Search Teacher', 'teacher')", teacherUserID); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, "INSERT INTO teacher_profiles (user_id, full_name, employee_number, joined_date, nic_number) VALUES ($1, 'Search Shared Teacher', 'SEARCH-T-1', $2, '199000000021')", teacherUserID, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, "INSERT INTO student_profiles (full_name, index_number) VALUES ('Search Shared Student', 'SEARCH-S-1')"); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, "INSERT INTO guardians (full_name, relationship, phone, nic_number) VALUES ('Search Shared Guardian', 'guardian', '0770000021', '199000000022')"); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, "INSERT INTO non_academic_staff (full_name, employee_number, designation, joined_date) VALUES ('Search Shared Staff', 'SEARCH-N-1', 'office_staff', $1)", time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)); err != nil {
		t.Fatal(err)
	}

	router := gin.New()
	RegisterRoutes(router.Group(""), NewService(NewRepository(pool)))
	empty := searchRequest(router, "")
	if empty.Code != http.StatusOK || empty.Body.String() != `{"students":[],"teachers":[],"guardians":[],"non_academic_staff":[]}` {
		t.Fatalf("empty search: code=%d body=%s", empty.Code, empty.Body.String())
	}
	result := searchRequest(router, "Shared")
	if result.Code != http.StatusOK {
		t.Fatalf("search: code=%d body=%s", result.Code, result.Body.String())
	}
	for _, value := range []string{"Search Shared Student", "Search Shared Teacher", "Search Shared Guardian", "Search Shared Staff"} {
		if !strings.Contains(result.Body.String(), value) {
			t.Fatalf("search result missing %q: %s", value, result.Body.String())
		}
	}
}

func searchRequest(router http.Handler, term string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodGet, "/admin/search?q="+term, nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}
