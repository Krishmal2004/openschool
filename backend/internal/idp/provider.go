// Package idp defines the provider-neutral external identity contract and configuration.
package idp

import (
	"context"
	"errors"
	"os"
	"strings"
)

// ErrDuplicateUser is returned by Provider.CreateUser when a unique attribute is already taken by another account.
var ErrDuplicateUser = errors.New("a user with these details already exists")

// User is the provider-neutral account shape returned by every Provider method.
type User struct {
	ID       string
	Username string
	Email    string
}

// Provider is the identity-provider seam every concrete client (e.g. ThunderID) implements.
type Provider interface {
	CreateUser(ctx context.Context, userType string, attrs map[string]any) (*User, error)
	UpdateUser(ctx context.Context, userID string, userType string, attrs map[string]any) error
	DeleteUser(ctx context.Context, userID string) error
	AssignRole(ctx context.Context, roleID string, userID string) error
	// ListUsers returns every account the provider knows of, used by the orphaned-identity reconciliation job to find local `users` rows with no match.
	ListUsers(ctx context.Context) ([]User, error)
}

// JWKSURL returns the identity provider's JWKS endpoint used to validate access tokens.
func JWKSURL() string {
	return os.Getenv("THUNDERID_JWKS_URL")
}

// Issuer returns the expected `iss` claim on access tokens issued by the identity provider.
func Issuer() string {
	return os.Getenv("THUNDERID_ISSUER")
}

// RoleID returns the identity provider's role ID configured for the given base role.
func RoleID(role string) string {
	return os.Getenv("THUNDERID_ROLE_" + strings.ToUpper(role))
}
