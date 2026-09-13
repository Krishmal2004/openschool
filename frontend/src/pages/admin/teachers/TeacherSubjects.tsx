import { useState } from "react";
import { Search, Tag, SkeletonText } from "@carbon/react";
import { useTeachers, useTeacherSubjects, useAssignTeacherSubject, useRemoveTeacherSubject } from "../../../queries/useTeachers";
import { useSubjects } from "../../../queries/useSubjects";
import EntityCombobox from "../../../components/common/EntityCombobox";
import type { Teacher, TeacherSubject } from "../../../services/teacher";
import type { Subject } from "../../../services/subject";
import LoadingSpinner from "../../../components/common/LoadingSpinner";
import ErrorMessage from "../../../components/common/ErrorMessage";
import ConfirmDeleteModal from "../../../components/common/ConfirmDeleteModal";

function TeacherSubjectRow({ teacher, allSubjects }: { teacher: Teacher; allSubjects: Subject[] }) {
  const { data: assignedSubjects, isLoading, isError } = useTeacherSubjects(teacher.id);
  const assignMutation = useAssignTeacherSubject(teacher.id);
  const removeMutation = useRemoveTeacherSubject(teacher.id);
  const [subjectToRemove, setSubjectToRemove] = useState<TeacherSubject | null>(null);

  const handleAssign = (subjectId: string) => {
    if (!subjectId) return;
    assignMutation.mutate(subjectId);
  };

  const confirmRemove = () => {
    if (!subjectToRemove) return;
    removeMutation.mutate(subjectToRemove.id, { onSettled: () => setSubjectToRemove(null) });
  };

  const assignedIds = new Set(assignedSubjects?.map((s) => s.id) ?? []);
  const assignableSubjects = allSubjects.filter((s) => !assignedIds.has(s.id));

  return (
    <tr>
      <td style={{ fontWeight: 500 }}>{teacher.full_name}</td>
      <td className="os-table__mono">{teacher.employee_number}</td>
      <td>
        {isLoading ? (
          <SkeletonText width="6rem" />
        ) : isError ? (
          <span style={{ color: "var(--os-danger)", fontSize: "0.875rem" }}>Error loading subjects</span>
        ) : !assignedSubjects || assignedSubjects.length === 0 ? (
          <span style={{ fontSize: "0.875rem", color: "var(--os-text-tertiary)" }}>No subjects assigned</span>
        ) : (
          <div style={{ display: "flex", flexWrap: "wrap", gap: "0.25rem" }}>
            {assignedSubjects.map((s) => (
              <Tag
                key={s.id}
                type="blue"
                size="sm"
                title="Click to remove"
                onClick={() => setSubjectToRemove(s)}
                style={{ cursor: "pointer" }}
              >
                {s.name} &times;
              </Tag>
            ))}
          </div>
        )}
      </td>
      <td style={{ minWidth: "12rem" }}>
        <EntityCombobox
          id={`assign-subject-${teacher.id}`}
          items={assignableSubjects}
          selectedId=""
          onSelect={handleAssign}
          getId={(s) => s.id}
          itemToString={(s) => `${s.name} (${s.code})`}
          labelText=""
          placeholder="Assign subject…"
        />
        {assignMutation.isError && (
          <div style={{ color: "var(--os-danger)", fontSize: "0.75rem", marginTop: "0.25rem" }}>
            Failed to assign
          </div>
        )}
        {removeMutation.isError && (
          <div style={{ color: "var(--os-danger)", fontSize: "0.75rem", marginTop: "0.25rem" }}>
            Failed to remove
          </div>
        )}
      </td>
      <ConfirmDeleteModal
        open={!!subjectToRemove}
        title="Remove subject"
        description={
          <>
            Remove <strong>{subjectToRemove?.name}</strong> from {teacher.full_name}&apos;s assigned subjects?
          </>
        }
        confirmLabel="Remove"
        pendingLabel="Removing…"
        isPending={removeMutation.isPending}
        onClose={() => setSubjectToRemove(null)}
        onConfirm={confirmRemove}
      />
    </tr>
  );
}

export default function TeacherSubjects() {
  const { data: teachers, isLoading: loadingTeachers, isError: teachersError, refetch: refetchTeachers } = useTeachers();
  const { data: subjects, isLoading: loadingSubjects, isError: subjectsError, refetch: refetchSubjects } = useSubjects();
  const [searchQuery, setSearchQuery] = useState("");

  if (loadingTeachers || loadingSubjects) {
    return <LoadingSpinner />;
  }

  if (teachersError) {
    return <ErrorMessage message="Could not load teachers." onRetry={refetchTeachers} />;
  }

  if (subjectsError) {
    return <ErrorMessage message="Could not load subjects." onRetry={refetchSubjects} />;
  }

  const filteredTeachers = (teachers ?? []).filter((t) =>
    t.full_name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    t.employee_number.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <div className="os-page">
      <div className="os-page__header">
        <div className="os-page__header-left">
          <h1 className="os-page__title">Teacher Subjects</h1>
          <p className="os-page__subtitle">
            Assign subjects to teachers globally. These assignments designate which subjects a teacher is qualified to teach.
          </p>
        </div>
      </div>

      <div style={{ marginBottom: "1.25rem", maxWidth: "24rem" }}>
        <Search
          id="teacher-search"
          placeholder="Search teachers by name or employee number…"
          labelText="Search"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
        />
      </div>

      <div className="os-section">
        <table className="os-table">
          <thead>
            <tr>
              <th>Teacher</th>
              <th>Employee #</th>
              <th>Assigned Subjects (Click to remove)</th>
              <th>Assign Subject</th>
            </tr>
          </thead>
          <tbody>
            {filteredTeachers.map((t) => (
              <TeacherSubjectRow key={t.id} teacher={t} allSubjects={subjects ?? []} />
            ))}
            {filteredTeachers.length === 0 && (
              <tr>
                <td colSpan={4} style={{ textAlign: "center", color: "var(--os-text-tertiary)", padding: "2rem" }}>
                  No teachers found matching your search.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
