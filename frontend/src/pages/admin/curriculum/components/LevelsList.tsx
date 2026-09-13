import { Link } from "react-router";
import { Add, ChevronRight, Copy, Edit, Layers } from "@carbon/icons-react";
import { Button, Tag } from "@carbon/react";
import type { useLevels, useDeleteLevel } from "../../../../queries/useCurriculum";
import type { useGrades } from "../../../../queries/useGrades";
import type { Level } from "../../../../services/curriculum";
import ErrorMessage from "../../../../components/common/ErrorMessage";
import EmptyState from "../../../../components/common/EmptyState";
import ListRowSkeleton from "../../../../components/common/ListRowSkeleton";
import MutationErrorNotification from "../../../../components/common/MutationErrorNotification";
import SectionHeader from "../../../../components/common/SectionHeader";

interface Props {
  levels: ReturnType<typeof useLevels>["data"];
  isLoading: boolean;
  isError: boolean;
  refetch: () => void;
  grades: ReturnType<typeof useGrades>["data"];
  deleteLevel: ReturnType<typeof useDeleteLevel>;
  onOpenCreate: () => void;
  onEdit: (level: Level) => void;
  onDuplicate: (level: Level) => void;
  onRequestDelete: (level: Level) => void;
}

export default function LevelsList({
  levels,
  isLoading,
  isError,
  refetch,
  grades,
  deleteLevel,
  onOpenCreate,
  onEdit,
  onDuplicate,
  onRequestDelete,
}: Props) {
  const gradeName = (id: string | null) => grades?.find((g) => g.id === id)?.name ?? null;

  return (
    <div className="os-section">
      <SectionHeader
        title="Levels"
        meta={levels && <span className="os-section__meta">{levels.length} total</span>}
      />

      {isLoading && (
        <div>
          {Array.from({ length: 4 }).map((_, i) => (
            <ListRowSkeleton key={i} titleWidth="30%" subtitleWidth="15%" trailingWidth="4rem" />
          ))}
        </div>
      )}
      {isError && <ErrorMessage message="Could not load levels." onRetry={refetch} />}

      <MutationErrorNotification
        isError={deleteLevel.isError}
        error={deleteLevel.error}
        title="Could not delete level"
        fallback="The level may have students enrolled through its groups."
        onClose={() => deleteLevel.reset()}
        style={{ margin: "0 1.5rem 1rem" }}
      />

      {!isLoading && !isError && levels?.length === 0 && (
        <EmptyState
          title="No levels yet"
          description="Create a level for each place a distinct set of subject rules applies - for example one per grade, or one per stream."
          action={
            <Button renderIcon={Add} kind="primary" onClick={onOpenCreate}>
              New Level
            </Button>
          }
        />
      )}

      {!isLoading && levels && levels.length > 0 && (
        <div>
          {levels.map((l) => (
            <div key={l.id} className="os-list-row">
              <Layers size={20} style={{ fill: "var(--os-accent)", flexShrink: 0 }} />
              <div style={{ flex: 1 }}>
                <p style={{ margin: "0 0 0.125rem", fontWeight: 600, fontSize: "0.9rem", color: "var(--os-text-primary)" }}>
                  {l.label}
                </p>
                <p style={{ margin: 0, fontSize: "0.75rem", color: "var(--os-text-secondary)" }}>Order {l.sort_order}</p>
              </div>
              {gradeName(l.grade_id) ? (
                <Tag type="teal" size="sm">
                  {gradeName(l.grade_id)}
                </Tag>
              ) : (
                <Tag type="gray" size="sm">
                  No grade
                </Tag>
              )}
              <Button kind="ghost" size="sm" renderIcon={ChevronRight} as={Link} to={`/curriculum/${l.id}`}>
                Configure
              </Button>
              <Button kind="ghost" size="sm" renderIcon={Edit} onClick={() => onEdit(l)}>
                Edit
              </Button>
              <Button kind="ghost" size="sm" renderIcon={Copy} onClick={() => onDuplicate(l)}>
                Duplicate
              </Button>
              <Button kind="danger--ghost" size="sm" onClick={() => onRequestDelete(l)}>
                Delete
              </Button>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
