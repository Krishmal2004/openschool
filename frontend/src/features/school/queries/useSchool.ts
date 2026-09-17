import { useMutation, useQuery } from "@tanstack/react-query";
import { schoolApi } from "@/features/school/api/school";
import type { CreateSchoolRequest } from "@/features/school/api/school";
import { schoolKeys } from "@/features/school/keys";
import { useInvalidate } from "@/shared/api/useInvalidate";

export const useSchool = () => useQuery({ queryKey: schoolKeys.info(), queryFn: schoolApi.get });

export const useCreateSchool = () => {
  const invalidate = useInvalidate();
  return useMutation({ mutationFn: (data: CreateSchoolRequest) => schoolApi.create(data), onSuccess: () => invalidate(schoolKeys.info()) });
};

export const useUpdateSchool = () => {
  const invalidate = useInvalidate();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: CreateSchoolRequest }) => schoolApi.update(id, data),
    onSuccess: () => invalidate(schoolKeys.info()),
  });
};
