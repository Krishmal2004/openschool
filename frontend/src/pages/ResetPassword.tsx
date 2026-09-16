import { useState } from "react";
import { Link, useSearchParams } from "react-router";
import { Button, PasswordInput, InlineNotification } from "@carbon/react";
import { CheckmarkFilled } from "@carbon/icons-react";
import { useResetPassword } from "../queries/useAuth";
import { getErrorMessage } from "../lib/errorMessage";

export default function ResetPassword() {
  const [searchParams] = useSearchParams();
  const token = searchParams.get("token") ?? "";
  const resetPassword = useResetPassword();

  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [done, setDone] = useState(false);

  const canReset = token.length > 0 && newPassword.length >= 8 && newPassword === confirmPassword;

  const handleReset = () => {
    if (!canReset) return;
    resetPassword.mutate({ token, new_password: newPassword }, { onSuccess: () => setDone(true) });
  };

  if (done) {
    return (
      <div className="os-signin-wrapper">
        <div className="os-auth-card os-auth-card--center">
          <div className="os-auth-card__success-icon">
            <CheckmarkFilled size={28} />
          </div>
          <h1 className="os-auth-card__title">Password updated</h1>
          <p className="os-auth-card__subtitle">You can now sign in with your new password.</p>
          <Button href="/signin" className="os-full-width-btn">
            Go to Sign In
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="os-signin-wrapper">
      <div className="os-auth-card">
        <div className="os-auth-card__brand">
          <img
            src="/favicon.webp"
            alt="OpenSchool"
            width={36}
            height={36}
            className="os-auth-card__logo"
          />
          <span className="os-auth-card__brand-name">OpenSchool</span>
        </div>
        <h1 className="os-auth-card__title">Reset password</h1>
        <p className="os-auth-card__subtitle">Choose a new password for your account.</p>

        {!token && (
          <InlineNotification
            kind="error"
            title="Invalid reset link"
            subtitle="This link is missing its reset token. Request a new one below."
            lowContrast
            hideCloseButton
            style={{ marginBottom: "1.25rem", maxWidth: "100%" }}
          />
        )}
        {resetPassword.isError && (
          <InlineNotification
            kind="error"
            title="Could not reset password"
            subtitle={getErrorMessage(resetPassword.error, "The reset link is invalid, already used, or has expired.")}
            lowContrast
            hideCloseButton
            style={{ marginBottom: "1.25rem", maxWidth: "100%" }}
          />
        )}
        <div className="os-auth-card__form">
          <PasswordInput
            id="reset-password-new"
            labelText="New Password"
            value={newPassword}
            onChange={(e) => setNewPassword(e.target.value)}
            invalid={newPassword.length > 0 && newPassword.length < 8}
            invalidText="Must be at least 8 characters."
          />
          <PasswordInput
            id="reset-password-confirm"
            labelText="Confirm New Password"
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
            invalid={confirmPassword.length > 0 && confirmPassword !== newPassword}
            invalidText="Passwords do not match."
          />
        </div>
        <div className="os-auth-card__actions">
          <Button
            className="os-full-width-btn"
            onClick={handleReset}
            disabled={!canReset || resetPassword.isPending}
          >
            {resetPassword.isPending ? "Saving…" : "Set new password"}
          </Button>
        </div>

        <div className="os-auth-card__footer">
          <Link to="/forgot-password">Request a new link</Link> · <Link to="/signin">Back to sign in</Link>
        </div>
      </div>
    </div>
  );
}
