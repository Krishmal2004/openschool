import { useMutation, useQuery } from "@tanstack/react-query";
import { timetableApi } from "@/features/timetable/api/timetable";
import type { TimetableEntryInput } from "@/features/timetable/api/timetable";
import { timetableKeys } from "@/features/timetable/keys";
import { useInvalidate } from "@/shared/api/useInvalidate";

export const useTimetablesByYear = (academicYearId: string) =>
  useQuery({ queryKey: timetableKeys.byYear(academicYearId), queryFn: () => timetableApi.listByYear(academicYearId), enabled: !!academicYearId });

// 404 means "not yet scheduled", so callers treat isError as empty.
export const usePublishedTimetableForClass = (classId: string, academicYearId: string) =>
  useQuery({
    queryKey: timetableKeys.publishedForClass(classId, academicYearId),
    queryFn: () => timetableApi.publishedForClass(classId, academicYearId),
    enabled: !!classId && !!academicYearId,
    retry: false,
  });

// Every class's timetable for the current year, for the Principal / Vice Principal.
export const useTimetablesForLeadership = () =>
  useQuery({ queryKey: timetableKeys.leadership(), queryFn: () => timetableApi.listForLeadership() });

export const useTimetable = (id: string) =>
  useQuery({ queryKey: timetableKeys.detail(id), queryFn: () => timetableApi.get(id), enabled: !!id });

export const useTimetableEntries = (id: string) =>
  useQuery({ queryKey: timetableKeys.entries(id), queryFn: () => timetableApi.getEntries(id), enabled: !!id });

export const useTimetableValidation = (id: string) =>
  useQuery({ queryKey: timetableKeys.validation(id), queryFn: () => timetableApi.validate(id), enabled: !!id });

export const useTimetableStatusHistory = (id: string) =>
  useQuery({ queryKey: timetableKeys.history(id), queryFn: () => timetableApi.statusHistory(id), enabled: !!id });

export const useReviewQueue = (academicYearId: string, enabled = true) =>
  useQuery({
    queryKey: timetableKeys.reviewQueue(academicYearId),
    queryFn: () => timetableApi.reviewQueue(academicYearId),
    enabled: !!academicYearId && enabled,
  });

export const useMyTeacherSchedule = (academicYearId: string) =>
  useQuery({ queryKey: timetableKeys.myTeacherSchedule(academicYearId), queryFn: () => timetableApi.myTeacherSchedule(academicYearId), enabled: !!academicYearId });

export const useMyClassTimetable = () => useQuery({ queryKey: timetableKeys.myClass(), queryFn: timetableApi.myClassTimetable, retry: false });

export const useChildTimetable = (studentId: string) =>
  useQuery({ queryKey: timetableKeys.child(studentId), queryFn: () => timetableApi.childTimetable(studentId), enabled: !!studentId, retry: false });

// Every write touches lists, review queues and the detail views, so invalidate the whole feature.
const useTimetableMutation = <TVariables = void, TResult = unknown>(mutationFn: (variables: TVariables) => Promise<TResult>) => {
  const invalidate = useInvalidate();
  return useMutation({ mutationFn, onSuccess: () => invalidate(timetableKeys.all) });
};

export const useCreateTimetable = () =>
  useTimetableMutation((data: { academic_year_id: string; class_id: string }) => timetableApi.create(data));

export const useCopyTimetable = () =>
  useTimetableMutation((data: { academic_year_id: string; class_id: string; source_timetable_id: string }) => timetableApi.copyFrom(data));

export const useReviseTimetable = () => useTimetableMutation((publishedId: string) => timetableApi.revise(publishedId));

export const useDeleteTimetable = () => useTimetableMutation((id: string) => timetableApi.remove(id));

export const useSaveTimetableEntries = (id: string) =>
  useTimetableMutation((entries: TimetableEntryInput[]) => timetableApi.saveEntries(id, entries));

export const useDeleteTimetableEntry = (id: string) =>
  useTimetableMutation(({ day, period }: { day: number; period: number }) => timetableApi.deleteEntry(id, day, period));

export const useSubmitTimetable = (id: string) => useTimetableMutation(() => timetableApi.submit(id));

export const usePublishTimetable = (id: string) => useTimetableMutation(() => timetableApi.publish(id));

export const useApproveTimetable = (id: string) => useTimetableMutation((comment?: string) => timetableApi.approve(id, comment));

export const useRejectTimetable = (id: string) => useTimetableMutation((comment: string) => timetableApi.reject(id, comment));
