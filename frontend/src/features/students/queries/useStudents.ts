import { queryOptions, useMutation, useQuery } from "@tanstack/react-query";
import { studentApi } from "@/features/students/api/student";
import type { CreateStudentRequest, UpdateStudentRequest, StudentEnrollmentStatus } from "@/features/students/api/student";
import { studentKeys } from "@/features/students/keys";
import { useInvalidate } from "@/shared/api/useInvalidate";

export const useStudents = () => useQuery({ queryKey: studentKeys.list(), queryFn: studentApi.list });

export const useStudentWithClass = (id: string) =>
  useQuery({ queryKey: studentKeys.withClass(id), queryFn: () => studentApi.getWithClass(id), enabled: !!id });

// Exposed as options so other features can batch it with useQueries.
export const studentsByClassOptions = (classId: string) =>
  queryOptions({
    queryKey: studentKeys.byClass(classId),
    queryFn: () => studentApi.listByClass(classId),
    enabled: !!classId,
  });

export const useStudentsByClass = (classId: string) => useQuery(studentsByClassOptions(classId));

export const useCreateStudent = () => {
  const invalidate = useInvalidate();
  return useMutation({
    mutationFn: (data: CreateStudentRequest) => studentApi.create(data),
    onSuccess: () => invalidate(studentKeys.all),
  });
};

export const useUpdateStudent = () => {
  const invalidate = useInvalidate();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateStudentRequest }) => studentApi.update(id, data),
    onSuccess: () => invalidate(studentKeys.all),
  });
};

export const useUpdateStudentHouse = () => {
  const invalidate = useInvalidate();
  return useMutation({
    mutationFn: ({ id, houseId }: { id: string; houseId: string }) => studentApi.updateHouse(id, houseId),
    onSuccess: () => invalidate(studentKeys.all),
  });
};

export const useUpdateStudentEnrollmentStatus = () => {
  const invalidate = useInvalidate();
  return useMutation({
    mutationFn: ({ id, status }: { id: string; status: StudentEnrollmentStatus }) => studentApi.updateEnrollmentStatus(id, status),
    onSuccess: () => invalidate(studentKeys.all),
  });
};

export const useDeleteStudent = () => {
  const invalidate = useInvalidate();
  return useMutation({
    mutationFn: (id: string) => studentApi.remove(id),
    onSuccess: () => invalidate(studentKeys.all),
  });
};
