import { Tag } from "@carbon/react";
import { EventSchedule, Report, Table, Notification } from "@carbon/icons-react";
import { useMemo } from "react";
import { useMyStudentProfile, useMyAttendance, useMyMarks } from "@/features/students/queries/useStudentSelf";
import { useMyClassTimetable } from "@/features/timetable/queries/useTimetables";
import { useUnreadNotificationCount } from "@/features/notifications/queries/useNotifications";
import { useCurrentTerm } from "@/features/school/queries/useTerms";
import LoadingSpinner from "@/shared/ui/LoadingSpinner";
import ErrorMessage from "@/shared/ui/ErrorMessage";
import TodayStatCard from "@/shared/ui/TodayStatCard";
import { getInitials } from "@/shared/lib/name";
import { todayISODate } from "@/shared/lib/date";

export default function StudentDashboard() {
  const { data: profile, isLoading, isError, refetch } = useMyStudentProfile();
  const { data: attendance, isLoading: attendanceLoading } = useMyAttendance();
  const { data: currentTerm } = useCurrentTerm();
  const { data: marks, isLoading: marksLoading } = useMyMarks(currentTerm?.id ?? "");
  const { data: timetable, isLoading: timetableLoading } = useMyClassTimetable();
  const { data: unread, isLoading: unreadLoading } = useUnreadNotificationCount();

  const attendanceThisMonth = useMemo(() => {
    const month = todayISODate().slice(0, 7);
    const rows = (attendance ?? []).filter((r) => r.session_date.startsWith(month));
    if (rows.length === 0) return null;
    const present = rows.filter((r) => r.status === "present" || r.status === "late").length;
    return Math.round((present / rows.length) * 100);
  }, [attendance]);

  const marksAverage = useMemo(() => {
    const graded = (marks ?? []).filter((m) => !m.is_absent && m.max_marks > 0);
    if (graded.length === 0) return null;
    const pct = graded.reduce((sum, m) => sum + m.marks / m.max_marks, 0) / graded.length;
    return Math.round(pct * 100);
  }, [marks]);

  const todayClasses = useMemo(
    () => (timetable?.entries ?? []).filter((e) => e.day_of_week === new Date().getDay()),
    [timetable],
  );

  if (isLoading) return <LoadingSpinner />;
  if (isError || !profile) {
    return (
      <div className="os-p-8">
        <ErrorMessage message="Failed to load your profile" onRetry={refetch} />
      </div>
    );
  }

  return (
    <div className="os-bg-layer-hover os-min-h-content">
      <div className="os-profile__banner">
        <div className="os-profile__avatar">
          {getInitials(profile.full_name)}
        </div>
        <div className="os-flex-1">
          <p className="os-profile__name">{profile.full_name}</p>
          <p className="os-profile__meta">
            {profile.index_number}
            {profile.class_name ? ` · ${profile.class_name}` : ""}
            {profile.grade_name ? ` · ${profile.grade_name}` : ""}
          </p>
        </div>
        {profile.house_name && (
          <div className="os-profile__actions">
            <Tag type="blue" size="sm">
              {profile.house_name} House
            </Tag>
          </div>
        )}
      </div>

      <div className="os-py-6 os-px-8">
        <div className="os-stat-grid">
          <TodayStatCard
            label="Attendance this month"
            value={attendanceThisMonth === null ? "-" : `${attendanceThisMonth}%`}
            loading={attendanceLoading}
            Icon={EventSchedule}
            path="/s/attendance"
          />
          <TodayStatCard
            label="Latest marks"
            value={marksAverage === null ? "-" : `${marksAverage}%`}
            meta={currentTerm?.name}
            loading={marksLoading || !currentTerm}
            Icon={Report}
            path="/s/marks"
          />
          <TodayStatCard
            label="Today's classes"
            value={todayClasses.length}
            meta={todayClasses.length > 0 ? todayClasses[0].subject_name ?? undefined : undefined}
            loading={timetableLoading}
            Icon={Table}
            path="/s/timetable"
          />
          <TodayStatCard
            label="Unread notices"
            value={unread ?? 0}
            loading={unreadLoading}
            Icon={Notification}
            path="/notification-center"
          />
        </div>

        <div className="os-section">
          <div className="os-section__header">
            <h2 className="os-section__title">Student Details</h2>
          </div>
          <div className="os-section__body os-grid os-grid-auto-200 os-gap-4">
            <div className="os-kv-item os-border os-bg-layer">
              <p className="os-kv-item__label">Full Name</p>
              <p className="os-kv-item__value">{profile.full_name}</p>
            </div>
            <div className="os-kv-item os-border os-bg-layer">
              <p className="os-kv-item__label">Index Number</p>
              <p className="os-kv-item__value">{profile.index_number}</p>
            </div>
            <div className="os-kv-item os-border os-bg-layer">
              <p className="os-kv-item__label">Class</p>
              <p className="os-kv-item__value">{profile.class_name || "-"}</p>
            </div>
            <div className="os-kv-item os-border os-bg-layer">
              <p className="os-kv-item__label">Grade</p>
              <p className="os-kv-item__value">{profile.grade_name || "-"}</p>
            </div>
            <div className="os-kv-item os-border os-bg-layer">
              <p className="os-kv-item__label">House</p>
              <p className="os-kv-item__value">{profile.house_name || "-"}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
