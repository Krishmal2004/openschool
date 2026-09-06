import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { staffAttendanceApi } from "../services/staffAttendance";
import type { MarkStaffAttendanceRequest } from "../services/staffAttendance";

const byDateKey = (date: string) => ["staff-attendance", "by-date", date];
const monthlySummaryKey = (year: number, month: number) => ["staff-attendance", "monthly-summary", year, month];

export const useStaffAttendanceByDate = (date: string) =>
  useQuery({
    queryKey: byDateKey(date),
    queryFn: () => staffAttendanceApi.byDate(date),
    enabled: !!date,
  });

export const useStaffAttendanceMonthlySummary = (year: number, month: number) =>
  useQuery({
    queryKey: monthlySummaryKey(year, month),
    queryFn: () => staffAttendanceApi.monthlySummary(year, month),
  });

const myHistoryKey = (year: number, month: number) => ["me", "teacher", "attendance", year, month];

// The signed-in teacher's own staff-attendance history for a month.
export const useMyStaffAttendanceHistory = (year: number, month: number) =>
  useQuery({
    queryKey: myHistoryKey(year, month),
    queryFn: () => staffAttendanceApi.myHistory(year, month),
  });

export const useMarkStaffAttendance = (date: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: MarkStaffAttendanceRequest) => staffAttendanceApi.mark(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: byDateKey(date) });
      queryClient.invalidateQueries({ queryKey: ["staff-attendance", "monthly-summary"] });
    },
  });
};
