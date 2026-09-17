import {
  Home, ChartLine, UserMultiple, Education, Building, Book, Language, EventSchedule, Settings, Calendar,
  Notification, UserFollow, UserAdmin, UserRole, Trophy, Renew, Table, Location, Rule, SettingsAdjust,
  Group, DocumentPdf, Idea, Bot, Rocket,
} from "@carbon/icons-react";
import type { NavGroup } from "@/layouts/nav/types";

export const ADMIN_NAV: NavGroup[] = [
  { label: "Overview", items: [{ path: "/", label: "Overview", Icon: Home }] },
  {
    label: "People",
    items: [
      { path: "/positions", label: "Leadership Positions", Icon: UserRole },
      { path: "/teachers", label: "Teachers", Icon: Education },
      { path: "/non-academic-staff", label: "Non-Academic Staff", Icon: Group },
      { path: "/students", label: "Students", Icon: UserMultiple },
      { path: "/guardians", label: "Guardians", Icon: UserAdmin },
    ],
  },
  {
    label: "Student Life",
    items: [
      { path: "/prefects", label: "School Prefects", Icon: Trophy },
      { path: "/societies", label: "Societies", Icon: Idea },
    ],
  },
  {
    label: "Academics",
    items: [
      { path: "/classes", label: "Grades & Classes", Icon: Building },
      { path: "/streams", label: "Streams", Icon: UserFollow },
      { path: "/subjects", label: "Subjects & Curriculum", Icon: Book },
      { path: "/mediums", label: "Mediums", Icon: Language },
      { path: "/teacher-subjects", label: "Teacher Subjects", Icon: Education },
    ],
  },
  {
    label: "Scheduling",
    items: [
      { path: "/academic-years", label: "Academic Years", Icon: Calendar },
      { path: "/timetables", label: "Timetables", Icon: Table },
      { path: "/timetables/generate", label: "Generate Timetable", Icon: Rocket },
      { path: "/classrooms", label: "Classrooms & Facilities", Icon: Location },
      { path: "/subject-requirements", label: "Subject Requirements", Icon: Rule },
      { path: "/timetable-settings", label: "Timetable Settings", Icon: SettingsAdjust },
    ],
  },
  {
    label: "Attendance",
    items: [
      { path: "/attendance", label: "Student Attendance", Icon: EventSchedule },
      { path: "/staff-attendance", label: "Staff Attendance", Icon: Group },
    ],
  },
  {
    label: "Operations",
    items: [
      { path: "/promotion", label: "Promotion", Icon: Renew },
      { path: "/notifications", label: "Notifications", Icon: Notification },
      { path: "/reports", label: "Reports", Icon: DocumentPdf },
      { path: "/analytics", label: "Analytics", Icon: ChartLine },
    ],
  },
  {
    label: "System",
    items: [
      { path: "/settings", label: "Settings", Icon: Settings },
      { path: "/automation", label: "Automation", Icon: Bot },
    ],
  },
];
