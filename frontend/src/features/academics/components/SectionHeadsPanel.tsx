import { useMemo, useState } from "react";
import { UserFollow } from "@carbon/icons-react";
import { Button } from "@carbon/react";
import { useCurrentClasses, useStreams } from "@/features/academics/queries/useClasses";
import { useCurrentAcademicYear } from "@/features/school/queries/useAcademicYears";
import { useTeachers } from "@/features/teachers/queries/useTeachers";
import { useSectionHeads, useAssignSectionHead, useRemoveSectionHead } from "@/features/teachers/queries/useSectionHeads";
import type { SectionHead } from "@/features/teachers/api/sectionHead";
import EmptyState from "@/shared/ui/EmptyState";
import EntityCombobox from "@/shared/ui/EntityCombobox";
import ConfirmDeleteModal from "@/shared/ui/ConfirmDeleteModal";
import MutationErrorNotification from "@/shared/ui/MutationErrorNotification";

interface Row {
  key: string;
  gradeId: string;
  gradeName: string;
  streamId: string | null;
  streamName: string | null;
}

// One row per grade, or per stream in a grade that has streamed classes.
function useSectionHeadRows() {
  const { data: classes } = useCurrentClasses();
  const { data: streams } = useStreams();
  return useMemo(() => {
    const byGrade = new Map<string, { gradeName: string; streamIds: Map<string, string> }>();
    for (const c of classes ?? []) {
      if (!byGrade.has(c.grade_id)) byGrade.set(c.grade_id, { gradeName: c.grade_name, streamIds: new Map() });
      if (c.stream_id) byGrade.get(c.grade_id)!.streamIds.set(c.stream_id, streams?.find((s) => s.id === c.stream_id)?.name ?? "Stream");
    }
    const rows: Row[] = [];
    for (const [gradeId, { gradeName, streamIds }] of byGrade) {
      if (streamIds.size === 0) rows.push({ key: gradeId, gradeId, gradeName, streamId: null, streamName: null });
      else for (const [streamId, streamName] of streamIds) rows.push({ key: `${gradeId}-${streamId}`, gradeId, gradeName, streamId, streamName });
    }
    return rows.sort((a, b) => a.gradeName.localeCompare(b.gradeName) || (a.streamName ?? "").localeCompare(b.streamName ?? ""));
  }, [classes, streams]);
}

export default function SectionHeadsPanel() {
  const { data: currentYear } = useCurrentAcademicYear();
  const { data: teachers } = useTeachers();
  const { data: sectionHeads } = useSectionHeads(currentYear?.id ?? "");
  const assign = useAssignSectionHead();
  const remove = useRemoveSectionHead();
  const rows = useSectionHeadRows();
  const [toRemove, setToRemove] = useState<SectionHead | null>(null);

  const headFor = (row: Row) => sectionHeads?.find((sh) => sh.grade_id === row.gradeId && sh.stream_id === row.streamId);

  return (
    <div className="os-section">
      <div className="os-section__header">
        <h2 className="os-section__title">Section Heads (Teachers in Charge)</h2>
        {currentYear && <span className="os-text-xs os-c-tertiary">{currentYear.label}</span>}
      </div>

      {!currentYear ? (
        <EmptyState title="No current academic year" description="Set an academic year as current before assigning section heads." />
      ) : rows.length === 0 ? (
        <EmptyState title="No classes yet" description="Section heads are derived from the grades and streams your classes actually use." />
      ) : (
        rows.map((row) => {
          const head = headFor(row);
          return (
            <div key={row.key} className="os-list-row os-py-3 os-px-6">
              <UserFollow size={16} className="os-fill-tertiary os-shrink-0" />
              <div className="os-flex-1 os-min-w-0">
                <p className="os-m-0 os-text-md os-fw-500 os-c-primary">{row.gradeName}{row.streamName ? ` - ${row.streamName}` : ""}</p>
                {head && <p className="os-m-0 os-text-xs os-c-secondary">Current: {head.teacher_name}</p>}
              </div>
              <div className="os-w-16">
                <EntityCombobox
                  id={`tic-${row.key}`}
                  items={teachers ?? []}
                  selectedId={head?.teacher_id ?? ""}
                  onSelect={(teacher_id) => teacher_id && assign.mutate({ academic_year_id: currentYear.id, grade_id: row.gradeId, stream_id: row.streamId, teacher_id })}
                  getId={(t) => t.id}
                  itemToString={(t) => `${t.full_name} - ${t.employee_number}`}
                  placeholder="Search teachers…"
                />
              </div>
              {/* Vacating a post is only possible by removing the appointment. */}
              <Button kind="danger--ghost" size="sm" onClick={() => head && setToRemove(head)} disabled={!head || remove.isPending}>Remove</Button>
            </div>
          );
        })
      )}

      <MutationErrorNotification isError={assign.isError} error={assign.error} title="Could not assign section head" fallback="Please try again." onClose={() => assign.reset()} className="os-mx-6" />
      <MutationErrorNotification isError={remove.isError} error={remove.error} title="Could not remove section head" fallback="Please try again." onClose={() => remove.reset()} className="os-mx-6" />

      <ConfirmDeleteModal
        open={!!toRemove}
        title="Remove section head"
        description={
          <>
            Remove <strong>{toRemove?.teacher_name}</strong> as teacher-in-charge of {toRemove?.grade_name}
            {toRemove?.stream_name ? ` - ${toRemove.stream_name}` : ""}? The post is left vacant, and they lose the notification reach the role grants.
          </>
        }
        isPending={remove.isPending}
        onClose={() => setToRemove(null)}
        onConfirm={() => toRemove && currentYear && remove.mutate({ id: toRemove.id, academicYearId: currentYear.id }, { onSettled: () => setToRemove(null) })}
      />
    </div>
  );
}
