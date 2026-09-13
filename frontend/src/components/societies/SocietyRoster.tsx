import { useState } from "react";
import { Add } from "@carbon/icons-react";
import {
  Button,
  Select,
  SelectItem,
  ComposedModal,
  ModalHeader,
  ModalBody,
  ModalFooter,
  SkeletonText,
} from "@carbon/react";
import { Link } from "react-router";
import { useStudents } from "../../queries/useStudents";
import { useSocietyMembers, useAssignSocietyMember, useRemoveSocietyMember } from "../../queries/useSocieties";
import EmptyState from "../common/EmptyState";
import ErrorMessage from "../common/ErrorMessage";
import EntityCombobox from "../common/EntityCombobox";
import ConfirmDeleteModal from "../common/ConfirmDeleteModal";
import MutationErrorNotification from "../common/MutationErrorNotification";
import RemoveIconButton from "../common/RemoveIconButton";
import type { SocietyMember, SocietyRole } from "../../services/society";

const ROLES: { value: SocietyRole; label: string }[] = [
  { value: "leader", label: "Leaders" },
  { value: "deputy_leader", label: "Deputy Leaders" },
  { value: "secretary", label: "Secretaries" },
  { value: "treasurer", label: "Treasurers" },
  { value: "member", label: "Members" },
];

interface Props {
  societyId: string;
  readOnly?: boolean;
}

// A society's roster: add/remove students and assign one of the five roles.
// Shared by the admin Societies page (managing any society) and the
// teacher-portal My Society page (managing only the caller's own) — the
// backend enforces who's allowed to call the mutations either way.
export default function SocietyRoster({ societyId, readOnly }: Props) {
  const { data: members, isLoading, isError, refetch } = useSocietyMembers(societyId);
  const { data: students } = useStudents();
  const assignMember = useAssignSocietyMember(societyId);
  const removeMember = useRemoveSocietyMember(societyId);

  const [assignOpen, setAssignOpen] = useState(false);
  const [studentChoice, setStudentChoice] = useState("");
  const [roleChoice, setRoleChoice] = useState<SocietyRole>("member");
  const [memberToRemove, setMemberToRemove] = useState<SocietyMember | null>(null);

  const openAssign = () => {
    assignMember.reset();
    setStudentChoice("");
    setRoleChoice("member");
    setAssignOpen(true);
  };

  const handleAssign = () => {
    if (!studentChoice) return;
    assignMember.mutate(
      { student_id: studentChoice, role: roleChoice },
      { onSuccess: () => setAssignOpen(false) },
    );
  };

  const confirmRemove = () => {
    if (!memberToRemove) return;
    removeMember.mutate(
      { memberId: memberToRemove.id, studentId: memberToRemove.student_id },
      { onSettled: () => setMemberToRemove(null) },
    );
  };

  const byRole = (role: SocietyRole) => (members ?? []).filter((m) => m.role === role);
  const memberStudentIds = new Set((members ?? []).map((m) => m.student_id));
  const availableStudents = (students ?? []).filter((s) => !memberStudentIds.has(s.id));

  if (isError) {
    return <ErrorMessage message="Could not load the roster." onRetry={refetch} />;
  }

  return (
    <div>
      <div style={{ display: "flex", justifyContent: "flex-end", marginBottom: "1rem" }}>
        {!readOnly && (
          <Button renderIcon={Add} kind="primary" size="sm" onClick={openAssign}>
            Add Member
          </Button>
        )}
      </div>

      <MutationErrorNotification
        isError={removeMember.isError}
        error={removeMember.error}
        title="Could not remove member"
        fallback="Please try again."
        onClose={() => removeMember.reset()}
      />

      {isLoading ? (
        <SkeletonText width="40%" />
      ) : (members ?? []).length === 0 ? (
        <EmptyState title="No members yet" description="Add a student to this society's roster." />
      ) : (
        ROLES.map(({ value, label }) =>
          byRole(value).length === 0 ? null : (
            <div key={value} style={{ marginBottom: "1rem" }}>
              <div style={{ display: "flex", alignItems: "center", gap: "0.5rem", marginBottom: "0.25rem" }}>
                <h3 style={{ fontSize: "0.75rem", fontWeight: 600, textTransform: "uppercase", color: "var(--os-text-tertiary)", margin: 0 }}>
                  {label}
                </h3>
                <span style={{ fontSize: "0.75rem", color: "var(--os-text-tertiary)" }}>{byRole(value).length}</span>
              </div>
              {byRole(value).map((m) => (
                <div key={m.id} className="os-list-row os-list-row--compact">
                  <div style={{ flex: 1, minWidth: 0 }}>
                    <Link to={`/students/${m.student_id}`} className="os-table__link" style={{ fontSize: "0.875rem", fontWeight: 500 }}>
                      {m.student_name}
                    </Link>
                    <p style={{ margin: "0.1rem 0 0", fontSize: "0.75rem", color: "var(--os-text-secondary)" }}>
                      {[m.grade_name, m.student_index].filter(Boolean).join(" · ")}
                    </p>
                  </div>
                  {!readOnly && (
                    <RemoveIconButton disabled={removeMember.isPending} onClick={() => setMemberToRemove(m)} />
                  )}
                </div>
              ))}
            </div>
          ),
        )
      )}

      <ComposedModal open={assignOpen} size="sm" onClose={() => setAssignOpen(false)}>
        <ModalHeader title="Add society member" />
        <ModalBody>
          <MutationErrorNotification
            isError={assignMember.isError}
            error={assignMember.error}
            fallback="Failed to add member"
          />
          <div style={{ display: "grid", gap: "1rem" }}>
            <EntityCombobox
              id="society-member-student"
              labelText="Student"
              items={availableStudents}
              selectedId={studentChoice}
              onSelect={setStudentChoice}
              getId={(s) => s.id}
              itemToString={(s) => `${s.full_name} — ${s.index_number}`}
              placeholder="Search students by name or index number…"
            />
            <Select
              id="society-member-role"
              labelText="Role"
              value={roleChoice}
              onChange={(e) => setRoleChoice(e.target.value as SocietyRole)}
            >
              {ROLES.map((r) => (
                <SelectItem key={r.value} value={r.value} text={r.label.replace(/s$/, "")} />
              ))}
            </Select>
          </div>
        </ModalBody>
        <ModalFooter>
          <Button kind="secondary" onClick={() => setAssignOpen(false)}>
            Cancel
          </Button>
          <Button kind="primary" onClick={handleAssign} disabled={!studentChoice || assignMember.isPending}>
            {assignMember.isPending ? "Saving…" : "Add"}
          </Button>
        </ModalFooter>
      </ComposedModal>

      <ConfirmDeleteModal
        open={!!memberToRemove}
        title="Remove member"
        description={
          <>
            Remove <strong>{memberToRemove?.student_name}</strong> from this society?
          </>
        }
        confirmLabel="Remove"
        pendingLabel="Removing…"
        isPending={removeMember.isPending}
        onClose={() => setMemberToRemove(null)}
        onConfirm={confirmRemove}
      />
    </div>
  );
}
