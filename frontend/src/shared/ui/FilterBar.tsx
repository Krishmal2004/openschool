import type { ReactNode } from "react";
import { TableToolbarSearch } from "@carbon/react";

interface Props {
  search: { value: string; onChange: (value: string) => void; placeholder?: string };
  children?: ReactNode;
}

// Search box plus any filter controls; each child is wrapped as one control.
export default function FilterBar({ search, children }: Props) {
  return (
    <div className="os-filter-bar">
      <div className="os-filter-bar__search">
        <TableToolbarSearch
          persistent
          placeholder={search.placeholder ?? "Search…"}
          value={search.value}
          onChange={(e) => search.onChange(typeof e === "string" ? e : e.target.value)}
        />
      </div>
      {Array.isArray(children)
        ? children.map((child, i) => child && <div key={i} className="os-filter-bar__control">{child}</div>)
        : children && <div className="os-filter-bar__control">{children}</div>}
    </div>
  );
}
