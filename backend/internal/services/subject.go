package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	db "github.com/openschool-org/openschool/db/sqlc"
	"github.com/openschool-org/openschool/internal/models"
	"github.com/openschool-org/openschool/internal/repositories"
)

var (
	ErrSubjectNotFound = errors.New("subject not found")
	ErrSubjectInUse    = errors.New("subject is in use and cannot be deleted")
)

type SubjectService struct {
	repo *repositories.SubjectRepository
}

func NewSubjectService(repo *repositories.SubjectRepository) *SubjectService {
	return &SubjectService{repo: repo}
}

func (s *SubjectService) CreateSubject(ctx context.Context, req models.CreateSubjectRequest) (models.SubjectResponse, error) {
	maxMarks := 100.00
	if req.MaxMarks != nil {
		maxMarks = *req.MaxMarks
	}
	maxMarksNumeric, err := pgNumeric(maxMarks)
	if err != nil {
		return models.SubjectResponse{}, err
	}
	subject, err := s.repo.Create(ctx, db.CreateSubjectParams{
		Name:     req.Name,
		Code:     req.Code,
		Type:     optionalText(req.Type),
		MaxMarks: maxMarksNumeric,
	})
	return toSubjectResponse(subject), err
}

func (s *SubjectService) GetSubject(ctx context.Context, id uuid.UUID) (models.SubjectResponse, error) {
	subject, err := s.repo.GetByID(ctx, id)
	return toSubjectResponse(subject), err
}

func (s *SubjectService) ListSubjects(ctx context.Context) ([]models.SubjectResponse, error) {
	subjects, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	resp := make([]models.SubjectResponse, len(subjects))
	for i, subject := range subjects {
		resp[i] = toSubjectResponse(subject)
	}
	return resp, nil
}

func (s *SubjectService) UpdateSubject(ctx context.Context, id uuid.UUID, req models.UpdateSubjectRequest) (models.SubjectResponse, error) {
	maxMarks := 100.00
	if req.MaxMarks != nil {
		maxMarks = *req.MaxMarks
	}
	maxMarksNumeric, err := pgNumeric(maxMarks)
	if err != nil {
		return models.SubjectResponse{}, err
	}
	subject, err := s.repo.Update(ctx, db.UpdateSubjectParams{
		ID:       id,
		Name:     req.Name,
		Code:     req.Code,
		Type:     optionalText(req.Type),
		MaxMarks: maxMarksNumeric,
	})
	return toSubjectResponse(subject), err
}

func toSubjectResponse(s db.Subject) models.SubjectResponse {
	var subjectType *string
	if s.Type.Valid {
		subjectType = &s.Type.String
	}
	maxMarks := float64(0)
	if value, err := s.MaxMarks.Float64Value(); err == nil && value.Valid {
		maxMarks = value.Float64
	}
	return models.SubjectResponse{ID: s.ID.String(), Name: s.Name, Code: s.Code, Type: subjectType, MaxMarks: maxMarks, CreatedAt: s.CreatedAt.Time.String()}
}

func (s *SubjectService) DeleteSubject(ctx context.Context, id uuid.UUID) error {
	rows, err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	if rows == 0 {
		if _, err := s.repo.GetByID(ctx, id); err != nil {
			return ErrSubjectNotFound
		}
		return ErrSubjectInUse
	}
	return nil
}
