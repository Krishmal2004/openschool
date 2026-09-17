import { NON_ACADEMIC_DESIGNATIONS } from "@/features/staff/api/nonAcademicStaff";

export { EMPLOYMENT_STATUSES } from "@/shared/lib/constants/people";

export function designationLabel(value: string) {
  return NON_ACADEMIC_DESIGNATIONS.find((d) => d.value === value)?.label ?? value;
}
