import { queryOptions, useMutation, useQuery } from "@tanstack/react-query";
import { attendanceApi } from "@/features/attendance/api/attendance";
import type { CreateSessionRequest, MarkAttendanceRequest } from "@/features/attendance/api/attendance";
import { attendanceKeys } from "@/features/attendance/keys";
import { useInvalidate } from "@/shared/api/useInvalidate";

// Exposed as options so other features can batch it with useQueries.
export const classSessionsOptions = (classId: string) =>
  queryOptions({
    queryKey: attendanceKeys.classSessions(classId),
    queryFn: () => attendanceApi.listSessionsByClass(classId),
    enabled: !!classId,
  });

export const useClassSessions = (classId: string) => useQuery(classSessionsOptions(classId));

// Every record a student has across sessions, for the portfolio attendance rollup.
export const useStudentAttendanceHistory = (studentId: string) =>
  useQuery({ queryKey: attendanceKeys.byStudent(studentId), queryFn: () => attendanceApi.listByStudent(studentId), enabled: !!studentId });

export const useDailySessions = (date: string) =>
  useQuery({ queryKey: attendanceKeys.dailySessions(date), queryFn: () => attendanceApi.listSessionsByDate(date), enabled: !!date });

export const useSession = (id: string) =>
  useQuery({ queryKey: attendanceKeys.session(id), queryFn: () => attendanceApi.getSession(id), enabled: !!id });

export const useSessionRecords = (id: string) =>
  useQuery({ queryKey: attendanceKeys.sessionRecords(id), queryFn: () => attendanceApi.listRecords(id), enabled: !!id });

export const useCreateSession = () => {
  const invalidate = useInvalidate();
  return useMutation({ mutationFn: (data: CreateSessionRequest) => attendanceApi.createSession(data), onSuccess: () => invalidate(attendanceKeys.all) });
};

export const useMarkAttendance = (sessionId: string) => {
  const invalidate = useInvalidate();
  return useMutation({
    mutationFn: (data: MarkAttendanceRequest) => attendanceApi.markAttendance(sessionId, data),
    onSuccess: () => invalidate(attendanceKeys.all),
  });
};

export const useDeleteSession = () => {
  const invalidate = useInvalidate();
  return useMutation({ mutationFn: (id: string) => attendanceApi.deleteSession(id), onSuccess: () => invalidate(attendanceKeys.all) });
};
