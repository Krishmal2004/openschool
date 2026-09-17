import type { ReactNode } from "react";
import {
  Pagination,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableHeader,
  TableRow,
  TableToolbar,
  TableToolbarContent,
} from "@carbon/react";
import { usePagination } from "@/shared/hooks/usePagination";

export interface GridColumn<T> {
  key: string;
  header: string;
  render: (row: T) => ReactNode;
  align?: "start" | "end";
}

interface ServerPaging {
  page: number;
  pageSize: number;
  totalItems: number;
  onChange: (next: { page: number; pageSize: number }) => void;
}

interface Props<T> {
  rows: T[];
  columns: GridColumn<T>[];
  getRowId: (row: T) => string;
  pageSize?: number;
  pageSizes?: number[];
  // Pass when the backend paginates; rows are then the current page only.
  server?: ServerPaging;
  countLabel?: (shown: number, total: number) => string;
  toolbar?: ReactNode;
  noHover?: boolean;
  // false renders every row with no pager, for short lists inside tabs and panels.
  pagination?: boolean;
  className?: string;
  onRowClick?: (row: T) => void;
}

// Carbon table with pagination, count label and optional toolbar. Client-side paging by default.
export default function DataGrid<T>({
  rows,
  columns,
  getRowId,
  pageSize = 10,
  pageSizes = [10, 20, 50],
  server,
  countLabel,
  toolbar,
  noHover,
  pagination = true,
  className,
  onRowClick,
}: Props<T>) {
  const client = usePagination(rows, pageSize);
  const pageRows = server || !pagination ? rows : client.pageItems;
  const paging = server ?? client;
  const total = server ? server.totalItems : rows.length;

  return (
    <>
      <TableContainer className="os-table-container">
        {(countLabel || toolbar) && (
          <TableToolbar className="os-grid__toolbar">
            <TableToolbarContent>
              {countLabel && <span className="os-grid__count">{countLabel(pageRows.length, total)}</span>}
              {toolbar}
            </TableToolbarContent>
          </TableToolbar>
        )}
        <Table className={`os-table${noHover ? " os-table--no-hover" : ""}${className ? ` ${className}` : ""}`}>
          <TableHead>
            <TableRow>
              {columns.map((c) => (
                <TableHeader key={c.key} className={c.align === "end" ? "os-grid__cell--end" : undefined}>
                  {c.header}
                </TableHeader>
              ))}
            </TableRow>
          </TableHead>
          <TableBody>
            {pageRows.map((row) => (
              <TableRow key={getRowId(row)} onClick={onRowClick ? () => onRowClick(row) : undefined} className={onRowClick ? "os-pointer" : undefined}>
                {columns.map((c) => (
                  <TableCell key={c.key} className={c.align === "end" ? "os-grid__cell--end" : undefined}>
                    {c.render(row)}
                  </TableCell>
                ))}
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
      {pagination && (
        <Pagination totalItems={total} page={paging.page} pageSize={paging.pageSize} pageSizes={pageSizes} onChange={paging.onChange} />
      )}
    </>
  );
}
