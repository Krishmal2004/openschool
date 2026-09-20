import { describe, expect, it } from "vitest";
import { AxiosError, AxiosHeaders } from "axios";
import { getErrorMessage, isNotFoundError } from "@/shared/api/errors";

function axiosError(status: number, data?: unknown): AxiosError {
  return new AxiosError("Request failed", "ERR", undefined, undefined, {
    status,
    statusText: "",
    headers: {},
    config: { headers: new AxiosHeaders() },
    data,
  });
}

describe("getErrorMessage", () => {
  it("returns the backend error string when present", () => {
    expect(getErrorMessage(axiosError(400, { error: "bad input" }))).toBe("bad input");
  });

  it("falls back for non-axios errors and missing bodies", () => {
    expect(getErrorMessage(new Error("x"), "fallback")).toBe("fallback");
    expect(getErrorMessage(axiosError(500), "fallback")).toBe("fallback");
  });

  it("falls back instead of showing text that looks like a leaked internal error", () => {
    expect(getErrorMessage(axiosError(400, { error: "Key: 'CreateSessionRequest.ClassID' Error:Field validation for 'ClassID' failed on the 'required' tag" }), "fallback")).toBe("fallback");
    expect(getErrorMessage(axiosError(500, { error: "pq: duplicate key value violates unique constraint" }), "fallback")).toBe("fallback");
    expect(getErrorMessage(axiosError(500, { error: "SQLSTATE 23505: unique_violation" }), "fallback")).toBe("fallback");
  });
});

describe("isNotFoundError", () => {
  it("is true only for axios 404 responses", () => {
    expect(isNotFoundError(axiosError(404))).toBe(true);
    expect(isNotFoundError(axiosError(500))).toBe(false);
    expect(isNotFoundError(new Error("x"))).toBe(false);
  });
});
