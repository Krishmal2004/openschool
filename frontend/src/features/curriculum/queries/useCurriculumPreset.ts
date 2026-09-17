import { useMutation } from "@tanstack/react-query";
import { curriculumPresetApi } from "@/features/curriculum/api/curriculumPreset";
import { curriculumKeys } from "@/features/curriculum/keys";
import { useInvalidate } from "@/shared/api/useInvalidate";

export const useRunCurriculumPreset = () => {
  const invalidate = useInvalidate();
  return useMutation({ mutationFn: curriculumPresetApi.run, onSuccess: () => invalidate(curriculumKeys.all) });
};
