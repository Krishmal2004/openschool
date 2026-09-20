import { Home, Notification, EventSchedule, Report, Table, UserMultiple, Document, Idea } from "@carbon/icons-react";
import type { NavGroup } from "@/layouts/nav/types";

export const STUDENT_NAV: NavGroup[] = [
  { label: "Overview", items: [{ path: "/", label: "Dashboard", Icon: Home }] },
  {
    label: "Academics",
    items: [
      { path: "/s/attendance", label: "Attendance", Icon: EventSchedule },
      { path: "/s/marks", label: "Marks", Icon: Report },
      { path: "/s/timetable", label: "Timetable", Icon: Table },
      { path: "/s/enrollment", label: "Subject Enrolment", Icon: Document },
    ],
  },
  {
    label: "Portfolio",
    items: [
      { path: "/s/progress", label: "Progress Reports", Icon: Report },
      { path: "/s/portfolio", label: "Activities & Leadership", Icon: Idea },
      { path: "/s/guardians", label: "My Guardians", Icon: UserMultiple },
    ],
  },
  { label: "System", items: [{ path: "/notification-center", label: "Notifications", Icon: Notification }] },
];
