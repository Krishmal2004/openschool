package selfservice

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

type studentProfileStoreStub struct {
	resolved uuid.UUID
	profile  StudentProfile
	err      error
}

func (s studentProfileStoreStub) studentIDForUser(context.Context, uuid.UUID) (uuid.UUID, error) {
	return s.resolved, s.err
}

func (s studentProfileStoreStub) studentWithClass(context.Context, uuid.UUID) (StudentProfile, error) {
	return s.profile, s.err
}

type teacherProfileStoreStub struct {
	resolved uuid.UUID
	year     uuid.UUID
	profile  TeacherProfile
	err      error
}

func (s teacherProfileStoreStub) teacherIDForUser(context.Context, uuid.UUID) (uuid.UUID, error) {
	return s.resolved, s.err
}

func (s teacherProfileStoreStub) teacherForUser(context.Context, uuid.UUID) (TeacherProfile, error) {
	return s.profile, s.err
}

func (s teacherProfileStoreStub) currentAcademicYearID(context.Context) (uuid.UUID, error) {
	return s.year, s.err
}

func TestStudentProfilesResolveAndProfile(t *testing.T) {
	studentID := uuid.New()
	want := StudentProfile{ID: studentID, FullName: "Student One"}
	service := NewStudentProfiles(studentProfileStoreStub{resolved: studentID, profile: want})

	resolved, err := service.Resolve(context.Background(), uuid.New())
	if err != nil || resolved != studentID {
		t.Fatalf("Resolve() = %s, %v; want %s, nil", resolved, err, studentID)
	}
	profile, err := service.Profile(context.Background(), studentID)
	if err != nil || profile != want {
		t.Fatalf("Profile() = %+v, %v; want %+v, nil", profile, err, want)
	}
}

func TestTeacherProfilesExposeResolverAndCurrentYear(t *testing.T) {
	teacherID, yearID := uuid.New(), uuid.New()
	service := NewTeacherProfiles(teacherProfileStoreStub{resolved: teacherID, year: yearID})

	resolved, err := service.Resolve(context.Background(), uuid.New())
	if err != nil || resolved != teacherID {
		t.Fatalf("Resolve() = %s, %v; want %s, nil", resolved, err, teacherID)
	}
	year, err := service.CurrentAcademicYearID(context.Background())
	if err != nil || year != yearID {
		t.Fatalf("CurrentAcademicYearID() = %s, %v; want %s, nil", year, err, yearID)
	}
}

func TestProfileServicesPropagateRepositoryErrors(t *testing.T) {
	wantErr := errors.New("query failed")
	students := NewStudentProfiles(studentProfileStoreStub{err: wantErr})
	if _, err := students.Resolve(context.Background(), uuid.New()); !errors.Is(err, wantErr) {
		t.Fatalf("Resolve() error = %v, want %v", err, wantErr)
	}

	teachers := NewTeacherProfiles(teacherProfileStoreStub{err: wantErr})
	if _, err := teachers.Profile(context.Background(), uuid.New()); !errors.Is(err, wantErr) {
		t.Fatalf("Profile() error = %v, want %v", err, wantErr)
	}
}

func TestTeacherProfileJSONPreservesDateAndNullableFields(t *testing.T) {
	profile := TeacherProfile{
		ID:       uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		UserID:   uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		FullName: "Teacher One", EmployeeNumber: "T001", JoinedDate: "2026-09-14",
		CreatedAt:        time.Date(2026, 9, 14, 8, 30, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2026, 9, 14, 8, 30, 0, 0, time.UTC),
		EmploymentStatus: "active", NICNumber: "123456789V",
	}

	payload, err := json.Marshal(profile)
	if err != nil {
		t.Fatal(err)
	}
	jsonText := string(payload)
	for _, fragment := range []string{`"joined_date":"2026-09-14"`, `"phone":null`, `"house_id":null`} {
		if !strings.Contains(jsonText, fragment) {
			t.Fatalf("profile JSON %s does not contain %s", jsonText, fragment)
		}
	}
}
