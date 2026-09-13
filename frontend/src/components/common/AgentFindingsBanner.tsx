import { useState } from "react";
import { ActionableNotification } from "@carbon/react";
import { useNavigate } from "react-router";
import { useMyNotifications } from "../../queries/notifications/useNotifications";
import type { MyNotification } from "../../services/notifications/notification";

interface Props {
  /** Notification titles this page cares about; must match a backend agent's `notifyAdmins(...)` title exactly. */
  titles: string[];
}

// Dismissible nudge surfacing a background agent's latest finding on the page it's actionable from.
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
