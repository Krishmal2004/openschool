package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/repositories"
)

// TeacherSelfService owns the signed-in teacher profile lookup used by self-service endpoints.
type TeacherSelfService struct {
	repo *repositories.TeacherRepository
}

func NewTeacherSelfService(repo *repositories.TeacherRepository) *TeacherSelfService {
	return &TeacherSelfService{repo: repo}
}

func (s *TeacherSelfService) Resolve(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	teacher, err := s.repo.GetByUserID(ctx, userID)
	return teacher.ID, err
}

func (s *TeacherSelfService) Profile(ctx context.Context, userID uuid.UUID) (interface{}, error) {
	return s.repo.GetByUserID(ctx, userID)
}

// StudentSelfService owns the signed-in student profile lookup used by self-service endpoints.
type StudentSelfService struct {
	repo *repositories.StudentRepository
}

func NewStudentSelfService(repo *repositories.StudentRepository) *StudentSelfService {
	return &StudentSelfService{repo: repo}
}

func (s *StudentSelfService) Resolve(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	student, err := s.repo.GetByUserID(ctx, userID)
	return student.ID, err
}

func (s *StudentSelfService) Profile(ctx context.Context, studentID uuid.UUID) (interface{}, error) {
	return s.repo.GetWithClass(ctx, studentID)
}
