import { useQuery } from "@tanstack/react-query";
import { dashboardAnalyticsApi } from "@/features/reports/api/dashboardAnalytics";
import { reportKeys } from "@/features/reports/keys";

export const useDashboardAnalytics = (enabled = true) =>
  useQuery({ queryKey: reportKeys.dashboardAnalytics(), queryFn: () => dashboardAnalyticsApi.get(), enabled });

// Same data through the Principal / Vice Principal scoped endpoint.
export const useLeadershipAnalytics = (enabled = true) =>
  useQuery({ queryKey: reportKeys.leadershipAnalytics(), queryFn: () => dashboardAnalyticsApi.getForLeadership(), enabled });
