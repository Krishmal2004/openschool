package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/authz"
)

type studentAccessStub struct {
	studentID   uuid.UUID
	guardian    bool
	studentErr  error
	guardianErr error
}

func (s studentAccessStub) StudentIDForUser(context.Context, uuid.UUID) (uuid.UUID, error) {
	return s.studentID, s.studentErr
}

func (s studentAccessStub) IsGuardianOfStudent(context.Context, uuid.UUID, uuid.UUID) (bool, error) {
	return s.guardian, s.guardianErr
}

func TestRequireStudentAccessUsesInjectedRelationshipReader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	studentID := uuid.New()
	tests := []struct {
		name   string
		roles  []string
		access studentAccessStub
		want   int
	}{
		{name: "admin bypass", roles: []string{authz.RoleAdmin}, want: http.StatusOK},
		{name: "student owns profile", roles: []string{authz.RoleStudent}, access: studentAccessStub{studentID: studentID}, want: http.StatusOK},
		{name: "student does not own profile", roles: []string{authz.RoleStudent}, access: studentAccessStub{studentID: uuid.New()}, want: http.StatusForbidden},
		{name: "student lookup fails closed", roles: []string{authz.RoleStudent}, access: studentAccessStub{studentErr: errors.New("lookup failed")}, want: http.StatusForbidden},
		{name: "guardian owns student", roles: []string{authz.RoleParent}, access: studentAccessStub{guardian: true}, want: http.StatusOK},
		{name: "guardian does not own student", roles: []string{authz.RoleParent}, want: http.StatusForbidden},
		{name: "guardian lookup fails closed", roles: []string{authz.RoleParent}, access: studentAccessStub{guardianErr: errors.New("lookup failed")}, want: http.StatusForbidden},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router := gin.New()
			router.Use(func(c *gin.Context) {
				c.Set("userID", uuid.NewString())
				c.Set("roles", test.roles)
				c.Next()
			})
			router.GET("/students/:id", RequireStudentAccess(test.access), func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			response := httptest.NewRecorder()
			router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/students/"+studentID.String(), nil))
			if response.Code != test.want {
				t.Fatalf("status = %d, want %d; body=%s", response.Code, test.want, response.Body.String())
			}
		})
	}
}
