import { useMutation, useQuery } from "@tanstack/react-query";
import { houseApi } from "@/features/school/api/house";
import type { CreateHouseRequest } from "@/features/school/api/house";
import { schoolKeys } from "@/features/school/keys";
import { studentKeys } from "@/features/students/keys";
import { teacherKeys } from "@/features/teachers/keys";
import { useInvalidate } from "@/shared/api/useInvalidate";

export const useHouses = () => useQuery({ queryKey: schoolKeys.houses(), queryFn: houseApi.list });

export const useCreateHouse = () => {
  const invalidate = useInvalidate();
  return useMutation({ mutationFn: (data: CreateHouseRequest) => houseApi.create(data), onSuccess: () => invalidate(schoolKeys.houses()) });
};

export const useUpdateHouse = () => {
  const invalidate = useInvalidate();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: CreateHouseRequest }) => houseApi.update(id, data),
    onSuccess: () => invalidate(schoolKeys.houses()),
  });
};

export const useDeleteHouse = () => {
  const invalidate = useInvalidate();
  return useMutation({ mutationFn: (id: string) => houseApi.remove(id), onSuccess: () => invalidate(schoolKeys.houses()) });
};

// Bulk-assigns houses to students that have none; their cached rows change.
export const useReassignMissingHouses = () => {
  const invalidate = useInvalidate();
  return useMutation({ mutationFn: () => houseApi.reassignMissing(), onSuccess: () => invalidate(studentKeys.all) });
};

export const useReassignMissingStaffHouses = () => {
  const invalidate = useInvalidate();
  return useMutation({ mutationFn: () => houseApi.reassignMissingStaff(), onSuccess: () => invalidate(teacherKeys.all) });
};
