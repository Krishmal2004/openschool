import { useQuery } from "@tanstack/react-query";
import { searchApi } from "@/features/system/api/search";
import { systemKeys } from "@/features/system/keys";

// The caller debounces; this only fires from two characters.
export const useGlobalSearch = (q: string) =>
  useQuery({ queryKey: systemKeys.globalSearch(q), queryFn: () => searchApi.global(q), enabled: q.trim().length >= 2 });
