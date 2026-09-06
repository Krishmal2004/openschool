import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { termApi } from "../services/term";
import type { CreateTermRequest, UpdateTermRequest } from "../services/term";

export const termsKey = (academicYearId: string) => ["terms", academicYearId];
export const CURRENT_TERM_KEY = ["terms", "current"];

export const useTerms = (academicYearId: string | undefined) =>
  useQuery({
    queryKey: termsKey(academicYearId ?? ""),
    queryFn: () => termApi.listByAcademicYear(academicYearId!),
    enabled: !!academicYearId,
  });

// The term marked is_current — mirrors useCurrentAcademicYear. 404s (no
// current term set yet) surface as isError; callers should fall back to
// manual term selection in that case rather than blocking the page.
export const useCurrentTerm = () =>
  useQuery({
    queryKey: CURRENT_TERM_KEY,
    queryFn: termApi.getCurrent,
    retry: false,
  });

export const useCreateTerm = (academicYearId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateTermRequest) => termApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: termsKey(academicYearId) });
    },
  });
};

export const useUpdateTerm = (academicYearId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateTermRequest }) =>
      termApi.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: termsKey(academicYearId) });
    },
  });
};

export const useSetCurrentTerm = (academicYearId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => termApi.setCurrent(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: termsKey(academicYearId) });
      queryClient.invalidateQueries({ queryKey: CURRENT_TERM_KEY });
    },
  });
};

export const useDeleteTerm = (academicYearId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => termApi.remove(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: termsKey(academicYearId) });
    },
  });
};
