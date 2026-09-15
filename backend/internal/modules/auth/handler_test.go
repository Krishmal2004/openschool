package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHandlerSecurityResponses(t *testing.T) {
	gin.SetMode(gin.TestMode)
	store := &authStoreStub{consumeErr: errTestNoRows}
	handler := NewAuthHandler(newTestService(store, guardianStub{}, &passwordUpdaterStub{}, &mailerStub{}))

	tests := []struct {
		name   string
		path   string
		body   any
		handle gin.HandlerFunc
		status int
	}{
		{name: "invalid reset token", path: "/auth/reset-password", body: ResetPasswordRequest{Token: "invalid", NewPassword: "long-enough"}, handle: handler.ResetPassword, status: http.StatusUnauthorized},
		{name: "short password", path: "/auth/reset-password", body: ResetPasswordRequest{Token: "token", NewPassword: "short"}, handle: handler.ResetPassword, status: http.StatusBadRequest},
		{name: "admin forgot role rejected", path: "/auth/forgot-password", body: ForgotPasswordRequest{Role: "admin", Identifier: "admin@example.com", Secret: "secret"}, handle: handler.ForgotPassword, status: http.StatusBadRequest},
		{name: "missing authenticated identity", path: "/auth/change-password", body: ChangePasswordRequest{NewPassword: "long-enough"}, handle: handler.ChangePassword, status: http.StatusUnauthorized},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			payload, err := json.Marshal(test.body)
			if err != nil {
				t.Fatal(err)
			}
			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			ctx.Request = httptest.NewRequest(http.MethodPost, test.path, bytes.NewReader(payload))
			ctx.Request.Header.Set("Content-Type", "application/json")

			test.handle(ctx)
			if recorder.Code != test.status {
				t.Fatalf("status = %d, want %d; body=%s", recorder.Code, test.status, recorder.Body.String())
			}
		})
	}
}

var errTestNoRows = &testError{"no rows"}

type testError struct{ message string }

func (e *testError) Error() string { return e.message }
