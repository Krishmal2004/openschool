export const staffKeys = {
  all: ["non-academic-staff"] as const,
  list: (search: string, designation: string) => ["non-academic-staff", "list", { search, designation }] as const,
};
