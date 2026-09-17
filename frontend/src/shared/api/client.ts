import axios, { type InternalAxiosRequestConfig } from "axios";

type TokenProvider = () => Promise<string | null | undefined>;

let tokenProvider: TokenProvider | null = null;
let unauthorizedHandler: (() => void) | null = null;
let markReady: () => void = () => {};
// Requests wait here until the auth bridge has mounted, so none go out without a token.
const providerReady = new Promise<void>((resolve) => {
  markReady = resolve;
});

export function setAccessTokenProvider(provider: TokenProvider) {
  tokenProvider = provider;
  markReady();
}

export function setUnauthorizedHandler(handler: (() => void) | null) {
  unauthorizedHandler = handler;
}

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL ?? "http://localhost:8080/api/v1",
  headers: { "Content-Type": "application/json" },
});

api.interceptors.request.use(async (config: InternalAxiosRequestConfig) => {
  await providerReady;
  try {
    const token = await tokenProvider?.();
    if (token) config.headers.Authorization = `Bearer ${token}`;
  } catch {
    // Not signed in yet; public endpoints still work without a token.
  }
  return config;
});

let handling401 = false;
api.interceptors.response.use(
  (response) => response,
  (error) => {
    const status = error?.response?.status;
    const url: string = error?.config?.url ?? "";
    // Self-service auth endpoints return 401 for bad secrets, not for a dead session.
    if (status === 401 && !url.startsWith("/auth/") && unauthorizedHandler && !handling401) {
      handling401 = true;
      try {
        unauthorizedHandler();
      } finally {
        setTimeout(() => {
          handling401 = false;
        }, 5000);
      }
    }
    return Promise.reject(error);
  },
);

export default api;
