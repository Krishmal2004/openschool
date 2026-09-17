import { SkeletonText, Tag } from "@carbon/react";
import { useAuditLogs } from "@/features/system/queries/useAuditLogs";
import ErrorMessage from "@/shared/ui/ErrorMessage";
import EmptyState from "@/shared/ui/EmptyState";
import AgentFindingsBanner from "@/features/notifications/components/AgentFindingsBanner";

function formatEntity(entityType: string) {
  return entityType.replace(/_/g, " ");
}

function formatAction(action: string) {
  return action.replace(/_/g, " ");
}

export default function AuditLog() {
  const { data: logs, isLoading, isError, refetch } = useAuditLogs();

  return (
    <div>
      <div className="os-page__header">
        <div>
          <h2 className="os-section__title os-m-0">
            Audit Log
          </h2>
          <p className="os-page__subtitle os-mt-1">
            Manual house re-assignments and admin edits made to attendance
            after its 24-hour lock - most recent first.
          </p>
        </div>
      </div>

      <AgentFindingsBanner titles={["Unusual audit-log activity", "Unusual off-hours account activity"]} />

      {isError && <ErrorMessage message="Could not load the audit log." onRetry={refetch} />}

      <div className="os-section">
        {isLoading && (
          <div>
            {Array.from({ length: 4 }).map((_, i) => (
              <div
                key={i} className="os-flex os-gap-4 os-py-3h os-px-6 os-border-b"
              >
                <SkeletonText width="10rem" />
                <SkeletonText width="20%" />
              </div>
            ))}
          </div>
        )}

        {!isLoading && !isError && (logs?.length ?? 0) === 0 && (
          <EmptyState
            title="No audit entries yet"
            description="Manual house changes and attendance edits made after the 24-hour lock will show up here."
          />
        )}

        {!isLoading && logs && logs.length > 0 && (
          <div>
            {logs.map((entry, i) => (
              <div
                key={entry.id} className={`os-py-3h os-px-6 ${i < logs.length - 1 ? "os-border-b" : ""}`}
              >
                <div className="os-flex os-items-center os-gap-3 os-wrap">
                  <Tag type="cool-gray" size="sm">
                    {formatEntity(entry.entity_type)}
                  </Tag>
                  <span className="os-fw-600 os-text-md">
                    {formatAction(entry.action)}
                  </span>
                  <div className="os-flex-1" />
                  <span className="os-text-xs os-c-tertiary">
                    {new Date(entry.created_at).toLocaleString()}
                  </span>
                </div>
                <p className="os-mt-1h os-mx-0 os-mb-0 os-text-sm os-c-secondary">
                  By {entry.actor_name ?? "unknown"}
                  {entry.reason ? ` - "${entry.reason}"` : ""}
                </p>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
