import { Tag } from "@carbon/react";
import type { PrefectAppointment, StudentActivity, LeadershipRole, StudentAward, DisciplinaryRecord } from "@/features/portfolio/api/studentPortfolio";
import type { SocietyMembership } from "@/features/portfolio/api/society";
import DataGrid from "@/shared/ui/DataGrid";
import { SEVERITY_TAG } from "@/shared/lib/constants/tags";
import { capitalize } from "@/shared/lib/text";
import { formatDate } from "@/shared/lib/date";

// Read-only portfolio tables shared by the admin, student and parent portals.
const grid = { pagination: false, noHover: true } as const;

export const PrefectsTable = ({ rows }: { rows: PrefectAppointment[] }) => (
  <DataGrid {...grid} rows={rows} getRowId={(p) => p.id} columns={[
    { key: "year", header: "Year", render: (p) => <span className="os-fw-500">{p.academic_year_label}</span> },
    { key: "rank", header: "Appointment / Rank", render: (p) => p.rank },
  ]} />
);

export const SocietiesTable = ({ rows }: { rows: SocietyMembership[] }) => (
  <DataGrid {...grid} rows={rows} getRowId={(s) => s.id} columns={[
    { key: "society", header: "Society", render: (s) => <span className="os-fw-500">{s.society_name}</span> },
    { key: "role", header: "Role", render: (s) => capitalize(s.role, "Member") },
    { key: "joined", header: "Joined Date", render: (s) => <span className="os-table__mono">{formatDate(s.created_at)}</span> },
  ]} />
);

export const ActivitiesTable = ({ rows }: { rows: StudentActivity[] }) => (
  <DataGrid {...grid} rows={rows} getRowId={(a) => a.id} columns={[
    { key: "name", header: "Activity", render: (a) => <span className="os-fw-500">{a.name}</span> },
    { key: "category", header: "Category", render: (a) => capitalize(a.category) },
    { key: "role", header: "Role", render: (a) => a.role || "Participant" },
    { key: "achievement", header: "Achievement", render: (a) => a.achievement || "-" },
  ]} />
);

export const LeadershipTable = ({ rows, className }: { rows: LeadershipRole[]; className?: string }) => (
  <DataGrid {...grid} className={className} rows={rows} getRowId={(l) => l.id} columns={[
    { key: "title", header: "Leadership Role", render: (l) => <span className="os-fw-500">{l.title}</span> },
    { key: "scope", header: "Scope", render: (l) => l.scope || "School" },
    { key: "assigned", header: "Assigned", render: (l) => <span className="os-table__mono">{formatDate(l.created_at)}</span> },
  ]} />
);

export const AwardsTable = ({ rows }: { rows: StudentAward[] }) => (
  <DataGrid {...grid} rows={rows} getRowId={(a) => a.id} columns={[
    { key: "title", header: "Award Recognition", render: (a) => <span className="os-fw-500">{a.title}</span> },
    { key: "category", header: "Category", render: (a) => a.category || "-" },
    { key: "date", header: "Date", render: (a) => <span className="os-table__mono">{formatDate(a.awarded_date)}</span> },
    { key: "description", header: "Description", render: (a) => a.description || "-" },
  ]} />
);

export const DisciplineTable = ({ rows }: { rows: DisciplinaryRecord[] }) => (
  <DataGrid {...grid} rows={rows} getRowId={(d) => d.id} columns={[
    { key: "date", header: "Date", render: (d) => <span className="os-table__mono">{formatDate(d.incident_date)}</span> },
    { key: "incident", header: "Incident", render: (d) => <span className="os-fw-500">{d.description}</span> },
    { key: "severity", header: "Severity", render: (d) => <Tag type={SEVERITY_TAG[d.severity] ?? "gray"} size="sm">{capitalize(d.severity, "Minor")}</Tag> },
    { key: "action", header: "Action Taken", render: (d) => d.action_taken || "-" },
  ]} />
);
