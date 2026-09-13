import { Time } from "@carbon/icons-react";
import RoleBadge from "./RoleBadge";
import type { PositionRankLabel } from "../../../services/position";
import type { Teacher } from "../../../services/teacher";
import { getInitials } from "../../../lib/name";

export default function WelcomeBanner({
  profile,
  subjectSummary,
  currentYearLabel,
  currentTermName,
  pendingCount,
  rankLabel,
}: {
  profile: Teacher;
  subjectSummary: string;
  currentYearLabel?: string;
  currentTermName?: string;
  pendingCount: number;
  rankLabel: PositionRankLabel;
}) {
  const hour = new Date().getHours();
  const greeting = hour < 12 ? "Good morning" : hour < 17 ? "Good afternoon" : "Good evening";

  return (
    <div
      style={{
        background: "var(--os-layer)",
        border: "1px solid var(--os-border-subtle)",
        borderTop: "3px solid var(--os-accent)",
        padding: "1.25rem 1.5rem",
        marginBottom: "1.5rem",
        display: "flex",
        alignItems: "center",
        gap: "1rem",
        flexWrap: "wrap",
      }}
    >
      <div
        style={{
          width: "2.75rem",
          height: "2.75rem",
          borderRadius: "50%",
          background: "var(--os-accent)",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          color: "var(--os-layer)",
          fontWeight: 700,
          fontSize: "1rem",
          flexShrink: 0,
        }}
      >
        {getInitials(profile.full_name)}
      </div>
      <div style={{ flex: 1 }}>
        <div style={{ display: "flex", alignItems: "center", gap: "0.5rem", marginBottom: "0.15rem", flexWrap: "wrap" }}>
          <p style={{ margin: 0, fontSize: "1.1rem", fontWeight: 500, color: "var(--os-text-primary)" }}>
            {greeting}, {profile.title ? `${profile.title} ` : ""}{profile.full_name}
          </p>
          <RoleBadge rankLabel={rankLabel} />
        </div>
        <p style={{ margin: 0, fontSize: "0.8125rem", color: "var(--os-text-secondary)" }}>
          {subjectSummary} · {profile.employee_number}
          {currentYearLabel ? ` · ${currentYearLabel}` : ""}
          {currentTermName ? ` ${currentTermName}` : ""}
        </p>
      </div>
      {pendingCount > 0 && (
        <div style={{ display: "flex", alignItems: "center", gap: "0.5rem", padding: "0.5rem 0.875rem", background: "var(--os-status-late-bg)", border: "1px solid var(--os-warning)", fontSize: "0.8125rem", color: "var(--os-warning-text)" }}>
          <Time size={14} style={{ fill: "var(--os-warning)" }} />
          {pendingCount} session{pendingCount > 1 ? "s" : ""} pending today
        </div>
      )}
    </div>
  );
}
