import { useCallback } from "react";
import { useQueryClient } from "@tanstack/react-query";

// Invalidate one or more key prefixes after a mutation; replaces repeated queryClient boilerplate.
export function useInvalidate() {
  const queryClient = useQueryClient();
  return useCallback(
    (...prefixes: readonly (readonly unknown[])[]) => {
      for (const queryKey of prefixes) void queryClient.invalidateQueries({ queryKey });
    },
    [queryClient],
  );
}
