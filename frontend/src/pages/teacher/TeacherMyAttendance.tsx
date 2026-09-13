import { useState } from "react";
import { Button, Tag } from "@carbon/react";
import { ChevronLeft, ChevronRight, CheckmarkFilled, CloseFilled, Time, Renew } from "@carbon/icons-react";
import { useMyStaffAttendanceHistory } from "../../queries/useStaffAttendance";
import type { StaffAttendanceStatus } from "../../services/staffAttendance";
import LoadingSpinner from "../../components/common/LoadingSpinner";
import EmptyState from "../../components/common/EmptyState";

// "Leave" reuses the "excused" attendance-status colors — same palette, different label for the staff-attendance domain.
const STATUS_STYLES: Record<StaffAttendanceStatus, { bg: string; border: string; color: string; label: string; Icon: typeof CheckmarkFilled }> = {
  present: { bg: "var(--os-status-present-bg)", border: "var(--os-status-present-border)", color: "var(--os-status-present-text)", label: "Present", Icon: CheckmarkFilled },
  absent: { bg: "var(--os-status-absent-bg)", border: "var(--os-status-absent-border)", color: "var(--os-status-absent-text)", label: "Absent", Icon: CloseFilled },
  late: { bg: "var(--os-status-late-bg)", border: "var(--os-status-late-border)", color: "var(--os-status-late-text)", label: "Late", Icon: Time },
  leave: { bg: "var(--os-status-excused-bg)", border: "var(--os-status-excused-border)", color: "var(--os-status-excused-text)", label: "Leave", Icon: Renew },
};

const MONTH_NAMES = [
  "January", "February", "March", "April", "May", "June",
  "July", "August", "September", "October", "November", "December",
];

export default function TeacherMyAttendance() {
  const now = new Date();
  const [year, setYear] = useState(now.getFullYear());
  const [month, setMonth] = useState(now.getMonth() + 1); // 1-12

  const { data: records, isLoading, isError } = useMyStaffAttendanceHistory(year, month);

  const goPrevMonth = () => {
    if (month === 1) { setYear((y) => y - 1); setMonth(12); } else { setMonth((m) => m - 1); }
  };
  const goNextMonth = () => {
    if (month === 12) { setYear((y) => y + 1); setMonth(1); } else { setMonth((m) => m + 1); }
  };

  const counts = (records ?? []).reduce(
    (acc, r) => {
      acc[r.status] = (acc[r.status] ?? 0) + 1;
      return acc;
    },
    {} as Record<StaffAttendanceStatus, number>,
  );

  const sorted = [...(records ?? [])].sort((a, b) => b.date.localeCompare(a.date));

  return (
    <div className="os-page">
      <div className="os-page__header">
        <div className="os-page__header-left">
          <h1 className="os-page__title">My Attendance</h1>
          <p className="os-page__subtitle">Your own attendance record — marked by an administrator</p>
        </div>
      </div>

      <div className="os-section" style={{ marginBottom: "1.5rem" }}>
        <div className="os-section__header" style={{ display: "flex", alignItems: "center", gap: "0.75rem" }}>
          <Button hasIconOnly iconDescription="Previous month" renderIcon={ChevronLeft} kind="ghost" size="sm" onClick={goPrevMonth} />
          <h2 className="os-section__title" style={{ minWidth: "10rem", textAlign: "center" }}>
            {MONTH_NAMES[month - 1]} {year}
          </h2>
          <Button hasIconOnly iconDescription="Next month" renderIcon={ChevronRight} kind="ghost" size="sm" onClick={goNextMonth} />
        </div>

        <div className="os-section__body" style={{ display: "flex", gap: "1rem", flexWrap: "wrap" }}>
          {(Object.keys(STATUS_STYLES) as StaffAttendanceStatus[]).map((status) => {
            const cfg = STATUS_STYLES[status];
            return (
              <div
                key={status}
                style={{
                  flex: "1 1 8rem",
                  padding: "0.875rem 1rem",
                  background: cfg.bg,
                  border: `1px solid ${cfg.border}`,
                  borderRadius: "2px",
                }}
              >
                <p style={{ margin: "0 0 0.25rem", fontSize: "0.6875rem", fontWeight: 600, textTransform: "uppercase", letterSpacing: "0.06em", color: cfg.color }}>
                  {cfg.label}
                </p>
                <p style={{ margin: 0, fontSize: "1.5rem", fontWeight: 300, color: cfg.color }}>{counts[status] ?? 0}</p>
              </div>
            );
          })}
        </div>
      </div>

      <div className="os-section">
        {isLoading ? (
          <LoadingSpinner />
        ) : isError ? (
          <EmptyState title="Could not load attendance" description="Please try again later." />
        ) : sorted.length === 0 ? (
          <EmptyState title="No attendance marked yet" description="No attendance has been recorded for you this month." />
        ) : (
          <table className="os-table">
            <thead>
              <tr>
                <th>Date</th>
                <th>Status</th>
                <th>Note</th>
              </tr>
            </thead>
            <tbody>
              {sorted.map((r) => {
                const cfg = STATUS_STYLES[r.status];
                return (
                  <tr key={r.id}>
                    <td className="os-table__mono">{r.date}</td>
                    <td>
                      <Tag renderIcon={cfg.Icon} size="sm" style={{ background: cfg.bg, color: cfg.color }}>
                        {cfg.label}
                      </Tag>
                    </td>
                    <td className="os-table__muted">{r.note ?? "—"}</td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}
