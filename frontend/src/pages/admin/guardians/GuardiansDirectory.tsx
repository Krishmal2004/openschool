import { useMemo, useState } from "react";
import { Close } from "@carbon/icons-react";
import { Checkbox, Pagination, Tile, Tag, TableToolbarSearch, ClickableTile, Button } from "@carbon/react";
import { useGuardians, useSearchGuardians } from "../../../queries/useGuardians";
import { usePagination } from "../../../hooks/usePagination";
import EmptyState from "../../../components/common/EmptyState";
import ErrorMessage from "../../../components/common/ErrorMessage";
import Avatar from "../../../components/common/Avatar";
import ListRowSkeleton from "../../../components/common/ListRowSkeleton";
import { relationshipLabel } from "./constants";
import GuardianDetail from "./components/GuardianDetail";

export default function GuardiansDirectory() {
  const [search, setSearch] = useState("");
  const [orphansOnly, setOrphansOnly] = useState(false);
  const [selectedId, setSelectedId] = useState<string | null>(null);

  const allGuardians = useGuardians(orphansOnly);
  const searchResults = useSearchGuardians(search, orphansOnly);
  const isSearching = search.trim().length > 0;

  const { data: guardians, isLoading, isError, refetch } = isSearching ? searchResults : allGuardians;

  const ordered = useMemo(() => {
    if (!guardians) return [];
    return [...guardians].sort((a, b) => a.full_name.localeCompare(b.full_name));
  }, [guardians]);

  const selected = ordered.find((g) => g.id === selectedId) ?? null;

  const { page, pageSize, pageItems, totalItems, onChange } = usePagination(ordered, 10);

  return (
    <div className="os-page">
      <div className="os-page__header">
        <div className="os-page__header-left">
          <h1 className="os-page__title">Guardians</h1>
          <p className="os-page__subtitle">
            Directory of parents and guardians linked to student profiles.
          </p>
        </div>
      </div>

      <div className="os-section" style={{ background: "#ffffff", marginBottom: "1.5rem", padding: "1rem 1.5rem" }}>
        <div style={{ display: "flex", alignItems: "center", gap: "1.5rem", flexWrap: "wrap" }}>
          <div style={{ flex: "1 1 18rem", minWidth: "14rem" }}>
            <TableToolbarSearch
              persistent
              placeholder="Search guardians by name, phone, or email…"
              value={search}
              onChange={(e) => setSearch(typeof e === "string" ? e : e.target.value)}
            />
          </div>
          <Checkbox
            id="orphans-only"
            labelText="Unlinked Guardians Only (Orphans)"
            checked={orphansOnly}
            onChange={(_e, { checked }) => setOrphansOnly(checked)}
          />
        </div>

        {(search || orphansOnly) && (
          <div style={{ display: "flex", alignItems: "center", gap: "0.5rem", flexWrap: "wrap", marginTop: "0.75rem" }}>
            <span style={{ fontSize: "0.75rem", fontWeight: 600, color: "var(--os-text-tertiary)" }}>Active Filters:</span>
            {orphansOnly && <Tag type="magenta" filter onClose={() => setOrphansOnly(false)}>Filter: Unlinked Only</Tag>}
            {search && <Tag type="blue" filter onClose={() => setSearch("")}>Search: "{search}"</Tag>}
            <Button kind="ghost" size="sm" renderIcon={Close} onClick={() => { setSearch(""); setOrphansOnly(false); }}>
              Clear All
            </Button>
          </div>
        )}
      </div>

      {isError && <ErrorMessage message="Could not load guardian records." onRetry={refetch} />}

      <div style={{ display: "grid", gridTemplateColumns: "22rem 1fr", gap: "1.5rem", alignItems: "start" }}>
        <div className="os-section" style={{ background: "#ffffff", marginTop: 0, padding: 0, overflow: "hidden" }}>
          {isLoading && (
            <div style={{ padding: "1rem" }}>
              {Array.from({ length: 5 }).map((_, i) => (
                <ListRowSkeleton key={i} leadingWidth="2.25rem" titleWidth="70%" subtitleWidth={null} trailingWidth={null} />
              ))}
            </div>
          )}

          {!isLoading && !isError && ordered.length === 0 && (
            <div style={{ padding: "1.5rem" }}>
              <EmptyState
                title={isSearching ? "No matching guardians" : "No guardians found"}
                description={
                  isSearching
                    ? "Try searching by a different name, phone number, or email address."
                    : "Guardians are automatically linked when adding or editing student profiles."
                }
              />
            </div>
          )}

          {!isLoading &&
            pageItems.map((g) => {
              const isSelected = selected?.id === g.id;
              return (
                <ClickableTile
                  key={g.id}
                  onClick={() => setSelectedId(g.id)}
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
                  <Avatar name={g.full_name} size="sm" />
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
                      {g.full_name}
                    </div>
                    <div style={{ fontSize: "0.75rem", color: "var(--os-text-tertiary)" }}>
                      {relationshipLabel(g.relationship)} · {g.phone}
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
          <GuardianDetail guardian={selected} onDeleted={() => setSelectedId(null)} />
        ) : (
          <Tile className="os-section" style={{ background: "#ffffff", marginTop: 0, textAlign: "center", padding: "3rem 1.5rem" }}>
            <EmptyState title="Select a guardian" description="Choose a guardian from the directory list to view contact info and linked students." />
          </Tile>
        )}
      </div>
    </div>
  );
}


