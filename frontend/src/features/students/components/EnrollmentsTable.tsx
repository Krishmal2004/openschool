import DataGrid from "@/shared/ui/DataGrid";
import { capitalize } from "@/shared/lib/text";

interface Row {
  subject_id: string;
  subject_name: string;
  subject_code: string;
  subject_type?: string | null;
}

export default function EnrollmentsTable({ rows }: { rows: Row[] }) {
  return (
    <DataGrid
      rows={rows}
      getRowId={(e) => e.subject_id}
      pagination={false}
      noHover
      columns={[
        { key: "name", header: "Subject Name", render: (e) => <span className="os-fw-500">{e.subject_name}</span> },
        { key: "code", header: "Subject Code", render: (e) => <span className="os-table__mono">{e.subject_code}</span> },
        { key: "type", header: "Type", render: (e) => capitalize(e.subject_type, "Core") },
      ]}
    />
  );
}
