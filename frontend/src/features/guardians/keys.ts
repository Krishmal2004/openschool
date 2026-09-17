export const guardianKeys = {
  all: ["guardians"] as const,
  list: (orphansOnly: boolean) => ["guardians", "list", { orphansOnly }] as const,
  search: (search: string, orphansOnly: boolean) => ["guardians", "search", search, { orphansOnly }] as const,
  byStudent: (studentId: string) => ["guardians", "by-student", studentId] as const,
  students: (guardianId: string) => ["guardians", "detail", guardianId, "students"] as const,
  notifications: (guardianId: string) => ["guardians", "detail", guardianId, "notifications"] as const,
};
