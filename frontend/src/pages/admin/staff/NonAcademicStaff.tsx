import { useMemo, useState } from "react";
import { Add, Close } from "@carbon/icons-react";
import { Button, Select, SelectItem, Pagination, Tile, Tag, TableToolbarSearch, ClickableTile } from "@carbon/react";
import { useNonAcademicStaffList } from "../../../queries/useNonAcademicStaff";
import { NON_ACADEMIC_DESIGNATIONS } from "../../../services/nonAcademicStaff";
import { usePagination } from "../../../hooks/usePagination";
import EmptyState from "../../../components/common/EmptyState";
import ErrorMessage from "../../../components/common/ErrorMessage";
import Avatar from "../../../components/common/Avatar";
import ListRowSkeleton from "../../../components/common/ListRowSkeleton";
import { designationLabel } from "./constants";
import StaffFormModal from "./components/StaffFormModal";
import StaffDetail from "./components/StaffDetail";

export default function NonAcademicStaff() {
  const [search, setSearch] = useState("");
  const [designation, setDesignation] = useState("");
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [creating, setCreating] = useState(false);

  const { data: staff, isLoading, isError, refetch } = useNonAcademicStaffList(search, designation);

  const ordered = useMemo(() => {
    if (!staff) return [];
    return [...staff].sort((a, b) => a.full_name.localeCompare(b.full_name));
  }, [staff]);

  const selected = ordered.find((s) => s.id === selectedId) ?? null;

  const { page, pageSize, pageItems, totalItems, onChange } = usePagination(ordered, 10);

  return (
    <div className="os-page">
      <div className="os-page__header">
        <div className="os-page__header-left">
          <h1 className="os-page__title">Non-Academic Staff</h1>
          <p className="os-page__subtitle">
            Lab assistants, librarians, office staff, and other staff without a portal login.
          </p>
        </div>
        <Button renderIcon={Add} kind="primary" size="md" onClick={() => setCreating(true)}>
          Add Staff
        </Button>
      </div>

      <div className="os-section" style={{ background: "#ffffff", marginBottom: "1.5rem", padding: "1rem 1.5rem" }}>
        <div style={{ display: "flex", alignItems: "center", gap: "1rem", flexWrap: "wrap" }}>
          <div style={{ flex: "1 1 16rem", minWidth: "12rem" }}>
            <TableToolbarSearch
              persistent
              placeholder="Search staff by name or employee number…"
              value={search}
              onChange={(e) => setSearch(typeof e === "string" ? e : e.target.value)}
            />
          </div>
          <div style={{ minWidth: "14rem" }}>
            <Select
              id="staff-designation-filter"
              labelText="Designation"
              hideLabel
              size="md"
              value={designation}
              onChange={(e) => setDesignation(e.target.value)}
            >
              <SelectItem value="" text="All designations" />
              {NON_ACADEMIC_DESIGNATIONS.map((d) => (
                <SelectItem key={d.value} value={d.value} text={d.label} />
              ))}
            </Select>
          </div>
        </div>

        {(search || designation) && (
          <div style={{ display: "flex", alignItems: "center", gap: "0.5rem", flexWrap: "wrap", marginTop: "0.75rem" }}>
            <span style={{ fontSize: "0.75rem", fontWeight: 600, color: "var(--os-text-tertiary)" }}>Active Filters:</span>
            {designation && (
              <Tag type="teal" filter onClose={() => setDesignation("")}>
                Role: {designationLabel(designation)}
              </Tag>
            )}
            {search && <Tag type="blue" filter onClose={() => setSearch("")}>Search: "{search}"</Tag>}
            <Button kind="ghost" size="sm" renderIcon={Close} onClick={() => { setSearch(""); setDesignation(""); }}>
              Clear All
            </Button>
          </div>
        )}
      </div>

      {isError && <ErrorMessage message="Could not load staff members." onRetry={refetch} />}

      <div style={{ display: "grid", gridTemplateColumns: "22rem 1fr", gap: "1.5rem", alignItems: "start" }}>
        <div className="os-section" style={{ background: "#ffffff", marginTop: 0, padding: 0, overflow: "hidden" }}>
          {isLoading && (
            <div style={{ padding: "1rem" }}>
              {Array.from({ length: 5 }).map((_, i) => (
                <ListRowSkeleton key={i} leadingWidth="1.5rem" titleWidth="70%" subtitleWidth="40%" trailingWidth={null} />
              ))}
            </div>
          )}

          {!isLoading && !isError && ordered.length === 0 && (
            <div style={{ padding: "1.5rem" }}>
              <EmptyState
                title="No staff found"
                description="Add lab assistants, librarians, office staff, or adjust your filter."
              />
            </div>
          )}

          {!isLoading &&
            pageItems.map((s) => {
              const isSelected = selected?.id === s.id;
              return (
                <ClickableTile
                  key={s.id}
                  onClick={() => setSelectedId(s.id)}
                  className={`os-list-row${isSelected ? " is-selected" : ""}`}
                  style={{
                    display: "flex",
                    alignItems: "center",
                    gap: "0.75rem",
                    padding: "0.875rem 1.25rem",
                    borderRadius: 0,
                    borderBottom: "1px solid var(--os-border-subtle)",
                    backgroundColor: isSelected ? "var(--os-accent-light)" : "#ffffff",
                  }}
                >
                  <Avatar name={s.full_name} size="sm" />
                  <div style={{ minWidth: 0, flex: 1 }}>
                    <div
                      style={{
                        fontWeight: 600,
                        fontSize: "0.875rem",
                        color: "var(--os-text-primary)",
                        overflow: "hidden",
                        textOverflow: "ellipsis",
                        whiteSpace: "nowrap",
                      }}
                    >
                      {s.full_name}
                    </div>
                    <div style={{ fontSize: "0.75rem", color: "var(--os-text-tertiary)" }}>
                      {designationLabel(s.designation)} · {s.employee_number}
                    </div>
                  </div>
                </ClickableTile>
              );
            })}

          {!isLoading && ordered.length > 0 && (
            <Pagination
              totalItems={totalItems}
              page={page}
              pageSize={pageSize}
              pageSizes={[10, 20, 50]}
              onChange={onChange}
              size="sm"
            />
          )}
        </div>

        {selected ? (
          <StaffDetail staff={selected} onDeleted={() => setSelectedId(null)} />
        ) : (
          <Tile className="os-section" style={{ background: "#ffffff", marginTop: 0, textAlign: "center", padding: "3rem 1.5rem" }}>
            <EmptyState title="Select a staff member" description="Choose a staff member from the left directory to view full profile details." />
          </Tile>
        )}
      </div>

      {creating && <StaffFormModal staff={null} onClose={() => setCreating(false)} />}
    </div>
  );
}


