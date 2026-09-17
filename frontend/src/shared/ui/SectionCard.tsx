import type { ReactNode } from "react";

interface Props {
  title: ReactNode;
  meta?: ReactNode;
  children: ReactNode;
  // Tables sit flush with the card edges; text and forms get body padding.
  flush?: boolean;
  className?: string;
}

// A titled `.os-section` card, the standard container for page content.
export default function SectionCard({ title, meta, children, flush = false, className }: Props) {
  return (
    <div className={`os-section${className ? ` ${className}` : ""}`}>
      <div className="os-section__header">
        <h2 className="os-section__title">{title}</h2>
        {meta}
      </div>
      {flush ? children : <div className="os-section__body">{children}</div>}
    </div>
  );
}
