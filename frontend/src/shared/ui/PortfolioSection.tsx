import type { ReactNode } from "react";

// One titled block in a portfolio view: its table, or an empty-state message.
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
      <h3 className="os-text-md os-fw-600 os-mt-0 os-mx-0 os-mb-2">{title}</h3>
      {isEmpty ? (
        <p className="os-text-sm os-c-tertiary">{emptyMessage}</p>
      ) : (
        children
      )}
    </div>
  );
}
