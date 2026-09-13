import { useState } from "react";
import { Button, TextInput, DatePicker, DatePickerInput } from "@carbon/react";
import { Add } from "@carbon/icons-react";
import { useCurrentAcademicYear } from "../../../queries/useAcademicYears";
import {
  useLeadershipRoles,
  useCreateLeadershipRole,
  useDeleteLeadershipRole,
  useStudentAwards,
  useCreateStudentAward,
  useDeleteStudentAward,
} from "../../../queries/useStudentPortfolio";
import { todayISODate, toYmd } from "../../../lib/date";
import EmptyState from "../../../components/common/EmptyState";
import ConfirmDeleteModal from "../../../components/common/ConfirmDeleteModal";
import RemoveIconButton from "../../../components/common/RemoveIconButton";
import MutationErrorNotification from "../../../components/common/MutationErrorNotification";

export default function StudentLeadershipAwards({ studentId }: { studentId: string }) {
  const { data: currentYear } = useCurrentAcademicYear();

  const { data: roles, isLoading: rolesLoading } = useLeadershipRoles(studentId);
  const createRole = useCreateLeadershipRole(studentId);
  const deleteRole = useDeleteLeadershipRole(studentId);
  const [title, setTitle] = useState("");
  const [scope, setScope] = useState("");
  const [pendingDeleteRoleId, setPendingDeleteRoleId] = useState<string | null>(null);

  const { data: awards, isLoading: awardsLoading } = useStudentAwards(studentId);
  const createAward = useCreateStudentAward(studentId);
  const deleteAward = useDeleteStudentAward(studentId);
  const [awardTitle, setAwardTitle] = useState("");
  const [awardDate, setAwardDate] = useState(todayISODate());
  const [pendingDeleteAwardId, setPendingDeleteAwardId] = useState<string | null>(null);

  const addRole = () => {
    if (!title.trim() || !currentYear) return;
    createRole.mutate(
      { academic_year_id: currentYear.id, title: title.trim(), scope: scope.trim() || undefined },
      { onSuccess: () => { setTitle(""); setScope(""); } },
    );
  };

  const addAward = () => {
    if (!awardTitle.trim() || !currentYear) return;
    createAward.mutate(
      { academic_year_id: currentYear.id, title: awardTitle.trim(), awarded_date: new Date(awardDate).toISOString() },
      { onSuccess: () => setAwardTitle("") },
    );
  };

  return (
    <>
      <div className="os-section" style={{ marginTop: "1rem" }}>
        <div className="os-section__header">
          <h2 className="os-section__title">Leadership Roles</h2>
        </div>
        <div className="os-section__body">
          <MutationErrorNotification
            isError={createRole.isError}
            error={createRole.error}
            fallback="Failed to add role"
            style={{ marginBottom: "1rem" }}
          />
          <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr auto", gap: "0.75rem", alignItems: "end", marginBottom: "1.5rem" }}>
            <TextInput id="leadership-title" labelText="Title" placeholder="e.g. Debate Club President" value={title} onChange={(e) => setTitle(e.target.value)} />
            <TextInput id="leadership-scope" labelText="Scope (optional)" value={scope} onChange={(e) => setScope(e.target.value)} />
            <Button renderIcon={Add} kind="primary" size="md" onClick={addRole} disabled={!title.trim() || createRole.isPending}>Add</Button>
          </div>
          {!rolesLoading && (roles?.length ?? 0) === 0 && <EmptyState title="No leadership roles yet" description="Roles outside the prefect board (club president, sports captain, etc.)." />}
          {roles?.map((r) => (
            <div key={r.id} style={{ display: "flex", alignItems: "center", justifyContent: "space-between", padding: "0.625rem 0", borderBottom: "1px solid #e0e0e0" }}>
              <div>
                <span style={{ fontWeight: 500, fontSize: "0.875rem" }}>{r.title}</span>
                {r.scope && <span style={{ marginLeft: "0.5rem", fontSize: "0.8125rem", color: "#8d8d8d" }}>{r.scope}</span>}
              </div>
              <RemoveIconButton label="Delete" onClick={() => setPendingDeleteRoleId(r.id)} />
            </div>
          ))}
        </div>
      </div>

      <ConfirmDeleteModal
        open={pendingDeleteRoleId !== null}
        title="Delete leadership role"
        description="This will permanently remove this leadership role. This action cannot be undone."
        isPending={deleteRole.isPending}
        onClose={() => setPendingDeleteRoleId(null)}
        onConfirm={() => {
          if (pendingDeleteRoleId) deleteRole.mutate(pendingDeleteRoleId, { onSuccess: () => setPendingDeleteRoleId(null) });
        }}
      />

      <div className="os-section">
        <div className="os-section__header">
          <h2 className="os-section__title">Awards &amp; Achievements</h2>
        </div>
        <div className="os-section__body">
          <MutationErrorNotification
            isError={createAward.isError}
            error={createAward.error}
            fallback="Failed to add award"
            style={{ marginBottom: "1rem" }}
          />
          <div style={{ display: "grid", gridTemplateColumns: "1fr 12rem auto", gap: "0.75rem", alignItems: "end", marginBottom: "1.5rem" }}>
            <TextInput id="award-title" labelText="Title" value={awardTitle} onChange={(e) => setAwardTitle(e.target.value)} />
            <DatePicker datePickerType="single" dateFormat="Y-m-d" value={awardDate} onChange={(dates) => {
              const ymd = toYmd(dates[0]);
              if (ymd) setAwardDate(ymd);
            }}>
              <DatePickerInput id="award-date" labelText="Date" placeholder="YYYY-MM-DD" />
            </DatePicker>
            <Button renderIcon={Add} kind="primary" size="md" onClick={addAward} disabled={!awardTitle.trim() || createAward.isPending}>Add</Button>
          </div>
          {!awardsLoading && (awards?.length ?? 0) === 0 && <EmptyState title="No awards yet" description="Add awards and achievements for this student." />}
          {awards?.map((a) => (
            <div key={a.id} style={{ display: "flex", alignItems: "center", justifyContent: "space-between", padding: "0.625rem 0", borderBottom: "1px solid #e0e0e0" }}>
              <div>
                <span style={{ fontWeight: 500, fontSize: "0.875rem" }}>{a.title}</span>
                <span style={{ marginLeft: "0.5rem", fontSize: "0.8125rem", color: "#8d8d8d" }}>{a.awarded_date}</span>
              </div>
              <RemoveIconButton label="Delete" onClick={() => setPendingDeleteAwardId(a.id)} />
            </div>
          ))}
        </div>
      </div>

      <ConfirmDeleteModal
        open={pendingDeleteAwardId !== null}
        title="Delete award"
        description="This will permanently remove this award. This action cannot be undone."
        isPending={deleteAward.isPending}
        onClose={() => setPendingDeleteAwardId(null)}
        onConfirm={() => {
          if (pendingDeleteAwardId) deleteAward.mutate(pendingDeleteAwardId, { onSuccess: () => setPendingDeleteAwardId(null) });
        }}
      />
    </>
  );
}
