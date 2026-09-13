import { Location } from "@carbon/icons-react";
import type { LeadershipOverviewSummary } from "../../../services/position";
import InfoRow from "../../../components/common/InfoRow";

// The §9.3 "deepened" panel — real scoped counts (not just label text) for
// Principal/Vice Principal/Section Head. Sits alongside LeadershipPanel
// (notification reach), doesn't replace it.
export default function LeadershipOverviewPanel({ overview }: { overview: LeadershipOverviewSummary }) {
  const scopeLabel = overview.scope === "school" ? "Whole school" : overview.grade_names.join(", ") || "Your grades";

  const rows = [
    { label: "Classes", value: overview.class_count, color: "var(--os-text-primary)" },
    { label: "Students", value: overview.student_count, color: "var(--os-accent)" },
    { label: "Sessions Marked Today", value: overview.sessions_marked_today, color: "var(--os-success)" },
    {
      label: "Sessions Pending Today",
      value: overview.sessions_pending_today,
      color: overview.sessions_pending_today > 0 ? "var(--os-warning)" : "var(--os-text-tertiary)",
    },
  ];

  return (
    <div className="os-section">
      <div className="os-section__header">
        <h2 className="os-section__title" style={{ display: "flex", alignItems: "center", gap: "0.5rem" }}>
          <Location size={16} style={{ fill: "var(--os-accent)" }} /> Overview — {scopeLabel}
        </h2>
      </div>
      <div className="os-section__body" style={{ padding: "0.75rem 1.5rem" }}>
        {rows.map(({ label, value, color }, i) => (
          <InfoRow key={label} label={label} value={<span style={{ color }}>{value}</span>} bold divider={i < rows.length - 1} />
        ))}
      </div>
    </div>
  );
}
