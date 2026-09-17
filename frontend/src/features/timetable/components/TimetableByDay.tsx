import { WEEKDAYS } from "@/shared/lib/timetable";
import DataGrid, { type GridColumn } from "@/shared/ui/DataGrid";

interface Entry {
  day_of_week: number;
  period_number: number;
  subject_name?: string | null;
  classroom_name?: string | null;
}

interface Props<T extends Entry> {
  entries: T[];
  // The column between Subject and Classroom differs per portal (teacher vs class).
  middle: GridColumn<T>;
  getRowId: (e: T, index: number) => string;
  // Wrap each day in its own section card instead of a heading.
  asSections?: boolean;
}

// One read-only table per weekday, used by the student, teacher and parent timetable views.
export default function TimetableByDay<T extends Entry>({ entries, middle, getRowId, asSections = false }: Props<T>) {
  const columns: GridColumn<T>[] = [
    { key: "period", header: "Period", render: (e) => `P${e.period_number}` },
    { key: "subject", header: "Subject", render: (e) => e.subject_name ?? <span className="os-table__muted">-</span> },
    middle,
    { key: "room", header: "Classroom", render: (e) => e.classroom_name ?? <span className="os-table__muted">-</span> },
  ];
  return (
    <>
      {WEEKDAYS.map((day) => {
        const rows = entries.filter((e) => e.day_of_week === day.value).sort((a, b) => a.period_number - b.period_number);
        if (rows.length === 0) return null;
        const grid = <DataGrid rows={rows} columns={columns} getRowId={(e) => getRowId(e, rows.indexOf(e))} pagination={false} noHover />;
        return asSections ? (
          <div key={day.value} className="os-section">
            <div className="os-section__header">
              <h2 className="os-section__title">{day.label}</h2>
            </div>
            {grid}
          </div>
        ) : (
          <div key={day.value}>
            <h3 className="os-text-md os-fw-600 os-mt-0 os-mx-0 os-mb-2">{day.label}</h3>
            {grid}
          </div>
        );
      })}
    </>
  );
}
