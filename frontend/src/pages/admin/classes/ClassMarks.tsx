import { useState } from "react";
import { Select, SelectItem, NumberInput, Button, InlineNotification, Checkbox, Tag } from "@carbon/react";
import { Save } from "@carbon/icons-react";
import { useTerms, useCurrentTerm } from "../../../queries/useTerms";
import { useSubjects } from "../../../queries/useSubjects";
import { useClassSubjectTeachers } from "../../../queries/useClasses";
import { useClassMarks, useSaveClassMarks } from "../../../queries/useTermMarks";
import { getErrorMessage } from "../../../lib/errorMessage";
import LoadingSpinner from "../../../components/common/LoadingSpinner";
import EmptyState from "../../../components/common/EmptyState";
import AgentFindingsBanner from "../../../components/common/AgentFindingsBanner";
import EntityCombobox from "../../../components/common/EntityCombobox";

export default function ClassMarks({
  classId,
  academicYearId,
}: {
  classId: string;
  academicYearId: string;
}) {
  const { data: terms } = useTerms(academicYearId);
  const { data: currentTerm } = useCurrentTerm();
  const { data: subjects } = useSubjects();
  // Scope the subject picker to subjects actually assigned to this class
  // (class_subject_teachers) rather than every subject in the school.
  const { data: classSubjectTeachers } = useClassSubjectTeachers(classId);

  const [termId, setTermId] = useState("");
  const [termTouched, setTermTouched] = useState(false);
  const [subjectId, setSubjectId] = useState("");

  // Auto-select the current term once it loads, unless the admin has
  // already picked one manually (e.g. to review/enter a past term).
  const effectiveTermId = termTouched ? termId : termId || currentTerm?.id || "";

  const { data: rows, isLoading } = useClassMarks(classId, effectiveTermId, subjectId);
  const saveMarks = useSaveClassMarks(classId);

  // Draft edits keyed by student_id, seeded from the fetched grid whenever
  // the term/subject selection (i.e. a fresh set of rows) changes.
  const [draft, setDraft] = useState<Record<string, { marks: number; max_marks: number; is_absent: boolean }>>({});
  const [syncedFor, setSyncedFor] = useState("");
  const rowsKey = `${effectiveTermId}:${subjectId}`;

  const classSubjects = classSubjectTeachers?.map((st) => ({ id: st.subject_id, name: st.subject_name })) ?? [];
  const currentSubject = subjects?.find((s) => s.id === subjectId);
  const defaultMaxMarks = currentSubject?.max_marks ?? 100;

  if (rows && syncedFor !== rowsKey) {
    setDraft(
      Object.fromEntries(
        rows.map((r) => [
          r.student_id,
          { marks: r.marks ?? 0, max_marks: r.max_marks ?? defaultMaxMarks, is_absent: r.is_absent ?? false },
        ]),
      ),
    );
    setSyncedFor(rowsKey);
  }

  const handleSave = () => {
    if (!effectiveTermId || !subjectId) return;
    saveMarks.mutate({
      term_id: effectiveTermId,
      subject_id: subjectId,
      entries: Object.entries(draft).map(([student_id, v]) => ({ student_id, ...v })),
    });
  };

  return (
    <div>
      <AgentFindingsBanner titles={["Terms nearing deadline with no marks entered", "Terms falling behind on marks-entry pace"]} />
      <div className="os-section">
      <div className="os-section__header" style={{ flexWrap: "wrap", rowGap: "0.75rem" }}>
        <h2 className="os-section__title">Term Marks</h2>
        <div style={{ display: "flex", gap: "0.75rem", flexWrap: "wrap" }}>
          <Select
            id="marks-term"
            labelText=""
            size="sm"
            value={effectiveTermId}
            onChange={(e) => {
              setTermId(e.target.value);
              setTermTouched(true);
            }}
          >
            <SelectItem value="" text="Choose a term…" />
            {terms?.map((t) => (
              <SelectItem key={t.id} value={t.id} text={t.is_current ? `${t.name} (current)` : t.name} />
            ))}
          </Select>
          <div style={{ minWidth: "12rem" }}>
            <EntityCombobox
              id="marks-subject"
              items={classSubjects}
              selectedId={subjectId}
              onSelect={setSubjectId}
              getId={(s) => s.id}
              itemToString={(s) => s.name}
              placeholder="Choose a subject…"
            />
          </div>
        </div>
      </div>

      {!effectiveTermId || !subjectId ? (
        <EmptyState
          title="Pick a term and subject"
          description="Choose which term test and subject you're recording marks for."
        />
      ) : isLoading ? (
        <LoadingSpinner />
      ) : rows && rows.length > 0 ? (
        <>
          {saveMarks.isError && (
            <InlineNotification
              kind="error"
              title="Could not save marks"
              subtitle={getErrorMessage(saveMarks.error, "Please try again.")}
              lowContrast
              onClose={() => saveMarks.reset()}
              style={{ maxWidth: "100%", margin: "0 1.5rem 1rem" }}
            />
          )}
          {saveMarks.isSuccess && (
            <InlineNotification
              kind="success"
              title="Marks saved"
              lowContrast
              onClose={() => saveMarks.reset()}
              style={{ maxWidth: "100%", margin: "0 1.5rem 1rem" }}
            />
          )}
          <table className="os-table os-table--no-hover">
            <thead>
              <tr>
                <th>Index No.</th>
                <th>Student</th>
                <th style={{ width: "8rem" }}>Marks</th>
                <th style={{ width: "8rem" }}>Out of</th>
                <th style={{ width: "5rem" }}>Absent</th>
              </tr>
            </thead>
            <tbody>
              {rows.map((r) => {
                const isAbsent = draft[r.student_id]?.is_absent ?? false;
                return (
                <tr key={r.student_id}>
                  <td className="os-table__mono">{r.index_number}</td>
                  <td>{r.student_name}{isAbsent && <Tag type="red" size="sm" style={{ marginLeft: "0.5rem" }}>AB</Tag>}</td>
                  <td>
                    <NumberInput
                      id={`marks-${r.student_id}`}
                      label=""
                      hideLabel
                      size="sm"
                      min={0}
                      max={draft[r.student_id]?.max_marks ?? 100}
                      value={draft[r.student_id]?.marks ?? 0}
                      disabled={isAbsent}
                      onChange={(_e, { value }) =>
                        setDraft((d) => ({
                          ...d,
                          [r.student_id]: {
                            ...d[r.student_id],
                            marks: value === "" ? 0 : Number(value),
                          },
                        }))
                      }
                    />
                  </td>
                  <td>
                    <NumberInput
                      id={`max-marks-${r.student_id}`}
                      label=""
                      hideLabel
                      size="sm"
                      min={1}
                      value={draft[r.student_id]?.max_marks ?? 100}
                      onChange={(_e, { value }) =>
                        setDraft((d) => ({
                          ...d,
                          [r.student_id]: {
                            ...d[r.student_id],
                            max_marks: value === "" ? 100 : Number(value),
                          },
                        }))
                      }
                    />
                  </td>
                  <td>
                    <Checkbox
                      id={`absent-${r.student_id}`}
                      labelText="AB"
                      checked={isAbsent}
                      onChange={(_e, { checked }) =>
                        setDraft((d) => ({
                          ...d,
                          [r.student_id]: {
                            ...d[r.student_id],
                            marks: checked ? 0 : d[r.student_id]?.marks ?? 0,
                            is_absent: checked,
                          },
                        }))
                      }
                    />
                  </td>
                </tr>
                );
              })}
            </tbody>
          </table>
          <div style={{ padding: "1rem 1.5rem", borderTop: "1px solid #e0e0e0" }}>
            <Button
              renderIcon={Save}
              kind="primary"
              size="sm"
              onClick={handleSave}
              disabled={saveMarks.isPending}
            >
              {saveMarks.isPending ? "Saving…" : "Save Marks"}
            </Button>
          </div>
        </>
      ) : (
        <EmptyState
          title="No students enrolled"
          description="Enrol students into this class before recording marks."
        />
      )}
      </div>
    </div>
  );
}
