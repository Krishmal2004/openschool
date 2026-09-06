// This file defines query and mutation hooks for retrieving and updating teacher profiles, workloads, assigned subjects, and form classes.

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { teacherApi } from "../services/teacher";
import { useCurrentClasses } from "./useClasses";
import type {
  CreateTeacherRequest,
  UpdateTeacherRequest,
  TeacherEmploymentStatus,
} from "../services/teacher";

export const TEACHERS_KEY = ["teachers"];
export const teacherKey = (id: string) => ["teachers", id];
export const teacherSubjectsKey = (id: string) => ["teachers", id, "subjects"];
export const teacherWorkloadKey = (id: string) => ["teachers", id, "workload"];
export const MY_TEACHER_PROFILE_KEY = ["me", "teacher"];

export const useTeachers = () =>
  useQuery({
    queryKey: TEACHERS_KEY,
    queryFn: teacherApi.list,
  });

export const useTeacher = (id: string) =>
  useQuery({
    queryKey: teacherKey(id),
    queryFn: () => teacherApi.get(id),
    enabled: !!id,
  });

export const useTeacherSubjects = (id: string) =>
  useQuery({
    queryKey: teacherSubjectsKey(id),
    queryFn: () => teacherApi.listSubjects(id),
    enabled: !!id,
  });

export const useMyTeacherProfile = () =>
  useQuery({
    queryKey: MY_TEACHER_PROFILE_KEY,
    queryFn: teacherApi.me,
  });

export const teachersBySubjectKey = (subjectId: string) => ["teachers", "by-subject", subjectId];

// Teachers qualified (teacher_subjects) for a given subject — used to scope
// the class-subject-teacher assignment picker.
export const useTeachersBySubject = (subjectId: string) =>
  useQuery({
    queryKey: teachersBySubjectKey(subjectId),
    queryFn: () => teacherApi.listBySubject(subjectId),
    enabled: !!subjectId,
  });

export const useTeacherWorkload = (id: string) =>
  useQuery({
    queryKey: teacherWorkloadKey(id),
    queryFn: () => teacherApi.workload(id),
    enabled: !!id,
  });

export interface MyClass {
  class_id: string;
  class_name: string;
  grade_name: string;
  subjects: string[];
  isFormTeacher: boolean;
}

// Every class a teacher is actually involved with this year: classes they
// are the form teacher of, UNION classes where they teach a subject
// (class_subject_teachers, via workload) — previously this only looked at
// form-teacher classes and always reported an empty subjects list, so a
// subject-only teacher never saw their classes anywhere in the portal.
export const useMyClasses = () => {
  const teacher = useMyTeacherProfile();
  const teacherId = teacher.data?.id ?? "";
  const { data: allClasses, isLoading: classesLoading, isError: classesError } = useCurrentClasses();
  const { data: workload, isLoading: workloadLoading, isError: workloadError } = useTeacherWorkload(teacherId);

  const classMap = new Map<string, MyClass>();

  for (const c of allClasses ?? []) {
    if (c.form_teacher_id === teacherId) {
      classMap.set(c.id, {
        class_id: c.id,
        class_name: c.name,
        grade_name: c.grade_name,
        subjects: [],
        isFormTeacher: true,
      });
    }
  }

  for (const w of workload ?? []) {
    if (!w.academic_year_is_current) continue;
    const existing = classMap.get(w.class_id);
    if (existing) {
      if (!existing.subjects.includes(w.subject_name)) existing.subjects.push(w.subject_name);
    } else {
      classMap.set(w.class_id, {
        class_id: w.class_id,
        class_name: w.class_name,
        grade_name: w.grade_name,
        subjects: [w.subject_name],
        isFormTeacher: false,
      });
    }
  }

  return {
    teacher: teacher.data,
    classes: [...classMap.values()],
    isLoading: teacher.isLoading || classesLoading || workloadLoading,
    isError: teacher.isError || classesError || workloadError,
    refetch: teacher.refetch,
  };
};

export const useCreateTeacher = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateTeacherRequest) => teacherApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: TEACHERS_KEY });
    },
  });
};

export const useUpdateTeacher = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateTeacherRequest }) =>
      teacherApi.update(id, data),
    onSuccess: (_data, { id }) => {
      queryClient.invalidateQueries({ queryKey: TEACHERS_KEY });
      queryClient.invalidateQueries({ queryKey: teacherKey(id) });
    },
  });
};

export const useUpdateTeacherEmploymentStatus = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, status }: { id: string; status: TeacherEmploymentStatus }) =>
      teacherApi.updateEmploymentStatus(id, status),
    onSuccess: (_data, { id }) => {
      queryClient.invalidateQueries({ queryKey: TEACHERS_KEY });
      queryClient.invalidateQueries({ queryKey: teacherKey(id) });
    },
  });
};

export const useUpdateTeacherHouse = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, houseId }: { id: string; houseId: string }) =>
      teacherApi.updateHouse(id, houseId),
    onSuccess: (_data, { id }) => {
      queryClient.invalidateQueries({ queryKey: TEACHERS_KEY });
      queryClient.invalidateQueries({ queryKey: teacherKey(id) });
    },
  });
};

export const useDeleteTeacher = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => teacherApi.remove(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: TEACHERS_KEY });
    },
  });
};

export const useAssignTeacherSubject = (teacherId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (subjectId: string) => teacherApi.assignSubject(teacherId, subjectId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: teacherSubjectsKey(teacherId) });
    },
  });
};

export const useRemoveTeacherSubject = (teacherId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (subjectId: string) => teacherApi.removeSubject(teacherId, subjectId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: teacherSubjectsKey(teacherId) });
    },
  });
};

