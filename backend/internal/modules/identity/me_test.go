package identity

import (
	"context"
	"testing"
)

type provisionerStub struct {
	called  bool
	command ensureUserCommand
	result  provisionedUser
	err     error
}

func (s *provisionerStub) ensureExists(_ context.Context, command ensureUserCommand) (provisionedUser, error) {
	s.called = true
	s.command = command
	return s.result, s.err
}

func TestEnsureProvisionedSkipsUnknownRoles(t *testing.T) {
	repository := &provisionerStub{}
	service := newMeService(repository)

	result, err := service.ensureProvisioned(context.Background(), ensureUserCommand{})
	if err != nil || result.MustChangePassword || repository.called {
		t.Fatalf("unknown role should not provision a user")
	}
}

func TestEnsureProvisionedDelegatesKnownRole(t *testing.T) {
	repository := &provisionerStub{result: provisionedUser{MustChangePassword: true}}
	service := newMeService(repository)
	command := ensureUserCommand{Role: "teacher", Email: "teacher@example.test"}

	result, err := service.ensureProvisioned(context.Background(), command)
	if err != nil {
		t.Fatal(err)
	}
	if !repository.called || repository.command != command || !result.MustChangePassword {
		t.Fatalf("known role was not forwarded to the repository")
	}
}
