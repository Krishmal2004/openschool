import { useMemo, useState } from "react";
import { Layers, Add, UserFollow } from "@carbon/icons-react";
import {
  Button,
  TextInput,
  Tag,
  ComposedModal,
  ModalHeader,
  ModalBody,
  ModalFooter,
  SkeletonText,
} from "@carbon/react";
import {
  useStreams,
  useStreamGroups,
  useCreateStream,
  useCreateStreamGroup,
  useCurrentClasses,
} from "../../../queries/useClasses";
import { useCurrentAcademicYear } from "../../../queries/useAcademicYears";
import { useTeachers } from "../../../queries/useTeachers";
import {
  useSectionHeads,
  useAssignSectionHead,
  useRemoveSectionHead,
} from "../../../queries/useSectionHeads";
import ErrorMessage from "../../../components/common/ErrorMessage";
import EmptyState from "../../../components/common/EmptyState";
import AgentFindingsBanner from "../../../components/common/AgentFindingsBanner";
import EntityCombobox from "../../../components/common/EntityCombobox";
import ConfirmDeleteModal from "../../../components/common/ConfirmDeleteModal";
import type { Stream } from "../../../services/stream";
import type { SectionHead } from "../../../services/sectionHead";
import MutationErrorNotification from "../../../components/common/MutationErrorNotification";

function StreamGroups({ stream }: { stream: Stream }) {
  const { data: groups, isLoading } = useStreamGroups(stream.id);
  const createGroup = useCreateStreamGroup();
  const [name, setName] = useState("");

  const handleAdd = () => {
    if (!name.trim()) return;
    createGroup.mutate(
      { streamId: stream.id, data: { name: name.trim() } },
      { onSuccess: () => setName("") },
    );
  };

  return (
    <div style={{ paddingLeft: "1.5rem" }}>
      <div style={{ display: "flex", flexWrap: "wrap", gap: "0.5rem", marginBottom: "0.5rem" }}>
        {isLoading ? (
          <SkeletonText width="30%" />
        ) : groups && groups.length > 0 ? (
          groups.map((g) => (
            <Tag key={g.id} type="teal" size="sm">
              {g.name}
            </Tag>
          ))
        ) : (
          <span style={{ fontSize: "0.75rem", color: "var(--os-text-tertiary)" }}>No sub-groups</span>
        )}
      </div>
      <div style={{ display: "flex", gap: "0.5rem", alignItems: "flex-end" }}>
        <TextInput
          id={`new-group-${stream.id}`}
          labelText=""
          placeholder="e.g. Physical Science"
          size="sm"
          value={name}
          onChange={(e) => setName(e.target.value)}
          style={{ maxWidth: "14rem" }}
        />
        <Button kind="ghost" size="sm" renderIcon={Add} disabled={!name.trim() || createGroup.isPending} onClick={handleAdd}>
          Add group
        </Button>
      </div>
    </div>
  );
}

export default function Streams() {
  const { data: streams, isLoading, isError, refetch } = useStreams();
  const createStream = useCreateStream();
  const { data: currentYear } = useCurrentAcademicYear();
  const { data: classes } = useCurrentClasses();
  const { data: teachers } = useTeachers();
  const { data: sectionHeads } = useSectionHeads(currentYear?.id ?? "");
  const assignSectionHead = useAssignSectionHead();
  const removeSectionHead = useRemoveSectionHead();

  const [createOpen, setCreateOpen] = useState(false);
  const [newStreamName, setNewStreamName] = useState("");
  const [toRemove, setToRemove] = useState<SectionHead | null>(null);

  const handleCreateStream = () => {
    if (!newStreamName.trim()) return;
    createStream.mutate(
      { name: newStreamName.trim() },
      { onSuccess: () => { setNewStreamName(""); setCreateOpen(false); } },
    );
  };

  // Derive section-head rows from the classes that actually exist this
  // year: a grade with any class carrying a stream_id gets one row per
  // stream in use; every other grade gets a single whole-grade row.
  const sectionHeadRows = useMemo(() => {
    type Row = { key: string; gradeId: string; gradeName: string; streamId: string | null; streamName: string | null };
    const byGrade = new Map<string, { gradeName: string; streamIds: Map<string, string> }>();
    for (const c of classes ?? []) {
      if (!byGrade.has(c.grade_id)) byGrade.set(c.grade_id, { gradeName: c.grade_name, streamIds: new Map() });
      if (c.stream_id) {
        const streamName = streams?.find((s) => s.id === c.stream_id)?.name ?? "Stream";
        byGrade.get(c.grade_id)!.streamIds.set(c.stream_id, streamName);
      }
    }
    const rows: Row[] = [];
    for (const [gradeId, { gradeName, streamIds }] of byGrade) {
      if (streamIds.size === 0) {
        rows.push({ key: gradeId, gradeId, gradeName, streamId: null, streamName: null });
      } else {
        for (const [streamId, streamName] of streamIds) {
          rows.push({ key: `${gradeId}-${streamId}`, gradeId, gradeName, streamId, streamName });
        }
      }
    }
    return rows.sort((a, b) => a.gradeName.localeCompare(b.gradeName) || (a.streamName ?? "").localeCompare(b.streamName ?? ""));
  }, [classes, streams]);

  const currentHeadFor = (gradeId: string, streamId: string | null) =>
    sectionHeads?.find((sh) => sh.grade_id === gradeId && sh.stream_id === streamId);

  const handleAssign = (gradeId: string, streamId: string | null, teacherId: string) => {
    if (!currentYear || !teacherId) return;
    assignSectionHead.mutate({
      academic_year_id: currentYear.id,
      grade_id: gradeId,
      stream_id: streamId,
      teacher_id: teacherId,
    });
  };

  return (
    <div className="os-page">
      <div className="os-page__header">
        <div className="os-page__header-left">
          <h1 className="os-page__title">Streams & Section Heads</h1>
          <p className="os-page__subtitle">
            A/L streams (with their sub-groups) and the teacher-in-charge for each grade or stream.
          </p>
        </div>
        <Button renderIcon={Add} kind="primary" size="md" onClick={() => setCreateOpen(true)}>
          New Stream
        </Button>
      </div>

      <AgentFindingsBanner titles={["Streams with no current-year classes"]} />

      <div className="os-section">
        <div className="os-section__header">
          <h2 className="os-section__title">Streams</h2>
        </div>

        {isLoading && (
          <div style={{ padding: "1.25rem 1.5rem" }}>
            <SkeletonText width="40%" />
          </div>
        )}
        {isError && <ErrorMessage message="Could not load streams." onRetry={refetch} />}

        {!isLoading && !isError && (!streams || streams.length === 0) && (
          <EmptyState
            title="No streams yet"
            description="Add A/L streams like Science, Commerce, Arts or Technology."
            action={
              <Button renderIcon={Add} kind="primary" onClick={() => setCreateOpen(true)}>
                New Stream
              </Button>
            }
          />
        )}

        {!isLoading && streams && streams.length > 0 && (
          <div>
            {streams.map((s) => (
              <div key={s.id} className="os-list-row" style={{ alignItems: "flex-start", flexDirection: "column", padding: "1rem 1.5rem" }}>
                <div style={{ display: "flex", alignItems: "center", gap: "0.5rem", marginBottom: "0.5rem" }}>
                  <Layers size={16} style={{ fill: "var(--os-accent)" }} />
                  <span style={{ fontWeight: 600, fontSize: "0.9rem", color: "var(--os-text-primary)" }}>{s.name}</span>
                </div>
                <StreamGroups stream={s} />
              </div>
            ))}
          </div>
        )}
      </div>

      <div className="os-section">
        <div className="os-section__header">
          <h2 className="os-section__title">Section Heads (Teachers in Charge)</h2>
          {currentYear && (
            <span style={{ fontSize: "0.75rem", color: "var(--os-text-tertiary)" }}>{currentYear.label}</span>
          )}
        </div>

        {!currentYear ? (
          <EmptyState
            title="No current academic year"
            description="Set an academic year as current before assigning section heads."
          />
        ) : sectionHeadRows.length === 0 ? (
          <EmptyState
            title="No classes yet"
            description="Section heads are derived from the grades and streams your classes actually use."
          />
        ) : (
          <div>
            {sectionHeadRows.map((row) => {
              const head = currentHeadFor(row.gradeId, row.streamId);
              return (
                <div key={row.key} className="os-list-row" style={{ padding: "0.75rem 1.5rem" }}>
                  <UserFollow size={16} style={{ fill: "var(--os-text-tertiary)", flexShrink: 0 }} />
                  <div style={{ flex: 1, minWidth: 0 }}>
                    <p style={{ margin: 0, fontSize: "0.875rem", fontWeight: 500, color: "var(--os-text-primary)" }}>
                      {row.gradeName}
                      {row.streamName ? ` - ${row.streamName}` : ""}
                    </p>
                    {head && (
                      <p style={{ margin: 0, fontSize: "0.75rem", color: "var(--os-text-secondary)" }}>
                        Current: {head.teacher_name}
                      </p>
                    )}
                  </div>
                  <div style={{ width: "16rem" }}>
                    <EntityCombobox
                      id={`tic-${row.key}`}
                      items={teachers ?? []}
                      selectedId={head?.teacher_id ?? ""}
                      onSelect={(id) => handleAssign(row.gradeId, row.streamId, id)}
                      getId={(t) => t.id}
                      itemToString={(t) => `${t.full_name} - ${t.employee_number}`}
                      placeholder="Search teachers…"
                    />
                  </div>
                  {/* Assigning requires a teacher, so leaving the post vacant
                      is only possible by removing the appointment outright. */}
                  <Button
                    kind="danger--ghost"
                    size="sm"
                    onClick={() => head && setToRemove(head)}
                    disabled={!head || removeSectionHead.isPending}
                  >
                    Remove
                  </Button>
                </div>
              );
            })}
          </div>
        )}

        <MutationErrorNotification
          isError={assignSectionHead.isError}
          error={assignSectionHead.error}
          title="Could not assign section head"
          fallback="Please try again."
          onClose={() => assignSectionHead.reset()}
          style={{ margin: "0 1.5rem 1rem" }}
        />

        <MutationErrorNotification
          isError={removeSectionHead.isError}
          error={removeSectionHead.error}
          title="Could not remove section head"
          fallback="Please try again."
          onClose={() => removeSectionHead.reset()}
          style={{ margin: "0 1.5rem 1rem" }}
        />
      </div>

      <ConfirmDeleteModal
        open={!!toRemove}
        title="Remove section head"
        description={
          <>
            Remove <strong>{toRemove?.teacher_name}</strong> as teacher-in-charge
            of {toRemove?.grade_name}
            {toRemove?.stream_name ? ` - ${toRemove.stream_name}` : ""}? The post
            is left vacant, and they lose the notification reach the role grants.
          </>
        }
        isPending={removeSectionHead.isPending}
        onClose={() => setToRemove(null)}
        onConfirm={() => {
          if (!toRemove || !currentYear) return;
          removeSectionHead.mutate(
            { id: toRemove.id, academicYearId: currentYear.id },
            { onSettled: () => setToRemove(null) },
          );
        }}
      />

      <ComposedModal open={createOpen} size="sm" onClose={() => setCreateOpen(false)}>
        <ModalHeader title="New stream" />
        <ModalBody>
          <MutationErrorNotification
            isError={createStream.isError}
            error={createStream.error}
            fallback="Failed to create stream"
            style={{ marginBottom: "1rem" }}
          />
          <TextInput
            id="new-stream-name"
            labelText="Stream name"
            placeholder="e.g. Science, Commerce, Arts, Technology"
            value={newStreamName}
            onChange={(e) => setNewStreamName(e.target.value)}
          />
        </ModalBody>
        <ModalFooter>
          <Button kind="secondary" onClick={() => setCreateOpen(false)}>
            Cancel
          </Button>
          <Button kind="primary" onClick={handleCreateStream} disabled={!newStreamName.trim() || createStream.isPending}>
            {createStream.isPending ? "Creating…" : "Create"}
          </Button>
        </ModalFooter>
      </ComposedModal>
    </div>
  );
}
