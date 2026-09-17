import DataGrid from "@/shared/ui/DataGrid";
import { formatDate } from "@/shared/lib/date";

interface Row {
  id: string;
  term_name?: string | null;
  narrative: string;
  created_at?: string | null;
}

export default function ProgressReportsTable({ rows }: { rows: Row[] }) {
  return (
    <DataGrid
      rows={rows}
      getRowId={(r) => r.id}
      pagination={false}
      noHover
      columns={[
        { key: "term", header: "Term", render: (r) => <span className="os-fw-500">{r.term_name || "-"}</span> },
        { key: "narrative", header: "Narrative Remarks", render: (r) => <span className="os-pre-wrap">{r.narrative}</span> },
        { key: "created", header: "Created At", render: (r) => <span className="os-table__mono os-text-xs">{formatDate(r.created_at)}</span> },
      ]}
    />
  );
}
