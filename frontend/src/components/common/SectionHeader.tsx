import type { CSSProperties, ReactNode } from "react";

interface Props {
  title: ReactNode;
  meta?: ReactNode;
  style?: CSSProperties;
}

// Standard `.os-section` header: a title with an optional trailing count/action.
export default function SectionHeader({ title, meta, style }: Props) {
  return (
    <div className="os-section__header" style={style}>
      <h2 className="os-section__title">{title}</h2>
      {meta}
    </div>
  );
}
