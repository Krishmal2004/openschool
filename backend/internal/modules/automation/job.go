// Package automation owns the scheduled maintenance and operational checks
// that support OpenSchool without being part of a user-facing workflow.
// These are deterministic checking algorithms, not AI agents.
//
// Each algorithm has one clearly named file:
//   - system_health.go       — backup, retention, and migration checks
//   - structural_integrity.go — academic structure and roster invariants
//   - people_compliance.go   — guardian, employment, and onboarding checks
//   - academic_delivery.go   — attendance and marks-delivery checks
//   - security_audit.go      — audit anomalies and token housekeeping
//
// The package also owns route handling, scheduling, persisted settings and
// run history. internal/app wires it once and cmd/api controls Start and Stop.
package automation

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
