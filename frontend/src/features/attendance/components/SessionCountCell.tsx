import { useSessionRecords } from "@/features/attendance/queries/useAttendance";

// Present / absent count for one session; each cell loads its own records.
export default function SessionCountCell({ sessionId, status }: { sessionId: string; status: "present" | "absent" }) {
  const { data: records, isLoading } = useSessionRecords(sessionId);
  if (isLoading) return <span className="os-c-tertiary">…</span>;
  const n = records?.filter((r) => r.status === status).length ?? 0;
  const tone = status === "present" ? "os-c-success os-fw-600" : n > 0 ? "os-c-danger os-fw-600" : "os-c-tertiary";
  return <span className={tone}>{n}</span>;
}
