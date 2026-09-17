import type { Classroom } from "@/features/timetable/api/classroom";

// Homerooms are usually named after the class, so suggest the matching room; the admin can override.
export function suggestHomeClassroom(classrooms: Classroom[] | undefined, className: string): Classroom | undefined {
  const trimmed = className.trim().toLowerCase();
  if (!trimmed) return undefined;
  return classrooms?.find((c) => c.room_type === "regular" && c.name.trim().toLowerCase() === trimmed);
}
