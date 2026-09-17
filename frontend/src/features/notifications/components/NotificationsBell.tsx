import { NOTIFICATION_PRIORITY_TAG as PRIORITY_TAG } from "@/shared/lib/constants/tags";
import { useState, useRef, useEffect } from "react";
import { Link } from "react-router";
import { HeaderGlobalAction, Tag } from "@carbon/react";
import { Notification as NotificationIcon } from "@carbon/icons-react";
import {
  useMyNotifications,
  useUnreadNotificationCount,
  useMarkNotificationRead,
} from "@/features/notifications/queries/useNotifications";

export default function NotificationsBell() {
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);
  const { data: unreadCount } = useUnreadNotificationCount();
  const { data: notifications } = useMyNotifications({ enabled: open });
  const markRead = useMarkNotificationRead();

  useEffect(() => {
    if (!open) return;
    const handler = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false);
    };
    document.addEventListener("mousedown", handler);
    return () => document.removeEventListener("mousedown", handler);
  }, [open]);

  const preview = (notifications ?? []).slice(0, 8);

  return (
    <div ref={ref} className="os-relative">
      <HeaderGlobalAction aria-label="Notifications" onClick={() => setOpen((o) => !o)} isActive={open}>
        <span className="os-relative os-inline-flex">
          <NotificationIcon size={20} className="os-header-icon" />
          {!!unreadCount && (
            <span className="os-bell__dot os-absolute os-w-px-8 os-h-px-8 os-rounded-full"
            />
          )}
        </span>
      </HeaderGlobalAction>

      {open && (
        <div className="os-bell__popover os-absolute os-right-0 os-w-24 os-flex os-col os-bg-layer os-border os-shadow-sm os-z-top"
        >
          <div className="os-py-3 os-px-4 os-border-b os-fw-600 os-text-md">
            Notifications
          </div>
          <div className="os-overflow-y-auto os-flex-1">
            {preview.length === 0 ? (
              <div className="os-py-6 os-px-4 os-text-center os-c-tertiary os-text-sm">
                No notifications yet
              </div>
            ) : (
              preview.map((n) => (
                <div
                  key={n.recipient_id}
                  onClick={() => !n.is_read && markRead.mutate(n.notification_id)} className={`os-py-3 os-px-4 os-border-layer-hover-b ${n.is_read ? "os-cursor-default" : "os-pointer"} ${n.is_read ? "os-bg-transparent" : "os-bg-accent-light"}`}
                >
                  <div className="os-flex os-items-center os-gap-1h os-mb-1">
                    <span className="os-text-sm os-fw-600 os-c-primary os-flex-1">{n.title}</span>
                    {n.priority !== "normal" && (
                      <Tag type={PRIORITY_TAG[n.priority]} size="sm">
                        {n.priority}
                      </Tag>
                    )}
                  </div>
                  <p className="os-m-0 os-text-sm os-c-secondary os-overflow-hidden os-truncate os-clamp-2">
                    {n.message}
                  </p>
                  <p className="os-mt-1 os-mx-0 os-mb-0 os-text-xs os-c-tertiary">
                    {n.sender_name} &middot; {new Date(n.sent_at).toLocaleString()}
                  </p>
                </div>
              ))
            )}
          </div>
          <Link
            to="/notification-center"
            onClick={() => setOpen(false)} className="os-block os-text-center os-p-2h os-border-t os-text-sm os-fw-500 os-c-accent-dark os-no-underline"
          >
            View all
          </Link>
        </div>
      )}
    </div>
  );
}
