import type { Classroom } from "../services/timetable/classroom";

// Sri Lankan schools usually name a class's homeroom the same as the class itself
// (e.g. class "13-M1" sits in room "13-M1") — suggest that match automatically,
// but the caller always lets the admin override it.
export function suggestHomeClassroom(classrooms: Classroom[] | undefined, className: string): Classroom | undefined {
  const trimmed = className.trim().toLowerCase();
  if (!trimmed) return undefined;
  return classrooms?.find((c) => c.room_type === "regular" && c.name.trim().toLowerCase() === trimmed);
}
