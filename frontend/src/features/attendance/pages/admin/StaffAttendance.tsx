import { STAFF_ATTENDANCE_STYLES } from "@/shared/lib/constants/attendance";
import { useState } from "react";
import { DatePicker, DatePickerInput, Tag, Toggle } from "@carbon/react";
import { CheckmarkFilled, CloseFilled, Time, Renew } from "@carbon/icons-react";
import {
  useStaffAttendanceByDate,
  useStaffAttendanceMonthlySummary,
  useMarkStaffAttendance,
} from "@/features/attendance/queries/useStaffAttendance";
import type { StaffAttendanceRow, StaffAttendanceStatus } from "@/features/attendance/api/staffAttendance";
import { todayISODate, toYmd, formatLongDate } from "@/shared/lib/date";
import LoadingSpinner from "@/shared/ui/LoadingSpinner";
import ErrorMessage from "@/shared/ui/ErrorMessage";

// "Leave" reuses the "excused" attendance-status colors, same palette, different label for the staff-attendance domain.

const STATUS_ICONS: Record<StaffAttendanceStatus, typeof CheckmarkFilled> = {
  present: CheckmarkFilled,
  absent: CloseFilled,
  late: Time,
  leave: Renew,
};

function StatusButton({
  value,
  selected,
  onClick,
}: {
  value: StaffAttendanceStatus;
  selected: boolean;
  onClick: () => void;
}) {
  const cfg = STAFF_ATTENDANCE_STYLES[value];
  const Icon = STATUS_ICONS[value];
  return (
    <button
      onClick={onClick} className={`os-py-1 os-px-3 os-text-xs ${selected ? "os-fw-600" : "os-fw-400"} os-pointer os-rounded-sm os-nowrap`} style={{ border: `1px solid ${selected ? cfg.border : "var(--os-border-subtle)"}`, background: selected ? cfg.bg : "var(--os-layer)", color: selected ? cfg.color : "var(--os-text-secondary)" }}
    >
      <Icon size={12} className="os-mr-1 os-align-middle" style={{ fill: selected ? cfg.color : "var(--os-text-tertiary)" }} />
      {cfg.label}
    </button>
  );
}

function AttendanceTable({
  title,
  rows,
  onMark,
  markingId,
}: {
  title: string;
  rows: StaffAttendanceRow[];
  onMark: (row: StaffAttendanceRow, status: StaffAttendanceStatus) => void;
  markingId: string | null;
}) {
  return (
    <div className="os-section">
      <div className="os-section__header">
        <h2 className="os-section__title">{title}</h2>
        <span className="os-text-xs os-c-tertiary">{rows.length} staff</span>
      </div>
      <table className="os-table">
        <thead>
          <tr>
            <th>Name</th>
            <th>Employee No.</th>
            <th>Attendance</th>
          </tr>
        </thead>
        <tbody>
          {rows.map((row) => (
            <tr key={row.staff_id}>
              <td className="os-fw-500">{row.full_name}</td>
              <td className="os-table__mono">{row.employee_number}</td>
              <td>
                <div className={`os-flex os-gap-1h ${markingId === row.staff_id ? "os-opacity-50" : ""}`}>
                  {(["present", "late", "absent", "leave"] as const).map((s) => (
                    <StatusButton key={s} value={s} selected={row.status === s} onClick={() => onMark(row, s)} />
                  ))}
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      {rows.length === 0 && (
        <div className="os-placeholder">
          <p>No active staff.</p>
        </div>
      )}
    </div>
  );
}

function MonthlySummaryView({ year, month }: { year: number; month: number }) {
  const { data, isLoading, isError, refetch } = useStaffAttendanceMonthlySummary(year, month);

  if (isLoading) return <LoadingSpinner />;
  if (isError || !data) return <ErrorMessage message="Could not load monthly summary." onRetry={refetch} />;

  const renderTable = (title: string, rows: typeof data.teachers) => (
    <div className="os-section">
      <div className="os-section__header">
        <h2 className="os-section__title">{title}</h2>
      </div>
      <table className="os-table">
        <thead>
          <tr>
            <th>Name</th>
            <th>Present</th>
            <th>Late</th>
            <th>Absent</th>
            <th>Leave</th>
          </tr>
        </thead>
        <tbody>
          {rows.map((r) => (
            <tr key={r.staff_id}>
              <td className="os-fw-500">{r.full_name}</td>
              <td>{r.present_count}</td>
              <td>{r.late_count}</td>
              <td>{r.absent_count}</td>
              <td>{r.leave_count}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );

  return (
    <>
      {renderTable("Teachers", data.teachers)}
      {renderTable("Non-Academic Staff", data.non_academic_staff)}
    </>
  );
}

export default function StaffAttendance() {
  const [date, setDate] = useState(todayISODate());
  const [showMonthly, setShowMonthly] = useState(false);
  const [markingId, setMarkingId] = useState<string | null>(null);

  const { data, isLoading, isError, refetch } = useStaffAttendanceByDate(date);
  const markAttendance = useMarkStaffAttendance();

  const dateObj = new Date(date);
  const year = dateObj.getFullYear();
  const month = dateObj.getMonth() + 1;

  const handleMark = (row: StaffAttendanceRow, status: StaffAttendanceStatus, isTeacher: boolean) => {
    setMarkingId(row.staff_id);
    markAttendance.mutate(
      {
        ...(isTeacher ? { teacher_id: row.staff_id } : { non_academic_staff_id: row.staff_id }),
        date: new Date(date).toISOString(),
        status,
      },
      { onSettled: () => setMarkingId(null) },
    );
  };

  return (
    <div className="os-page">
      <div className="os-page__header">
        <div className="os-page__header-left">
          <h1 className="os-page__title">Staff Attendance</h1>
          <p className="os-page__subtitle">
            Mark daily attendance for teachers and non-academic staff, or view monthly summaries.
          </p>
        </div>
      </div>

      <div className="os-flex os-items-center os-gap-6 os-mb-6 os-wrap">
        {!showMonthly && (
          <DatePicker datePickerType="single" dateFormat="Y-m-d" value={date} onChange={(dates) => {
            const ymd = toYmd(dates[0]);
            if (ymd) setDate(ymd);
          }}>
            <DatePickerInput id="staff-attendance-date" labelText="Date" placeholder="YYYY-MM-DD" />
          </DatePicker>
        )}
        <Toggle
          id="staff-attendance-view-toggle"
          labelText="View"
          labelA="Daily"
          labelB="Monthly summary"
          toggled={showMonthly}
          onToggle={(checked) => setShowMonthly(checked)}
        />
        {!showMonthly && <Tag type="gray">{formatLongDate(date)}</Tag>}
      </div>

      {showMonthly ? (
        <MonthlySummaryView year={year} month={month} />
      ) : isLoading ? (
        <LoadingSpinner />
      ) : isError || !data ? (
        <ErrorMessage message="Could not load staff attendance." onRetry={refetch} />
      ) : (
        <>
          <AttendanceTable
            title="Teachers"
            rows={data.teachers}
            markingId={markingId}
            onMark={(row, status) => handleMark(row, status, true)}
          />
          <AttendanceTable
            title="Non-Academic Staff"
            rows={data.non_academic_staff}
            markingId={markingId}
            onMark={(row, status) => handleMark(row, status, false)}
          />
        </>
      )}
    </div>
  );
}
