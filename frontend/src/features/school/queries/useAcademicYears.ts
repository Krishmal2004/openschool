import { useMutation, useQuery } from "@tanstack/react-query";
import { academicYearApi } from "@/features/school/api/academicYear";
import type { CreateAcademicYearRequest } from "@/features/school/api/academicYear";
import { academicYearKeys } from "@/features/school/keys";
import { useInvalidate } from "@/shared/api/useInvalidate";
import { REFERENCE_DATA_STALE_TIME_MS } from "@/shared/api/staleTime";

export const useAcademicYears = () =>
  useQuery({ queryKey: academicYearKeys.list(), queryFn: academicYearApi.list, staleTime: REFERENCE_DATA_STALE_TIME_MS });

export const useCurrentAcademicYear = () => useQuery({ queryKey: academicYearKeys.current(), queryFn: academicYearApi.getCurrent });

export const useCreateAcademicYear = () => {
  const invalidate = useInvalidate();
  return useMutation({ mutationFn: (data: CreateAcademicYearRequest) => academicYearApi.create(data), onSuccess: () => invalidate(academicYearKeys.all) });
};

export const useSetCurrentAcademicYear = () => {
  const invalidate = useInvalidate();
  return useMutation({ mutationFn: (id: string) => academicYearApi.setCurrent(id), onSuccess: () => invalidate(academicYearKeys.all) });
};

export const useDeleteAcademicYear = () => {
  const invalidate = useInvalidate();
  return useMutation({ mutationFn: (id: string) => academicYearApi.remove(id), onSuccess: () => invalidate(academicYearKeys.all) });
};
