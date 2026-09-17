import type { ReactNode } from "react";

interface Props {
  label: string;
  value: ReactNode;
  bold?: boolean;
  divider?: boolean;
  // Highlights the value with the accent colour, for a count or stat.
  accent?: boolean;
}

// One label/value line in a compact summary panel (e.g. a sidebar "Quick Info" card).
export default function InfoRow({ label, value, bold, divider = true, accent }: Props) {
  return (
    <div className={`os-flex os-justify-between os-py-2 os-text-sm${divider ? " os-border-b" : ""}`}>
      <span className="os-c-secondary">{label}</span>
      <span className={`${bold ? "os-fw-600" : "os-fw-500"} ${accent ? "os-c-accent" : "os-c-primary"}`}>{value}</span>
    </div>
  );
}
