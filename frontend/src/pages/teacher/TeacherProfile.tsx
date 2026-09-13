import { useQueries } from "@tanstack/react-query";
import { Tag } from "@carbon/react";
import { Book } from "@carbon/icons-react";
import { useMyTeacherProfile, useTeacherSubjects, useMyClasses } from "../../queries/useTeachers";
import { classStudentsKey } from "../../queries/useClasses";
import { classSessionsKey } from "../../queries/useAttendance";
import { attendanceApi } from "../../services/attendance";
import { studentApi } from "../../services/student";
import LoadingSpinner from "../../components/common/LoadingSpinner";
import ErrorMessage from "../../components/common/ErrorMessage";
import InfoRow from "../../components/common/InfoRow";
import { getInitials } from "../../lib/name";

export default function TeacherProfile() {
  const { data: profile, isLoading, isError, refetch } = useMyTeacherProfile();
  const { data: subjects } = useTeacherSubjects(profile?.id ?? "");
  const { classes: myClasses } = useMyClasses();

  const classIds = myClasses.map((c) => c.class_id);
  const studentQueries = useQueries({
    queries: classIds.map((id) => ({
      queryKey: classStudentsKey(id),
      queryFn: () => studentApi.listByClass(id),
      enabled: !!id,
    })),
  });
  const sessionQueries = useQueries({
    queries: classIds.map((id) => ({
      queryKey: classSessionsKey(id),
      queryFn: () => attendanceApi.listSessionsByClass(id),
      enabled: !!id,
    })),
  });
  const studentsTaught = new Set(studentQueries.flatMap((q) => (q.data ?? []).map((s) => s.id))).size;
  const sessionsTaken = sessionQueries.reduce((sum, q) => sum + (q.data?.length ?? 0), 0);

  if (isLoading) return <LoadingSpinner />;
  if (isError || !profile) {
    return (
      <div style={{ padding: "2rem" }}>
        <ErrorMessage message="Failed to load your profile" onRetry={refetch} />
      </div>
    );
  }

  return (
    <div style={{ background: "var(--os-layer-hover)", minHeight: "calc(100vh - 3rem)" }}>
      <div className="os-profile__banner">
        <div className="os-profile__avatar">{getInitials(profile.full_name)}</div>
        <div style={{ flex: 1 }}>
          <p className="os-profile__name">{profile.title ? `${profile.title} ` : ""}{profile.full_name}</p>
          <p className="os-profile__meta">{profile.employee_number}</p>
        </div>
        <div className="os-profile__actions">
          <Tag type={profile.is_active ? "green" : "gray"} size="sm">{profile.is_active ? "Active" : "Inactive"}</Tag>
        </div>
      </div>

      <div style={{ padding: "1.5rem 2rem" }}>
        <div style={{ display: "grid", gridTemplateColumns: "2fr 1fr", gap: "1.5rem", alignItems: "start" }}>
          {/* Main */}
          <div>
            <div className="os-section">
              <div className="os-section__header"><h2 className="os-section__title">Personal Details</h2></div>
              <div className="os-kv-grid">
                {[
                  ["Full Name", profile.full_name],
                  ["Title", profile.title ?? "—"],
                  ["Gender", profile.gender ? profile.gender[0].toUpperCase() + profile.gender.slice(1) : "—"],
                  ["Phone", profile.phone ?? "—"],
                  ["Employee Number", profile.employee_number],
                  ["Joined Date", profile.joined_date ?? "—"],
                ].map(([label, value]) => (
                  <div key={label} className="os-kv-item">
                    <p className="os-kv-item__label">{label}</p>
                    <p className="os-kv-item__value">{value}</p>
                  </div>
                ))}
              </div>
            </div>

            <div className="os-section">
              <div className="os-section__header"><h2 className="os-section__title">Subjects</h2></div>
              <div className="os-section__body" style={{ display: "flex", gap: "0.5rem", flexWrap: "wrap" }}>
                {!subjects || subjects.length === 0 ? (
                  <p style={{ color: "var(--os-text-tertiary)", fontSize: "0.8125rem" }}>No subjects assigned yet.</p>
                ) : (
                  subjects.map((s) => (
                    <div key={s.id} style={{ display: "flex", alignItems: "center", gap: "0.5rem", padding: "0.5rem 0.875rem", border: "1px solid var(--os-border-subtle)", background: "var(--os-layer-hover)" }}>
                      <Book size={14} style={{ fill: "var(--os-accent)" }} />
                      <span style={{ fontSize: "0.875rem", fontWeight: 500 }}>{s.name}</span>
                    </div>
                  ))
                )}
              </div>
            </div>
          </div>

          {/* Sidebar */}
          <div>
            <div className="os-section">
              <div className="os-section__header"><h2 className="os-section__title">Quick Info</h2></div>
              <div className="os-section__body" style={{ padding: "0.75rem 1.5rem" }}>
                <InfoRow label="Employee ID" value={profile.employee_number} />
                <InfoRow label="Status" value={profile.is_active ? "Active" : "Inactive"} />
                <InfoRow label="Subjects" value={subjects?.length ?? 0} />
                <InfoRow label="Classes" value={myClasses.length} />
                <InfoRow label="Joined" value={profile.joined_date ?? "—"} divider={false} />
              </div>
            </div>

            <div className="os-section">
              <div className="os-section__header"><h2 className="os-section__title">This Year</h2></div>
              <div className="os-section__body" style={{ padding: "0.75rem 1.5rem" }}>
                <InfoRow label="Sessions Taken" value={sessionsTaken} bold accent />
                <InfoRow label="Students Taught" value={studentsTaught} bold accent divider={false} />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
