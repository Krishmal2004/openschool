import { useMutation, useQuery } from "@tanstack/react-query";
import { subjectPeriodRequirementApi } from "@/features/timetable/api/subjectPeriodRequirement";
import type { UpsertSubjectPeriodRequirementRequest } from "@/features/timetable/api/subjectPeriodRequirement";
import { timetableKeys } from "@/features/timetable/keys";
import { useInvalidate } from "@/shared/api/useInvalidate";

export const useSubjectPeriodRequirements = (academicYearId: string, gradeId: string) =>
  useQuery({
    queryKey: timetableKeys.subjectRequirements(academicYearId, gradeId),
    queryFn: () => subjectPeriodRequirementApi.listByGrade(academicYearId, gradeId),
    enabled: !!academicYearId && !!gradeId,
  });

export const useUpsertSubjectPeriodRequirement = () => {
  const invalidate = useInvalidate();
  return useMutation({
    mutationFn: (data: UpsertSubjectPeriodRequirementRequest) => subjectPeriodRequirementApi.upsert(data),
    onSuccess: (_r, variables) => invalidate(timetableKeys.subjectRequirements(variables.academic_year_id, variables.grade_id)),
  });
};

export const useDeleteSubjectPeriodRequirement = (academicYearId: string, gradeId: string) => {
  const invalidate = useInvalidate();
  return useMutation({
    mutationFn: (id: string) => subjectPeriodRequirementApi.remove(id),
    onSuccess: () => invalidate(timetableKeys.subjectRequirements(academicYearId, gradeId)),
  });
};
