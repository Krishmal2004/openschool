import { useState } from "react";
import { Button, Select, SelectItem, TextArea, DatePicker, DatePickerInput, Tag } from "@carbon/react";
import { Add } from "@carbon/icons-react";
import { useCurrentAcademicYear } from "../../../queries/useAcademicYears";
import {
  useDisciplinaryRecords,
  useCreateDisciplinaryRecord,
  useDeleteDisciplinaryRecord,
} from "../../../queries/useStudentPortfolio";
import { DISCIPLINARY_SEVERITIES } from "../../../services/studentPortfolio";
import type { DisciplinarySeverity } from "../../../services/studentPortfolio";
import { todayISODate, toYmd } from "../../../lib/date";
import EmptyState from "../../../components/common/EmptyState";
import ConfirmDeleteModal from "../../../components/common/ConfirmDeleteModal";
import RemoveIconButton from "../../../components/common/RemoveIconButton";
import MutationErrorNotification from "../../../components/common/MutationErrorNotification";

const SEVERITY_TAG: Record<DisciplinarySeverity, "gray" | "warm-gray" | "red"> = {
  minor: "gray",
  major: "warm-gray",
  severe: "red",
};

export default function StudentDisciplinary({ studentId }: { studentId: string }) {
  const { data: currentYear } = useCurrentAcademicYear();
  const { data: records, isLoading } = useDisciplinaryRecords(studentId);
  const createRecord = useCreateDisciplinaryRecord(studentId);
  const deleteRecord = useDeleteDisciplinaryRecord(studentId);

  const [date, setDate] = useState(todayISODate());
  const [description, setDescription] = useState("");
  const [actionTaken, setActionTaken] = useState("");
  const [severity, setSeverity] = useState<DisciplinarySeverity | "">("");
  const [pendingDeleteId, setPendingDeleteId] = useState<string | null>(null);

  const handleAdd = () => {
    if (!description.trim() || !severity || !currentYear) return;
    createRecord.mutate(
      {
        academic_year_id: currentYear.id,
        incident_date: new Date(date).toISOString(),
        description: description.trim(),
        action_taken: actionTaken.trim() || undefined,
        severity,
      },
      { onSuccess: () => { setDescription(""); setActionTaken(""); setSeverity(""); } },
    );
  };

  return (
    <div className="os-section" style={{ marginTop: "1rem" }}>
      <div className="os-section__header">
        <h2 className="os-section__title">Disciplinary Records</h2>
      </div>
      <div className="os-section__body">
        <MutationErrorNotification
          isError={createRecord.isError}
          error={createRecord.error}
          fallback="Failed to add record"
          style={{ marginBottom: "1rem" }}
        />

        <div style={{ display: "grid", gridTemplateColumns: "10rem 10rem 1fr 1fr auto", gap: "0.75rem", alignItems: "end", marginBottom: "1.5rem" }}>
          <DatePicker datePickerType="single" dateFormat="Y-m-d" value={date} onChange={(dates) => {
            const ymd = toYmd(dates[0]);
            if (ymd) setDate(ymd);
          }}>
            <DatePickerInput id="disciplinary-date" labelText="Date" placeholder="YYYY-MM-DD" />
          </DatePicker>
          <Select id="disciplinary-severity" labelText="Severity" value={severity} onChange={(e) => setSeverity(e.target.value as DisciplinarySeverity)}>
            <SelectItem value="" text="Select…" />
            {DISCIPLINARY_SEVERITIES.map((s) => (
              <SelectItem key={s.value} value={s.value} text={s.label} />
            ))}
          </Select>
          <TextArea id="disciplinary-description" labelText="Description" rows={1} value={description} onChange={(e) => setDescription(e.target.value)} />
          <TextArea id="disciplinary-action" labelText="Action taken (optional)" rows={1} value={actionTaken} onChange={(e) => setActionTaken(e.target.value)} />
          <Button renderIcon={Add} kind="primary" size="md" onClick={handleAdd} disabled={!description.trim() || !severity || createRecord.isPending}>Add</Button>
        </div>

        {!isLoading && (records?.length ?? 0) === 0 && <EmptyState title="No disciplinary records" description="Nothing on file for this student." />}

        {records?.map((r) => (
          <div key={r.id} className="os-list-row" style={{ justifyContent: "space-between", alignItems: "flex-start" }}>
            <div>
              <div style={{ display: "flex", alignItems: "center", gap: "0.5rem", marginBottom: "0.25rem" }}>
                <Tag size="sm" type={SEVERITY_TAG[r.severity]}>{r.severity}</Tag>
                <span style={{ fontSize: "0.75rem", color: "var(--os-text-tertiary)" }}>{r.incident_date}</span>
              </div>
              <p style={{ margin: 0, fontSize: "0.875rem" }}>{r.description}</p>
              {r.action_taken && <p style={{ margin: "0.25rem 0 0", fontSize: "0.8125rem", color: "var(--os-text-secondary)" }}>Action: {r.action_taken}</p>}
            </div>
            <RemoveIconButton label="Delete" onClick={() => setPendingDeleteId(r.id)} />
          </div>
        ))}
      </div>

      <ConfirmDeleteModal
        open={pendingDeleteId !== null}
        title="Delete disciplinary record"
        description="This will permanently remove this record. This action cannot be undone."
        isPending={deleteRecord.isPending}
        onClose={() => setPendingDeleteId(null)}
        onConfirm={() => {
          if (pendingDeleteId) deleteRecord.mutate(pendingDeleteId, { onSuccess: () => setPendingDeleteId(null) });
        }}
      />
    </div>
  );
}
