import { useMutation, useQuery } from "@tanstack/react-query";
import { timetableSettingsApi } from "@/features/timetable/api/timetableSettings";
import type { TimetableSettings } from "@/features/timetable/api/timetableSettings";
import { timetableKeys } from "@/features/timetable/keys";
import { useInvalidate } from "@/shared/api/useInvalidate";

export const useTimetableSettings = (academicYearId: string) =>
  useQuery({
    queryKey: timetableKeys.settings(academicYearId),
    queryFn: () => timetableSettingsApi.getByYear(academicYearId),
    enabled: !!academicYearId,
    retry: false,
  });

export const useUpsertTimetableSettings = () => {
  const invalidate = useInvalidate();
  return useMutation({
    mutationFn: (data: TimetableSettings) => timetableSettingsApi.upsert(data),
    onSuccess: (_r, variables) => invalidate(timetableKeys.settings(variables.academic_year_id)),
  });
};
