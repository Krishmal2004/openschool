const PASSWORD_MIN_LENGTH = 8;

// One place for the password policy shown in every password form.
export function validateNewPassword(password: string, confirm: string) {
  const tooShort = password.length > 0 && password.length < PASSWORD_MIN_LENGTH;
  const mismatch = confirm.length > 0 && confirm !== password;
  return {
    valid: password.length >= PASSWORD_MIN_LENGTH && confirm === password,
    passwordError: tooShort ? `Must be at least ${PASSWORD_MIN_LENGTH} characters.` : undefined,
    confirmError: mismatch ? "Passwords do not match." : undefined,
  };
}
