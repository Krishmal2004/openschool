import type { ReactNode } from "react";

interface Props {
  title: ReactNode;
  meta?: ReactNode;
  className?: string;
}

// Standard `.os-section` header: a title with an optional trailing count/action.
export default function SectionHeader({ title, meta, className }: Props) {
  return (
    <div className={`os-section__header${className ? ` ${className}` : ""}`}>
      <h2 className="os-section__title">{title}</h2>
      {meta}
    </div>
  );
}
