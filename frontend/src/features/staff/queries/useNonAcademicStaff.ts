import { keepPreviousData, useMutation, useQuery } from "@tanstack/react-query";
import { nonAcademicStaffApi } from "@/features/staff/api/nonAcademicStaff";
import type {
  CreateNonAcademicStaffRequest,
  UpdateNonAcademicStaffRequest,
  NonAcademicEmploymentStatus,
  StaffListParams,
} from "@/features/staff/api/nonAcademicStaff";
import { staffKeys } from "@/features/staff/keys";
import { useInvalidate } from "@/shared/api/useInvalidate";

// Server-paginated (docs/SECURITY_AND_PERFORMANCE_PLAYBOOK.md section 4).
export const useNonAcademicStaffList = (params: StaffListParams = {}) =>
  useQuery({
    queryKey: staffKeys.list(params),
    queryFn: () => nonAcademicStaffApi.list(params),
    placeholderData: keepPreviousData,
  });

export const useCreateNonAcademicStaff = () => {
  const invalidate = useInvalidate();
  return useMutation({ mutationFn: (data: CreateNonAcademicStaffRequest) => nonAcademicStaffApi.create(data), onSuccess: () => invalidate(staffKeys.all) });
};

export const useUpdateNonAcademicStaff = () => {
  const invalidate = useInvalidate();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateNonAcademicStaffRequest }) => nonAcademicStaffApi.update(id, data),
    onSuccess: () => invalidate(staffKeys.all),
  });
};

export const useUpdateNonAcademicStaffEmploymentStatus = () => {
  const invalidate = useInvalidate();
  return useMutation({
    mutationFn: ({ id, status }: { id: string; status: NonAcademicEmploymentStatus }) => nonAcademicStaffApi.updateEmploymentStatus(id, status),
    onSuccess: () => invalidate(staffKeys.all),
  });
};

export const useUpdateNonAcademicStaffHouse = () => {
  const invalidate = useInvalidate();
  return useMutation({
    mutationFn: ({ id, houseId }: { id: string; houseId: string }) => nonAcademicStaffApi.updateHouse(id, houseId),
    onSuccess: () => invalidate(staffKeys.all),
  });
};

export const useDeleteNonAcademicStaff = () => {
  const invalidate = useInvalidate();
  return useMutation({ mutationFn: (id: string) => nonAcademicStaffApi.remove(id), onSuccess: () => invalidate(staffKeys.all) });
};
