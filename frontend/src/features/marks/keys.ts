export const markKeys = {
  all: ["marks"] as const,
  byClass: (classId: string, termId: string, subjectId: string) => ["marks", "by-class", classId, termId, subjectId] as const,
  byStudent: (studentId: string, termId: string) => ["marks", "by-student", studentId, termId] as const,
};
