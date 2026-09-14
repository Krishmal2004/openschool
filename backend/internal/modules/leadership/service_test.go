package leadership

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/openschool-org/openschool/internal/models"
)

type leadershipStore struct {
	store
	principal       bool
	position        Position
	positionErr     error
	headedGrades    []uuid.UUID
	formTeacher     bool
	subjectTeacher  bool
	scopeGrades     []scopeGrade
	replacedScopes  []uuid.UUID
	streamSection   SectionHead
	streamArguments []uuid.UUID
}

func (s *leadershipStore) upsertVicePrincipal(_ context.Context, teacherID uuid.UUID, wholeSchool bool) (Position, error) {
	s.position = Position{ID: uuid.New(), TeacherID: teacherID, Position: "vice_principal", NotifyWholeSchool: wholeSchool}
	return s.position, nil
}

func (s *leadershipStore) replaceVicePrincipalScopes(_ context.Context, _ uuid.UUID, gradeIDs []uuid.UUID) error {
	s.replacedScopes = gradeIDs
	return nil
}

func (s *leadershipStore) isPrincipal(context.Context, uuid.UUID) (bool, error) {
	return s.principal, nil
}

func (s *leadershipStore) getPositionForTeacher(context.Context, uuid.UUID) (Position, error) {
	if s.positionErr != nil {
		return Position{}, s.positionErr
	}
	if s.position.Position == "" {
		return Position{}, pgx.ErrNoRows
	}
	return s.position, nil
}

func (s *leadershipStore) listGradeIDsHeadedByTeacher(context.Context, uuid.UUID, uuid.UUID) ([]uuid.UUID, error) {
	return s.headedGrades, nil
}

func (s *leadershipStore) isFormTeacher(context.Context, uuid.UUID, uuid.UUID) (bool, error) {
	return s.formTeacher, nil
}

func (s *leadershipStore) isSubjectTeacher(context.Context, uuid.UUID, uuid.UUID) (bool, error) {
	return s.subjectTeacher, nil
}

func (s *leadershipStore) listVicePrincipalScopeGrades(context.Context, uuid.UUID) ([]scopeGrade, error) {
	return s.scopeGrades, nil
}

func (s *leadershipStore) upsertStreamSectionHead(_ context.Context, yearID, gradeID, streamID, teacherID uuid.UUID) (SectionHead, error) {
	s.streamArguments = []uuid.UUID{yearID, gradeID, streamID, teacherID}
	return s.streamSection, nil
}

type auditRecord struct {
	entityType string
	action     string
	actorID    uuid.UUID
}

func (a *auditRecord) Record(_ context.Context, entityType string, _ uuid.UUID, action string, actorID uuid.UUID, _, _ interface{}, _ string) error {
	a.entityType, a.action, a.actorID = entityType, action, actorID
	return nil
}

func TestAssignVicePrincipalReplacesGradeScopeAndAudits(t *testing.T) {
	store := &leadershipStore{}
	audit := &auditRecord{}
	teacherID, actorID := uuid.New(), uuid.New()
	gradeOne, gradeTwo := uuid.New(), uuid.New()

	position, err := NewService(store, audit).AssignVicePrincipal(context.Background(), models.AssignVicePrincipalRequest{
		TeacherID: teacherID.String(), GradeIDs: []string{gradeOne.String(), gradeTwo.String()},
	}, actorID)
	if err != nil {
		t.Fatal(err)
	}
	if position.TeacherID != teacherID || len(store.replacedScopes) != 2 || store.replacedScopes[0] != gradeOne || store.replacedScopes[1] != gradeTwo {
		t.Fatalf("unexpected scoped appointment: position=%#v scopes=%v", position, store.replacedScopes)
	}
	if audit.entityType != "teacher_position" || audit.action != "assigned_vice_principal" || audit.actorID != actorID {
		t.Fatalf("unexpected audit record: %#v", audit)
	}
}

func TestAssignVicePrincipalWholeSchoolIgnoresGradeInput(t *testing.T) {
	store := &leadershipStore{}
	_, err := NewService(store, nil).AssignVicePrincipal(context.Background(), models.AssignVicePrincipalRequest{
		TeacherID: uuid.NewString(), NotifyWholeSchool: true, GradeIDs: []string{"not-a-uuid"},
	}, uuid.New())
	if err != nil {
		t.Fatal(err)
	}
	if len(store.replacedScopes) != 0 {
		t.Fatalf("whole-school appointment retained grade scope: %v", store.replacedScopes)
	}
}

func TestRankForTeacherUsesHighestMatchingRank(t *testing.T) {
	yearID, teacherID := uuid.New(), uuid.New()
	tests := []struct {
		name  string
		store *leadershipStore
		want  PositionRank
	}{
		{"principal", &leadershipStore{principal: true, formTeacher: true}, RankPrincipal},
		{"vice principal", &leadershipStore{position: Position{Position: "vice_principal"}, headedGrades: []uuid.UUID{uuid.New()}}, RankVicePrincipal},
		{"section head", &leadershipStore{headedGrades: []uuid.UUID{uuid.New()}, formTeacher: true}, RankSectionHead},
		{"class teacher", &leadershipStore{formTeacher: true, subjectTeacher: true}, RankClassTeacher},
		{"subject teacher", &leadershipStore{subjectTeacher: true}, RankSubjectTeacher},
		{"teacher", &leadershipStore{}, RankTeacher},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rank, _, err := NewService(test.store, nil).RankForTeacher(context.Background(), teacherID, yearID)
			if err != nil {
				t.Fatal(err)
			}
			if rank != test.want {
				t.Fatalf("rank=%d, want %d", rank, test.want)
			}
		})
	}
}

func TestLeadershipScopeMapsVicePrincipalGrades(t *testing.T) {
	positionID, first, second := uuid.New(), uuid.New(), uuid.New()
	store := &leadershipStore{
		position:    Position{ID: positionID, Position: "vice_principal"},
		scopeGrades: []scopeGrade{{ID: first}, {ID: second}},
	}
	wholeSchool, grades, err := NewService(store, nil).LeadershipScope(context.Background(), uuid.New(), uuid.New())
	if err != nil {
		t.Fatal(err)
	}
	if wholeSchool || len(grades) != 2 || grades[0] != first || grades[1] != second {
		t.Fatalf("unexpected scope: wholeSchool=%v grades=%v", wholeSchool, grades)
	}
}

func TestLeadershipScopeRejectsOrdinaryTeacher(t *testing.T) {
	_, _, err := NewService(&leadershipStore{}, nil).LeadershipScope(context.Background(), uuid.New(), uuid.New())
	if !errors.Is(err, ErrInsufficientRank) {
		t.Fatalf("error=%v, want ErrInsufficientRank", err)
	}
}

func TestAssignStreamSectionHeadParsesIdentifiers(t *testing.T) {
	store := &leadershipStore{streamSection: SectionHead{ID: uuid.New()}}
	yearID, gradeID, streamID, teacherID := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	rawStream := streamID.String()
	result, err := NewService(store, nil).AssignSectionHead(context.Background(), models.AssignSectionHeadRequest{
		AcademicYearID: yearID.String(), GradeID: gradeID.String(), StreamID: &rawStream, TeacherID: teacherID.String(),
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != store.streamSection.ID || len(store.streamArguments) != 4 || store.streamArguments[2] != streamID {
		t.Fatalf("unexpected section-head assignment: result=%#v args=%v", result, store.streamArguments)
	}
}
