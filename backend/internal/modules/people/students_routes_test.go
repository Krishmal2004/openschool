package people

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/platform/httpx"
)

type studentListReaderStub struct {
	gotParams StudentListParams
}

func (*studentListReaderStub) Get(context.Context, uuid.UUID) (any, error)          { return nil, nil }
func (*studentListReaderStub) GetWithClass(context.Context, uuid.UUID) (any, error) { return nil, nil }
func (*studentListReaderStub) ListByClass(context.Context, uuid.UUID) (any, error)  { return nil, nil }
func (s *studentListReaderStub) ListPage(_ context.Context, p StudentListParams) (any, error) {
	s.gotParams = p
	return httpx.Page[string]{Items: []string{"one"}, Total: 1, Limit: p.Limit, Offset: p.Offset}, nil
}

func TestListStudentsRouteParsesPageParamsAndFilters(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	group := router.Group("")
	reader := &studentListReaderStub{}
	RegisterStudentRoutes(group, group, nil, reader, nil)

	request := httptest.NewRequest(http.MethodGet, "/students?limit=10&offset=20&search=perera&grade=Grade+5&gender=female", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	if reader.gotParams.Limit != 10 || reader.gotParams.Offset != 20 || reader.gotParams.Search != "perera" ||
		reader.gotParams.Grade != "Grade 5" || reader.gotParams.Gender != "female" {
		t.Fatalf("unexpected params forwarded to ListPage: %+v", reader.gotParams)
	}

	var body httpx.Page[string]
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Total != 1 || body.Limit != 10 || body.Offset != 20 || len(body.Items) != 1 {
		t.Fatalf("unexpected page envelope: %+v", body)
	}
}
