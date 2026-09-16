import { useState } from "react";
import { Link, useNavigate } from "react-router";
import { Add, Edit, TrashCan, Close } from "@carbon/icons-react";
import {
  Button,
  IconButton,
  Pagination,
  Tag,
  Select,
  SelectItem,
  DataTable,
  TableContainer,
  TableToolbar,
  TableToolbarContent,
  TableToolbarSearch,
  Table,
  TableHead,
  TableRow,
  TableHeader,
  TableBody,
  TableCell,
} from "@carbon/react";
import { useTeachers, useDeleteTeacher } from "../../../queries/useTeachers";
import type { Teacher } from "../../../services/teacher";
import { usePagination } from "../../../hooks/usePagination";
import TableSkeleton from "../../../components/common/TableSkeleton";
import ErrorMessage from "../../../components/common/ErrorMessage";
import EmptyState from "../../../components/common/EmptyState";
import ConfirmDeleteModal from "../../../components/common/ConfirmDeleteModal";
import AgentFindingsBanner from "../../../components/common/AgentFindingsBanner";
import MutationErrorNotification from "../../../components/common/MutationErrorNotification";

const TEACHER_TABLE_HEADERS = [
  { key: "employee_number", header: "Employee No." },
  { key: "full_name", header: "Full Name" },
  { key: "phone", header: "Phone" },
  { key: "joined_date", header: "Joined Date" },
  { key: "employment_status", header: "Status" },
  { key: "actions", header: "Actions" },
];

export default function Teachers() {
  const navigate = useNavigate();
  const { data: teachers, isLoading, isError, refetch } = useTeachers();
  const deleteTeacher = useDeleteTeacher();
  const [query, setQuery] = useState("");
  const [statusFilter, setStatusFilter] = useState("");
  const [toDelete, setToDelete] = useState<Teacher | null>(null);

  const filtered = (teachers ?? []).filter((t) => {
    const q = query.toLowerCase();
    const matchesQuery =
      t.full_name.toLowerCase().includes(q) ||
      t.employee_number.toLowerCase().includes(q);
    const matchesStatus = !statusFilter || t.employment_status === statusFilter;
    return matchesQuery && matchesStatus;
  });

  const { page, pageSize, pageItems, totalItems, onChange } = usePagination(filtered, 10);

  const handleDelete = () => {
    if (!toDelete) return;
    deleteTeacher.mutate(toDelete.id, {
      onSettled: () => setToDelete(null),
    });
  };

  return (
    <div className="os-page">
      <div className="os-page__header">
        <div className="os-page__header-left">
          <h1 className="os-page__title">Teachers</h1>
          <p className="os-page__subtitle">Manage teacher profiles</p>
        </div>
        <Button renderIcon={Add} kind="primary" size="md" as={Link} to="/teachers/new">
          Add Teacher
        </Button>
      </div>

      <AgentFindingsBanner titles={["Inactive teachers still assigned to classes", "Teacher accounts stuck in first-login setup"]} />

      <div className="os-section" style={{ background: "#ffffff" }}>
        <div className="os-toolbar" style={{ padding: "1rem 1.5rem 0.5rem" }}>
          <div style={{ flex: "1 1 16rem", minWidth: "12rem" }}>
            <TableToolbarSearch
              persistent
              placeholder="Search by name or employee number…"
              value={query}
              onChange={(e) => setQuery(typeof e === "string" ? e : e.target.value)}
            />
          </div>
          <div style={{ minWidth: "10rem" }}>
            <Select
              id="filter-teacher-status"
              labelText=""
              size="md"
              value={statusFilter}
              onChange={(e) => setStatusFilter(e.target.value)}
            >
              <SelectItem value="" text="All Statuses" />
              <SelectItem value="active" text="Active" />
              <SelectItem value="resigned" text="Resigned" />
              <SelectItem value="transferred" text="Transferred" />
            </Select>
          </div>
        </div>

        {(query || statusFilter) && (
          <div style={{ display: "flex", alignItems: "center", gap: "0.5rem", flexWrap: "wrap", padding: "0 1.5rem 0.75rem" }}>
            <span style={{ fontSize: "0.75rem", fontWeight: 600, color: "var(--os-text-tertiary)" }}>Active Filters:</span>
            {statusFilter && <Tag type="teal" filter onClose={() => setStatusFilter("")}>Status: {statusFilter}</Tag>}
            {query && <Tag type="blue" filter onClose={() => setQuery("")}>Search: "{query}"</Tag>}
            <Button kind="ghost" size="sm" renderIcon={Close} onClick={() => { setQuery(""); setStatusFilter(""); }}>
              Clear All
            </Button>
          </div>
        )}

        <MutationErrorNotification
          isError={deleteTeacher.isError}
          error={deleteTeacher.error}
          title="Could not delete teacher"
          fallback="The teacher may be assigned to a class or have attendance records."
          onClose={() => deleteTeacher.reset()}
          style={{ margin: "0 1.5rem 1rem" }}
        />

        {isLoading ? (
          <TableSkeleton headers={TEACHER_TABLE_HEADERS.map(h => h.header)} />
        ) : isError ? (
          <ErrorMessage message="Failed to load teachers" onRetry={refetch} />
        ) : filtered.length === 0 ? (
          <EmptyState
            title="No teachers found"
            description="Add your first teacher or adjust your search filter to get started."
          />
        ) : (
          <>
            <DataTable
              rows={pageItems.map((t) => ({ ...t, id: t.id }))}
              headers={TEACHER_TABLE_HEADERS}
              render={({
                rows,
                headers,
                getHeaderProps,
                getRowProps,
              }) => (
                <TableContainer className="os-table-container">
                  <TableToolbar style={{ background: "#ffffff" }}>
                    <TableToolbarContent>
                      <span style={{ fontSize: "0.75rem", fontWeight: 600, color: "var(--os-text-tertiary)", paddingRight: "1rem" }}>
                        Showing {pageItems.length} of {filtered.length} teachers
                      </span>
                    </TableToolbarContent>
                  </TableToolbar>

                  <Table className="os-table">
                    <TableHead>
                      <TableRow>
                        {headers.map((header) => (
                          <TableHeader {...getHeaderProps({ header })}>
                            {header.header}
                          </TableHeader>
                        ))}
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {rows.map((row) => {
                        const teacher = pageItems.find((t) => t.id === row.id);
                        if (!teacher) return null;
                        return (
                          <TableRow {...getRowProps({ row })}>
                            <TableCell>
                              <span className="os-table__mono">{teacher.employee_number}</span>
                            </TableCell>
                            <TableCell>
                              <Link to={`/teachers/${teacher.id}`} className="os-table__link">
                                {teacher.full_name}
                              </Link>
                            </TableCell>
                            <TableCell className="os-table__muted">{teacher.phone ?? "-"}</TableCell>
                            <TableCell className="os-table__muted">{teacher.joined_date ?? "-"}</TableCell>
                            <TableCell>
                              {teacher.employment_status === "active" ? (
                                <Tag type="green" size="sm">Active</Tag>
                              ) : teacher.employment_status === "resigned" ? (
                                <Tag type="red" size="sm">Resigned</Tag>
                              ) : (
                                <Tag type="magenta" size="sm">Transferred</Tag>
                              )}
                            </TableCell>
                            <TableCell>
                              <div style={{ display: "flex", gap: "0.25rem", justifyContent: "flex-end" }}>
                                <IconButton
                                  label="Edit"
                                  kind="ghost"
                                  size="sm"
                                  onClick={() =>
                                    navigate(`/teachers/${teacher.id}`, {
                                      state: { edit: true },
                                    })
                                  }
                                >
                                  <Edit />
                                </IconButton>
                                <IconButton
                                  label="Delete"
                                  kind="ghost"
                                  size="sm"
                                  onClick={() => {
                                    deleteTeacher.reset();
                                    setToDelete(teacher);
                                  }}
                                >
                                  <TrashCan />
                                </IconButton>
                              </div>
                            </TableCell>
                          </TableRow>
                        );
                      })}
                    </TableBody>
                  </Table>
                </TableContainer>
              )}
            />
            <Pagination
              totalItems={totalItems}
              page={page}
              pageSize={pageSize}
              pageSizes={[10, 20, 50]}
              onChange={onChange}
            />
          </>
        )}
      </div>

      <ConfirmDeleteModal
        open={!!toDelete}
        title="Delete teacher"
        description={
          <>
            Delete <strong>{toDelete?.full_name}</strong>? This removes their account and cannot be undone.
          </>
        }
        isPending={deleteTeacher.isPending}
        onClose={() => setToDelete(null)}
        onConfirm={handleDelete}
      />
    </div>
  );
}



