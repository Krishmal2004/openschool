import { useState } from "react";
import { Button, Select, SelectItem, TextInput, Tag } from "@carbon/react";
import { Add } from "@carbon/icons-react";
import { useCurrentAcademicYear } from "../../../queries/useAcademicYears";
import {
  useStudentActivities,
  useCreateStudentActivity,
  useDeleteStudentActivity,
} from "../../../queries/useStudentPortfolio";
import { useStudentSocietyMemberships } from "../../../queries/useSocieties";
import { ACTIVITY_CATEGORIES } from "../../../services/studentPortfolio";
import type { ActivityCategory } from "../../../services/studentPortfolio";
import EmptyState from "../../../components/common/EmptyState";
import ErrorMessage from "../../../components/common/ErrorMessage";
import ConfirmDeleteModal from "../../../components/common/ConfirmDeleteModal";
import RemoveIconButton from "../../../components/common/RemoveIconButton";
import MutationErrorNotification from "../../../components/common/MutationErrorNotification";

const SOCIETY_ROLE_LABELS: Record<string, string> = {
  leader: "Leader",
  deputy_leader: "Deputy Leader",
  secretary: "Secretary",
  treasurer: "Treasurer",
  member: "Member",
};

function StudentSocietyMemberships({ studentId }: { studentId: string }) {
  const { data: memberships, isLoading, isError, refetch } = useStudentSocietyMemberships(studentId);

  if (isLoading) return null;

  if (isError) {
    return (
      <div style={{ marginBottom: "1.5rem" }}>
        <ErrorMessage message="Could not load society memberships." onRetry={refetch} />
      </div>
    );
  }

  if ((memberships?.length ?? 0) === 0) return null;

  return (
    <div style={{ marginBottom: "1.5rem" }}>
      <h3 style={{ fontSize: "0.75rem", fontWeight: 600, textTransform: "uppercase", color: "var(--os-text-tertiary)", margin: "0 0 0.5rem" }}>
        Society Memberships
      </h3>
      {memberships?.map((m) => (
        <div key={m.id} className="os-list-row os-list-row--compact" style={{ gap: "0.625rem" }}>
          <Tag size="sm" type="purple">{SOCIETY_ROLE_LABELS[m.role] ?? m.role}</Tag>
          <span style={{ fontWeight: 500, fontSize: "0.875rem" }}>{m.society_name}</span>
          <span style={{ fontSize: "0.8125rem", color: "var(--os-text-tertiary)" }}>{m.academic_year_label}</span>
        </div>
      ))}
    </div>
  );
}

export default function StudentActivities({ studentId }: { studentId: string }) {
  const { data: currentYear } = useCurrentAcademicYear();
  const { data: activities, isLoading } = useStudentActivities(studentId);
  const createActivity = useCreateStudentActivity(studentId);
  const deleteActivity = useDeleteStudentActivity(studentId);

  const [category, setCategory] = useState<ActivityCategory | "">("");
  const [name, setName] = useState("");
  const [role, setRole] = useState("");
  const [pendingDeleteId, setPendingDeleteId] = useState<string | null>(null);

  const handleAdd = () => {
    if (!category || !name.trim() || !currentYear) return;
    createActivity.mutate(
      {
        academic_year_id: currentYear.id,
        category,
        name: name.trim(),
        role: role.trim() || undefined,
      },
      { onSuccess: () => { setName(""); setRole(""); } },
    );
  };

  return (
    <div className="os-section" style={{ marginTop: "1rem" }}>
      <div className="os-section__header">
        <h2 className="os-section__title">Activities</h2>
        <span style={{ fontSize: "0.75rem", color: "var(--os-text-tertiary)" }}>Clubs, sports, societies &amp; competitions</span>
      </div>
      <div className="os-section__body">
        <StudentSocietyMemberships studentId={studentId} />

        <MutationErrorNotification
          isError={createActivity.isError}
          error={createActivity.error}
          fallback="Failed to add activity"
          style={{ marginBottom: "1rem" }}
        />

        <div style={{ display: "grid", gridTemplateColumns: "10rem 1fr 10rem auto", gap: "0.75rem", alignItems: "end", marginBottom: "1.5rem" }}>
          <Select id="activity-category" labelText="Category" value={category} onChange={(e) => setCategory(e.target.value as ActivityCategory)}>
            <SelectItem value="" text="Select…" />
            {ACTIVITY_CATEGORIES.map((c) => (
              <SelectItem key={c.value} value={c.value} text={c.label} />
            ))}
          </Select>
          <TextInput id="activity-name" labelText="Name" value={name} onChange={(e) => setName(e.target.value)} />
          <TextInput id="activity-role" labelText="Role (optional)" value={role} onChange={(e) => setRole(e.target.value)} />
          <Button renderIcon={Add} kind="primary" size="md" onClick={handleAdd} disabled={!category || !name.trim() || createActivity.isPending}>
            Add
          </Button>
        </div>

        {!isLoading && (activities?.length ?? 0) === 0 && (
          <EmptyState title="No activities yet" description="Add clubs, sports, societies, or competitions this student takes part in." />
        )}

        {activities?.map((a) => (
          <div key={a.id} className="os-list-row os-list-row--compact" style={{ justifyContent: "space-between" }}>
            <div style={{ display: "flex", alignItems: "center", gap: "0.625rem" }}>
              <Tag size="sm" type="gray">{ACTIVITY_CATEGORIES.find((c) => c.value === a.category)?.label ?? a.category}</Tag>
              <span style={{ fontWeight: 500, fontSize: "0.875rem" }}>{a.name}</span>
              {a.role && <span style={{ fontSize: "0.8125rem", color: "var(--os-text-tertiary)" }}>{a.role}</span>}
            </div>
            <RemoveIconButton label="Delete" onClick={() => setPendingDeleteId(a.id)} />
          </div>
        ))}
      </div>

      <ConfirmDeleteModal
        open={pendingDeleteId !== null}
        title="Delete activity"
        description="This will permanently remove this activity record. This action cannot be undone."
        isPending={deleteActivity.isPending}
        onClose={() => setPendingDeleteId(null)}
        onConfirm={() => {
          if (pendingDeleteId) deleteActivity.mutate(pendingDeleteId, { onSuccess: () => setPendingDeleteId(null) });
        }}
      />
    </div>
  );
}
