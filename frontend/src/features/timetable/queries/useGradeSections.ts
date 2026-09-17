import { useMutation, useQuery } from "@tanstack/react-query";
import { gradeSectionApi } from "@/features/timetable/api/gradeSection";
import type { CreateGradeSectionRequest, UpdateGradeSectionRequest, TimetablePeriod } from "@/features/timetable/api/gradeSection";
import { timetableKeys } from "@/features/timetable/keys";
import { useInvalidate } from "@/shared/api/useInvalidate";

export const useGradeSections = (academicYearId: string) =>
  useQuery({ queryKey: timetableKeys.gradeSections(academicYearId), queryFn: () => gradeSectionApi.listByYear(academicYearId), enabled: !!academicYearId });

export const useCreateGradeSection = () => {
  const invalidate = useInvalidate();
  return useMutation({
    mutationFn: (data: CreateGradeSectionRequest) => gradeSectionApi.create(data),
    onSuccess: (_r, variables) => invalidate(timetableKeys.gradeSections(variables.academic_year_id)),
  });
};

export const useUpdateGradeSection = (academicYearId: string) => {
  const invalidate = useInvalidate();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateGradeSectionRequest }) => gradeSectionApi.update(id, data),
    onSuccess: () => invalidate(timetableKeys.gradeSections(academicYearId)),
  });
};

export const useDeleteGradeSection = (academicYearId: string) => {
  const invalidate = useInvalidate();
  return useMutation({ mutationFn: (id: string) => gradeSectionApi.remove(id), onSuccess: () => invalidate(timetableKeys.gradeSections(academicYearId)) });
};

export const useAssignGradesToSection = (academicYearId: string) => {
  const invalidate = useInvalidate();
  return useMutation({
    mutationFn: ({ id, gradeIds }: { id: string; gradeIds: string[] }) => gradeSectionApi.assignGrades(id, gradeIds),
    onSuccess: () => invalidate(timetableKeys.gradeSections(academicYearId)),
  });
};

export const useGradeSectionPeriods = (id: string) =>
  useQuery({ queryKey: timetableKeys.gradeSectionPeriods(id), queryFn: () => gradeSectionApi.getPeriods(id), enabled: !!id });

export const useSaveGradeSectionPeriods = (id: string) => {
  const invalidate = useInvalidate();
  return useMutation({
    mutationFn: (periods: Omit<TimetablePeriod, "id">[]) => gradeSectionApi.savePeriods(id, periods),
    onSuccess: () => invalidate(timetableKeys.gradeSectionPeriods(id)),
  });
};

export const useRegenerateGradeSectionPeriods = (id: string) => {
  const invalidate = useInvalidate();
  return useMutation({ mutationFn: () => gradeSectionApi.regeneratePeriods(id), onSuccess: () => invalidate(timetableKeys.gradeSectionPeriods(id)) });
};
