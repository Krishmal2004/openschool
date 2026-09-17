import { useMutation } from "@tanstack/react-query";
import { generateApi } from "@/features/timetable/api/generate";
import type { GenerateTimetablesRequest } from "@/features/timetable/api/generate";
import { timetableKeys } from "@/features/timetable/keys";
import { useInvalidate } from "@/shared/api/useInvalidate";

export const useGenerateTimetables = () => {
  const invalidate = useInvalidate();
  return useMutation({ mutationFn: (data: GenerateTimetablesRequest) => generateApi.generate(data), onSuccess: () => invalidate(timetableKeys.all) });
};
