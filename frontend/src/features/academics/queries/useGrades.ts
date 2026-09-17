import { useMutation, useQuery } from "@tanstack/react-query";
import { gradeApi } from "@/features/academics/api/grade";
import type { Grade, CreateGradeRequest } from "@/features/academics/api/grade";
import { gradeKeys } from "@/features/academics/keys";
import { useInvalidate } from "@/shared/api/useInvalidate";

export const useGrades = () => useQuery({ queryKey: gradeKeys.all, queryFn: gradeApi.list });

export const useCreateGrade = () => {
  const invalidate = useInvalidate();
  return useMutation({ mutationFn: (data: CreateGradeRequest) => gradeApi.create(data), onSuccess: () => invalidate(gradeKeys.all) });
};

export const useUpdateGrade = () => {
  const invalidate = useInvalidate();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: CreateGradeRequest }) => gradeApi.update(id, data),
    onSuccess: () => invalidate(gradeKeys.all),
  });
};

export const useDeleteGrade = () => {
  const invalidate = useInvalidate();
  return useMutation({ mutationFn: (id: string) => gradeApi.remove(id), onSuccess: () => invalidate(gradeKeys.all) });
};

// Persists only the grades whose position changed.
export const useReorderGrades = () => {
  const invalidate = useInvalidate();
  return useMutation({
    mutationFn: async (ordered: Grade[]) => {
      const changed = ordered.map((g, index) => ({ g, index })).filter(({ g, index }) => g.sort_order !== index);
      await Promise.all(changed.map(({ g, index }) => gradeApi.update(g.id, { name: g.name, sort_order: index })));
      return changed.length;
    },
    onSuccess: () => invalidate(gradeKeys.all),
  });
};
