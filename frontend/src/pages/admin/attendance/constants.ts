import type { AttendanceRecordRow } from "../../../services/attendance";

export type Status = "present" | "absent" | "late" | "excused" | null;

export const STATUS_STYLES: Record<
  NonNullable<Status>,
  { bg: string; border: string; color: string; label: string }
> = {
  present: { bg: "var(--os-status-present-bg)", border: "var(--os-status-present-border)", color: "var(--os-status-present-text)", label: "Present" },
  absent: { bg: "var(--os-status-absent-bg)", border: "var(--os-status-absent-border)", color: "var(--os-status-absent-text)", label: "Absent" },
  late: { bg: "var(--os-status-late-bg)", border: "var(--os-status-late-border)", color: "var(--os-status-late-text)", label: "Late" },
  excused: { bg: "var(--os-status-excused-bg)", border: "var(--os-status-excused-border)", color: "var(--os-status-excused-text)", label: "Excused" },
};

export function recordsToState(records: AttendanceRecordRow[]) {
  const statuses: Record<string, Status> = {};
  const notes: Record<string, string> = {};
  for (const r of records) {
    if (r.status === "present" || r.status === "absent" || r.status === "late" || r.status === "excused") {
      statuses[r.student_id] = r.status;
    }
    if (r.note) notes[r.student_id] = r.note;
  }
  return { statuses, notes };
}
