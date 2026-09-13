import { Link } from "react-router";
import { CheckmarkFilled } from "@carbon/icons-react";
import type { MyClass } from "../../../queries/useTeachers";
import type { DailySession } from "../../../services/attendance";

export default function MyClassesPanel({
  myClasses,
  studentCountByClass,
  todaySessionByClass,
}: {
  myClasses: MyClass[];
  studentCountByClass: Map<string, number>;
  todaySessionByClass: Map<string, DailySession>;
}) {
  return (
    <div className="os-section">
      <div className="os-section__header">
        <h2 className="os-section__title">My Classes</h2>
        <Link to="/t/classes" style={{ fontSize: "0.75rem", color: "var(--os-accent)", textDecoration: "none" }}>View →</Link>
      </div>
      <div>
        {myClasses.length === 0 ? (
          <p style={{ padding: "1.5rem", color: "var(--os-text-tertiary)", fontSize: "0.8125rem" }}>No classes assigned yet.</p>
        ) : (
          myClasses.map((cls) => {
            const session = todaySessionByClass.get(cls.class_id);
            const isMarked = !!session && session.marked_count > 0;
            return (
              <div key={cls.class_id} className="os-list-row os-list-row--compact" style={{ padding: "0.75rem 1.5rem" }}>
                <div style={{ flex: 1 }}>
                  <p style={{ margin: "0 0 0.1rem", fontWeight: 600, fontSize: "0.8125rem", color: "var(--os-text-primary)" }}>{cls.grade_name} — {cls.class_name}</p>
                  <p style={{ margin: 0, fontSize: "0.75rem", color: "var(--os-text-secondary)" }}>{cls.subjects.join(", ")} · {studentCountByClass.get(cls.class_id) ?? 0} students</p>
                </div>
                <Link to="/t/classes" style={{ fontSize: "0.75rem", color: "var(--os-accent)", textDecoration: "none" }}>
                  <CheckmarkFilled size={14} style={{ fill: isMarked ? "var(--os-success)" : "var(--os-border-subtle)" }} />
                </Link>
              </div>
            );
          })
        )}
      </div>
    </div>
  );
}
