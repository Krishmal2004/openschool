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
	ErrMediumNotFound         = errors.New("medium not found")
	ErrMediumInUse            = errors.New("medium is in use and cannot be deleted")
	ErrLevelNotFound          = errors.New("level not found")
	ErrLevelInUse             = errors.New("level has enrolled students and cannot be deleted")
	ErrSelectionGroupNotFound = errors.New("selection group not found")
	ErrSelectionGroupInUse    = errors.New("selection group has enrolled students and cannot be deleted")
	ErrInvalidSelectRange     = errors.New("max_select must be greater than or equal to min_select")
)

type CurriculumService struct {
	repo *repositories.CurriculumRepository
}

func NewCurriculumService(repo *repositories.CurriculumRepository) *CurriculumService {
	return &CurriculumService{repo: repo}
}

func (s *CurriculumService) CreateMedium(ctx context.Context, req models.CreateMediumRequest) (models.MediumResponse, error) {
	medium, err := s.repo.CreateMedium(ctx, req.Name)
	return toMediumResponse(medium), err
}

func (s *CurriculumService) ListMediums(ctx context.Context) ([]models.MediumResponse, error) {
	mediums, err := s.repo.ListMediums(ctx)
	if err != nil {
		return nil, err
	}
	resp := make([]models.MediumResponse, len(mediums))
	for i, medium := range mediums {
		resp[i] = toMediumResponse(medium)
	}
	return resp, nil
}

func (s *CurriculumService) UpdateMedium(ctx context.Context, id uuid.UUID, req models.UpdateMediumRequest) (models.MediumResponse, error) {
	medium, err := s.repo.UpdateMedium(ctx, db.UpdateMediumParams{ID: id, Name: req.Name})
	return toMediumResponse(medium), err
}

func (s *CurriculumService) DeleteMedium(ctx context.Context, id uuid.UUID) error {
	rows, err := s.repo.DeleteMedium(ctx, id)
	if err != nil {
		return err
	}
	if rows == 0 {
		if _, err := s.repo.GetMediumByID(ctx, id); err != nil {
			return ErrMediumNotFound
		}
		return ErrMediumInUse
	}
	return nil
}

func (s *CurriculumService) CreateLevel(ctx context.Context, req models.CreateLevelRequest) (models.LevelResponse, error) {
	gradeID, err := parseOptionalUUID(req.GradeID)
	if err != nil {
		return models.LevelResponse{}, errors.New("invalid grade_id")
	}

	level, err := s.repo.CreateLevel(ctx, db.CreateLevelParams{
		Label:     req.Label,
		GradeID:   gradeID,
		SortOrder: req.SortOrder,
	})
	return ToLevelResponse(level), err
}

func (s *CurriculumService) GetLevel(ctx context.Context, id uuid.UUID) (models.LevelResponse, error) {
	level, err := s.repo.GetLevelByID(ctx, id)
	return ToLevelResponse(level), err
}

func (s *CurriculumService) ListLevels(ctx context.Context) ([]models.LevelResponse, error) {
	levels, err := s.repo.ListLevels(ctx)
	return toLevelResponses(levels), err
}

func (s *CurriculumService) ListLevelsByGrade(ctx context.Context, gradeID uuid.UUID) ([]models.LevelResponse, error) {
	levels, err := s.repo.ListLevelsByGrade(ctx, pgUUID(gradeID))
	return toLevelResponses(levels), err
}

func (s *CurriculumService) UpdateLevel(ctx context.Context, id uuid.UUID, req models.UpdateLevelRequest) (models.LevelResponse, error) {
	gradeID, err := parseOptionalUUID(req.GradeID)
	if err != nil {
		return models.LevelResponse{}, errors.New("invalid grade_id")
	}

	level, err := s.repo.UpdateLevel(ctx, db.UpdateLevelParams{
		ID:        id,
		Label:     req.Label,
		GradeID:   gradeID,
		SortOrder: req.SortOrder,
	})
	return ToLevelResponse(level), err
}

func (s *CurriculumService) DuplicateLevel(ctx context.Context, sourceID uuid.UUID, req models.DuplicateLevelRequest) (models.LevelResponse, error) {
	if _, err := s.repo.GetLevelByID(ctx, sourceID); err != nil {
		return models.LevelResponse{}, ErrLevelNotFound
	}

	gradeID, err := parseOptionalUUID(req.GradeID)
	if err != nil {
		return models.LevelResponse{}, errors.New("invalid grade_id")
	}

	level, err := s.repo.DuplicateLevel(ctx, sourceID, db.CreateLevelParams{
		Label:     req.Label,
		GradeID:   gradeID,
		SortOrder: req.SortOrder,
	})
	return ToLevelResponse(level), err
}

func (s *CurriculumService) DeleteLevel(ctx context.Context, id uuid.UUID) error {
	rows, err := s.repo.DeleteLevel(ctx, id)
	if err != nil {
		return err
	}
	if rows == 0 {
		if _, err := s.repo.GetLevelByID(ctx, id); err != nil {
			return ErrLevelNotFound
		}
		return ErrLevelInUse
	}
	return nil
}

func (s *CurriculumService) CreateSelectionGroup(ctx context.Context, levelID uuid.UUID, req models.CreateSelectionGroupRequest) (models.SelectionGroupResponse, error) {
	if req.MaxSelect < req.MinSelect {
		return models.SelectionGroupResponse{}, ErrInvalidSelectRange
	}
	if _, err := s.repo.GetLevelByID(ctx, levelID); err != nil {
		return models.SelectionGroupResponse{}, ErrLevelNotFound
	}

	group, err := s.repo.CreateSelectionGroup(ctx, db.CreateSelectionGroupParams{
		LevelID:   levelID,
		Label:     req.Label,
		MinSelect: req.MinSelect,
		MaxSelect: req.MaxSelect,
		SortOrder: req.SortOrder,
	})
	return toSelectionGroupResponse(group), err
}

func (s *CurriculumService) ListSelectionGroupsByLevel(ctx context.Context, levelID uuid.UUID) ([]models.SelectionGroupResponse, error) {
	groups, err := s.repo.ListSelectionGroupsByLevel(ctx, levelID)
	if err != nil {
		return nil, err
	}
	resp := make([]models.SelectionGroupResponse, len(groups))
	for i, group := range groups {
		resp[i] = toSelectionGroupResponse(group)
	}
	return resp, nil
}

func (s *CurriculumService) UpdateSelectionGroup(ctx context.Context, id uuid.UUID, req models.UpdateSelectionGroupRequest) (models.SelectionGroupResponse, error) {
	if req.MaxSelect < req.MinSelect {
		return models.SelectionGroupResponse{}, ErrInvalidSelectRange
	}

	group, err := s.repo.UpdateSelectionGroup(ctx, db.UpdateSelectionGroupParams{
		ID:        id,
		Label:     req.Label,
		MinSelect: req.MinSelect,
		MaxSelect: req.MaxSelect,
		SortOrder: req.SortOrder,
	})
	return toSelectionGroupResponse(group), err
}

func (s *CurriculumService) DeleteSelectionGroup(ctx context.Context, id uuid.UUID) error {
	rows, err := s.repo.DeleteSelectionGroup(ctx, id)
	if err != nil {
		return err
	}
	if rows == 0 {
		if _, err := s.repo.GetSelectionGroupByID(ctx, id); err != nil {
			return ErrSelectionGroupNotFound
		}
		return ErrSelectionGroupInUse
	}
	return nil
}

func (s *CurriculumService) AddGroupSubject(ctx context.Context, groupID uuid.UUID, req models.AddGroupSubjectRequest) error {
	subjectID, err := uuid.Parse(req.SubjectID)
	if err != nil {
		return errors.New("invalid subject_id")
	}

	mediumID, err := parseOptionalUUID(req.MediumID)
	if err != nil {
		return errors.New("invalid medium_id")
	}

	if _, err := s.repo.GetSelectionGroupByID(ctx, groupID); err != nil {
		return ErrSelectionGroupNotFound
	}

	_, err = s.repo.AddGroupSubject(ctx, db.AddGroupSubjectParams{
		GroupID:          groupID,
		SubjectID:        subjectID,
		MediumID:         mediumID,
		PrerequisiteNote: optionalText(req.PrerequisiteNote),
		SortOrder:        req.SortOrder,
	})
	return err
}

func (s *CurriculumService) ListGroupSubjects(ctx context.Context, groupID uuid.UUID) ([]models.GroupSubjectResponse, error) {
	rows, err := s.repo.ListGroupSubjects(ctx, groupID)
	if err != nil {
		return nil, err
	}

	resp := make([]models.GroupSubjectResponse, len(rows))
	for i, r := range rows {
		resp[i] = models.GroupSubjectResponse{
			SubjectID:        r.SubjectID.String(),
			SubjectName:      r.SubjectName,
			SubjectCode:      r.SubjectCode,
			SubjectType:      textString(r.SubjectType),
			MediumID:         uuidString(r.MediumID),
			MediumName:       textString(r.MediumName),
			PrerequisiteNote: textString(r.PrerequisiteNote),
			SortOrder:        r.SortOrder,
		}
	}
	return resp, nil
}

func (s *CurriculumService) RemoveGroupSubject(ctx context.Context, groupID, subjectID uuid.UUID) error {
	return s.repo.RemoveGroupSubject(ctx, db.RemoveGroupSubjectParams{
		GroupID:   groupID,
		SubjectID: subjectID,
	})
}

func (s *CurriculumService) GetCurriculumTree(ctx context.Context, levelID uuid.UUID) (models.CurriculumTreeResponse, error) {
	level, err := s.repo.GetLevelByID(ctx, levelID)
	if err != nil {
		return models.CurriculumTreeResponse{}, ErrLevelNotFound
	}

	rows, err := s.repo.GetCurriculumTreeByLevel(ctx, levelID)
	if err != nil {
		return models.CurriculumTreeResponse{}, err
	}

	groups := make([]models.CurriculumGroupResponse, 0)
	position := make(map[uuid.UUID]int)

	for _, row := range rows {
		i, seen := position[row.GroupID]
		if !seen {
			groups = append(groups, models.CurriculumGroupResponse{
				ID:        row.GroupID.String(),
				Label:     row.GroupLabel,
				MinSelect: row.MinSelect,
				MaxSelect: row.MaxSelect,
				SortOrder: row.GroupSortOrder,
				Subjects:  []models.GroupSubjectResponse{},
			})
			i = len(groups) - 1
			position[row.GroupID] = i
		}

		// a group with no subjects still returns one row, with a NULL subject
		if !row.SubjectID.Valid {
			continue
		}

		groups[i].Subjects = append(groups[i].Subjects, models.GroupSubjectResponse{
			SubjectID:        uuid.UUID(row.SubjectID.Bytes).String(),
			SubjectName:      row.SubjectName.String,
			SubjectCode:      row.SubjectCode.String,
			SubjectType:      textString(row.SubjectType),
			MediumID:         uuidString(row.MediumID),
			MediumName:       textString(row.MediumName),
			PrerequisiteNote: textString(row.PrerequisiteNote),
			SortOrder:        row.SubjectSortOrder.Int32,
		})
	}

	return models.CurriculumTreeResponse{
		Level:  ToLevelResponse(level),
		Groups: groups,
	}, nil
}

func toMediumResponse(m db.Medium) models.MediumResponse {
	return models.MediumResponse{ID: m.ID.String(), Name: m.Name, CreatedAt: m.CreatedAt.Time.String()}
}

func toSelectionGroupResponse(g db.SelectionGroup) models.SelectionGroupResponse {
	return models.SelectionGroupResponse{ID: g.ID.String(), LevelID: g.LevelID.String(), Label: g.Label, MinSelect: g.MinSelect, MaxSelect: g.MaxSelect, SortOrder: g.SortOrder, CreatedAt: g.CreatedAt.Time.String()}
}

func ToLevelResponse(l db.Level) models.LevelResponse {
	return models.LevelResponse{
		ID:        l.ID.String(),
		Label:     l.Label,
		GradeID:   uuidString(l.GradeID),
		SortOrder: l.SortOrder,
		CreatedAt: l.CreatedAt.Time.String(),
	}
}

func toLevelResponses(levels []db.Level) []models.LevelResponse {
	resp := make([]models.LevelResponse, len(levels))
	for i, level := range levels {
		resp[i] = ToLevelResponse(level)
	}
	return resp
}
