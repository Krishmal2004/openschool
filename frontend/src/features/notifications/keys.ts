export const notificationKeys = {
  all: ["notifications"] as const,
  mine: () => ["notifications", "mine"] as const,
  archived: () => ["notifications", "mine", "archived"] as const,
  unread: () => ["notifications", "unread-count"] as const,
  sent: () => ["notifications", "sent"] as const,
  drafts: () => ["notifications", "drafts"] as const,
  stats: (id: string) => ["notifications", "stats", id] as const,
};
