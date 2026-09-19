import type { StaffListParams } from "@/features/staff/api/nonAcademicStaff";

export const staffKeys = {
  all: ["non-academic-staff"] as const,
  list: (params: StaffListParams = {}) => ["non-academic-staff", "list", params] as const,
};
