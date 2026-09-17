import { useMemo, useState } from "react";
import type { AttendanceRecordRow } from "@/features/attendance/api/attendance";
import type { Student } from "@/features/students/api/student";
import { recordsToState, type Status } from "@/features/attendance/constants";

// Local marking state for one session: per-student status and note, plus the running summary.
export function useAttendanceMarking(sessionId: string, records: AttendanceRecordRow[] | undefined, students: Student[] | undefined) {
  const [statuses, setStatuses] = useState<Record<string, Status>>({});
  const [notes, setNotes] = useState<Record<string, string>>({});
  const [loadedFor, setLoadedFor] = useState<string | null>(null);

  if (records && loadedFor !== sessionId) {
    const seeded = recordsToState(records);
    setStatuses(seeded.statuses);
    setNotes(seeded.notes);
    setLoadedFor(sessionId);
  }

  const summary = useMemo(() => {
    const values = Object.values(statuses);
    const count = (s: Status) => values.filter((v) => v === s).length;
    return {
      present: count("present"),
      absent: count("absent"),
      late: count("late"),
      excused: count("excused"),
      unmarked: (students ?? []).length - values.filter(Boolean).length,
    };
  }, [statuses, students]);

  return {
    statuses,
    notes,
    summary,
    // Clicking the active status again clears it.
    mark: (studentId: string, status: NonNullable<Status>) =>
      setStatuses((prev) => ({ ...prev, [studentId]: prev[studentId] === status ? null : status })),
    markAll: (status: NonNullable<Status>) => setStatuses(Object.fromEntries((students ?? []).map((s) => [s.id, status]))),
    clearAll: () => setStatuses({}),
    setNote: (studentId: string, value: string) => setNotes((prev) => ({ ...prev, [studentId]: value })),
    toRecords: () =>
      Object.entries(statuses)
        .filter((entry): entry is [string, NonNullable<Status>] => !!entry[1])
        .map(([student_id, status]) => ({ student_id, status, note: notes[student_id]?.trim() || undefined })),
  };
}
