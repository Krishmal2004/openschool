import { useMutation, useQuery } from "@tanstack/react-query";
import { staffAttendanceApi } from "@/features/attendance/api/staffAttendance";
import type { MarkStaffAttendanceRequest } from "@/features/attendance/api/staffAttendance";
import { staffAttendanceKeys } from "@/features/attendance/keys";
import { useInvalidate } from "@/shared/api/useInvalidate";

export const useStaffAttendanceByDate = (date: string) =>
  useQuery({ queryKey: staffAttendanceKeys.byDate(date), queryFn: () => staffAttendanceApi.byDate(date), enabled: !!date });

export const useStaffAttendanceMonthlySummary = (year: number, month: number) =>
  useQuery({ queryKey: staffAttendanceKeys.monthlySummary(year, month), queryFn: () => staffAttendanceApi.monthlySummary(year, month) });

// The signed-in teacher's own attendance for one month.
export const useMyStaffAttendanceHistory = (year: number, month: number) =>
  useQuery({ queryKey: staffAttendanceKeys.myHistory(year, month), queryFn: () => staffAttendanceApi.myHistory(year, month) });

export const useMarkStaffAttendance = () => {
  const invalidate = useInvalidate();
  return useMutation({ mutationFn: (data: MarkStaffAttendanceRequest) => staffAttendanceApi.mark(data), onSuccess: () => invalidate(staffAttendanceKeys.all) });
};
