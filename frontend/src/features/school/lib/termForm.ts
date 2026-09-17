import { isDateRangeInvalid } from "@/shared/lib/date";

export interface TermFormValues {
  name: string;
  start_date: string;
  end_date: string;
}

export type TermTouched = Partial<Record<keyof TermFormValues, boolean>>;

export const isTermFormValid = (f: TermFormValues) =>
  f.name.trim().length > 0 && !!f.start_date && !!f.end_date && !isDateRangeInvalid(f.start_date, f.end_date);
