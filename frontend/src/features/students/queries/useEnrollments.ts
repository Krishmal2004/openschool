import { useMutation, useQuery } from "@tanstack/react-query";
import { enrollmentApi } from "@/features/students/api/enrollment";
import type { SubmitEnrollmentRequest } from "@/features/students/api/enrollment";
import { studentKeys } from "@/features/students/keys";
import { useInvalidate } from "@/shared/api/useInvalidate";

export const useStudentEnrollments = (studentId: string, academicYearId: string) =>
  useQuery({
    queryKey: studentKeys.enrollments(studentId, academicYearId),
    queryFn: () => enrollmentApi.listByStudent(studentId, academicYearId),
    enabled: !!studentId && !!academicYearId,
  });

// Validates first; an invalid pick set is returned to the caller instead of submitted.
export const useSubmitEnrollments = (studentId: string, academicYearId: string) => {
  const invalidate = useInvalidate();
  return useMutation({
    mutationFn: async (data: SubmitEnrollmentRequest) => {
      const check = await enrollmentApi.validate(data);
      if (!check.valid) return check;
      return enrollmentApi.submit(studentId, data);
    },
    onSuccess: (result) => {
      if (result.valid) invalidate(studentKeys.enrollments(studentId, academicYearId));
    },
  });
};
