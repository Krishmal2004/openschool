import { Toggle, Tag, Button, SkeletonText } from "@carbon/react";
import { Play } from "@carbon/icons-react";
import { useJobs, useSetJobEnabled, useRunJobNow } from "../../../queries/useJobs";
import ErrorMessage from "../../../components/common/ErrorMessage";
import MutationErrorNotification from "../../../components/common/MutationErrorNotification";
import type { JobRunStatus } from "../../../services/jobs";

function humanizeJobName(name: string) {
  return name.replace(/_/g, " ").replace(/\b\w/g, (c) => c.toUpperCase());
}

// Agent schedules are fixed cron expressions set in code; this just reads friendlier than raw cron next to each name.
const SCHEDULE_LABELS: Record<string, string> = {
  "0 * * * *": "Hourly",
  "0 2 * * *": "Daily at 2:00 AM",
  "0 3 * * *": "Daily at 3:00 AM",
  "0 5 * * *": "Daily at 5:00 AM",
  "0 12 * * 1-5": "Weekdays at 12:00 PM",
};

function humanizeSchedule(cron: string) {
  return SCHEDULE_LABELS[cron] ?? cron;
}

// Can't be disabled — it's the school's only backup mechanism (see internal/handlers/jobs.go's SetEnabled).
const NON_DISABLEABLE_JOBS = new Set(["system_health_agent"]);

function statusTag(status: JobRunStatus) {
  switch (status) {
    case "ok":
      return <Tag type="green" size="sm">OK</Tag>;
    case "failed":
      return <Tag type="red" size="sm">Failed</Tag>;
    case "running":
      return <Tag type="blue" size="sm">Running</Tag>;
  }
}

export default function Automation() {
  const { data: jobs, isLoading, isError, refetch } = useJobs();
  const setEnabled = useSetJobEnabled();
  const runNow = useRunJobNow();

  return (
    <div className="os-page">
      <div className="os-page__header">
        <div className="os-page__header-left">
          <h1 className="os-page__title">Automation</h1>
          <p className="os-page__subtitle">
            Five scheduled background agents that support the system's
            operation — none of the app's other features depend on them, so
            any of these can be turned off safely, except System Health
            (backup).
          </p>
        </div>
      </div>

      {isError && (
        <div style={{ marginBottom: "1.5rem" }}>
          <ErrorMessage message="Could not load background agents." onRetry={refetch} />
        </div>
      )}

      <MutationErrorNotification
        isError={setEnabled.isError || runNow.isError}
        error={setEnabled.error ?? runNow.error}
        title="Action failed"
        fallback="Please try again."
        onClose={() => {
          setEnabled.reset();
          runNow.reset();
        }}
        style={{ marginBottom: "1.5rem" }}
      />

      <div className="os-section" style={{ marginTop: 0 }}>
        {isLoading && (
          <div style={{ padding: "1.25rem 1.5rem" }}>
            <SkeletonText width="60%" />
          </div>
        )}

        {!isLoading &&
          (jobs ?? []).map((job, i) => (
            <div
              key={job.name}
              style={{
                display: "flex",
                alignItems: "flex-start",
                justifyContent: "space-between",
                gap: "1.5rem",
                padding: "1rem 1.5rem",
                borderBottom: i < (jobs?.length ?? 0) - 1 ? "1px solid var(--os-border-subtle)" : "none",
                flexWrap: "wrap",
              }}
            >
              <div style={{ flex: 1, minWidth: "20rem" }}>
                <div style={{ display: "flex", alignItems: "center", gap: "0.625rem", flexWrap: "wrap" }}>
                  <span style={{ fontWeight: 600, fontSize: "0.875rem" }}>{humanizeJobName(job.name)}</span>
                  <Tag type="blue" size="sm">{humanizeSchedule(job.schedule)}</Tag>
                </div>
                <p style={{ margin: "0.25rem 0 0", fontSize: "0.8125rem", color: "var(--os-text-secondary)" }}>{job.description}</p>

                {job.last_run ? (
                  <div style={{ display: "flex", alignItems: "center", gap: "0.5rem", marginTop: "0.625rem", flexWrap: "wrap" }}>
                    {statusTag(job.last_run.status)}
                    {job.last_run.findings > 0 && (
                      <Tag type="magenta" size="sm">{job.last_run.findings} finding{job.last_run.findings === 1 ? "" : "s"}</Tag>
                    )}
                    <span style={{ fontSize: "0.75rem", color: "var(--os-text-tertiary)" }}>
                      Last ran {new Date(job.last_run.started_at).toLocaleString()}
                    </span>
                  </div>
                ) : (
                  <p style={{ margin: "0.625rem 0 0", fontSize: "0.75rem", color: "var(--os-text-tertiary)" }}>Never run yet</p>
                )}
                {job.last_run?.summary && (
                  <p style={{ margin: "0.375rem 0 0", fontSize: "0.8125rem", color: "var(--os-text-primary)" }}>{job.last_run.summary}</p>
                )}
              </div>

              <div style={{ display: "flex", alignItems: "center", gap: "1rem", flexShrink: 0 }}>
                <Button
                  kind="ghost"
                  size="sm"
                  renderIcon={Play}
                  disabled={!job.enabled || runNow.isPending}
                  onClick={() => runNow.mutate(job.name)}
                >
                  Run now
                </Button>
                {NON_DISABLEABLE_JOBS.has(job.name) ? (
                  <Tag type="gray" size="sm">Always on</Tag>
                ) : (
                  <Toggle
                    id={`job-toggle-${job.name}`}
                    size="sm"
                    labelText={`Enable ${humanizeJobName(job.name)}`}
                    hideLabel
                    toggled={job.enabled}
                    disabled={setEnabled.isPending}
                    onToggle={(checked) => setEnabled.mutate({ name: job.name, enabled: checked })}
                  />
                )}
              </div>
            </div>
          ))}
      </div>
    </div>
  );
}
