//go:build integration

package curriculum

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/openschool-org/openschool/internal/testutil/testdb"
)

func TestCurriculumConfigurationAPIWithPostgres(t *testing.T) {
	gin.SetMode(gin.TestMode)
	pool := testdb.Open(t)
	gradeID, subjectID := seedCurriculumFixture(t, pool)
	router := gin.New()
	group := router.Group("")
	RegisterMediumRoutes(group, group, pool)
	RegisterLevelRoutes(group, group, pool)
	RegisterPresetRoutes(group, pool)

	createdMedium := performCurriculumRequest(t, router, http.MethodPost, "/mediums", mediumRequest{Name: "English"})
	if createdMedium.Code != http.StatusCreated {
		t.Fatalf("create medium: code=%d body=%s", createdMedium.Code, createdMedium.Body.String())
	}
	var medium Medium
	if err := json.Unmarshal(createdMedium.Body.Bytes(), &medium); err != nil {
		t.Fatal(err)
	}
	updatedMedium := performCurriculumRequest(t, router, http.MethodPut, "/mediums/"+medium.ID, mediumRequest{Name: "English Medium"})
	if updatedMedium.Code != http.StatusOK || !bytes.Contains(updatedMedium.Body.Bytes(), []byte(`"name":"English Medium"`)) {
		t.Fatalf("update medium: code=%d body=%s", updatedMedium.Code, updatedMedium.Body.String())
	}

	createdLevel := performCurriculumRequest(t, router, http.MethodPost, "/levels", levelRequest{Label: "Custom Grade 6", GradeID: gradeID.String(), SortOrder: 1})
	if createdLevel.Code != http.StatusCreated {
		t.Fatalf("create level: code=%d body=%s", createdLevel.Code, createdLevel.Body.String())
	}
	var level Level
	if err := json.Unmarshal(createdLevel.Body.Bytes(), &level); err != nil {
		t.Fatal(err)
	}
	invalidGroup := performCurriculumRequest(t, router, http.MethodPost, "/levels/"+level.ID+"/groups", groupRequest{Label: "Invalid", MinSelect: 2, MaxSelect: 1})
	if invalidGroup.Code != http.StatusBadRequest {
		t.Fatalf("invalid selection range: code=%d body=%s", invalidGroup.Code, invalidGroup.Body.String())
	}
	createdGroup := performCurriculumRequest(t, router, http.MethodPost, "/levels/"+level.ID+"/groups", groupRequest{Label: "Core", MinSelect: 1, MaxSelect: 1, SortOrder: 1})
	if createdGroup.Code != http.StatusCreated {
		t.Fatalf("create selection group: code=%d body=%s", createdGroup.Code, createdGroup.Body.String())
	}
	var selectionGroup SelectionGroup
	if err := json.Unmarshal(createdGroup.Body.Bytes(), &selectionGroup); err != nil {
		t.Fatal(err)
	}
	linked := performCurriculumRequest(t, router, http.MethodPost, "/groups/"+selectionGroup.ID+"/subjects", groupSubjectRequest{SubjectID: subjectID.String(), MediumID: medium.ID, PrerequisiteNote: "Foundation", SortOrder: 1})
	if linked.Code != http.StatusCreated {
		t.Fatalf("link group subject: code=%d body=%s", linked.Code, linked.Body.String())
	}
	tree := performCurriculumRequest(t, router, http.MethodGet, "/levels/"+level.ID+"/tree", nil)
	if tree.Code != http.StatusOK || !bytes.Contains(tree.Body.Bytes(), []byte(`"Label":"Core"`)) || !bytes.Contains(tree.Body.Bytes(), []byte(`"SubjectCode":"CUR-MATH"`)) || !bytes.Contains(tree.Body.Bytes(), []byte(`"MediumName":"English Medium"`)) {
		t.Fatalf("curriculum tree: code=%d body=%s", tree.Code, tree.Body.String())
	}

	duplicated := performCurriculumRequest(t, router, http.MethodPost, "/levels/"+level.ID+"/duplicate", levelRequest{Label: "Custom Grade 6 Copy", GradeID: gradeID.String(), SortOrder: 2})
	if duplicated.Code != http.StatusCreated {
		t.Fatalf("duplicate curriculum level: code=%d body=%s", duplicated.Code, duplicated.Body.String())
	}
	var copyLevel Level
	if err := json.Unmarshal(duplicated.Body.Bytes(), &copyLevel); err != nil {
		t.Fatal(err)
	}
	copyTree := performCurriculumRequest(t, router, http.MethodGet, "/levels/"+copyLevel.ID+"/tree", nil)
	if copyTree.Code != http.StatusOK || !bytes.Contains(copyTree.Body.Bytes(), []byte(`"SubjectCode":"CUR-MATH"`)) {
		t.Fatalf("duplicated curriculum tree: code=%d body=%s", copyTree.Code, copyTree.Body.String())
	}
	mediumInUse := performCurriculumRequest(t, router, http.MethodDelete, "/mediums/"+medium.ID, nil)
	if mediumInUse.Code != http.StatusConflict {
		t.Fatalf("delete medium in use: code=%d body=%s", mediumInUse.Code, mediumInUse.Body.String())
	}

	before := curriculumCounts(t, pool)
	preview := performCurriculumRequest(t, router, http.MethodGet, "/curriculum/preset/preview", nil)
	if preview.Code != http.StatusOK || !bytes.Contains(preview.Body.Bytes(), []byte(`"dry_run":true`)) {
		t.Fatalf("preview curriculum preset: code=%d body=%s", preview.Code, preview.Body.String())
	}
	if afterPreview := curriculumCounts(t, pool); afterPreview != before {
		t.Fatalf("preset preview mutated database: before=%v after=%v", before, afterPreview)
	}
	applied := performCurriculumRequest(t, router, http.MethodPost, "/curriculum/preset", nil)
	if applied.Code != http.StatusOK || !bytes.Contains(applied.Body.Bytes(), []byte(`"dry_run":false`)) {
		t.Fatalf("apply curriculum preset: code=%d body=%s", applied.Code, applied.Body.String())
	}
	afterApply := curriculumCounts(t, pool)
	if afterApply == before {
		t.Fatalf("preset did not create curriculum rows: %v", afterApply)
	}
	reapplied := performCurriculumRequest(t, router, http.MethodPost, "/curriculum/preset", nil)
	if reapplied.Code != http.StatusOK {
		t.Fatalf("reapply curriculum preset: code=%d body=%s", reapplied.Code, reapplied.Body.String())
	}
	if afterReapply := curriculumCounts(t, pool); afterReapply != afterApply {
		t.Fatalf("preset was not idempotent: first=%v second=%v", afterApply, afterReapply)
	}
}

type curriculumRowCounts struct{ subjects, levels, groups, links int }

func curriculumCounts(t *testing.T, pool *pgxpool.Pool) curriculumRowCounts {
	t.Helper()
	result := curriculumRowCounts{}
	for table, target := range map[string]*int{
		"subjects": &result.subjects, "levels": &result.levels, "selection_groups": &result.groups, "group_subjects": &result.links,
	} {
		if err := pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM "+table).Scan(target); err != nil {
			t.Fatal(err)
		}
	}
	return result
}

func seedCurriculumFixture(t *testing.T, pool *pgxpool.Pool) (uuid.UUID, uuid.UUID) {
	t.Helper()
	ctx := context.Background()
	var gradeID, subjectID uuid.UUID
	if err := pool.QueryRow(ctx, "INSERT INTO grades (name, sort_order) VALUES ('Grade 6 Curriculum', 6) RETURNING id").Scan(&gradeID); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, "INSERT INTO subjects (name, code) VALUES ('Curriculum Mathematics', 'CUR-MATH') RETURNING id").Scan(&subjectID); err != nil {
		t.Fatal(err)
	}
	return gradeID, subjectID
}

func performCurriculumRequest(t *testing.T, router http.Handler, method, path string, body any) *httptest.ResponseRecorder {
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
