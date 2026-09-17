// k6 load test simulating a school day's admin/teacher/parent traffic at
// roughly 7,000 signed-in sessions (docs/SECURITY_AND_PERFORMANCE_PLAYBOOK.md
// sections 5 and 7): one dashboard view, one paginated list page (the
// Students directory — section 4's reference implementation), and one poll
// (the notification unread-count badge), per simulated session.
//
// Auth note: same constraint as loadtest_attendance_rush.js — ThunderID
// doesn't implement the OAuth2 "password" grant, so this script can't log
// in as each simulated user itself. It consumes a pool of already-issued
// access tokens instead.
//
// Usage:
//   1. Bring up the full stack and seed it with ~7,000 students (and a
//      realistic number of teachers/staff) so the Students page and
//      dashboard analytics reflect real scale.
//   2. Create TOKENS.json — a plain array of already-issued access tokens,
//      one per virtual session, e.g. ["eyJ...", "eyJ...", ...]. Sign in
//      through the real frontend (or drive @thunderid/react's flow
//      headlessly) to get them; tokens are short-lived, so regenerate this
//      file right before running.
//   3. k6 run -e BASE_URL=http://localhost:8080 scripts/loadtest_school_day.js

import http from "k6/http";
import { check, sleep } from "k6";
import { SharedArray } from "k6/data";

const BASE_URL = __ENV.BASE_URL || "http://localhost:8080";

const tokens = new SharedArray("tokens", function () {
  return JSON.parse(open("./TOKENS.json"));
});

export const options = {
  scenarios: {
    school_day: {
      executor: "constant-vus",
      vus: Number(__ENV.LOAD_VUS || 200),
      duration: __ENV.LOAD_DURATION || "5m",
    },
  },
  thresholds: {
    http_req_duration: ["p(95)<500"],
    http_req_failed: ["rate<0.01"],
  },
};

export default function () {
  const token = tokens[__VU % tokens.length];
  const headers = { Authorization: `Bearer ${token}` };

  // One dashboard view: the admin analytics endpoint the Dashboard/Analytics
  // pages chart (frontend/src/features/reports/pages/admin/Analytics.tsx).
  const dashboard = http.get(`${BASE_URL}/api/v1/dashboard/analytics`, { headers });
  check(dashboard, { "dashboard ok": (r) => r.status === 200 });

  sleep(1);

  // One paginated list page.
  const students = http.get(`${BASE_URL}/api/v1/students?limit=25&offset=0`, { headers });
  check(students, { "students page ok": (r) => r.status === 200 });

  sleep(1);

  // One poll (refetches every 2 minutes client-side — see
  // frontend/src/features/notifications/queries/useNotifications.ts).
  const unread = http.get(`${BASE_URL}/api/v1/me/notifications/unread-count`, { headers });
  check(unread, { "unread poll ok": (r) => r.status === 200 });

  sleep(Math.random() * 2);
}
