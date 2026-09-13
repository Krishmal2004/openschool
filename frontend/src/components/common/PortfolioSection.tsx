import type { ReactNode } from "react";

// One titled block in a portfolio view (prefect appointments, societies,
// activities, leadership & awards, disciplinary records) — either its table
// or a plain empty-state message.
export default function PortfolioSection({
  title,
  isEmpty,
  emptyMessage,
  children,
}: {
  title: string;
  isEmpty: boolean;
  emptyMessage: string;
  children: ReactNode;
}) {
  return (
    <div>
      <h3 style={{ fontSize: "0.875rem", fontWeight: 600, margin: "0 0 0.5rem" }}>{title}</h3>
      {isEmpty ? (
        <p style={{ fontSize: "0.8125rem", color: "var(--os-text-tertiary)" }}>{emptyMessage}</p>
      ) : (
        children
      )}
    </div>
  );
}
