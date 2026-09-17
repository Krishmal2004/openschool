import type { ReactNode } from "react";
import { TabPanel } from "@carbon/react";

interface Props {
  title: string;
  children: ReactNode;
  // Tables sit flush with the section edges; forms and text get body padding.
  padded?: boolean;
}

// A Carbon TabPanel holding one titled section card.
export default function TabSection({ title, children, padded = true }: Props) {
  return (
    <TabPanel className="os-p-0">
      <div className="os-section os-mt-4">
        <div className="os-section__header">
          <h2 className="os-section__title">{title}</h2>
        </div>
        {padded ? <div className="os-section__body">{children}</div> : children}
      </div>
    </TabPanel>
  );
}
