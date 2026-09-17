import DataGrid from "@/shared/ui/DataGrid";
import { capitalize } from "@/shared/lib/text";

interface Row {
  id: string;
  full_name: string;
  relationship?: string | null;
  phone?: string | null;
  is_primary_contact: boolean;
}

// Read-only guardian list for the student and parent portals.
export default function GuardiansTable({ rows }: { rows: Row[] }) {
  return (
    <DataGrid
      rows={rows}
      getRowId={(g) => g.id}
      pagination={false}
      noHover
      columns={[
        { key: "name", header: "Name", render: (g) => <span className="os-fw-500">{g.full_name}</span> },
        { key: "rel", header: "Relationship", render: (g) => capitalize(g.relationship, "Guardian") },
        { key: "phone", header: "Phone", render: (g) => <span className="os-table__mono">{g.phone || "-"}</span> },
        { key: "primary", header: "Primary Contact", render: (g) => (g.is_primary_contact ? <span className="os-c-success os-fw-600">Yes</span> : "No") },
      ]}
    />
  );
}
