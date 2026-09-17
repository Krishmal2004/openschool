export const reportKeys = {
  all: ["reports"] as const,
  dashboardAnalytics: () => ["reports", "dashboard-analytics"] as const,
  leadershipAnalytics: () => ["reports", "leadership-analytics"] as const,
};
