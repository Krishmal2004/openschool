import DataGrid from "@/shared/ui/DataGrid";

interface Row {
  id: string;
  subject_name: string;
  subject_code: string;
  teacher_name?: string | null;
  is_absent: boolean;
  marks: number | null;
  max_marks: number;
}

// Read-only term marks for one student; shared by the student and parent portals.
export default function TermMarksTable({ rows }: { rows: Row[] }) {
  return (
    <DataGrid
      rows={rows}
      getRowId={(m) => m.id}
      pagination={false}
      noHover
      columns={[
        { key: "subject", header: "Subject", render: (m) => <>{m.subject_name} <span className="os-table__muted">({m.subject_code})</span></> },
        { key: "teacher", header: "Teacher", render: (m) => <span className="os-table__muted">{m.teacher_name || "-"}</span> },
        { key: "marks", header: "Marks", align: "end", render: (m) => <span className="os-fw-600">{m.is_absent ? "AB" : `${m.marks} / ${m.max_marks}`}</span> },
      ]}
    />
  );
}
