// This file renders the TeacherMarks page: an overview of every subject/class a teacher teaches with a marks-entered status per term, and the entry grid for recording marks against a chosen class/subject/term.

import { useState, useMemo } from "react";
import { Select, SelectItem, NumberInput, Button, InlineNotification, Tag, ClickableTile, Checkbox } from "@carbon/react";
import { Save, ArrowLeft, CheckmarkFilled } from "@carbon/icons-react";
import { useMyTeacherProfile, useTeacherWorkload } from "../../queries/useTeachers";
import { useCurrentAcademicYear } from "../../queries/useAcademicYears";
import { useTerms, useCurrentTerm } from "../../queries/useTerms";
import { useSubjects } from "../../queries/useSubjects";
import { useClassMarks, useSaveClassMarks } from "../../queries/useTermMarks";
import { useClassStudents } from "../../queries/useClasses";
import { getErrorMessage } from "../../lib/errorMessage";
import LoadingSpinner from "../../components/common/LoadingSpinner";
import EmptyState from "../../components/common/EmptyState";

interface ClassCardProps {
  classId: string;
  className: string;
  gradeName: string;
  subjectId: string;
  termId: string;
  onOpen: () => void;
}

function ClassSubjectCard({ classId, className, gradeName, subjectId, termId, onOpen }: ClassCardProps) {
  const { data: rows, isLoading } = useClassMarks(classId, termId, subjectId);
  const total = rows?.length ?? 0;
  const entered = rows?.filter((r) => r.term_mark_id !== null).length ?? 0;
  const complete = total > 0 && entered === total;

  return (
    <ClickableTile onClick={onOpen} style={{ padding: "1rem 1.25rem" }}>
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "flex-start", gap: "0.5rem" }}>
        <div>
          <p style={{ margin: "0 0 0.25rem", fontWeight: 600, fontSize: "0.875rem" }}>{gradeName} — {className}</p>
          {isLoading ? (
            <p style={{ margin: 0, fontSize: "0.75rem", color: "#8d8d8d" }}>Loading…</p>
          ) : (
            <p style={{ margin: 0, fontSize: "0.75rem", color: complete ? "#24a148" : "#8d8d8d" }}>
              {entered}/{total} marks entered
            </p>
          )}
        </div>
        {complete && <CheckmarkFilled size={18} style={{ fill: "#24a148", flexShrink: 0 }} />}
      </div>
    </ClickableTile>
  );
}

export default function TeacherMarks() {
  const { data: teacher, isLoading: teacherLoading } = useMyTeacherProfile();
  const { data: workload, isLoading: workloadLoading } = useTeacherWorkload(teacher?.id ?? "");
  const { data: currentYear, isLoading: yearLoading } = useCurrentAcademicYear();
  const { data: terms } = useTerms(currentYear?.id ?? "");
  const { data: currentTerm } = useCurrentTerm();

  const [classId, setClassId] = useState("");
  const [subjectId, setSubjectId] = useState("");
  const [termId, setTermId] = useState("");
  const [termTouched, setTermTouched] = useState(false);
  const [mode, setMode] = useState<"overview" | "entry">("overview");

  // Auto-select the current term once it loads, unless the teacher has
  // already picked one manually (e.g. to review/enter a past term).
  const effectiveTermId = termTouched ? termId : termId || currentTerm?.id || "";

  const { data: students, isLoading: studentsLoading } = useClassStudents(classId);
  const { data: rows, isLoading: marksLoading } = useClassMarks(classId, effectiveTermId, subjectId);
  const { data: subjects } = useSubjects();
  const saveMarks = useSaveClassMarks(classId);

  const [draft, setDraft] = useState<Record<string, { marks: number; isAbsent: boolean }>>({});
  const [savedSnapshot, setSavedSnapshot] = useState<Record<string, { marks: number; isAbsent: boolean }>>({});
  const [syncedFor, setSyncedFor] = useState("");

  const uniqueClasses = useMemo(() => {
    const map = new Map<string, { id: string; name: string; grade_name: string }>();
    for (const r of workload ?? []) {
      if (r.academic_year_is_current) {
        map.set(r.class_id, { id: r.class_id, name: r.class_name, grade_name: r.grade_name });
      }
    }
    return [...map.values()];
  }, [workload]);

  const classSubjects = useMemo(() => {
    if (!classId) return [];
    const subjects = (workload ?? [])
      .filter((r) => r.academic_year_is_current && r.class_id === classId)
      .map((r) => ({ id: r.subject_id, name: r.subject_name }));
    const unique = new Map<string, typeof subjects[0]>();
    for (const s of subjects) {
      unique.set(s.id, s);
    }
    return [...unique.values()];
  }, [workload, classId]);

  // Overview grouping: subject -> the classes taught for that subject, for
  // the "My Subjects & Classes" landing view.
  const bySubject = useMemo(() => {
    const map = new Map<
      string,
      { subjectId: string; subjectName: string; classes: { id: string; name: string; gradeName: string }[] }
    >();
    for (const r of workload ?? []) {
      if (!r.academic_year_is_current) continue;
      if (!map.has(r.subject_id)) {
        map.set(r.subject_id, { subjectId: r.subject_id, subjectName: r.subject_name, classes: [] });
      }
      map.get(r.subject_id)!.classes.push({ id: r.class_id, name: r.class_name, gradeName: r.grade_name });
    }
    return [...map.values()];
  }, [workload]);

  const currentSubject = subjects?.find((s) => s.id === subjectId);
  const defaultMaxMarks = currentSubject?.max_marks ?? 100;

  const rowsKey = `${effectiveTermId}:${subjectId}:${classId}`;
  if (rows && syncedFor !== rowsKey) {
    const seeded = Object.fromEntries(
      rows.map((r) => [r.student_id, { marks: r.marks ?? 0, isAbsent: r.is_absent ?? false }]),
    );
    setDraft(seeded);
    setSavedSnapshot(seeded);
    setSyncedFor(rowsKey);
  }

  const isRowUnsaved = (studentId: string) => {
    const d = draft[studentId] ?? { marks: 0, isAbsent: false };
    const s = savedSnapshot[studentId] ?? { marks: 0, isAbsent: false };
    return d.marks !== s.marks || d.isAbsent !== s.isAbsent;
  };
  const hasUnsavedChanges = Object.keys(draft).some(isRowUnsaved);

  const openEntry = (nextClassId: string, nextSubjectId: string) => {
    setClassId(nextClassId);
    setSubjectId(nextSubjectId);
    setMode("entry");
  };

  const backToOverview = () => {
    setMode("overview");
    setClassId("");
    setSubjectId("");
  };

  const handleSave = () => {
    if (!effectiveTermId || !subjectId || !classId) return;
    const marksToSave = draft;
    saveMarks.mutate(
      {
        term_id: effectiveTermId,
        subject_id: subjectId,
        entries: Object.entries(draft).map(([student_id, v]) => ({
          student_id,
          marks: v.marks,
          max_marks: defaultMaxMarks,
          is_absent: v.isAbsent,
        })),
      },
      { onSuccess: () => setSavedSnapshot(marksToSave) },
    );
  };

  const isLoading = teacherLoading || workloadLoading || yearLoading;
  if (isLoading) return <LoadingSpinner />;

  const termSelector = (
    <Select
      id="marks-term"
      labelText="Term"
      size="sm"
      value={effectiveTermId}
      onChange={(e) => {
        setTermId(e.target.value);
        setTermTouched(true);
      }}
      style={{ minWidth: "12rem" }}
    >
      <SelectItem value="" text="Choose a term…" />
      {terms?.map((t) => (
        <SelectItem key={t.id} value={t.id} text={t.is_current ? `${t.name} (current)` : t.name} />
      ))}
    </Select>
  );

  if (mode === "overview") {
    return (
      <div className="os-page">
        <div className="os-page__header">
          <div className="os-page__header-left">
            <h1 className="os-page__title">My Subjects &amp; Classes</h1>
            <p className="os-page__subtitle">Every subject and class you teach, with marks status for the selected term</p>
          </div>
        </div>

        <div className="os-section" style={{ marginBottom: "1.5rem" }}>
          <div className="os-section__header">{termSelector}</div>
        </div>

        {!effectiveTermId ? (
          <EmptyState title="Choose a term" description="Select a term above to see marks status for your classes." />
        ) : bySubject.length === 0 ? (
          <EmptyState title="No subjects assigned yet" description="Classes you teach a subject in will appear here once an admin assigns you." />
        ) : (
          bySubject.map((subject) => (
            <div key={subject.subjectId} className="os-section" style={{ marginBottom: "1.5rem" }}>
              <div className="os-section__header">
                <h2 className="os-section__title">{subject.subjectName}</h2>
              </div>
              <div
                className="os-section__body"
                style={{ display: "grid", gridTemplateColumns: "repeat(auto-fill, minmax(16rem, 1fr))", gap: "0.875rem" }}
              >
                {subject.classes.map((c) => (
                  <ClassSubjectCard
                    key={c.id}
                    classId={c.id}
                    className={c.name}
                    gradeName={c.gradeName}
                    subjectId={subject.subjectId}
                    termId={effectiveTermId}
                    onOpen={() => openEntry(c.id, subject.subjectId)}
                  />
                ))}
              </div>
            </div>
          ))
        )}
      </div>
    );
  }

  return (
    <div className="os-page">
      <div className="os-page__header">
        <div className="os-page__header-left">
          <Button kind="ghost" size="sm" renderIcon={ArrowLeft} onClick={backToOverview} style={{ marginBottom: "0.5rem", paddingLeft: 0 }}>
            Back to overview
          </Button>
          <h1 className="os-page__title">Record Marks</h1>
          <p className="os-page__subtitle">Record term marks for classes and subjects you teach</p>
        </div>
      </div>

      <div className="os-section" style={{ marginBottom: "1.5rem" }}>
        <div className="os-section__header" style={{ flexWrap: "wrap", rowGap: "0.75rem" }}>
          <div style={{ display: "flex", gap: "0.75rem", flexWrap: "wrap", width: "100%" }}>
            <Select
              id="record-class"
              labelText="Class"
              size="sm"
              value={classId}
              onChange={(e) => {
                setClassId(e.target.value);
                setSubjectId("");
              }}
              style={{ minWidth: "12rem" }}
            >
              <SelectItem value="" text="Choose a class…" />
              {uniqueClasses.map((c) => (
                <SelectItem key={c.id} value={c.id} text={`${c.grade_name} — ${c.name}`} />
              ))}
            </Select>

            <Select
              id="record-subject"
              labelText="Subject"
              size="sm"
              value={subjectId}
              disabled={!classId}
              onChange={(e) => setSubjectId(e.target.value)}
              style={{ minWidth: "12rem" }}
            >
              <SelectItem value="" text="Choose a subject…" />
              {classSubjects.map((s) => (
                <SelectItem key={s.id} value={s.id} text={s.name} />
              ))}
            </Select>

            {termSelector}

            {classId && subjectId && effectiveTermId && (
              <div
                style={{
                  display: "flex",
                  alignItems: "center",
                  gap: "1rem",
                  marginLeft: "auto",
                  alignSelf: "flex-end",
                }}
              >
                <span style={{ fontSize: "0.8125rem", color: "#525252" }}>
                  Max Marks: <strong>{defaultMaxMarks}</strong>
                </span>
                <Button
                  renderIcon={Save}
                  size="sm"
                  onClick={handleSave}
                  disabled={
                    saveMarks.isPending ||
                    !students ||
                    students.length === 0 ||
                    !hasUnsavedChanges
                  }
                >
                  {saveMarks.isPending ? "Saving…" : "Save Marks"}
                </Button>
              </div>
            )}
          </div>
        </div>

        {!classId || !subjectId || !effectiveTermId ? (
          <EmptyState
            title="Pick class, subject, and term"
            description="Select options above to load student mark roster."
          />
        ) : studentsLoading || marksLoading ? (
          <LoadingSpinner />
        ) : students && students.length > 0 ? (
          <div className="os-section__body" style={{ padding: 0 }}>
            {saveMarks.isError && (
              <div style={{ padding: "1rem" }}>
                <InlineNotification
                  kind="error"
                  title="Could not save marks"
                  subtitle={getErrorMessage(saveMarks.error, "Please try again.")}
                  lowContrast
                  onClose={() => saveMarks.reset()}
                />
              </div>
            )}
            {saveMarks.isSuccess && (
              <div style={{ padding: "1rem" }}>
                <InlineNotification
                  kind="success"
                  title="Marks saved successfully"
                  lowContrast
                />
              </div>
            )}

            <table className="os-table">
              <thead>
                <tr>
                  <th style={{ width: "3rem" }}>#</th>
                  <th>Student Name</th>
                  <th>Index Number</th>
                  <th style={{ width: "8rem" }}>Marks</th>
                  <th style={{ width: "5rem" }}>Absent</th>
                  <th style={{ width: "7rem" }}>Status</th>
                </tr>
              </thead>
              <tbody>
                {students.map((student, i) => {
                  const row = draft[student.id] ?? { marks: 0, isAbsent: false };
                  const unsaved = isRowUnsaved(student.id);
                  return (
                    <tr key={student.id}>
                      <td className="os-table__muted">{i + 1}</td>
                      <td style={{ fontWeight: 500 }}>{student.full_name}</td>
                      <td className="os-table__mono">{student.index_number}</td>
                      <td>
                        <NumberInput
                          id={`marks-${student.id}`}
                          label=""
                          min={0}
                          max={defaultMaxMarks}
                          value={row.marks}
                          disabled={row.isAbsent}
                          onChange={(_e, { value }) => {
                            setDraft((prev) => ({
                              ...prev,
                              [student.id]: { ...row, marks: Number(value) },
                            }));
                          }}
                          size="sm"
                          hideSteppers
                        />
                      </td>
                      <td>
                        <Checkbox
                          id={`absent-${student.id}`}
                          labelText="AB"
                          checked={row.isAbsent}
                          onChange={(_e, { checked }) => {
                            setDraft((prev) => ({
                              ...prev,
                              [student.id]: { marks: checked ? 0 : row.marks, isAbsent: checked },
                            }));
                          }}
                        />
                      </td>
                      <td>
                        {row.isAbsent && (
                          <Tag type="red" size="sm">AB</Tag>
                        )}
                        {unsaved && (
                          <Tag type="high-contrast" size="sm">
                            Unsaved
                          </Tag>
                        )}
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        ) : (
          <EmptyState
            title="No students enrolled"
            description="There are no students enrolled in the selected class."
          />
        )}
      </div>
    </div>
  );
}
