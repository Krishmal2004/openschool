package people

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/openschool-org/openschool/internal/models"
)

var (
	ErrInvalidActivityCategory     = errors.New("invalid activity category")
	ErrInvalidDisciplinarySeverity = errors.New("invalid severity")
	ErrPortfolioRecordNotFound     = errors.New("record not found")
)

type portfolioStore interface {
	teacherProfileID(context.Context, uuid.UUID) (uuid.UUID, error)
	createProgress(context.Context, uuid.UUID, uuid.UUID, string, *uuid.UUID) (any, error)
	listProgress(context.Context, uuid.UUID) (any, error)
	updateProgress(context.Context, uuid.UUID, uuid.UUID, string) (any, error)
	deleteProgress(context.Context, uuid.UUID, uuid.UUID) (int64, error)
	createActivity(context.Context, uuid.UUID, uuid.UUID, string, string, string, string) (any, error)
	listActivities(context.Context, uuid.UUID) (any, error)
	updateActivity(context.Context, uuid.UUID, uuid.UUID, string, string, string, string) (any, error)
	deleteActivity(context.Context, uuid.UUID, uuid.UUID) (int64, error)
	createLeadership(context.Context, uuid.UUID, uuid.UUID, string, string) (any, error)
	listLeadership(context.Context, uuid.UUID) (any, error)
	deleteLeadership(context.Context, uuid.UUID, uuid.UUID) (int64, error)
	createAward(context.Context, uuid.UUID, uuid.UUID, string, string, time.Time, string) (any, error)
	listAwards(context.Context, uuid.UUID) (any, error)
	deleteAward(context.Context, uuid.UUID, uuid.UUID) (int64, error)
	createDiscipline(context.Context, uuid.UUID, uuid.UUID, time.Time, string, string, string, *uuid.UUID) (any, error)
	listDiscipline(context.Context, uuid.UUID) (any, error)
	deleteDiscipline(context.Context, uuid.UUID, uuid.UUID) (int64, error)
}

type StudentPortfolioService struct{ store portfolioStore }

func NewStudentPortfolioService(store portfolioStore) *StudentPortfolioService {
	return &StudentPortfolioService{store: store}
}
func (s *StudentPortfolioService) TeacherProfileIDForUser(ctx context.Context, user uuid.UUID) *uuid.UUID {
	id, err := s.store.teacherProfileID(ctx, user)
	if err != nil {
		return nil
	}
	return &id
}
func (s *StudentPortfolioService) CreateProgressReport(ctx context.Context, student uuid.UUID, req models.CreateProgressReportRequest, writer *uuid.UUID) (any, error) {
	term, err := parsePortfolioUUID(req.TermID, "term")
	if err != nil {
		return nil, err
	}
	return s.store.createProgress(ctx, student, term, req.Narrative, writer)
}
func (s *StudentPortfolioService) ListProgressReports(ctx context.Context, student uuid.UUID) (any, error) {
	return s.store.listProgress(ctx, student)
}
func (s *StudentPortfolioService) UpdateProgressReport(ctx context.Context, id, student uuid.UUID, req models.UpdateProgressReportRequest) (any, error) {
	value, err := s.store.updateProgress(ctx, id, student, req.Narrative)
	return portfolioValue(value, err)
}
func (s *StudentPortfolioService) DeleteProgressReport(ctx context.Context, id, student uuid.UUID) error {
	return portfolioDelete(s.store.deleteProgress(ctx, id, student))
}
func (s *StudentPortfolioService) CreateActivity(ctx context.Context, student uuid.UUID, req models.CreateActivityRequest) (any, error) {
	if !models.ValidActivityCategories[req.Category] {
		return nil, ErrInvalidActivityCategory
	}
	year, err := parsePortfolioUUID(req.AcademicYearID, "academic year")
	if err != nil {
		return nil, err
	}
	return s.store.createActivity(ctx, student, year, req.Category, req.Name, req.Role, req.Achievement)
}
func (s *StudentPortfolioService) ListActivities(ctx context.Context, student uuid.UUID) (any, error) {
	return s.store.listActivities(ctx, student)
}
func (s *StudentPortfolioService) UpdateActivity(ctx context.Context, id, student uuid.UUID, req models.UpdateActivityRequest) (any, error) {
	if !models.ValidActivityCategories[req.Category] {
		return nil, ErrInvalidActivityCategory
	}
	value, err := s.store.updateActivity(ctx, id, student, req.Category, req.Name, req.Role, req.Achievement)
	return portfolioValue(value, err)
}
func (s *StudentPortfolioService) DeleteActivity(ctx context.Context, id, student uuid.UUID) error {
	return portfolioDelete(s.store.deleteActivity(ctx, id, student))
}
func (s *StudentPortfolioService) CreateLeadershipRole(ctx context.Context, student uuid.UUID, req models.CreateLeadershipRoleRequest) (any, error) {
	year, err := parsePortfolioUUID(req.AcademicYearID, "academic year")
	if err != nil {
		return nil, err
	}
	return s.store.createLeadership(ctx, student, year, req.Title, req.Scope)
}
func (s *StudentPortfolioService) ListLeadershipRoles(ctx context.Context, student uuid.UUID) (any, error) {
	return s.store.listLeadership(ctx, student)
}
func (s *StudentPortfolioService) DeleteLeadershipRole(ctx context.Context, id, student uuid.UUID) error {
	return portfolioDelete(s.store.deleteLeadership(ctx, id, student))
}
func (s *StudentPortfolioService) CreateAward(ctx context.Context, student uuid.UUID, req models.CreateAwardRequest) (any, error) {
	year, err := parsePortfolioUUID(req.AcademicYearID, "academic year")
	if err != nil {
		return nil, err
	}
	return s.store.createAward(ctx, student, year, req.Title, req.Category, req.AwardedDate, req.Description)
}
func (s *StudentPortfolioService) ListAwards(ctx context.Context, student uuid.UUID) (any, error) {
	return s.store.listAwards(ctx, student)
}
func (s *StudentPortfolioService) DeleteAward(ctx context.Context, id, student uuid.UUID) error {
	return portfolioDelete(s.store.deleteAward(ctx, id, student))
}
func (s *StudentPortfolioService) CreateDisciplinaryRecord(ctx context.Context, student uuid.UUID, req models.CreateDisciplinaryRecordRequest, recorder *uuid.UUID) (any, error) {
	if !models.ValidDisciplinarySeverities[req.Severity] {
		return nil, ErrInvalidDisciplinarySeverity
	}
	year, err := parsePortfolioUUID(req.AcademicYearID, "academic year")
	if err != nil {
		return nil, err
	}
	return s.store.createDiscipline(ctx, student, year, req.IncidentDate, req.Description, req.ActionTaken, req.Severity, recorder)
}
func (s *StudentPortfolioService) ListDisciplinaryRecords(ctx context.Context, student uuid.UUID) (any, error) {
	return s.store.listDiscipline(ctx, student)
}
func (s *StudentPortfolioService) DeleteDisciplinaryRecord(ctx context.Context, id, student uuid.UUID) error {
	return portfolioDelete(s.store.deleteDiscipline(ctx, id, student))
}

func parsePortfolioUUID(raw, label string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid %s id", label)
	}
	return id, nil
}
func portfolioValue(value any, err error) (any, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrPortfolioRecordNotFound
	}
	return value, err
}
func portfolioDelete(rows int64, err error) error {
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrPortfolioRecordNotFound
	}
	return nil
}
