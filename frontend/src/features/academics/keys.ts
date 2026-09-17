export const classKeys = {
  all: ["classes"] as const,
  current: () => ["classes", "current"] as const,
  byYear: (academicYearId: string) => ["classes", "by-year", academicYearId] as const,
  detail: (id: string) => ["classes", "detail", id] as const,
  subjectTeachers: (id: string) => ["classes", "detail", id, "subject-teachers"] as const,
};

export const gradeKeys = {
  all: ["grades"] as const,
};

export const streamKeys = {
  all: ["streams"] as const,
  groups: (streamId: string) => ["streams", streamId, "groups"] as const,
};

export const promotionKeys = {
  all: ["promotion"] as const,
  preview: (sourceYearId: string, targetYearId: string, rankByTermId?: string) =>
    ["promotion", "preview", sourceYearId, targetYearId, rankByTermId ?? null] as const,
};
