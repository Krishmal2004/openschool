import { useState } from "react";
import { Link } from "react-router";
import { Button } from "@carbon/react";
import { Search, EventSchedule } from "@carbon/icons-react";
import { useMyClasses } from "../../queries/useTeachers";
import { useClassStudents } from "../../queries/useClasses";
import LoadingSpinner from "../../components/common/LoadingSpinner";
import ErrorMessage from "../../components/common/ErrorMessage";
import EmptyState from "../../components/common/EmptyState";
import InfoRow from "../../components/common/InfoRow";

export default function TeacherClasses() {
  const { classes: myClasses, isLoading, isError, refetch } = useMyClasses();
  const [activeClassId, setActiveClassId] = useState("");
  const [query, setQuery] = useState("");

  const activeClass = myClasses.find((c) => c.class_id === activeClassId) ?? myClasses[0];
  const { data: roster, isLoading: rosterLoading } = useClassStudents(activeClass?.class_id ?? "");

  if (isLoading) return <LoadingSpinner />;
  if (isError) {
    return (
      <div style={{ padding: "2rem" }}>
        <ErrorMessage message="Failed to load your classes" onRetry={refetch} />
      </div>
    );
  }
  if (myClasses.length === 0 || !activeClass) {
    return (
      <div className="os-page">
        <EmptyState
          title="No classes assigned yet"
          description="Classes you teach will appear here once an admin assigns you to a subject in a class."
        />
      </div>
    );
  }

  const filtered = (roster ?? []).filter(
    (s) =>
      s.full_name.toLowerCase().includes(query.toLowerCase()) ||
      s.index_number.toLowerCase().includes(query.toLowerCase())
  );

  const houseCounts = new Map<string, number>();
  for (const s of roster ?? []) {
    const house = s.house_name ?? "No house";
    houseCounts.set(house, (houseCounts.get(house) ?? 0) + 1);
  }

  return (
    <div className="os-page">
      <div className="os-page__header">
        <div className="os-page__header-left">
          <h1 className="os-page__title">My Classes</h1>
          <p className="os-page__subtitle">{myClasses.length} class{myClasses.length > 1 ? "es" : ""} this year</p>
        </div>
      </div>

      {/* Class tabs */}
      <div style={{ display: "flex", gap: "0.75rem", marginBottom: "1.5rem", flexWrap: "wrap" }}>
        {myClasses.map((c) => (
          <button
            key={c.class_id}
            onClick={() => { setActiveClassId(c.class_id); setQuery(""); }}
            style={{
              padding: "0.625rem 1.25rem", border: "1px solid", cursor: "pointer", fontSize: "0.875rem", fontWeight: 500, fontFamily: "inherit", transition: "all 0.15s",
              background: activeClass.class_id === c.class_id ? "var(--os-accent)" : "var(--os-layer)",
              borderColor: activeClass.class_id === c.class_id ? "var(--os-accent)" : "var(--os-border-subtle)",
              color: activeClass.class_id === c.class_id ? "var(--os-layer)" : "var(--os-text-primary)",
            }}
          >
            {c.grade_name} — {c.class_name}
          </button>
        ))}
      </div>

      <div style={{ display: "grid", gridTemplateColumns: "2fr 1fr", gap: "1.5rem", alignItems: "start" }}>
        {/* Student roster */}
        <div>
          <div className="os-section">
            <div className="os-section__header">
              <h2 className="os-section__title">Student Roster — {activeClass.class_name}</h2>
              <span style={{ fontSize: "0.75rem", color: "var(--os-text-tertiary)" }}>{roster?.length ?? 0} students</span>
            </div>
            <div className="os-toolbar">
              <div className="os-search" style={{ maxWidth: "22rem" }}>
                <Search size={16} className="os-search__icon" />
                <input
                  className="os-search__input"
                  placeholder="Search students…"
                  value={query}
                  onChange={e => setQuery(e.target.value)}
                />
              </div>
              <Button as={Link} to="/t/attendance" renderIcon={EventSchedule} kind="primary" size="md" style={{ marginLeft: "auto" }}>
                Mark Attendance
              </Button>
            </div>
            {rosterLoading ? (
              <LoadingSpinner />
            ) : (
              <table className="os-table">
                <thead>
                  <tr><th>#</th><th>Name</th><th>Index Number</th><th>Gender</th></tr>
                </thead>
                <tbody>
                  {filtered.map((s, i) => (
                    <tr key={s.id}>
                      <td className="os-table__muted" style={{ width: "2rem" }}>{i + 1}</td>
                      <td>{s.full_name}</td>
                      <td className="os-table__mono">{s.index_number}</td>
                      <td className="os-table__muted">{s.gender ? s.gender[0].toUpperCase() + s.gender.slice(1) : "—"}</td>
                    </tr>
                  ))}
                  {filtered.length === 0 && (
                    <tr><td colSpan={4} style={{ textAlign: "center", color: "var(--os-text-tertiary)", padding: "2rem" }}>No students found</td></tr>
                  )}
                </tbody>
              </table>
            )}
          </div>
        </div>

        {/* Class info */}
        <div>
          <div className="os-section">
            <div className="os-section__header"><h2 className="os-section__title">Class Details</h2></div>
            <div className="os-section__body" style={{ padding: "0.75rem 1.5rem" }}>
              <InfoRow label="Grade" value={activeClass.grade_name} />
              <InfoRow label="Class" value={activeClass.class_name} />
              <InfoRow label="Your role" value={activeClass.isFormTeacher ? "Form Teacher" : "Subject Teacher"} />
              <InfoRow label="Subjects you teach" value={activeClass.subjects.length > 0 ? activeClass.subjects.join(", ") : "—"} />
              <InfoRow label="Students" value={roster?.length ?? 0} divider={false} />
            </div>
          </div>

          <div className="os-section">
            <div className="os-section__header"><h2 className="os-section__title">Houses</h2></div>
            <div className="os-section__body" style={{ padding: "0.75rem 1.5rem" }}>
              {houseCounts.size === 0 ? (
                <p style={{ color: "var(--os-text-tertiary)", fontSize: "0.8125rem" }}>No students yet.</p>
              ) : (
                [...houseCounts.entries()].map(([house, count], i, arr) => (
                  <InfoRow key={house} label={house} value={count} bold accent divider={i < arr.length - 1} />
                ))
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
