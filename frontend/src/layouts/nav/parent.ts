import { Home, Notification } from "@carbon/icons-react";
import type { NavGroup } from "@/layouts/nav/types";

export const PARENT_NAV: NavGroup[] = [
  { label: "Overview", items: [{ path: "/", label: "My Children", Icon: Home }] },
  { label: "System", items: [{ path: "/notification-center", label: "Notifications", Icon: Notification }] },
];
