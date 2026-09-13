package services

import (
	"context"
	"errors"

	"github.com/openschool-org/openschool/internal/jobs"
	"github.com/openschool-org/openschool/internal/models"
	"github.com/openschool-org/openschool/internal/repositories"
)

var (
	ErrUnknownJob             = errors.New("unknown job")
	ErrSystemHealthCannotStop = errors.New("the system-health (backup) agent cannot be disabled")
)

// JobsService coordinates the scheduler's registered jobs with their persisted settings and history.
type JobsService struct {
	scheduler *jobs.Scheduler
	settings  *repositories.JobSchedulerRepository
}

func NewJobsService(scheduler *jobs.Scheduler, settings *repositories.JobSchedulerRepository) *JobsService {
	return &JobsService{scheduler: scheduler, settings: settings}
}

func (s *JobsService) List(ctx context.Context) ([]models.JobStatus, error) {
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
	lastRunByName := make(map[string]models.JobLastRun, len(runRows))
	for _, run := range runRows {
		lastRun := models.JobLastRun{Status: run.Status, Findings: run.Findings, StartedAt: run.StartedAt.Time}
		if run.Summary.Valid {
			lastRun.Summary = run.Summary.String
		}
		if run.FinishedAt.Valid {
			finishedAt := run.FinishedAt.Time
			lastRun.FinishedAt = &finishedAt
		}
		lastRunByName[run.JobName] = lastRun
	}

	out := make([]models.JobStatus, 0, len(s.scheduler.Jobs()))
	for _, job := range s.scheduler.Jobs() {
		enabled, ok := enabledByName[job.Name()]
		if !ok {
			enabled = true
		}
		status := models.JobStatus{Name: job.Name(), Description: job.Description(), Schedule: job.Schedule(), Enabled: enabled}
		if lastRun, ok := lastRunByName[job.Name()]; ok {
			status.LastRun = &lastRun
		}
		out = append(out, status)
	}
	return out, nil
}

func (s *JobsService) SetEnabled(ctx context.Context, name string, enabled bool) error {
	if !s.isKnown(name) {
		return ErrUnknownJob
	}
	if !enabled && name == jobs.SystemHealthAgentName {
		return ErrSystemHealthCannotStop
	}
	_, err := s.settings.SetEnabled(ctx, name, enabled)
	return err
}

func (s *JobsService) RunNow(ctx context.Context, name string) (jobs.Result, error) {
	if !s.isKnown(name) {
		return jobs.Result{}, ErrUnknownJob
	}
	return s.scheduler.RunNow(ctx, name)
}

func (s *JobsService) isKnown(name string) bool {
	for _, job := range s.scheduler.Jobs() {
		if job.Name() == name {
			return true
		}
	}
	return false
}
