import StatusTag from "../../../../components/common/StatusTag";
import { getInitials } from "../../../../lib/name";
import type { Student } from "../../../../services/student";
import { STATUS_STYLES, type Status } from "../constants";
import StatusButton from "./StatusButton";

export default function StudentAttendanceRow({
  student,
  idx,
  status,
  note,
  readOnly,
  onMark,
  onNoteChange,
}: {
  student: Student;
  idx: number;
  status: Status;
  note: string;
  readOnly: boolean;
  onMark: (status: NonNullable<Status>) => void;
  onNoteChange: (value: string) => void;
}) {
  return (
    <tr style={{ background: status ? `color-mix(in srgb, ${STATUS_STYLES[status].bg} 40%, white)` : "transparent" }}>
      <td style={{ color: "var(--os-text-tertiary)", fontFamily: "IBM Plex Mono, monospace", fontSize: "0.75rem" }}>
        {idx + 1}
      </td>
      <td>
        <div style={{ display: "flex", alignItems: "center", gap: "0.625rem" }}>
          <div
            style={{
              width: "1.75rem",
              height: "1.75rem",
              borderRadius: "50%",
              background: status ? STATUS_STYLES[status].bg : "var(--os-accent-light)",
              border: `1px solid ${status ? STATUS_STYLES[status].border : "var(--os-accent-border)"}`,
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              fontSize: "0.6rem",
              fontWeight: 700,
              color: status ? STATUS_STYLES[status].color : "var(--os-accent)",
              flexShrink: 0,
            }}
          >
            {getInitials(student.full_name)}
          </div>
          <span style={{ fontWeight: 500, fontSize: "0.875rem" }}>{student.full_name}</span>
        </div>
      </td>
      <td className="os-table__mono">{student.index_number}</td>
      <td>
        {readOnly ? (
          status ? (
            <StatusTag {...STATUS_STYLES[status]} />
          ) : (
            <span style={{ color: "var(--os-text-disabled)", fontSize: "0.75rem" }}>Not marked</span>
          )
        ) : (
          <div style={{ display: "flex", gap: "0.375rem" }}>
            {(["present", "absent", "late", "excused"] as const).map((s) => (
              <StatusButton key={s} value={s} selected={status === s} onClick={() => onMark(s)} />
            ))}
          </div>
        )}
      </td>
      <td>
        {readOnly ? (
          <span style={{ fontSize: "0.75rem", color: note ? "var(--os-text-secondary)" : "var(--os-text-disabled)" }}>{note || "-"}</span>
        ) : status === "absent" || status === "late" || status === "excused" ? (
          <input
            className="os-note-input"
            placeholder="Optional note…"
            aria-label={`Note for ${student.full_name}`}
            value={note}
            onChange={(e) => onNoteChange(e.target.value)}
          />
        ) : (
          <span style={{ color: "var(--os-text-disabled)", fontSize: "0.75rem" }}>-</span>
        )}
      </td>
    </tr>
  );
}
