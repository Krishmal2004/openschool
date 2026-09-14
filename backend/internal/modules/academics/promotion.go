package academics

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/openschool-org/openschool/internal/models"
)

var ErrCommitTargetClassMismatch = errors.New("one or more target classes do not belong to the target academic year")

type promotionStudent struct {
	ID, ClassID, GradeID              uuid.UUID
	Name, Index, ClassName, GradeName string
	MediumID, MediumName              *string
}
type promotionGrade struct {
	ID   uuid.UUID
	Name string
}
type promotionClass struct {
	ID   uuid.UUID
	Name string
}
type promotionMarks struct {
	ID         uuid.UUID
	Total, Max float64
}
type promotionStore interface {
	students(context.Context, uuid.UUID) ([]promotionStudent, error)
	nextGrade(context.Context, uuid.UUID) (promotionGrade, error)
	classByName(context.Context, uuid.UUID, uuid.UUID, string) (promotionClass, error)
	classByMedium(context.Context, uuid.UUID, uuid.UUID, string) (promotionClass, error)
	marks(context.Context, uuid.UUID, []uuid.UUID) ([]promotionMarks, error)
	validClasses(context.Context, uuid.UUID, []uuid.UUID) (int64, error)
	commit(context.Context, uuid.UUID, []uuid.UUID, []uuid.UUID) error
}
type promotionService struct{ store promotionStore }

func NewPromotionService(store promotionStore) PromotionRunner {
	return &promotionService{store: store}
}

func (s *promotionService) Preview(ctx context.Context, source, target uuid.UUID, term *uuid.UUID) ([]models.PromotionPreviewRow, error) {
	students, err := s.store.students(ctx, source)
	if err != nil {
		return nil, fmt.Errorf("failed to list students for source year: %w", err)
	}
	grades := map[uuid.UUID]*promotionGrade{}
	classes := map[string]*promotionClass{}
	result := make([]models.PromotionPreviewRow, 0, len(students))
	ids := make([]uuid.UUID, 0, len(students))
	for _, student := range students {
		row := models.PromotionPreviewRow{StudentID: student.ID.String(), StudentName: student.Name, StudentIndex: student.Index, CurrentClassID: student.ClassID.String(), CurrentClassName: student.ClassName, CurrentGradeID: student.GradeID.String(), CurrentGradeName: student.GradeName, MediumLocked: student.MediumID != nil, CurrentMediumID: student.MediumID, CurrentMediumName: student.MediumName}
		grade, ok := grades[student.GradeID]
		if !ok {
			value, e := s.store.nextGrade(ctx, student.GradeID)
			if errors.Is(e, pgx.ErrNoRows) {
				grades[student.GradeID] = nil
			} else if e != nil {
				return nil, fmt.Errorf("failed to resolve next grade: %w", e)
			} else {
				grades[student.GradeID] = &value
			}
			grade = grades[student.GradeID]
		}
		if grade == nil {
			row.Graduating = true
			result = append(result, row)
			ids = append(ids, student.ID)
			continue
		}
		nextID, nextName := grade.ID.String(), grade.Name
		row.NextGradeID, row.NextGradeName = &nextID, &nextName
		key := nextID + "|name|" + student.ClassName
		if student.MediumID != nil {
			key += "|medium|" + *student.MediumID
		}
		class, ok := classes[key]
		if !ok {
			var value promotionClass
			if student.MediumID != nil {
				value, err = s.store.classByMedium(ctx, grade.ID, target, *student.MediumID)
			} else {
				value, err = s.store.classByName(ctx, grade.ID, target, student.ClassName)
			}
			if errors.Is(err, pgx.ErrNoRows) {
				classes[key] = nil
			} else if err != nil {
				return nil, fmt.Errorf("failed to resolve suggested class: %w", err)
			} else {
				classes[key] = &value
			}
			class = classes[key]
		}
		if class != nil {
			id, name := class.ID.String(), class.Name
			row.SuggestedClassID, row.SuggestedClassName = &id, &name
		}
		result = append(result, row)
		ids = append(ids, student.ID)
	}
	if term != nil {
		marks, e := s.store.marks(ctx, *term, ids)
		if e != nil {
			return nil, fmt.Errorf("failed to load marks: %w", e)
		}
		by := map[uuid.UUID]promotionMarks{}
		for _, mark := range marks {
			by[mark.ID] = mark
		}
		for i := range result {
			if mark, ok := by[uuid.MustParse(result[i].StudentID)]; ok {
				result[i].TotalMarks, result[i].TotalMaxMarks = &mark.Total, &mark.Max
			}
		}
	}
	return result, nil
}

func (s *promotionService) CommitAssignments(ctx context.Context, request models.CommitAssignmentsRequest) (int, error) {
	year, err := uuid.Parse(request.AcademicYearID)
	if err != nil {
		return 0, errors.New("invalid academic year id")
	}
	students := make([]uuid.UUID, len(request.Assignments))
	classes := make([]uuid.UUID, len(request.Assignments))
	unique := map[uuid.UUID]bool{}
	for i, assignment := range request.Assignments {
		students[i], err = uuid.Parse(assignment.StudentID)
		if err != nil {
			return 0, fmt.Errorf("invalid student id: %s", assignment.StudentID)
		}
		classes[i], err = uuid.Parse(assignment.ClassID)
		if err != nil {
			return 0, fmt.Errorf("invalid class id: %s", assignment.ClassID)
		}
		unique[classes[i]] = true
	}
	ids := make([]uuid.UUID, 0, len(unique))
	for id := range unique {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i].String() < ids[j].String() })
	count, err := s.store.validClasses(ctx, year, ids)
	if err != nil {
		return 0, err
	}
	if int(count) != len(ids) {
		return 0, ErrCommitTargetClassMismatch
	}
	if err := s.store.commit(ctx, year, students, classes); err != nil {
		return 0, err
	}
	return len(students), nil
}
