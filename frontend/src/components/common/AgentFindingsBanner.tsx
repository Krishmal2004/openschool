import { useState } from "react";
import { ActionableNotification } from "@carbon/react";
import { useNavigate } from "react-router";
import { useMyNotifications } from "../../queries/notifications/useNotifications";
import type { MyNotification } from "../../services/notifications/notification";

interface Props {
  /**
   * Exact notification titles this page cares about — must match the
   * literal title string a backend agent's `notifyAdmins(...)` call uses
   * (see backend/internal/jobs/agent_*.go). Titles, not job/agent names:
   * five backend agents each run several concurrent checks and would
   * otherwise all report through one shared job_runs summary, which is too
   * coarse to show a page-relevant finding. Each check's notification
   * title stays stable and page-specific even though the checks were
   * consolidated, so matching on it here preserves the original "only show
   * this finding on the page it's actionable from" behavior.
   */
  titles: string[];
}

// A lightweight, dismissible nudge surfacing a background agent's most
// recent finding directly on the page it's actionable from (e.g. the
// no-guardian finding on the Students list) instead of only in the
// Notification Center. Built on the same admin notifications every agent
// check already sends — no separate backend surface — and clears itself
// once the admin reads the notification anywhere (bell icon, Notification
// Center, or dismissing it here).
export default function AgentFindingsBanner({ titles }: Props) {
  const { data: notifications } = useMyNotifications();
  const navigate = useNavigate();
  const [dismissed, setDismissed] = useState<Set<string>>(new Set());

  const latestByTitle = new Map<string, MyNotification>();
  for (const n of notifications ?? []) {
    if (!titles.includes(n.title) || n.is_read || n.is_archived || dismissed.has(n.notification_id)) continue;
    const existing = latestByTitle.get(n.title);
    if (!existing || new Date(n.sent_at) > new Date(existing.sent_at)) {
      latestByTitle.set(n.title, n);
    }
  }
  const flagged = Array.from(latestByTitle.values());

  if (flagged.length === 0) return null;

  return (
    <div style={{ display: "flex", flexDirection: "column", gap: "0.75rem", marginBottom: "1.5rem" }}>
      {flagged.map((n) => (
        <ActionableNotification
          key={n.notification_id}
          inline
          kind="warning"
          lowContrast
          hideCloseButton={false}
          title={n.title}
          subtitle={n.message}
          actionButtonLabel="View in Automation"
          onActionButtonClick={() => navigate("/automation")}
          onClose={() => setDismissed((prev) => new Set(prev).add(n.notification_id))}
          style={{ maxWidth: "100%" }}
        />
      ))}
    </div>
  );
}
