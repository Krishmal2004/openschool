import { useMutation, useQuery } from "@tanstack/react-query";
import { studentPortfolioApi } from "@/features/portfolio/api/studentPortfolio";
import type {
  CreateProgressReportRequest,
  CreateActivityRequest,
  CreateLeadershipRoleRequest,
  CreateAwardRequest,
  CreateDisciplinaryRecordRequest,
} from "@/features/portfolio/api/studentPortfolio";
import { portfolioKeys } from "@/features/portfolio/keys";
import { useInvalidate } from "@/shared/api/useInvalidate";

// Each record type follows the same list / create / delete shape.
const useCreate = <T,>(key: readonly unknown[], fn: (data: T) => Promise<unknown>) => {
  const invalidate = useInvalidate();
  return useMutation({ mutationFn: fn, onSuccess: () => invalidate(key) });
};

export const useProgressReports = (studentId: string) =>
  useQuery({ queryKey: portfolioKeys.progressReports(studentId), queryFn: () => studentPortfolioApi.listProgressReports(studentId), enabled: !!studentId });

export const useCreateProgressReport = (studentId: string) =>
  useCreate(portfolioKeys.progressReports(studentId), (data: CreateProgressReportRequest) => studentPortfolioApi.createProgressReport(studentId, data));

export const useDeleteProgressReport = (studentId: string) =>
  useCreate(portfolioKeys.progressReports(studentId), (recordId: string) => studentPortfolioApi.deleteProgressReport(studentId, recordId));

export const useStudentActivities = (studentId: string) =>
  useQuery({ queryKey: portfolioKeys.activities(studentId), queryFn: () => studentPortfolioApi.listActivities(studentId), enabled: !!studentId });

export const useCreateStudentActivity = (studentId: string) =>
  useCreate(portfolioKeys.activities(studentId), (data: CreateActivityRequest) => studentPortfolioApi.createActivity(studentId, data));

export const useDeleteStudentActivity = (studentId: string) =>
  useCreate(portfolioKeys.activities(studentId), (recordId: string) => studentPortfolioApi.deleteActivity(studentId, recordId));

export const useLeadershipRoles = (studentId: string) =>
  useQuery({ queryKey: portfolioKeys.leadershipRoles(studentId), queryFn: () => studentPortfolioApi.listLeadershipRoles(studentId), enabled: !!studentId });

export const useCreateLeadershipRole = (studentId: string) =>
  useCreate(portfolioKeys.leadershipRoles(studentId), (data: CreateLeadershipRoleRequest) => studentPortfolioApi.createLeadershipRole(studentId, data));

export const useDeleteLeadershipRole = (studentId: string) =>
  useCreate(portfolioKeys.leadershipRoles(studentId), (recordId: string) => studentPortfolioApi.deleteLeadershipRole(studentId, recordId));

export const useStudentAwards = (studentId: string) =>
  useQuery({ queryKey: portfolioKeys.awards(studentId), queryFn: () => studentPortfolioApi.listAwards(studentId), enabled: !!studentId });

export const useCreateStudentAward = (studentId: string) =>
  useCreate(portfolioKeys.awards(studentId), (data: CreateAwardRequest) => studentPortfolioApi.createAward(studentId, data));

export const useDeleteStudentAward = (studentId: string) =>
  useCreate(portfolioKeys.awards(studentId), (recordId: string) => studentPortfolioApi.deleteAward(studentId, recordId));

export const useDisciplinaryRecords = (studentId: string) =>
  useQuery({ queryKey: portfolioKeys.disciplinaryRecords(studentId), queryFn: () => studentPortfolioApi.listDisciplinaryRecords(studentId), enabled: !!studentId });

export const useCreateDisciplinaryRecord = (studentId: string) =>
  useCreate(portfolioKeys.disciplinaryRecords(studentId), (data: CreateDisciplinaryRecordRequest) => studentPortfolioApi.createDisciplinaryRecord(studentId, data));

export const useDeleteDisciplinaryRecord = (studentId: string) =>
  useCreate(portfolioKeys.disciplinaryRecords(studentId), (recordId: string) => studentPortfolioApi.deleteDisciplinaryRecord(studentId, recordId));

// Read-only rollup of prefect appointments.
export const usePrefectAppointmentsByStudent = (studentId: string) =>
  useQuery({ queryKey: portfolioKeys.prefectAppointments(studentId), queryFn: () => studentPortfolioApi.listPrefectAppointments(studentId), enabled: !!studentId });
