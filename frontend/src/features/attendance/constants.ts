import type { AttendanceRecordRow } from "@/features/attendance/api/attendance";

export type Status = "present" | "absent" | "late" | "excused" | null;

export { STUDENT_ATTENDANCE_STYLES as STATUS_STYLES } from "@/shared/lib/constants/attendance";

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
