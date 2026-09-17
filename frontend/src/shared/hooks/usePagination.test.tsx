import { describe, expect, it } from "vitest";
import { act, renderHook } from "@testing-library/react";
import { usePagination } from "@/shared/hooks/usePagination";

const items = Array.from({ length: 25 }, (_, i) => i + 1);

describe("usePagination", () => {
  it("slices the current page", () => {
    const { result } = renderHook(() => usePagination(items, 10));
    expect(result.current.pageItems).toEqual(items.slice(0, 10));
    act(() => result.current.onChange({ page: 3, pageSize: 10 }));
    expect(result.current.pageItems).toEqual([21, 22, 23, 24, 25]);
  });

  it("clamps the page when the list shrinks", () => {
    let list = items;
    const { result, rerender } = renderHook(() => usePagination(list, 10));
    act(() => result.current.onChange({ page: 3, pageSize: 10 }));
    list = items.slice(0, 5);
    rerender();
    expect(result.current.page).toBe(1);
    expect(result.current.totalItems).toBe(5);
  });
});
