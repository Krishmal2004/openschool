import { useState } from "react";
import { useNavigate } from "react-router";
import { useCreateClass, useStreamGroups } from "@/features/academics/queries/useClasses";
import { useClassrooms, useCreateClassroom } from "@/features/timetable/queries/useClassrooms";
import { useAcademicYears } from "@/features/school/queries/useAcademicYears";
import { getErrorMessage } from "@/shared/api/errors";
import { suggestHomeClassroom } from "@/features/timetable/lib/classroom";

const EMPTY_CLASS_FORM = {
  grade_id: "",
  academic_year_id: "",
  name: "",
  stream_id: "",
  stream_group_id: "",
  form_teacher_id: "",
  medium_id: "",
  home_classroom_id: "",
};

export type ClassForm = typeof EMPTY_CLASS_FORM;
type Touched = Partial<Record<"grade" | "name" | "year", boolean>>;

// Form state, validation and the two-step save (auto-create a homeroom, then the class).
export function useCreateClassForm(preselectedGradeId: string) {
  const navigate = useNavigate();
  const { data: years } = useAcademicYears();
  const { data: classrooms } = useClassrooms();
  const createClass = useCreateClass();
  const createClassroom = useCreateClassroom();

  const [form, setForm] = useState<ClassForm>({ ...EMPTY_CLASS_FORM, grade_id: preselectedGradeId });
  const [touched, setTouched] = useState<Touched>({});
  const [roomError, setRoomError] = useState<string | null>(null);
  const { data: streamGroups } = useStreamGroups(form.stream_id);

  const suggestedHomeClassroom = suggestHomeClassroom(classrooms, form.name);
  const effectiveHomeClassroomId = form.home_classroom_id || suggestedHomeClassroom?.id || "";
  const academicYearId = form.academic_year_id || years?.find((y) => y.is_current)?.id || "";
  const isValid = !!form.grade_id && !!academicYearId && !!form.name.trim();

  const set = (field: keyof ClassForm, value: string) =>
    setForm((f) => (field === "stream_id" ? { ...f, stream_id: value, stream_group_id: "" } : { ...f, [field]: value }));

  const save = async () => {
    setTouched({ grade: true, name: true, year: true });
    if (!isValid) return;
    setRoomError(null);

    // With no room chosen, a homeroom named after the class is created first.
    let homeClassroomId = effectiveHomeClassroomId;
    if (!homeClassroomId) {
      try {
        homeClassroomId = (await createClassroom.mutateAsync({ name: form.name.trim(), room_type: "regular" })).id;
      } catch (e) {
        setRoomError(getErrorMessage(e, "Failed to create this class's home classroom"));
        return;
      }
    }

    createClass.mutate(
      {
        grade_id: form.grade_id,
        academic_year_id: academicYearId,
        name: form.name.trim(),
        stream_id: form.stream_id || null,
        stream_group_id: form.stream_group_id || null,
        form_teacher_id: form.form_teacher_id || null,
        medium_id: form.medium_id || null,
        home_classroom_id: homeClassroomId || null,
      },
      { onSuccess: () => navigate("/classes") },
    );
  };

  return {
    form,
    set,
    touched,
    markTouched: (field: keyof Touched) => setTouched((t) => ({ ...t, [field]: true })),
    years,
    classrooms,
    streamGroups,
    suggestedHomeClassroom,
    effectiveHomeClassroomId,
    academicYearId,
    isValid,
    isSaving: createClassroom.isPending || createClass.isPending,
    error: roomError ?? (createClass.isError ? getErrorMessage(createClass.error, "Failed to create class") : null),
    clearError: () => { setRoomError(null); createClass.reset(); },
    save,
  };
}
