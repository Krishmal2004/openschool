// Package auth owns OpenSchool's password lifecycle and ThunderID password updates.
package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/authz"
	"github.com/openschool-org/openschool/internal/mailer"
	"github.com/openschool-org/openschool/internal/ports"
)

var (
	// ErrInvalidCredentials is returned for both "no such account" and "secret
	// didn't match" — never distinguished in the response, so a ForgotPassword
	// call can't be used to enumerate which identifiers exist.
	ErrInvalidCredentials = errors.New("no account matches those details")
	ErrResetTokenInvalid  = errors.New("reset link is invalid, already used, or has expired")
)

// passwordResetTokenTTL bounds how long an emailed reset link stays valid.
const passwordResetTokenTTL = 15 * time.Minute

type authStore interface {
	userByEmail(context.Context, string) (userAccount, error)
	userByID(context.Context, uuid.UUID) (userAccount, error)
	teacherCredentialsMatch(context.Context, uuid.UUID, string) bool
	studentCredentialsMatch(context.Context, uuid.UUID, string) bool
	createResetToken(context.Context, uuid.UUID, string, time.Time) error
	consumeResetToken(context.Context, string) (resetToken, error)
	setMustChangePassword(context.Context, uuid.UUID, bool) error
}

// PasswordUpdater is the narrow identity-provider operation needed by Auth.
type PasswordUpdater interface {
	UpdateUser(context.Context, string, string, map[string]any) error
}

// Service implements the password lifecycle without exposing raw reset tokens
// to persistence or generated database types to application code.
type Service struct {
	store     authStore
	guardians ports.GuardianAuthenticator
	idp       PasswordUpdater
	mailer    mailer.Mailer
	random    io.Reader
	now       func() time.Time
	frontend  func() string
}

func NewService(
	store authStore,
	guardians ports.GuardianAuthenticator,
	idp PasswordUpdater,
	mailSender mailer.Mailer,
) *Service {
	return &Service{store: store, guardians: guardians, idp: idp, mailer: mailSender, random: rand.Reader, now: time.Now, frontend: mailer.FrontendURL}
}

// ForgotPassword verifies the caller knows a user's login identifier and initial-password secret, then mints a short-lived one-time token and emails a reset link — hand-rolled since ThunderID exposes no reset primitive. The token is never returned in the response: NIC/index numbers appear on ID cards/report cards, so aren't secret enough to also hand over the takeover token (see docs audit C-1).
func (s *Service) ForgotPassword(ctx context.Context, req ForgotPasswordRequest) (ForgotPasswordResponse, error) {
	user, err := s.store.userByEmail(ctx, req.Identifier)
	if err != nil || user.Role != req.Role {
		return ForgotPasswordResponse{}, ErrInvalidCredentials
	}

	switch req.Role {
	case authz.RoleTeacher:
		if !s.store.teacherCredentialsMatch(ctx, user.ID, req.Secret) {
			return ForgotPasswordResponse{}, ErrInvalidCredentials
		}
	case authz.RoleStudent:
		if !s.store.studentCredentialsMatch(ctx, user.ID, req.Secret) {
			return ForgotPasswordResponse{}, ErrInvalidCredentials
		}
	case authz.RoleParent:
		if err := s.guardians.VerifyCredentials(ctx, user.ID, req.Secret); err != nil {
			return ForgotPasswordResponse{}, ErrInvalidCredentials
		}
	default:
		// Unreachable — ForgotPasswordRequest.Role is already
		// constrained to teacher/student/parent by its binding tag.
		return ForgotPasswordResponse{}, ErrInvalidCredentials
	}

	if err := s.issueAndEmailResetToken(ctx, user.ID, user.Email); err != nil {
		return ForgotPasswordResponse{}, err
	}

	return ForgotPasswordResponse{
		Message: "If those details match an account, a password reset link has been sent to the email on file.",
	}, nil
}

func (s *Service) issueAndEmailResetToken(ctx context.Context, userID uuid.UUID, email string) error {
	raw := make([]byte, 32)
	if _, err := io.ReadFull(s.random, raw); err != nil {
		return fmt.Errorf("failed to generate reset token: %w", err)
	}
	token := hex.EncodeToString(raw)
	hash := hashResetToken(token)
	expiresAt := s.now().Add(passwordResetTokenTTL)

	if err := s.store.createResetToken(ctx, userID, hash, expiresAt); err != nil {
		return fmt.Errorf("failed to create reset token: %w", err)
	}

	resetLink := fmt.Sprintf("%s/reset-password?token=%s", s.frontend(), token)
	body := fmt.Sprintf(
		"A password reset was requested for your OpenSchool account.\n\n"+
			"Reset your password using the link below. It expires in %d minutes and can only be used once.\n\n%s\n\n"+
			"If you didn't request this, you can safely ignore this email.",
		int(passwordResetTokenTTL.Minutes()), resetLink,
	)
	if err := s.mailer.Send(ctx, email, "Reset your OpenSchool password", body); err != nil {
		return fmt.Errorf("failed to send reset email: %w", err)
	}

	return nil
}

// ResetPassword is the unauthenticated counterpart to ChangePassword — it
// trusts the one-time token from ForgotPassword instead of a JWT.
func (s *Service) ResetPassword(ctx context.Context, req ResetPasswordRequest) error {
	// Consume first: only one concurrent request may proceed to ThunderID.
	record, err := s.store.consumeResetToken(ctx, hashResetToken(req.Token))
	if err != nil {
		return ErrResetTokenInvalid
	}
	return s.setPassword(ctx, record.UserID, req.NewPassword)
}

// ChangePassword is used by an already-authenticated caller (profile action or first-login "Set a new password") — a verified session already exists, so no reset token is needed.
func (s *Service) ChangePassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	return s.setPassword(ctx, userID, newPassword)
}

// KeepDefaultPassword clears the must-change flag without touching the password — the first-login "Keep this password" choice.
func (s *Service) KeepDefaultPassword(ctx context.Context, userID uuid.UUID) error {
	return s.store.setMustChangePassword(ctx, userID, false)
}

func (s *Service) setPassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	user, err := s.store.userByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if err := s.idp.UpdateUser(ctx, userID.String(), user.Role, map[string]any{
		"password": newPassword,
	}); err != nil {
		return fmt.Errorf("failed to update identity provider password: %w", err)
	}

	return s.store.setMustChangePassword(ctx, userID, false)
}

func hashResetToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
