import { useMutation, useQuery } from "@tanstack/react-query";
import { subjectApi } from "@/features/curriculum/api/subject";
import type { CreateSubjectRequest } from "@/features/curriculum/api/subject";
import { curriculumKeys } from "@/features/curriculum/keys";
import { useInvalidate } from "@/shared/api/useInvalidate";
import { REFERENCE_DATA_STALE_TIME_MS } from "@/shared/api/staleTime";

export const useSubjects = () =>
  useQuery({ queryKey: curriculumKeys.subjects(), queryFn: subjectApi.list, staleTime: REFERENCE_DATA_STALE_TIME_MS });

export const useCreateSubject = () => {
  const invalidate = useInvalidate();
  return useMutation({ mutationFn: (data: CreateSubjectRequest) => subjectApi.create(data), onSuccess: () => invalidate(curriculumKeys.subjects()) });
};

export const useUpdateSubject = () => {
  const invalidate = useInvalidate();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: CreateSubjectRequest }) => subjectApi.update(id, data),
    onSuccess: () => invalidate(curriculumKeys.subjects()),
  });
};

export const useDeleteSubject = () => {
  const invalidate = useInvalidate();
  return useMutation({ mutationFn: (id: string) => subjectApi.remove(id), onSuccess: () => invalidate(curriculumKeys.subjects()) });
};
