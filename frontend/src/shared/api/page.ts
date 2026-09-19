// Shared shape for a server-paginated list response
// (docs/SECURITY_AND_PERFORMANCE_PLAYBOOK.md section 4.1). Never a bare
// array, so `total` always travels with the page.
export interface Page<T> {
  items: T[];
  total: number;
  limit: number;
  offset: number;
}
