import { useMutation, useQuery } from "@tanstack/react-query";
import { sectionHeadApi } from "@/features/teachers/api/sectionHead";
import type { AssignSectionHeadRequest } from "@/features/teachers/api/sectionHead";
import { sectionHeadKeys } from "@/features/teachers/keys";
import { useInvalidate } from "@/shared/api/useInvalidate";

export const useSectionHeads = (academicYearId: string) =>
  useQuery({
    queryKey: sectionHeadKeys.byYear(academicYearId),
    queryFn: () => sectionHeadApi.listByYear(academicYearId),
    enabled: !!academicYearId,
  });

export const useAssignSectionHead = () => {
  const invalidate = useInvalidate();
  return useMutation({
    mutationFn: (data: AssignSectionHeadRequest) => sectionHeadApi.assign(data),
    onSuccess: (_result, variables) => invalidate(sectionHeadKeys.byYear(variables.academic_year_id)),
  });
};

export const useRemoveSectionHead = () => {
  const invalidate = useInvalidate();
  return useMutation({
    mutationFn: ({ id }: { id: string; academicYearId: string }) => sectionHeadApi.remove(id),
    onSuccess: (_result, variables) => invalidate(sectionHeadKeys.byYear(variables.academicYearId)),
  });
};
