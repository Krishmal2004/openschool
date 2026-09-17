import { useMutation, useQuery } from "@tanstack/react-query";
import { classroomApi } from "@/features/timetable/api/classroom";
import type { ClassroomRequest } from "@/features/timetable/api/classroom";
import { timetableKeys } from "@/features/timetable/keys";
import { useInvalidate } from "@/shared/api/useInvalidate";

export const useClassrooms = () => useQuery({ queryKey: timetableKeys.classrooms(), queryFn: classroomApi.list });

export const useCreateClassroom = () => {
  const invalidate = useInvalidate();
  return useMutation({ mutationFn: (data: ClassroomRequest) => classroomApi.create(data), onSuccess: () => invalidate(timetableKeys.classrooms()) });
};

export const useUpdateClassroom = () => {
  const invalidate = useInvalidate();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: ClassroomRequest }) => classroomApi.update(id, data),
    onSuccess: () => invalidate(timetableKeys.classrooms()),
  });
};

export const useDeleteClassroom = () => {
  const invalidate = useInvalidate();
  return useMutation({ mutationFn: (id: string) => classroomApi.remove(id), onSuccess: () => invalidate(timetableKeys.classrooms()) });
};
