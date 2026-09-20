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

type DateInput = string | Date | null | undefined;

const DATE_ONLY = /^(\d{4})-(\d{2})-(\d{2})$/;

function toDate(input: DateInput): Date | null {
  if (!input) return null;
  if (input instanceof Date) return Number.isNaN(input.getTime()) ? null : input;
  const dateOnly = DATE_ONLY.exec(input);
  if (dateOnly) {
    const [, y, m, d] = dateOnly;
    return new Date(Number(y), Number(m) - 1, Number(d));
  }
  const parsed = new Date(input);
  return Number.isNaN(parsed.getTime()) ? null : parsed;
}

export function formatDate(input: DateInput): string {
  const d = toDate(input);
  return d ? d.toLocaleDateString("en-LK") : "-";
}

export function formatDateTime(input: DateInput): string {
  const d = toDate(input);
  return d ? d.toLocaleString("en-LK") : "-";
}

export function formatMonth(input: DateInput): string {
  const d = toDate(input);
  return d ? d.toLocaleDateString("en-LK", { month: "short", year: "numeric" }) : "-";
}

export function formatLongDate(input: DateInput): string {
  const d = toDate(input);
  return d ? d.toLocaleDateString("en-LK", { weekday: "long", year: "numeric", month: "long", day: "numeric" }) : "-";
}

export function formatShortDayMonth(input: DateInput): string {
  const d = toDate(input);
  return d ? d.toLocaleDateString("en-LK", { month: "short", day: "numeric" }) : "-";
}

export function formatDayMonthYear(input: DateInput): string {
  const d = toDate(input);
  return d ? d.toLocaleDateString("en-LK", { month: "short", day: "numeric", year: "numeric" }) : "-";
}
