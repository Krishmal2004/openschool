import { useCallback, useMemo, useState } from "react";
import { useDebounced } from "@/shared/hooks/useDebounced";

type Filters = Record<string, string>;

// Holds a page's search and filter values in one place, with a debounced search for large lists.
export function useListFilters<T extends Filters>(
  initial: T,
  options: { searchKey?: keyof T; debounceMs?: number } = {},
) {
  // Captured once so clear() and activeKeys compare against a stable baseline.
  const [baseline] = useState(initial);
  const [filters, setFilters] = useState<T>(initial);
  const searchKey = options.searchKey ?? ("query" as keyof T);
  const debouncedSearch = useDebounced(filters[searchKey] ?? "", options.debounceMs ?? 250);

  const set = useCallback(<K extends keyof T>(key: K, value: T[K]) => {
    setFilters((current) => ({ ...current, [key]: value }));
  }, []);

  const clear = useCallback((key?: keyof T) => {
    setFilters((current) => (key ? { ...current, [key]: baseline[key] } : baseline));
  }, [baseline]);

  const activeKeys = useMemo(
    () => (Object.keys(filters) as (keyof T)[]).filter((k) => filters[k] !== baseline[k]),
    [filters, baseline],
  );

  return { filters, set, clear, activeKeys, hasActive: activeKeys.length > 0, debouncedSearch };
}
