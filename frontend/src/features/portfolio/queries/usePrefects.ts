import { useMutation, useQuery } from "@tanstack/react-query";
import { prefectApi } from "@/features/portfolio/api/prefect";
import type { AssignPrefectRequest } from "@/features/portfolio/api/prefect";
import { prefectKeys } from "@/features/portfolio/keys";
import { useInvalidate } from "@/shared/api/useInvalidate";

export const usePrefects = (academicYearId: string) =>
  useQuery({ queryKey: prefectKeys.byYear(academicYearId), queryFn: () => prefectApi.listByYear(academicYearId), enabled: !!academicYearId });

// Years with at least one appointment, for the archive selector.
export const usePrefectYears = () => useQuery({ queryKey: prefectKeys.years(), queryFn: () => prefectApi.listYears() });

export const useAssignPrefect = () => {
  const invalidate = useInvalidate();
  return useMutation({ mutationFn: (data: AssignPrefectRequest) => prefectApi.assign(data), onSuccess: () => invalidate(prefectKeys.all) });
};

export const useRemovePrefect = () => {
  const invalidate = useInvalidate();
  return useMutation({ mutationFn: ({ id }: { id: string; academicYearId: string }) => prefectApi.remove(id), onSuccess: () => invalidate(prefectKeys.all) });
};
