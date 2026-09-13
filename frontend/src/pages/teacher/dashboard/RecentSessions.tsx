import { Link } from "react-router";
import { Tag } from "@carbon/react";
import { useSessionRecords } from "../../../queries/useAttendance";
import type { AttendanceSession } from "../../../services/attendance";

function RecentSessionRow({ session, className }: { session: AttendanceSession; className: string }) {
  const { data: records, isLoading } = useSessionRecords(session.id);
  const present = records?.filter((r) => r.status === "present").length ?? 0;
  const absent = records?.filter((r) => r.status === "absent").length ?? 0;

  return (
    <tr>
      <td><span style={{ fontSize: "0.8125rem", color: "var(--os-text-secondary)" }}>{session.date}</span></td>
      <td className="os-table__link">{className}</td>
      <td>
        {isLoading ? <span style={{ color: "var(--os-text-tertiary)" }}>…</span> : <span style={{ color: "var(--os-success)", fontWeight: 600 }}>{present}</span>}
      </td>
      <td>
        {isLoading ? (
          <span style={{ color: "var(--os-text-tertiary)" }}>…</span>
        ) : (
          <span style={{ color: absent > 0 ? "var(--os-danger)" : "var(--os-text-tertiary)", fontWeight: absent > 0 ? 600 : 400 }}>{absent}</span>
        )}
      </td>
      <td><Tag type="blue" size="sm">Marked</Tag></td>
    </tr>
  );
}

export default function RecentSessions({
  sessions,
}: {
  sessions: { session: AttendanceSession; className: string }[];
}) {
  return (
    <div className="os-section">
      <div className="os-section__header">
        <h2 className="os-section__title">Recent Sessions</h2>
        <Link to="/t/attendance" style={{ fontSize: "0.75rem", color: "var(--os-accent)", textDecoration: "none" }}>View all →</Link>
      </div>
      <table className="os-table">
        <thead>
          <tr><th>Date</th><th>Class</th><th>Present</th><th>Absent</th><th>Status</th></tr>
        </thead>
        <tbody>
          {sessions.length === 0 ? (
            <tr><td colSpan={5} style={{ textAlign: "center", color: "var(--os-text-tertiary)", padding: "2rem" }}>No attendance sessions recorded yet</td></tr>
          ) : (
            sessions.map(({ session, className }) => (
              <RecentSessionRow key={session.id} session={session} className={className} />
            ))
          )}
        </tbody>
      </table>
    </div>
  );
}
