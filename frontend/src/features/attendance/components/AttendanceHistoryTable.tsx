import { Tag } from "@carbon/react";
import DataGrid from "@/shared/ui/DataGrid";
import { ATTENDANCE_STATUS_TAG } from "@/shared/lib/attendanceStatus";
import { capitalize } from "@/shared/lib/text";

interface Row {
  id: string;
  session_date: string;
  class_name: string;
  status: string;
  note?: string | null;
}

// Read-only attendance history, newest first; shared by the student and parent portals.
export default function AttendanceHistoryTable({ rows }: { rows: Row[] }) {
  const sorted = [...rows].sort((a, b) => b.session_date.localeCompare(a.session_date));
  return (
    <DataGrid
      rows={sorted}
      getRowId={(r) => r.id}
      noHover
      columns={[
        { key: "date", header: "Date", render: (r) => <span className="os-table__mono">{r.session_date}</span> },
        { key: "class", header: "Class", render: (r) => r.class_name },
        { key: "status", header: "Status", render: (r) => <Tag type={ATTENDANCE_STATUS_TAG[r.status] ?? "gray"} size="sm">{capitalize(r.status)}</Tag> },
        { key: "note", header: "Note", render: (r) => <span className="os-table__muted">{r.note || "-"}</span> },
      ]}
    />
  );
}
