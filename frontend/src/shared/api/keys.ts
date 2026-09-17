// Cross-cutting query keys only. Each feature owns its own keys.ts under one root segment.
export const keys = {
  me: {
    all: ["me"] as const,
    profile: () => ["me", "profile"] as const,
  },
};
