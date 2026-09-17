import { useMutation, useQuery } from "@tanstack/react-query";
import { promotionApi } from "@/features/academics/api/promotion";
import type { CommitAssignmentsRequest } from "@/features/academics/api/promotion";
import { classKeys, promotionKeys } from "@/features/academics/keys";
import { useInvalidate } from "@/shared/api/useInvalidate";

export const usePromotionPreview = (sourceYearId: string, targetYearId: string, rankByTermId?: string) =>
  useQuery({
    queryKey: promotionKeys.preview(sourceYearId, targetYearId, rankByTermId),
    queryFn: () =>
      promotionApi.preview({
        source_year_id: sourceYearId,
        target_year_id: targetYearId,
        rank_by_term_id: rankByTermId || undefined,
      }),
    enabled: !!sourceYearId && !!targetYearId,
  });

export const useCommitAssignments = () => {
  const invalidate = useInvalidate();
  return useMutation({
    mutationFn: (data: CommitAssignmentsRequest) => promotionApi.commit(data),
    // The target year's class rosters just changed.
    onSuccess: (_result, variables) => invalidate(classKeys.byYear(variables.academic_year_id), promotionKeys.all),
  });
};
