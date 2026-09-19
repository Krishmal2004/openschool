import type { TeacherListParams } from "@/features/teachers/api/teacher";

export const teacherKeys = {
  all: ["teachers"] as const,
  list: (params: TeacherListParams = {}) => ["teachers", "list", params] as const,
  detail: (id: string) => ["teachers", "detail", id] as const,
  subjects: (id: string) => ["teachers", "detail", id, "subjects"] as const,
  workload: (id: string) => ["teachers", "detail", id, "workload"] as const,
  bySubject: (subjectId: string) => ["teachers", "by-subject", subjectId] as const,
  me: () => ["teachers", "me"] as const,
};

export const sectionHeadKeys = {
  all: ["section-heads"] as const,
  byYear: (academicYearId: string) => ["section-heads", academicYearId] as const,
};
