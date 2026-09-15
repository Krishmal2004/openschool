//go:build integration

package school

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/openschool-org/openschool/internal/testutil/testdb"
)

func TestSchoolCalendarAndStructureAPIWithPostgres(t *testing.T) {
	gin.SetMode(gin.TestMode)
	pool := testdb.Open(t)
	router := gin.New()
	group := router.Group("")
	RegisterSchoolRoutes(group, group, group, pool)
	RegisterGradeRoutes(group, group, pool)
	RegisterTermRoutes(group, group, pool)
	RegisterHouseRoutes(group, group, NewHouseService(pool, nil))

	invalidSchool := performSchoolRequest(t, router, http.MethodPost, "/school", schoolCommand{Name: "Open School", Phone: "123"})
	if invalidSchool.Code != http.StatusBadRequest {
		t.Fatalf("invalid school phone: code=%d body=%s", invalidSchool.Code, invalidSchool.Body.String())
	}
	gradeFrom, gradeTo := int32(1), int32(13)
	createdSchool := performSchoolRequest(t, router, http.MethodPost, "/school", schoolCommand{Name: "Open School", Address: "Colombo", Phone: "0771234567", Email: "office@school.test", GradeFrom: &gradeFrom, GradeTo: &gradeTo})
	if createdSchool.Code != http.StatusCreated {
		t.Fatalf("create school: code=%d body=%s", createdSchool.Code, createdSchool.Body.String())
	}
	var school School
	if err := json.Unmarshal(createdSchool.Body.Bytes(), &school); err != nil {
		t.Fatal(err)
	}
	if school.ID == "" || school.SchoolType != "mixed" {
		t.Fatalf("unexpected school: %+v", school)
	}
	duplicateSchool := performSchoolRequest(t, router, http.MethodPost, "/school", schoolCommand{Name: "Duplicate", Phone: "0771234567"})
	if duplicateSchool.Code != http.StatusBadRequest {
		t.Fatalf("duplicate school: code=%d body=%s", duplicateSchool.Code, duplicateSchool.Body.String())
	}
	updatedSchool := performSchoolRequest(t, router, http.MethodPut, "/school/"+school.ID, schoolCommand{Name: "Open School Updated", Phone: "0712345678", SchoolType: "mixed"})
	if updatedSchool.Code != http.StatusOK || !bytes.Contains(updatedSchool.Body.Bytes(), []byte(`"name":"Open School Updated"`)) {
		t.Fatalf("update school: code=%d body=%s", updatedSchool.Code, updatedSchool.Body.String())
	}

	yearOne := createAcademicYear(t, router, "2026 School", 2026)
	yearTwo := createAcademicYear(t, router, "2027 School", 2027)
	setCurrent := performSchoolRequest(t, router, http.MethodPut, "/academic-years/"+yearTwo.ID+"/set-current", nil)
	if setCurrent.Code != http.StatusOK {
		t.Fatalf("set current academic year: code=%d body=%s", setCurrent.Code, setCurrent.Body.String())
	}
	currentYear := performSchoolRequest(t, router, http.MethodGet, "/academic-years/current", nil)
	if currentYear.Code != http.StatusOK || !bytes.Contains(currentYear.Body.Bytes(), []byte(yearTwo.ID)) {
		t.Fatalf("get current academic year: code=%d body=%s", currentYear.Code, currentYear.Body.String())
	}
	assertSingleCurrent(t, pool, "academic_years", yearTwo.ID)

	createdGrade := performSchoolRequest(t, router, http.MethodPost, "/grades", gradeCommand{Name: "Grade 6 School", SortOrder: 6})
	if createdGrade.Code != http.StatusCreated {
		t.Fatalf("create grade: code=%d body=%s", createdGrade.Code, createdGrade.Body.String())
	}
	var grade Grade
	if err := json.Unmarshal(createdGrade.Body.Bytes(), &grade); err != nil {
		t.Fatal(err)
	}
	updatedGrade := performSchoolRequest(t, router, http.MethodPut, "/grades/"+grade.ID.String(), gradeCommand{Name: "Grade 6 Updated", SortOrder: 7})
	if updatedGrade.Code != http.StatusOK || !bytes.Contains(updatedGrade.Body.Bytes(), []byte(`"sort_order":7`)) {
		t.Fatalf("update grade: code=%d body=%s", updatedGrade.Code, updatedGrade.Body.String())
	}

	termOne := createTerm(t, router, yearTwo.ID, "Term 1", 1, 4, 1)
	termTwo := createTerm(t, router, yearTwo.ID, "Term 2", 5, 8, 2)
	currentTerm := performSchoolRequest(t, router, http.MethodPut, "/terms/"+termTwo.ID+"/set-current", nil)
	if currentTerm.Code != http.StatusOK {
		t.Fatalf("set current term: code=%d body=%s", currentTerm.Code, currentTerm.Body.String())
	}
	assertSingleCurrent(t, pool, "terms", termTwo.ID)
	deletedTerm := performSchoolRequest(t, router, http.MethodDelete, "/terms/"+termOne.ID, nil)
	if deletedTerm.Code != http.StatusOK {
		t.Fatalf("delete unused term: code=%d body=%s", deletedTerm.Code, deletedTerm.Body.String())
	}

	createdHouse := performSchoolRequest(t, router, http.MethodPost, "/houses", houseCommand{Name: "Vijaya", Code: "VJ", Color: "#16a34a"})
	if createdHouse.Code != http.StatusCreated {
		t.Fatalf("create house: code=%d body=%s", createdHouse.Code, createdHouse.Body.String())
	}
	var house House
	if err := json.Unmarshal(createdHouse.Body.Bytes(), &house); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(context.Background(), "INSERT INTO classes (grade_id, academic_year_id, name) VALUES ($1, $2, 'A')", grade.ID, yearTwo.ID); err != nil {
		t.Fatal(err)
	}
	gradeInUse := performSchoolRequest(t, router, http.MethodDelete, "/grades/"+grade.ID.String(), nil)
	if gradeInUse.Code != http.StatusConflict {
		t.Fatalf("delete grade in use: code=%d body=%s", gradeInUse.Code, gradeInUse.Body.String())
	}
	yearInUse := performSchoolRequest(t, router, http.MethodDelete, "/academic-years/"+yearTwo.ID, nil)
	if yearInUse.Code != http.StatusConflict {
		t.Fatalf("delete academic year in use: code=%d body=%s", yearInUse.Code, yearInUse.Body.String())
	}
	deletedYear := performSchoolRequest(t, router, http.MethodDelete, "/academic-years/"+yearOne.ID, nil)
	if deletedYear.Code != http.StatusOK {
		t.Fatalf("delete unused academic year: code=%d body=%s", deletedYear.Code, deletedYear.Body.String())
	}

	seedHouseMembers(t, pool)
	studentsAssigned := performSchoolRequest(t, router, http.MethodPost, "/houses/reassign-missing", nil)
	staffAssigned := performSchoolRequest(t, router, http.MethodPost, "/houses/reassign-missing-staff", nil)
	if studentsAssigned.Code != http.StatusOK || !bytes.Contains(studentsAssigned.Body.Bytes(), []byte(`"assigned":1`)) {
		t.Fatalf("reassign students: code=%d body=%s", studentsAssigned.Code, studentsAssigned.Body.String())
	}
	if staffAssigned.Code != http.StatusOK || !bytes.Contains(staffAssigned.Body.Bytes(), []byte(`"assigned":1`)) {
		t.Fatalf("reassign teachers: code=%d body=%s", staffAssigned.Code, staffAssigned.Body.String())
	}
	assertHouseMembersAssigned(t, pool)
	houseInUse := performSchoolRequest(t, router, http.MethodDelete, "/houses/"+house.ID.String(), nil)
	if houseInUse.Code != http.StatusConflict {
		t.Fatalf("delete house in use: code=%d body=%s", houseInUse.Code, houseInUse.Body.String())
	}
}

func performSchoolRequest(t *testing.T, router http.Handler, method, path string, body any) *httptest.ResponseRecorder {
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

func createAcademicYear(t *testing.T, router http.Handler, label string, year int) AcademicYear {
	t.Helper()
	response := performSchoolRequest(t, router, http.MethodPost, "/academic-years", yearCommand{Label: label, StartDate: time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC), EndDate: time.Date(year, 12, 31, 0, 0, 0, 0, time.UTC)})
	if response.Code != http.StatusCreated {
		t.Fatalf("create academic year: code=%d body=%s", response.Code, response.Body.String())
	}
	var result AcademicYear
	if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	return result
}

func createTerm(t *testing.T, router http.Handler, yearID, name string, startMonth, endMonth int, sortOrder int32) Term {
	t.Helper()
	response := performSchoolRequest(t, router, http.MethodPost, "/terms", createTermCommand{AcademicYearID: yearID, Name: name, StartDate: time.Date(2027, time.Month(startMonth), 1, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2027, time.Month(endMonth), 28, 0, 0, 0, 0, time.UTC), SortOrder: sortOrder})
	if response.Code != http.StatusCreated {
		t.Fatalf("create term: code=%d body=%s", response.Code, response.Body.String())
	}
	var result Term
	if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	return result
}

func assertSingleCurrent(t *testing.T, pool *pgxpool.Pool, table, wantID string) {
	t.Helper()
	var count int
	var id uuid.UUID
	query := "SELECT COUNT(*) OVER (), id FROM " + table + " WHERE is_current = TRUE"
	if err := pool.QueryRow(context.Background(), query).Scan(&count, &id); err != nil {
		t.Fatal(err)
	}
	if count != 1 || id.String() != wantID {
		t.Fatalf("%s current count=%d id=%s, want %s", table, count, id, wantID)
	}
}

func seedHouseMembers(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	teacherUserID := uuid.New()
	if _, err := pool.Exec(ctx, "INSERT INTO users (id, email, full_name, role) VALUES ($1, 'house-teacher@example.test', 'House Teacher', 'teacher')", teacherUserID); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, "INSERT INTO teacher_profiles (user_id, full_name, employee_number, joined_date, nic_number) VALUES ($1, 'House Teacher', 'HOUSE-T-001', '2020-01-01', 'HOUSE-NIC-T-001')", teacherUserID); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, "INSERT INTO student_profiles (full_name, index_number) VALUES ('House Student', 'HOUSE-S-001')"); err != nil {
		t.Fatal(err)
	}
}

func assertHouseMembersAssigned(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	var missingStudents, missingTeachers int
	if err := pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM student_profiles WHERE house_id IS NULL").Scan(&missingStudents); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM teacher_profiles WHERE house_id IS NULL").Scan(&missingTeachers); err != nil {
		t.Fatal(err)
	}
	if missingStudents != 0 || missingTeachers != 0 {
		t.Fatalf("missing house assignments: students=%d teachers=%d", missingStudents, missingTeachers)
	}
}
