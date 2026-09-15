package automation

import (
	"context"
	"errors"
)

var (
	ErrUnknownJob             = errors.New("unknown job")
	ErrSystemHealthCannotStop = errors.New("the system-health (backup) agent cannot be disabled")
)

// Service coordinates the scheduler's registered jobs with their persisted settings and history.
type Service struct {
	scheduler *Scheduler
	settings  *Repository
}

func NewService(scheduler *Scheduler, settings *Repository) *Service {
	return &Service{scheduler: scheduler, settings: settings}
}

func (s *Service) List(ctx context.Context) ([]JobStatus, error) {
	settingsRows, err := s.settings.ListSettings(ctx)
	if err != nil {
		return nil, err
	}
	enabledByName := make(map[string]bool, len(settingsRows))
	for _, setting := range settingsRows {
		enabledByName[setting.JobName] = setting.Enabled
	}

	runRows, err := s.settings.ListLatestRuns(ctx)
	if err != nil {
		return nil, err
	}
	lastRunByName := make(map[string]JobLastRun, len(runRows))
	for _, run := range runRows {
		lastRun := JobLastRun{Status: run.Status, Summary: run.Summary, Findings: run.Findings, StartedAt: run.StartedAt, FinishedAt: run.FinishedAt}
		lastRunByName[run.JobName] = lastRun
	}

	out := make([]JobStatus, 0, len(s.scheduler.Jobs()))
	for _, job := range s.scheduler.Jobs() {
		enabled, ok := enabledByName[job.Name()]
		if !ok {
			enabled = true
		}
		status := JobStatus{Name: job.Name(), Description: job.Description(), Schedule: job.Schedule(), Enabled: enabled}
		if lastRun, ok := lastRunByName[job.Name()]; ok {
			status.LastRun = &lastRun
		}
		out = append(out, status)
	}
	return out, nil
}

func (s *Service) SetEnabled(ctx context.Context, name string, enabled bool) error {
	if !s.isKnown(name) {
		return ErrUnknownJob
	}
	if !enabled && name == SystemHealthAgentName {
		return ErrSystemHealthCannotStop
	}
	return s.settings.SetEnabled(ctx, name, enabled)
}

func (s *Service) RunNow(ctx context.Context, name string) (Result, error) {
	if !s.isKnown(name) {
		return Result{}, ErrUnknownJob
	}
	return s.scheduler.RunNow(ctx, name)
}

func (s *Service) isKnown(name string) bool {
	for _, job := range s.scheduler.Jobs() {
		if job.Name() == name {
			return true
		}
	}
	return false
}
