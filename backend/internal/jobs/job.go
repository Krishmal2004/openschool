// Package jobs holds the scheduled, read-mostly maintenance/ops agents
// described in docs/plan.md § Proposed — maintenance/ops agents: built-in
// background checks (backup, data-invariant scans, stale-record watchers,
// …) that support the system's operation without any user-facing feature
// depending on them. None of the app's other code imports this package —
// it's wired in exactly once, from routes.Setup.
//
// Structure, one agent per file:
//   - job.go       — the Job contract every agent implements (this file)
//   - registry.go  — BuildAll, the single place that lists every agent
//   - scheduler.go — the generic cron runner + job_runs bookkeeping
//   - system.go    — shared helpers (severity, notifyAdmins, runChecks)
//   - agent_*.go   — one agent's checks, self-contained
//
// Adding a wholly new agent means adding one agent_*.go file implementing
// Job, plus one line in registry.go's BuildAll — nothing else changes.
package jobs

import "context"

// Result is what a job's Run returns: a one-line Summary and a Findings count for the Automation panel.
type Result struct {
	Summary  string
	Findings int
}

// Job is a single scheduled agent; Run must tolerate running alongside any other Job, just never alongside itself.
type Job interface {
	// Name is the stable on-disk identifier stored in job_settings/job_runs — renaming it orphans existing history.
	Name() string

	// Schedule is a standard 5-field cron expression (robfig/cron/v3, no seconds field).
	Schedule() string

	// Description is a one-line explanation shown next to the Automation panel's toggle.
	Description() string

	Run(ctx context.Context) (Result, error)
}
