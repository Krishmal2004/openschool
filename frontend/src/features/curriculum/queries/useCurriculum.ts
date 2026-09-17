import { useMutation, useQuery } from "@tanstack/react-query";
import { mediumApi, levelApi, groupApi } from "@/features/curriculum/api/curriculum";
import type {
  CreateMediumRequest,
  CreateLevelRequest,
  CreateSelectionGroupRequest,
  AddGroupSubjectRequest,
} from "@/features/curriculum/api/curriculum";
import { curriculumKeys } from "@/features/curriculum/keys";
import { useInvalidate } from "@/shared/api/useInvalidate";

export const useMediums = () => useQuery({ queryKey: curriculumKeys.mediums(), queryFn: mediumApi.list });

export const useCreateMedium = () => {
  const invalidate = useInvalidate();
  return useMutation({ mutationFn: (data: CreateMediumRequest) => mediumApi.create(data), onSuccess: () => invalidate(curriculumKeys.mediums()) });
};

export const useUpdateMedium = () => {
  const invalidate = useInvalidate();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: CreateMediumRequest }) => mediumApi.update(id, data),
    onSuccess: () => invalidate(curriculumKeys.mediums()),
  });
};

export const useDeleteMedium = () => {
  const invalidate = useInvalidate();
  return useMutation({ mutationFn: (id: string) => mediumApi.remove(id), onSuccess: () => invalidate(curriculumKeys.mediums()) });
};

export const useLevels = () => useQuery({ queryKey: curriculumKeys.levels(), queryFn: levelApi.list });

export const useLevelTree = (id: string) =>
  useQuery({ queryKey: curriculumKeys.levelTree(id), queryFn: () => levelApi.tree(id), enabled: !!id });

export const useCreateLevel = () => {
  const invalidate = useInvalidate();
  return useMutation({ mutationFn: (data: CreateLevelRequest) => levelApi.create(data), onSuccess: () => invalidate(curriculumKeys.levels()) });
};

export const useUpdateLevel = () => {
  const invalidate = useInvalidate();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: CreateLevelRequest }) => levelApi.update(id, data),
    onSuccess: () => invalidate(curriculumKeys.levels()),
  });
};

export const useDuplicateLevel = () => {
  const invalidate = useInvalidate();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: CreateLevelRequest }) => levelApi.duplicate(id, data),
    onSuccess: () => invalidate(curriculumKeys.levels()),
  });
};

export const useDeleteLevel = () => {
  const invalidate = useInvalidate();
  return useMutation({ mutationFn: (id: string) => levelApi.remove(id), onSuccess: () => invalidate(curriculumKeys.levels()) });
};

export const useCreateSelectionGroup = (levelId: string) => {
  const invalidate = useInvalidate();
  return useMutation({
    mutationFn: (data: CreateSelectionGroupRequest) => groupApi.create(levelId, data),
    onSuccess: () => invalidate(curriculumKeys.levelTree(levelId)),
  });
};

export const useUpdateSelectionGroup = (levelId: string) => {
  const invalidate = useInvalidate();
  return useMutation({
    mutationFn: ({ groupId, data }: { groupId: string; data: CreateSelectionGroupRequest }) => groupApi.update(groupId, data),
    onSuccess: () => invalidate(curriculumKeys.levelTree(levelId)),
  });
};

export const useDeleteSelectionGroup = (levelId: string) => {
  const invalidate = useInvalidate();
  return useMutation({ mutationFn: (groupId: string) => groupApi.remove(groupId), onSuccess: () => invalidate(curriculumKeys.levelTree(levelId)) });
};

export const useAddGroupSubject = (levelId: string) => {
  const invalidate = useInvalidate();
  return useMutation({
    mutationFn: ({ groupId, data }: { groupId: string; data: AddGroupSubjectRequest }) => groupApi.addSubject(groupId, data),
    onSuccess: (_data, { groupId }) => invalidate(curriculumKeys.levelTree(levelId), curriculumKeys.groupSubjects(groupId)),
  });
};

export const useRemoveGroupSubject = (levelId: string) => {
  const invalidate = useInvalidate();
  return useMutation({
    mutationFn: ({ groupId, subjectId }: { groupId: string; subjectId: string }) => groupApi.removeSubject(groupId, subjectId),
    onSuccess: (_data, { groupId }) => invalidate(curriculumKeys.levelTree(levelId), curriculumKeys.groupSubjects(groupId)),
  });
};
