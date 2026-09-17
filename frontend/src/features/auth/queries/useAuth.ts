import { useMutation } from "@tanstack/react-query";
import { authApi } from "@/features/auth/api/auth";
import type { ForgotPasswordRequest, ResetPasswordRequest } from "@/features/auth/api/auth";
import { keys } from "@/shared/api/keys";
import { useInvalidate } from "@/shared/api/useInvalidate";

export const useForgotPassword = () => useMutation({ mutationFn: (data: ForgotPasswordRequest) => authApi.forgotPassword(data) });

export const useResetPassword = () => useMutation({ mutationFn: (data: ResetPasswordRequest) => authApi.resetPassword(data) });

export const useChangePassword = () => {
  const invalidate = useInvalidate();
  return useMutation({
    mutationFn: (newPassword: string) => authApi.changePassword({ new_password: newPassword }),
    onSuccess: () => invalidate(keys.me.all),
  });
};

export const useKeepDefaultPassword = () => {
  const invalidate = useInvalidate();
  return useMutation({ mutationFn: () => authApi.keepDefaultPassword(), onSuccess: () => invalidate(keys.me.all) });
};
