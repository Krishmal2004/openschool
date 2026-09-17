import type { RecipientRule } from "@/features/notifications/api/notification";

// Rules have no id; the id field each type sets gives a stable key so removing a chip does not reuse DOM nodes.
export function ruleKey(rule: RecipientRule): string {
  return [
    rule.type,
    rule.grade_id,
    rule.class_id,
    rule.grade_section_id,
    rule.subject_id,
    rule.subject_audience,
    rule.student_id,
    rule.guardian_id,
    rule.teacher_id,
  ]
    .filter(Boolean)
    .join(":");
}
