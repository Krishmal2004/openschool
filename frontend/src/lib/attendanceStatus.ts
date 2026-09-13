// Carbon Tag `type` per attendance status, for the read-only history views
// (student's own attendance, a parent's view of a child's attendance).
export const ATTENDANCE_STATUS_TAG: Record<string, "green" | "red" | "warm-gray" | "blue"> = {
  present: "green",
  absent: "red",
  late: "warm-gray",
  excused: "blue",
};
