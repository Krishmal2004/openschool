import { useState } from "react";
import { useCurrentAcademicYear } from "@/features/school/queries/useAcademicYears";
import { useLevels, useLevelTree } from "@/features/curriculum/queries/useCurriculum";
import { useStudentWithClass } from "@/features/students/queries/useStudents";
import { useGrades } from "@/features/academics/queries/useGrades";
import { useStudentEnrollments, useSubmitEnrollments } from "@/features/students/queries/useEnrollments";

interface Group {
  id: string;
  min_select: number;
  max_select: number;
  subjects: { subject_id: string; medium_id?: string | null }[];
}

// A group where every subject must be taken behaves as a fixed block.
export function isCompulsory(g: Group) {
  return g.subjects.length > 0 && g.min_select === g.subjects.length && g.max_select === g.subjects.length;
}

// Level preselection, per-group subject picks and submission for one student's enrolment.
export function useEnrollmentPicker(studentId: string) {
  const { data: currentYear } = useCurrentAcademicYear();
  const academicYearId = currentYear?.id ?? "";
  const { data: levels } = useLevels();
  const { data: student } = useStudentWithClass(studentId);
  const { data: grades } = useGrades();
  const { data: enrollments } = useStudentEnrollments(studentId, academicYearId);

  const [selectedLevel, setSelectedLevel] = useState("");
  const [presetDone, setPresetDone] = useState(false);
  const [syncedFor, setSyncedFor] = useState<string | null>(null);
  const [selected, setSelected] = useState<Record<string, string[]>>({});

  const { data: tree, isLoading: treeLoading } = useLevelTree(selectedLevel);
  const submit = useSubmitEnrollments(studentId, academicYearId);

  // Preselect once: an existing enrolment wins, else the level for the student's grade.
  if (!presetDone && enrollments) {
    if (enrollments.length > 0) {
      setSelectedLevel(enrollments[0].level_id);
      setPresetDone(true);
    } else if (levels && grades && student) {
      const gradeId = grades.find((g) => g.name === student.grade_name)?.id;
      const autoLevelId = gradeId ? levels.find((l) => l.grade_id === gradeId)?.id : undefined;
      if (autoLevelId) setSelectedLevel(autoLevelId);
      setPresetDone(true);
    }
  }

  // Seed picks from saved enrolments when the tree for the chosen level arrives.
  if (enrollments && tree && tree.level.id === selectedLevel && syncedFor !== selectedLevel) {
    const map: Record<string, string[]> = {};
    for (const e of enrollments) if (e.level_id === selectedLevel) (map[e.group_id] ??= []).push(e.subject_id);
    for (const g of tree.groups) if (isCompulsory(g) && !map[g.id]?.length) map[g.id] = g.subjects.map((s) => s.subject_id);
    setSelected(map);
    setSyncedFor(selectedLevel);
  }

  const changeLevel = (id: string) => {
    submit.reset();
    setSyncedFor(null);
    setSelectedLevel(id);
    if (!id) setSelected({});
  };

  const toggle = (groupId: string, subjectId: string, checked: boolean) =>
    setSelected((prev) => {
      const cur = prev[groupId] ?? [];
      return { ...prev, [groupId]: checked ? [...cur, subjectId] : cur.filter((s) => s !== subjectId) };
    });

  const save = () => {
    if (!tree) return;
    const picks = tree.groups.flatMap((g) =>
      (selected[g.id] ?? []).map((subjectId) => ({
        group_id: g.id,
        subject_id: subjectId,
        medium_id: g.subjects.find((s) => s.subject_id === subjectId)?.medium_id ?? undefined,
      })),
    );
    submit.mutate({ academic_year_id: academicYearId, level_id: selectedLevel, picks });
  };

  return { currentYear, academicYearId, levels, selectedLevel, changeLevel, tree, treeLoading, selected, toggle, submit, save };
}
