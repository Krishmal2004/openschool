import { useMutation, useQuery } from "@tanstack/react-query";
import { identityReconciliationApi } from "@/features/system/api/identityReconciliation";
import { systemKeys } from "@/features/system/keys";
import { useInvalidate } from "@/shared/api/useInvalidate";

// Walks every ThunderID user, so it only runs while the admin has the tab open.
export const useOrphanedAccounts = (enabled: boolean) =>
  useQuery({ queryKey: systemKeys.orphanedAccounts(), queryFn: identityReconciliationApi.list, enabled });

export const useDeleteOrphanedAccount = () => {
  const invalidate = useInvalidate();
  return useMutation({ mutationFn: (id: string) => identityReconciliationApi.remove(id), onSuccess: () => invalidate(systemKeys.orphanedAccounts()) });
};
