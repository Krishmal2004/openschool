import { useMutation, useQuery } from "@tanstack/react-query";
import { termMarkApi } from "@/features/marks/api/termMark";
import type { BulkUpsertMarksRequest } from "@/features/marks/api/termMark";
import { markKeys } from "@/features/marks/keys";
import { useInvalidate } from "@/shared/api/useInvalidate";

export const useClassMarks = (classId: string, termId: string, subjectId: string) =>
  useQuery({
    queryKey: markKeys.byClass(classId, termId, subjectId),
    queryFn: () => termMarkApi.listClassMarks(classId, termId, subjectId),
    enabled: !!classId && !!termId && !!subjectId,
  });

export const useStudentMarks = (studentId: string, termId: string) =>
  useQuery({
    queryKey: markKeys.byStudent(studentId, termId),
    queryFn: () => termMarkApi.listStudentMarks(studentId, termId),
    enabled: !!studentId && !!termId,
  });

export const useSaveClassMarks = (classId: string) => {
  const invalidate = useInvalidate();
  return useMutation({ mutationFn: (data: BulkUpsertMarksRequest) => termMarkApi.bulkUpsert(classId, data), onSuccess: () => invalidate(markKeys.all) });
};
