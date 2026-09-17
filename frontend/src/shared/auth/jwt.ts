// JWT payloads are base64url; atob only accepts standard base64.
function base64UrlDecode(input: string): string {
  const base64 = input.replace(/-/g, "+").replace(/_/g, "/");
  const padded = base64 + "=".repeat((4 - (base64.length % 4)) % 4);
  return atob(padded);
}

export function parseJwt(token: string): Record<string, unknown> | null {
  try {
    const payload = token.split(".")[1];
    if (!payload) return null;
    return JSON.parse(base64UrlDecode(payload));
  } catch {
    return null;
  }
}

export type Role = "admin" | "teacher" | "student" | "parent";

const ROLE_PRIORITY: Role[] = ["admin", "teacher", "student", "parent"];

// Picks the highest-privilege role from the token's `roles` claim.
export function resolveRole(payload: Record<string, unknown> | null): Role | null {
  const raw = payload?.roles;
  const roles = Array.isArray(raw) ? (raw as string[]) : typeof raw === "string" ? [raw] : [];
  return ROLE_PRIORITY.find((r) => roles.includes(r)) ?? null;
}
