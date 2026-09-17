import { NOTIFICATION_PRIORITY_TAG as PRIORITY_TAG } from "@/shared/lib/constants/tags";
import { useMemo, useState } from "react";
import { Tag, Dropdown, Button } from "@carbon/react";
import { Archive, ArrowUpRight, Search } from "@carbon/icons-react";
import {
  useMyNotifications,
  useMyArchivedNotifications,
  useMarkNotificationRead,
  useArchiveNotification,
  useUnarchiveNotification,
} from "@/features/notifications/queries/useNotifications";
import { CATEGORIES } from "@/features/notifications/api/notification";
import type { MyNotification, NotificationCategory } from "@/features/notifications/api/notification";
import EmptyState from "@/shared/ui/EmptyState";
import LoadingSpinner from "@/shared/ui/LoadingSpinner";

type Tab = "unread" | "read" | "archived";

const CATEGORY_LABEL: Record<string, string> = Object.fromEntries(CATEGORIES.map((c) => [c.value, c.label]));

function NotificationRow({ n }: { n: MyNotification }) {
  const markRead = useMarkNotificationRead();
  const archive = useArchiveNotification();
  const unarchive = useUnarchiveNotification();

  return (
    <div
      className={`${`os-list-row${!n.is_read ? " is-selected" : ""}`} os-gap-4 os-py-3h os-px-6`}
    >
      <div className="os-flex-1 os-min-w-0">
        <div className="os-flex os-items-center os-gap-2 os-mb-1">
          <span className="os-fw-600 os-text-md os-c-primary">{n.title}</span>
          {n.priority !== "normal" && (
            <Tag type={PRIORITY_TAG[n.priority]} size="sm">
              {n.priority}
            </Tag>
          )}
          <Tag type="blue" size="sm">
            {CATEGORY_LABEL[n.category] ?? n.category}
          </Tag>
          {!n.is_read && !n.is_archived && (
            <span className="os-w-px-8 os-h-px-8 os-rounded-full os-bg-accent os-inline-block" />
          )}
        </div>
        <p className="os-mt-0 os-mx-0 os-mb-1h os-text-sm os-c-secondary os-lh-normal">{n.message}</p>
        <p className="os-m-0 os-text-xs os-c-tertiary">
          {n.sender_name} &middot; {new Date(n.sent_at).toLocaleString()}
        </p>
      </div>
      <div className="os-flex os-col os-gap-1h os-items-end os-shrink-0">
        {!n.is_read && (
          <Button kind="ghost" size="sm" onClick={() => markRead.mutate(n.notification_id)}>
            Mark read
          </Button>
        )}
        {n.is_archived ? (
          <Button kind="ghost" size="sm" renderIcon={ArrowUpRight} onClick={() => unarchive.mutate(n.notification_id)}>
            Unarchive
          </Button>
        ) : (
          <Button kind="ghost" size="sm" renderIcon={Archive} onClick={() => archive.mutate(n.notification_id)}>
            Archive
          </Button>
        )}
      </div>
    </div>
  );
}

export default function NotificationCenter() {
  const { data: inbox, isLoading: inboxLoading } = useMyNotifications();
  const { data: archived, isLoading: archivedLoading } = useMyArchivedNotifications();

  const [tab, setTab] = useState<Tab>("unread");
  const [query, setQuery] = useState("");
  const [category, setCategory] = useState<NotificationCategory | "">("");

  const isLoading = tab === "archived" ? archivedLoading : inboxLoading;

  const filtered = useMemo(() => {
    let list = tab === "archived" ? (archived ?? []) : (inbox ?? []);
    if (tab === "unread") list = list.filter((n) => !n.is_read);
    if (tab === "read") list = list.filter((n) => n.is_read);
    if (category) list = list.filter((n) => n.category === category);
    const q = query.trim().toLowerCase();
    if (q) {
      list = list.filter((n) => n.title.toLowerCase().includes(q) || n.message.toLowerCase().includes(q));
    }
    return [...list].sort((a, b) => b.sent_at.localeCompare(a.sent_at));
  }, [inbox, archived, tab, category, query]);

  const unreadCount = (inbox ?? []).filter((n) => !n.is_read).length;
  const readCount = (inbox ?? []).filter((n) => n.is_read).length;

  return (
    <div className="os-page">
      <div className="os-page__header">
        <div className="os-page__header-left">
          <h1 className="os-page__title">Notification Center</h1>
          <p className="os-page__subtitle">Announcements and updates sent to you.</p>
        </div>
      </div>

      <div className="os-section">
        <div className="os-flex os-gap-2 os-pt-4 os-px-6 os-pb-0">
          {(["unread", "read", "archived"] as Tab[]).map((t) => (
            <button
              key={t}
              onClick={() => setTab(t)} className={`os-pill os-px-4${tab === t ? " is-active" : ""}`}
            >
              {t === "unread" ? `Unread (${unreadCount})` : t === "read" ? `Read (${readCount})` : "Archived"}
            </button>
          ))}
        </div>

        <div className="os-toolbar">
          <div className="os-search os-max-w-22">
            <Search size={16} className="os-search__icon" />
            <input
              className="os-search__input"
              placeholder="Search notifications…"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
            />
          </div>
          <div className="os-w-16">
            <Dropdown
              id="notification-category-filter"
              titleText=""
              label="All categories"
              items={["", ...CATEGORIES.map((c) => c.value)]}
              itemToString={(item) => (item ? CATEGORY_LABEL[item as string] : "All categories")}
              selectedItem={category}
              onChange={({ selectedItem }) => setCategory((selectedItem as NotificationCategory | "") ?? "")}
            />
          </div>
        </div>

        {isLoading ? (
          <LoadingSpinner />
        ) : filtered.length === 0 ? (
          <EmptyState
            title={tab === "archived" ? "No archived notifications" : "You're all caught up"}
            description={
              tab === "unread"
                ? "New notifications will appear here."
                : tab === "archived"
                  ? "Notifications you archive will show up here."
                  : "Notifications you've read will show up here."
            }
          />
        ) : (
          <div>
            {filtered.map((n) => (
              <NotificationRow key={n.recipient_id} n={n} />
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
