export const portfolioKeys = {
  all: ["portfolio"] as const,
  progressReports: (studentId: string) => ["portfolio", "students", studentId, "progress-reports"] as const,
  activities: (studentId: string) => ["portfolio", "students", studentId, "activities"] as const,
  leadershipRoles: (studentId: string) => ["portfolio", "students", studentId, "leadership-roles"] as const,
  awards: (studentId: string) => ["portfolio", "students", studentId, "awards"] as const,
  disciplinaryRecords: (studentId: string) => ["portfolio", "students", studentId, "disciplinary-records"] as const,
  prefectAppointments: (studentId: string) => ["portfolio", "students", studentId, "prefect-appointments"] as const,
  societyMemberships: (studentId: string) => ["portfolio", "students", studentId, "society-memberships"] as const,
};

export const prefectKeys = {
  all: ["prefects"] as const,
  byYear: (academicYearId: string) => ["prefects", "by-year", academicYearId] as const,
  years: () => ["prefects", "years"] as const,
};

export const societyKeys = {
  all: ["societies"] as const,
  byYear: (academicYearId: string) => ["societies", "by-year", academicYearId] as const,
  years: () => ["societies", "years"] as const,
  members: (societyId: string) => ["societies", "detail", societyId, "members"] as const,
  mine: () => ["societies", "me"] as const,
};
