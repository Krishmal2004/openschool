import { useState } from "react";
import { Link } from "react-router";
import { EventSchedule, CheckmarkFilled, WarningFilled } from "@carbon/icons-react";
import { Button, Tag, DatePicker, DatePickerInput } from "@carbon/react";
import { useDailySessions, useDeleteSession } from "@/features/attendance/queries/useAttendance";
import { useRole } from "@/shared/auth/useRole";
import type { DailySession } from "@/features/attendance/api/attendance";
import { toYmd, todayISODate, isLockedAfter24Hours } from "@/shared/lib/date";
import TableSkeleton from "@/shared/ui/TableSkeleton";
import DataGrid, { type GridColumn } from "@/shared/ui/DataGrid";
import StatCardSkeleton from "@/shared/ui/StatCardSkeleton";
import SectionHeader from "@/shared/ui/SectionHeader";
import ErrorMessage from "@/shared/ui/ErrorMessage";
import EmptyState from "@/shared/ui/EmptyState";
import ConfirmDeleteModal from "@/shared/ui/ConfirmDeleteModal";
import AgentFindingsBanner from "@/features/notifications/components/AgentFindingsBanner";
import MutationErrorNotification from "@/shared/ui/MutationErrorNotification";

const HEADERS = ["Class", "Grade", "Teacher", "Records", "Status", "Actions"];

function displayDate(ymd: string) {
  const [y, m, d] = ymd.split("-").map(Number);
  return new Date(y, m - 1, d).toLocaleDateString("en-LK", {
    weekday: "long",
    year: "numeric",
    month: "long",
    day: "numeric",
  });
}

export default function Attendance() {
  const [date, setDate] = useState(todayISODate());
  const { data: sessions, isLoading, isError, refetch } = useDailySessions(date);
  const deleteSession = useDeleteSession();
  const { role } = useRole();
  const isAdmin = role === "admin";
  const [toDelete, setToDelete] = useState<DailySession | null>(null);


  const columns: GridColumn<DailySession>[] = [
    { key: "class", header: "Class", render: (s) => <span className="os-fw-600">{s.class_name}</span> },
    { key: "grade", header: "Grade", render: (s) => <span className="os-table__muted">{s.grade_name}</span> },
    { key: "teacher", header: "Teacher", render: (s) => <span className="os-table__muted">{s.teacher_name}</span> },
    { key: "records", header: "Records", render: (s) => <span className="os-table__muted">{s.marked_count} / {s.enrolled_count}</span> },
    {
      key: "status",
      header: "Status",
      render: (s) => {
        const isMarked = s.marked_count > 0;
        return (
          <div className="os-flex os-gap-1h os-wrap">
            <Tag type={isMarked ? "blue" : "gray"} size="sm">
              {isMarked ? <CheckmarkFilled size={12} className="os-mr-1 os-align-middle" /> : <WarningFilled size={12} className="os-mr-1 os-align-middle" />}
              {isMarked ? "Marked" : "Pending"}
            </Tag>
            {isLockedAfter24Hours(s.created_at) && <Tag type="magenta" size="sm">Locked</Tag>}
          </div>
        );
      },
    },
    {
      key: "actions",
      header: "Actions",
      align: "end",
      render: (s) => {
        const readOnly = s.marked_count > 0 || (isLockedAfter24Hours(s.created_at) && !isAdmin);
        return (
          <div className="os-grid__actions os-nowrap-flex">
            <Button kind={readOnly ? "ghost" : "primary"} size="sm" as={Link} to={`/attendance/sessions/${s.id}/mark`} className={`os-nowrap${readOnly ? " os-c-accent" : ""}`}>
              {readOnly ? "View" : "Mark"}
            </Button>
            <Button kind="danger--ghost" size="sm" className="os-nowrap" onClick={() => setToDelete(s)}>Delete</Button>
          </div>
        );
      },
    },
  ];

  const marked = (sessions ?? []).filter((s) => s.marked_count > 0).length;
  const pending = (sessions ?? []).length - marked;
  const isToday = date === todayISODate();

  return (
    <div className="os-page">
      <div className="os-page__header">
        <div className="os-page__header-left">
          <h1 className="os-page__title">Attendance</h1>
          <p className="os-page__subtitle">{displayDate(date)}</p>
        </div>
        <div className="os-min-w-12">
          <DatePicker
            datePickerType="single"
            dateFormat="Y-m-d"
            value={date}
            onChange={(dates) => {
              if (dates[0]) setDate(toYmd(dates[0]));
            }}
          >
            <DatePickerInput id="attendance-date" labelText="" placeholder="YYYY-MM-DD" size="lg" />
          </DatePicker>
        </div>
      </div>

      <AgentFindingsBanner
        titles={[
          "Incomplete attendance sessions",
          "Classes missing today's attendance session",
          "Classes with inconsistent attendance-taking",
        ]}
      />

      {!isToday && (
        <div className="os-mb-4">
          <Button kind="ghost" size="sm" onClick={() => setDate(todayISODate())}>
            Jump to today
          </Button>
        </div>
      )}

      {isLoading ? (
        <>
          <div className="os-stat-grid os-grid-cols-3">
            {Array.from({ length: 3 }).map((_, i) => (
              <StatCardSkeleton key={i} />
            ))}
          </div>
          <div className="os-section">
            <SectionHeader title="Sessions" />
            <TableSkeleton headers={HEADERS} />
          </div>
        </>
      ) : isError ? (
        <ErrorMessage message="Could not load sessions for this date." onRetry={refetch} />
      ) : (
        <>
          <div className="os-stat-grid os-grid-cols-3">
            <div className="os-stat-card">
              <p className="os-stat-card__label">
                <EventSchedule size={14} className="os-fill-accent" /> Sessions
              </p>
              <p className="os-stat-card__value">{sessions?.length ?? 0}</p>
              <p className="os-stat-card__meta">Created for this date</p>
            </div>
            <div className="os-stat-card os-border-t-success">
              <p className="os-stat-card__label">
                <CheckmarkFilled size={14} className="os-fill-success" /> Marked
              </p>
              <p className="os-stat-card__value">{marked}</p>
              <p className="os-stat-card__meta">At least one record</p>
            </div>
            <div className="os-stat-card os-border-t-warning">
              <p className="os-stat-card__label">
                <WarningFilled size={14} className="os-fill-warning" /> Pending
              </p>
              <p className="os-stat-card__value">{pending}</p>
              <p className="os-stat-card__meta">No records yet</p>
            </div>
          </div>

          <div className="os-section">
            <SectionHeader title="Sessions" />

            <MutationErrorNotification
              isError={deleteSession.isError}
              error={deleteSession.error}
              title="Could not delete session"
              fallback="Please try again."
              onClose={() => deleteSession.reset()} className="os-mt-0 os-mx-6 os-mb-4"
            />

            {!sessions || sessions.length === 0 ? (
              <EmptyState
                title="No sessions for this date"
                description="Sessions are created from a class's Attendance tab - go to a class to start one."
                action={
                  <Button kind="primary" as={Link} to="/classes">
                    Go to Classes
                  </Button>
                }
              />
            ) : (
              <DataGrid rows={sessions} columns={columns} getRowId={(s) => s.id} noHover pageSize={20} />
            )}
          </div>
        </>
      )}

      <ConfirmDeleteModal
        open={!!toDelete}
        title="Delete attendance session"
        description={
          <>
            Delete the session for <strong>{toDelete?.class_name}</strong> on{" "}
            <strong>{toDelete?.date}</strong>? Every attendance record already
            marked for it is deleted too.
          </>
        }
        isPending={deleteSession.isPending}
        onClose={() => setToDelete(null)}
        onConfirm={() => toDelete && deleteSession.mutate(toDelete.id, { onSettled: () => setToDelete(null) })}
      />
    </div>
  );
}
