import { useState } from "react";
import type { ClassMarkRow } from "@/features/marks/api/termMark";

export interface MarkDraft {
  marks: number;
  isAbsent: boolean;
}

const EMPTY: MarkDraft = { marks: 0, isAbsent: false };

// Local edits for one class / subject / term roster, with unsaved-change tracking.
export function useMarksDraft(rows: ClassMarkRow[] | undefined, rowsKey: string) {
  const [draft, setDraft] = useState<Record<string, MarkDraft>>({});
  const [saved, setSaved] = useState<Record<string, MarkDraft>>({});
  const [syncedFor, setSyncedFor] = useState("");

  // Re-seed when the roster changes; done during render so the first paint is already correct.
  if (rows && syncedFor !== rowsKey) {
    const seeded = Object.fromEntries(rows.map((r) => [r.student_id, { marks: r.marks ?? 0, isAbsent: r.is_absent ?? false }]));
    setDraft(seeded);
    setSaved(seeded);
    setSyncedFor(rowsKey);
  }

  const get = (studentId: string) => draft[studentId] ?? EMPTY;
  const isUnsaved = (studentId: string) => {
    const d = get(studentId);
    const s = saved[studentId] ?? EMPTY;
    return d.marks !== s.marks || d.isAbsent !== s.isAbsent;
  };

  return {
    draft,
    get,
    isUnsaved,
    hasUnsaved: Object.keys(draft).some(isUnsaved),
    setMarks: (studentId: string, marks: number) => setDraft((p) => ({ ...p, [studentId]: { ...get(studentId), marks } })),
    setAbsent: (studentId: string, isAbsent: boolean) =>
      setDraft((p) => ({ ...p, [studentId]: { marks: isAbsent ? 0 : get(studentId).marks, isAbsent } })),
    markSaved: () => setSaved(draft),
  };
}
