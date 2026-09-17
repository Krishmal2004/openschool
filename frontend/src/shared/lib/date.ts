// Local YYYY-MM-DD; accepts undefined because Carbon DatePicker can pass an empty selection.
export function toYmd(d: Date | undefined): string {
  if (!d) return "";
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}-${String(d.getDate()).padStart(2, "0")}`;
}

// Call at use time, never cache at module load, or "today" freezes when midnight passes in an open tab.
export function todayISODate(): string {
  return toYmd(new Date());
}

// True once both dates are filled in and end isn't strictly after start.
export function isDateRangeInvalid(startDate: string, endDate: string): boolean {
  return !!startDate && !!endDate && endDate <= startDate;
}

// Attendance sessions lock 24 hours after creation, teachers lose edit access, admins keep an override path.
export function isLockedAfter24Hours(createdAt: string | null | undefined): boolean {
  return !!createdAt && Date.now() - new Date(createdAt).getTime() > 24 * 60 * 60 * 1000;
}

// Local short date for display; empty input shows a dash.
export function formatDate(iso: string | null | undefined): string {
  return iso ? new Date(iso).toLocaleDateString() : "-";
}
