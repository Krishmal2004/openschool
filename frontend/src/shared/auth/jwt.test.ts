import { describe, expect, it } from "vitest";
import { parseJwt, resolveRole } from "@/shared/auth/jwt";

function encode(payload: object): string {
  const json = JSON.stringify(payload);
  const b64 = btoa(json).replace(/\+/g, "-").replace(/\//g, "_").replace(/=+$/, "");
  return `hdr.${b64}.sig`;
}

describe("parseJwt", () => {
  it("decodes a base64url payload containing - and _ characters", () => {
    // Long payloads with binary-ish bytes produce '-' and '_' in base64url.
    const payload = { roles: ["teacher"], sub: "??>>??>>", extra: "~~~~~~" };
    expect(parseJwt(encode(payload))).toEqual(payload);
  });

  it("returns null for malformed tokens", () => {
    expect(parseJwt("not-a-jwt")).toBeNull();
    expect(parseJwt("a.!!!.c")).toBeNull();
  });
});

describe("resolveRole", () => {
  it("prefers admin over any other role", () => {
    expect(resolveRole({ roles: ["parent", "admin"] })).toBe("admin");
  });

  it("accepts a single string role", () => {
    expect(resolveRole({ roles: "student" })).toBe("student");
  });

  it("returns null when no known role is present", () => {
    expect(resolveRole({ roles: ["guest"] })).toBeNull();
    expect(resolveRole(null)).toBeNull();
  });
});
