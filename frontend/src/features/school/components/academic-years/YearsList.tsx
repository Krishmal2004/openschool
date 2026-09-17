import { Calendar, Checkmark } from "@carbon/icons-react";
import { Button, Tag } from "@carbon/react";
import type { useAcademicYears, useSetCurrentAcademicYear } from "@/features/school/queries/useAcademicYears";
import type { AcademicYear } from "@/features/school/api/academicYear";
import ErrorMessage from "@/shared/ui/ErrorMessage";
import EmptyState from "@/shared/ui/EmptyState";
import SectionHeader from "@/shared/ui/SectionHeader";
import ListRowSkeleton from "@/shared/ui/ListRowSkeleton";

function formatDate(iso: string | null) {
  if (!iso) return "-";
  // Built from the parsed Y/M/D as a local date, not `new Date(iso)`, the latter parses a date-only string as UTC midnight, which shifts to the previous day in negative-UTC timezones.
  const [y, m, d] = iso.split("-").map(Number);
  return new Date(y, m - 1, d).toLocaleDateString("en-LK", {
    month: "short",
    year: "numeric",
  });
}

interface Props {
  years: ReturnType<typeof useAcademicYears>["data"];
  isLoading: boolean;
  isError: boolean;
  refetch: () => void;
  setCurrent: ReturnType<typeof useSetCurrentAcademicYear>;
  onOpenTerms: (year: AcademicYear) => void;
  onRequestDelete: (year: AcademicYear) => void;
}

export default function YearsList({
  years,
  isLoading,
  isError,
  refetch,
  setCurrent,
  onOpenTerms,
  onRequestDelete,
}: Props) {
  return (
    <div className="os-section">
      <SectionHeader
        title="Academic Years"
        meta={years && <span className="os-section__meta">{years.length} total</span>}
      />

      {isLoading && (
        <div>
          {Array.from({ length: 3 }).map((_, i) => (
            <ListRowSkeleton key={i} titleWidth="25%" subtitleWidth="40%" trailingWidth="5rem" />
          ))}
        </div>
      )}
      {isError && <ErrorMessage message="Could not load academic years." onRetry={refetch} />}

      {!isLoading && !isError && years?.length === 0 && (
        <EmptyState title="No academic years" description="Create the first academic year to get started." />
      )}

      {!isLoading && years && years.length > 0 && (
        <div>
          {years.map((y) => (
            <div key={y.id} className="os-list-row">
              <Calendar size={20} className={`${y.is_current ? "os-fill-accent" : "os-fill-tertiary"} os-shrink-0`} />
              <div className="os-flex-1">
                <p className="os-mt-0 os-mx-0 os-mb-h os-fw-600 os-text-md os-c-primary">
                  {y.label}
                </p>
                <p className="os-m-0 os-text-xs os-c-secondary">
                  {formatDate(y.start_date)} - {formatDate(y.end_date)}
                </p>
              </div>
              <div className="os-flex os-items-center os-gap-2">
                <Tag type={y.is_current ? "teal" : "gray"} size="sm">
                  {y.is_current && <Checkmark size={12} className="os-mr-1" />}
                  {y.is_current ? "Current" : "Closed"}
                </Tag>
                <Button kind="ghost" size="sm" onClick={() => onOpenTerms(y)}>
                  Terms
                </Button>
                {!y.is_current && (
                  <Button
                    kind="ghost"
                    size="sm"
                    onClick={() => setCurrent.mutate(y.id)}
                    disabled={setCurrent.isPending}
                  >
                    Set Current
                  </Button>
                )}
                {!y.is_current && (
                  <Button kind="danger--ghost" size="sm" onClick={() => onRequestDelete(y)}>
                    Delete
                  </Button>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
