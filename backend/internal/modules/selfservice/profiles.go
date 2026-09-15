package selfservice

import (
	"context"

	"github.com/google/uuid"
)

type teacherProfileStore interface {
	teacherIDForUser(context.Context, uuid.UUID) (uuid.UUID, error)
	teacherForUser(context.Context, uuid.UUID) (TeacherProfile, error)
	currentAcademicYearID(context.Context) (uuid.UUID, error)
}

type studentProfileStore interface {
	studentIDForUser(context.Context, uuid.UUID) (uuid.UUID, error)
	studentWithClass(context.Context, uuid.UUID) (StudentProfile, error)
}

// TeacherProfiles resolves the authenticated account to its teacher profile.
// It also satisfies Attendance's narrow TeacherResolver contract.
type TeacherProfiles struct {
	repo teacherProfileStore
}

func NewTeacherProfiles(repo teacherProfileStore) *TeacherProfiles {
	return &TeacherProfiles{repo: repo}
}

func (s *TeacherProfiles) Resolve(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	return s.repo.teacherIDForUser(ctx, userID)
}

func (s *TeacherProfiles) Profile(ctx context.Context, userID uuid.UUID) (TeacherProfile, error) {
	return s.repo.teacherForUser(ctx, userID)
}

func (s *TeacherProfiles) CurrentAcademicYearID(ctx context.Context) (uuid.UUID, error) {
	return s.repo.currentAcademicYearID(ctx)
}

// StudentProfiles resolves an authenticated account and returns its enriched
// student profile, including current class, grade, house, and academic year.
type StudentProfiles struct {
	repo studentProfileStore
}

func NewStudentProfiles(repo studentProfileStore) *StudentProfiles {
	return &StudentProfiles{repo: repo}
}

func (s *StudentProfiles) Resolve(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	return s.repo.studentIDForUser(ctx, userID)
}

func (s *StudentProfiles) Profile(ctx context.Context, studentID uuid.UUID) (StudentProfile, error) {
	return s.repo.studentWithClass(ctx, studentID)
}
