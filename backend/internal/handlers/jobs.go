package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/openschool-org/openschool/db/sqlc"
	"github.com/openschool-org/openschool/internal/jobs"
	"github.com/openschool-org/openschool/internal/models"
	"github.com/openschool-org/openschool/internal/repositories"
)

// JobsHandler exposes the admin Automation panel's endpoints for the background agent scheduler.
type JobsHandler struct {
	scheduler *jobs.Scheduler
	settings  *repositories.JobSchedulerRepository
}

// NewJobsHandler constructs a JobsHandler with its scheduler and settings dependencies.
func NewJobsHandler(scheduler *jobs.Scheduler, settings *repositories.JobSchedulerRepository) *JobsHandler {
	return &JobsHandler{scheduler: scheduler, settings: settings}
}

// knownJobNames returns the set of currently-registered agent names, used to reject settings for an agent that no longer exists.
func (h *JobsHandler) knownJobNames() map[string]bool {
	names := make(map[string]bool, len(h.scheduler.Jobs()))
	for _, j := range h.scheduler.Jobs() {
		names[j.Name()] = true
	}
	return names
}

// toLastRun converts a generated job-run row into its JSON response shape.
func toLastRun(r db.JobRun) *models.JobLastRun {
	lastRun := &models.JobLastRun{
		Status:    r.Status,
		Findings:  r.Findings,
		StartedAt: r.StartedAt.Time,
	}
	if r.Summary.Valid {
		lastRun.Summary = r.Summary.String
	}
	if r.FinishedAt.Valid {
		finishedAt := r.FinishedAt.Time
		lastRun.FinishedAt = &finishedAt
	}
	return lastRun
}

// List lists every registered background job and its current state.
func (h *JobsHandler) List(c *gin.Context) {
	ctx := c.Request.Context()
	jobList := h.scheduler.Jobs()

	settingsRows, err := h.settings.ListSettings(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	enabledByName := make(map[string]bool, len(settingsRows))
	for _, s := range settingsRows {
		enabledByName[s.JobName] = s.Enabled
	}

	runRows, err := h.settings.ListLatestRuns(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	lastRunByName := make(map[string]db.JobRun, len(runRows))
	for _, r := range runRows {
		lastRunByName[r.JobName] = r
	}

	out := make([]models.JobStatus, 0, len(jobList))
	for _, j := range jobList {
		enabled, ok := enabledByName[j.Name()]
		if !ok {
			enabled = true // no row yet — see JobSchedulerRepository.IsEnabled
		}
		status := models.JobStatus{
			Name: j.Name(), Description: j.Description(), Schedule: j.Schedule(), Enabled: enabled,
		}
		if r, ok := lastRunByName[j.Name()]; ok {
			status.LastRun = toLastRun(r)
		}
		out = append(out, status)
	}

	c.JSON(http.StatusOK, out)
}

// SetEnabled enables or disables a background agent.
func (h *JobsHandler) SetEnabled(c *gin.Context) {
	name := c.Param("name")
	if !h.knownJobNames()[name] {
		c.JSON(http.StatusNotFound, gin.H{"error": "unknown job"})
		return
	}

	var req models.SetJobEnabledRequest
	if err := bindStrict(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// The system-health agent is the one exception to "every agent is
	// safely optional": disabling it silently stops the school's only
	// backup mechanism, with no other symptom until data loss during a
	// real incident. Blocked outright rather than just warned about, both
	// here and in the Automation UI (which never renders its toggle).
	if !req.Enabled && name == jobs.SystemHealthAgentName {
		c.JSON(http.StatusBadRequest, gin.H{"error": "the system-health (backup) agent cannot be disabled"})
		return
	}

	if _, err := h.settings.SetEnabled(c.Request.Context(), name, req.Enabled); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// RunNow runs a background agent immediately, outside its schedule, honoring its enabled/disabled setting.
func (h *JobsHandler) RunNow(c *gin.Context) {
	name := c.Param("name")
	if !h.knownJobNames()[name] {
		c.JSON(http.StatusNotFound, gin.H{"error": "unknown job"})
		return
	}

	result, err := h.scheduler.RunNow(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "summary": result.Summary, "findings": result.Findings})
		return
	}
	c.JSON(http.StatusOK, gin.H{"summary": result.Summary, "findings": result.Findings})
}
