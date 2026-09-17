import type { ReactNode } from "react";
import { DocumentBlank } from "@carbon/icons-react";

interface Props {
  title: string;
  description: string;
  action?: ReactNode;
}

export default function EmptyState({ title, description, action }: Props) {
  return (
    <div className="os-empty">
      <DocumentBlank size={40} className="os-empty__icon" />
      <p className="os-empty__title">{title}</p>
      <p className="os-empty__desc">{description}</p>
      {action}
    </div>
  );
}
