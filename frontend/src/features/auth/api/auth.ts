import api from "@/shared/api/client";

export type SelfServiceResetRole = "teacher" | "student" | "parent";

// identifier is the account email; secret is the initial default password (NIC for teacher/parent, index number for student). Admins have no secret on file.
export interface ForgotPasswordRequest {
  role: SelfServiceResetRole;
  identifier: string;
  secret: string;
}

// Deliberately carries no token, the reset link is delivered out-of-band by
// email, not handed back to whoever called this endpoint.
export interface ForgotPasswordResponse {
  message: string;
}

export interface ResetPasswordRequest {
  token: string;
  new_password: string;
}

export interface ChangePasswordRequest {
  new_password: string;
}

export const authApi = {
  forgotPassword: (data: ForgotPasswordRequest) =>
    api.post<ForgotPasswordResponse>("/auth/forgot-password", data).then((r) => r.data),

  resetPassword: (data: ResetPasswordRequest) =>
    api.post("/auth/reset-password", data).then((r) => r.data),

  // Requires an authenticated session, used for both the "Change password"
  // profile action and the first-login "Set a new password" choice.
  changePassword: (data: ChangePasswordRequest) =>
    api.post("/auth/change-password", data).then((r) => r.data),

  // The first-login "Keep this password" choice, clears the must-change
  // flag without touching the password.
  keepDefaultPassword: () => api.post("/auth/keep-default-password").then((r) => r.data),
};
