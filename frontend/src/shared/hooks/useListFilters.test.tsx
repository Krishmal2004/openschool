import { describe, expect, it, vi } from "vitest";
import { act, renderHook } from "@testing-library/react";
import { useListFilters } from "@/shared/hooks/useListFilters";

describe("useListFilters", () => {
  it("tracks active keys and clears one or all", () => {
    const { result } = renderHook(() => useListFilters({ query: "", grade: "" }));
    act(() => result.current.set("grade", "10"));
    expect(result.current.activeKeys).toEqual(["grade"]);
    act(() => result.current.set("query", "ann"));
    expect(result.current.hasActive).toBe(true);
    act(() => result.current.clear("grade"));
    expect(result.current.filters).toEqual({ query: "ann", grade: "" });
    act(() => result.current.clear());
    expect(result.current.hasActive).toBe(false);
  });

  it("debounces the search value", () => {
    vi.useFakeTimers();
    const { result } = renderHook(() => useListFilters({ query: "" }, { debounceMs: 100 }));
    act(() => result.current.set("query", "abc"));
    expect(result.current.debouncedSearch).toBe("");
    act(() => vi.advanceTimersByTime(100));
    expect(result.current.debouncedSearch).toBe("abc");
    vi.useRealTimers();
  });
});
