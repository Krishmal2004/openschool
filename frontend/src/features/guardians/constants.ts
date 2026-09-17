import type { GuardianRelationship } from "@/features/guardians/api/guardian";

export const GUARDIAN_RELATIONSHIPS: { value: GuardianRelationship; label: string }[] = [
  { value: "father", label: "Father" },
  { value: "mother", label: "Mother" },
  { value: "guardian", label: "Guardian" },
  { value: "other", label: "Other" },
];

export function relationshipLabel(value: string) {
  return GUARDIAN_RELATIONSHIPS.find((r) => r.value === value)?.label ?? value;
}
