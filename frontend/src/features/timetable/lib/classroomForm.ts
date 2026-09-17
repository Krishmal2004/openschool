import type { ClassroomType } from "@/features/timetable/api/classroom";

export interface ClassroomForm {
  name: string;
  code: string;
  capacity: string;
  room_type: ClassroomType;
  subject_id: string;
}

export const EMPTY_CLASSROOM_FORM: ClassroomForm = { name: "", code: "", capacity: "", room_type: "regular", subject_id: "" };
export const isLabMissingSubject = (f: ClassroomForm) => f.room_type === "lab" && !f.subject_id;
