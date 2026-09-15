//go:build integration

package dashboard

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/testutil/testdb"
)

func TestDashboardAnalyticsWithPostgres(t *testing.T) {
	gin.SetMode(gin.TestMode)
	pool := testdb.Open(t)
	ctx := context.Background()
	var yearID, gradeID, classID, studentID uuid.UUID
	if err := pool.QueryRow(ctx, "INSERT INTO academic_years (label, start_date, end_date, is_current) VALUES ('2026 Dashboard', '2026-01-01', '2026-12-31', TRUE) RETURNING id").Scan(&yearID); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, "INSERT INTO grades (name, sort_order) VALUES ('Dashboard Grade 10', 10) RETURNING id").Scan(&gradeID); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, "INSERT INTO classes (grade_id, academic_year_id, name) VALUES ($1, $2, 'A') RETURNING id", gradeID, yearID).Scan(&classID); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, "INSERT INTO student_profiles (full_name, index_number) VALUES ('Dashboard Student', 'DASH-STUDENT-1') RETURNING id").Scan(&studentID); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, "INSERT INTO class_students (class_id, student_id) VALUES ($1, $2)", classID, studentID); err != nil {
		t.Fatal(err)
	}

	router := gin.New()
	RegisterRoutes(router.Group(""), NewService(NewRepository(pool)))
	request := httptest.NewRequest(http.MethodGet, "/dashboard/analytics", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("dashboard analytics: code=%d body=%s", response.Code, response.Body.String())
	}
	for _, value := range []string{`"total":1`, `"label":"Dashboard Grade 10"`, `"label":"Dashboard Grade 10 A"`, `"total_classes":1`} {
		if !strings.Contains(response.Body.String(), value) {
			t.Fatalf("dashboard response missing %s: %s", value, response.Body.String())
		}
	}
}
