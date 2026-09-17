import { useState } from "react";
import { Tag, SkeletonText } from "@carbon/react";
import { useNotificationStats } from "@/features/notifications/queries/useNotifications";
import type { Notification } from "@/features/notifications/api/notification";

export default function SentHistoryRow({ notification }: { notification: Notification }) {
  const [expanded, setExpanded] = useState(false);
  const { data: stats } = useNotificationStats(expanded ? notification.id : "");

  return (
    <div
      role="button"
      tabIndex={0}
      aria-expanded={expanded}
      className="os-list-row os-col os-items-stretch os-py-3h os-px-6 os-pointer"
      onClick={() => setExpanded((e) => !e)}
      onKeyDown={(e) => {
        if (e.key === "Enter" || e.key === " ") {
          e.preventDefault();
          setExpanded((v) => !v);
        }
      }}
    >
      <div className="os-flex os-items-center os-gap-2 os-mb-1">
        <Tag type="blue" size="sm">
          {notification.category}
        </Tag>
        <span className="os-fw-500 os-text-sm">{notification.title}</span>
        <span className="os-ml-auto os-text-xs os-c-tertiary">
          {notification.sent_at ? new Date(notification.sent_at).toLocaleString() : ""}
        </span>
      </div>
      <p className="os-m-0 os-text-xs os-c-secondary os-lh-normal os-clamp-2 os-overflow-hidden"
      >
        {notification.message}
      </p>
      {expanded && (
        <div className="os-mt-2 os-text-xs os-c-primary">
          {stats ? (
            <span>
              Recipients: <strong>{stats.total}</strong> &middot; Read: <strong>{stats.read}</strong> &middot; Unread: <strong>{stats.unread}</strong>
            </span>
          ) : (
            <SkeletonText width="40%" />
          )}
        </div>
      )}
    </div>
  );
}
