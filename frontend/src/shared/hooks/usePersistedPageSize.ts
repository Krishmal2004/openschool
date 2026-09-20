import { useState } from "react";

function readStored(key: string, fallback: number): number {
  try {
    const raw = window.localStorage.getItem(key);
    const parsed = raw ? Number(raw) : NaN;
    return Number.isFinite(parsed) && parsed > 0 ? parsed : fallback;
  } catch {
    return fallback;
  }
}

export function usePersistedPageSize(key: string, fallback = 25) {
  const storageKey = `os-page-size:${key}`;
  const [pageSize, setPageSizeState] = useState(() => readStored(storageKey, fallback));

  const setPageSize = (next: number) => {
    setPageSizeState(next);
    try {
      window.localStorage.setItem(storageKey, String(next));
    } catch {
      void 0;
    }
  };

  return [pageSize, setPageSize] as const;
}
