import { useMutation, useQuery } from "@tanstack/react-query";
import { setupApi } from "@/features/school/api/setup";
import type { RegisterAdminRequest } from "@/features/school/api/setup";
import { schoolKeys } from "@/features/school/keys";

// Only matters before the first admin exists, so it never goes stale.
export const useSetupStatus = () =>
  useQuery({ queryKey: schoolKeys.setupStatus(), queryFn: setupApi.status, staleTime: Infinity, retry: false });

export const useRegisterAdmin = () => useMutation({ mutationFn: (data: RegisterAdminRequest) => setupApi.registerAdmin(data) });
