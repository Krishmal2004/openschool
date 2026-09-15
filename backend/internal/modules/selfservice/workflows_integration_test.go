//go:build integration

package selfservice

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

func TestStudentSelfServiceProfileWithPostgres(t *testing.T) {
	gin.SetMode(gin.TestMode)
	pool := testdb.Open(t)
	ctx := context.Background()
	studentUserID := uuid.New()
	if _, err := pool.Exec(ctx, "INSERT INTO users (id, email, full_name, role) VALUES ($1, 'portal-student@example.test', 'Portal Student', 'student')", studentUserID); err != nil {
		t.Fatal(err)
	}
	var yearID, gradeID, classID, studentID uuid.UUID
	if err := pool.QueryRow(ctx, "INSERT INTO academic_years (label, start_date, end_date, is_current) VALUES ('2026 Portal', '2026-01-01', '2026-12-31', TRUE) RETURNING id").Scan(&yearID); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, "INSERT INTO grades (name, sort_order) VALUES ('Portal Grade 7', 7) RETURNING id").Scan(&gradeID); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, "INSERT INTO classes (grade_id, academic_year_id, name) VALUES ($1, $2, 'A') RETURNING id", gradeID, yearID).Scan(&classID); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, "INSERT INTO student_profiles (user_id, full_name, index_number, phone) VALUES ($1, 'Portal Student', 'PORTAL-S-1', '0770000041') RETURNING id", studentUserID).Scan(&studentID); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, "INSERT INTO class_students (class_id, student_id) VALUES ($1, $2)", classID, studentID); err != nil {
		t.Fatal(err)
	}

	router := gin.New()
	student := router.Group("")
	student.Use(func(c *gin.Context) { c.Set("userID", studentUserID.String()); c.Next() })
	registerStudentRoutes(student, NewStudentProfiles(NewRepository(pool)), pool)
	request := httptest.NewRequest(http.MethodGet, "/me/student", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), `"full_name":"Portal Student"`) || !strings.Contains(response.Body.String(), `"class_name":"A"`) || !strings.Contains(response.Body.String(), `"grade_name":"Portal Grade 7"`) {
		t.Fatalf("student portal profile: code=%d body=%s", response.Code, response.Body.String())
	}
}
