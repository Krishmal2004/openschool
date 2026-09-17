// Reference data (grades, subjects, houses, mediums, academic years) rarely
// changes within a session, so it doesn't need TanStack Query's default
// background refetch-on-window-focus behaviour every few seconds
// (docs/SECURITY_AND_PERFORMANCE_PLAYBOOK.md section 5). An explicit
// mutation still invalidates and refetches immediately regardless of this.
export const REFERENCE_DATA_STALE_TIME_MS = 60 * 60 * 1000;
