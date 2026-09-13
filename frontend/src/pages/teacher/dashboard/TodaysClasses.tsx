import { Link } from "react-router";
import { Tag } from "@carbon/react";
import LoadingSpinner from "../../../components/common/LoadingSpinner";
import type { MyClass } from "../../../queries/useTeachers";
import type { DailySession } from "../../../services/attendance";

export default function TodaysClasses({
  loading,
  myClasses,
  studentCountByClass,
  todaySessionByClass,
}: {
  loading: boolean;
  myClasses: MyClass[];
  studentCountByClass: Map<string, number>;
  todaySessionByClass: Map<string, DailySession>;
}) {
  return (
    <div className="os-section">
      <div className="os-section__header">
        <h2 className="os-section__title">Today's Classes</h2>
        <span style={{ fontSize: "0.6875rem", color: "var(--os-text-tertiary)" }}>
          {new Date().toLocaleDateString("en-US", { weekday: "long", year: "numeric", month: "long", day: "numeric" })}
        </span>
      </div>
      <div>
        {loading ? (
          <LoadingSpinner />
        ) : myClasses.length === 0 ? (
          <p style={{ padding: "1.5rem", color: "var(--os-text-tertiary)", fontSize: "0.8125rem" }}>No classes assigned yet.</p>
        ) : (
          myClasses.map((cls) => {
            const session = todaySessionByClass.get(cls.class_id);
            const isMarked = !!session && session.marked_count > 0;
            return (
              <div key={cls.class_id} className="os-list-row" style={{ padding: "1rem 1.5rem", flexWrap: "wrap" }}>
                <div style={{ width: "2.25rem", height: "2.25rem", background: "var(--os-accent-light)", display: "flex", alignItems: "center", justifyContent: "center", flexShrink: 0, fontWeight: 700, fontSize: "0.75rem", color: "var(--os-accent)" }}>
                  {cls.class_name}
                </div>
                <div style={{ flex: 1, minWidth: 0 }}>
                  <p style={{ margin: "0 0 0.15rem", fontWeight: 600, fontSize: "0.875rem", color: "var(--os-text-primary)" }}>
                    {cls.grade_name} — {cls.class_name}
                  </p>
                  <p style={{ margin: 0, fontSize: "0.75rem", color: "var(--os-text-secondary)" }}>
                    {cls.subjects.join(", ")} · {studentCountByClass.get(cls.class_id) ?? 0} students
                  </p>
                </div>
                <div style={{ display: "flex", alignItems: "center", gap: "0.5rem" }}>
                  <Tag type={isMarked ? "blue" : "gray"} size="sm">
                    {isMarked ? "Marked" : "Pending"}
                  </Tag>
                  {session ? (
                    <Link to={`/attendance/sessions/${session.id}/mark`} style={{ fontSize: "0.8125rem", color: "var(--os-text-tertiary)", textDecoration: "none", whiteSpace: "nowrap" }}>
                      View →
                    </Link>
                  ) : (
                    <Link to="/t/attendance" style={{ fontSize: "0.8125rem", color: "var(--os-accent)", textDecoration: "none", fontWeight: 500, whiteSpace: "nowrap" }}>
                      Mark now →
                    </Link>
                  )}
                </div>
              </div>
            );
          })
        )}
      </div>
    </div>
  );
}
