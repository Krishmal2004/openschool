import type { ReactNode } from "react";

interface Props {
  label: string;
  value: ReactNode;
  bold?: boolean;
  divider?: boolean;
}

// One label/value line in a compact summary panel (e.g. a sidebar "Quick Info" card).
export default function InfoRow({ label, value, bold, divider = true }: Props) {
  return (
    <div
      style={{
        display: "flex",
        justifyContent: "space-between",
        padding: "0.5rem 0",
        borderBottom: divider ? "1px solid var(--os-border-subtle)" : "none",
        fontSize: "0.8125rem",
      }}
    >
      <span style={{ color: "var(--os-text-secondary)" }}>{label}</span>
      <span style={{ fontWeight: bold ? 600 : 500, color: "var(--os-text-primary)" }}>{value}</span>
    </div>
  );
}
