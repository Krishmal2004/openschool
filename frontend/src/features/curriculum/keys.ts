export const curriculumKeys = {
  all: ["curriculum"] as const,
  mediums: () => ["curriculum", "mediums"] as const,
  levels: () => ["curriculum", "levels"] as const,
  levelTree: (id: string) => ["curriculum", "levels", id, "tree"] as const,
  groupSubjects: (groupId: string) => ["curriculum", "groups", groupId, "subjects"] as const,
  subjects: () => ["curriculum", "subjects"] as const,
};
