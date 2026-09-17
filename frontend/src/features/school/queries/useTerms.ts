import { useMutation, useQuery } from "@tanstack/react-query";
import { termApi } from "@/features/school/api/term";
import type { CreateTermRequest, UpdateTermRequest } from "@/features/school/api/term";
import { termKeys } from "@/features/school/keys";
import { useInvalidate } from "@/shared/api/useInvalidate";

export const useTerms = (academicYearId: string | undefined) =>
  useQuery({
    queryKey: termKeys.byYear(academicYearId ?? ""),
    queryFn: () => termApi.listByAcademicYear(academicYearId!),
    enabled: !!academicYearId,
  });

// 404 when no current term is set; callers fall back to manual term selection.
export const useCurrentTerm = () => useQuery({ queryKey: termKeys.current(), queryFn: termApi.getCurrent, retry: false });

export const useCreateTerm = () => {
  const invalidate = useInvalidate();
  return useMutation({ mutationFn: (data: CreateTermRequest) => termApi.create(data), onSuccess: () => invalidate(termKeys.all) });
};

export const useUpdateTerm = () => {
  const invalidate = useInvalidate();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateTermRequest }) => termApi.update(id, data),
    onSuccess: () => invalidate(termKeys.all),
  });
};

export const useSetCurrentTerm = () => {
  const invalidate = useInvalidate();
  return useMutation({ mutationFn: (id: string) => termApi.setCurrent(id), onSuccess: () => invalidate(termKeys.all) });
};

export const useDeleteTerm = () => {
  const invalidate = useInvalidate();
  return useMutation({ mutationFn: (id: string) => termApi.remove(id), onSuccess: () => invalidate(termKeys.all) });
};
