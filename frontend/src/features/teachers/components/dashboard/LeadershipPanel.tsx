import { UserRole } from "@carbon/icons-react";
import type { PositionSummary } from "@/features/positions/api/position";

function reachText(summary: PositionSummary): string {
  if (summary.rank_label === "Principal") return "You can notify the whole school.";
  if (summary.rank_label === "Vice Principal") {
    return summary.notify_whole_school
      ? "You've been granted whole-school notification reach."
      : "You can notify your assigned grade(s).";
  }
  return "You can notify your assigned grade and its staff.";
}

// Only for Principal, Vice Principal and Section Head; a class teacher's reach is already obvious.
export default function LeadershipPanel({ summary }: { summary: PositionSummary }) {
  return (
    <div className="os-section">
      <div className="os-section__header">
        <h2 className="os-section__title os-flex os-items-center os-gap-2">
          <UserRole size={16} className="os-fill-accent" /> Leadership
        </h2>
      </div>
      <div className="os-py-4 os-px-6">
        <p className="os-m-0 os-text-sm os-c-secondary os-lh-normal">{reachText(summary)}</p>
      </div>
    </div>
  );
}
